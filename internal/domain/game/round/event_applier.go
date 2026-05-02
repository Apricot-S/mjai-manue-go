package round

import (
	"fmt"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type EventApplier interface {
	Apply(ev event.Event) error
}

func (s *State) Apply(ev event.Event) error {
	switch ev := ev.(type) {
	case *event.Draw:
		return s.applyDraw(ev)
	case *event.Discard:
		return s.applyDiscard(ev)
	case *event.Chii:
		return s.applyChii(ev)
	case *event.Pon:
		return s.applyPon(ev)
	case *event.CalledKan:
		return s.applyCalledKan(ev)
	case *event.ConcealedKan:
		return s.applyConcealedKan(ev)
	case *event.PromotedKan:
		return s.applyPromotedKan(ev)
	case *event.Dora:
		return s.applyDora(ev)
	case *event.Riichi:
		return s.applyRiichi(ev)
	case *event.RiichiAccepted:
		return s.applyRiichiAccepted(ev)
	case *event.Win:
		return s.applyWin(ev)
	case *event.DrawRound:
		return s.applyDrawRound(ev)
	default:
		return fmt.Errorf("unknown event: %T", ev)
	}
}

func (s *State) applyDraw(ev *event.Draw) error {
	if ev.Actor() != s.nextDraw {
		return fmt.Errorf("cannot Draw: actor %d is not next draw player %d", ev.Actor().Index(), s.nextDraw.Index())
	}
	if s.numLeftTiles <= 0 {
		return fmt.Errorf("cannot Draw: no tiles left")
	}

	actorSeat := ev.Actor()
	p := s.players[actorSeat.Index()]
	if err := p.Draw(ev.Tile()); err != nil {
		return err
	}

	s.numLeftTiles--
	s.pendingDiscard = &actorSeat
	return nil
}

func (s *State) applyDiscard(ev *event.Discard) error {
	actorSeat := ev.Actor()
	if s.pendingDiscard == nil || *s.pendingDiscard != actorSeat {
		return fmt.Errorf("cannot Discard: actor %d is not pending discard player", actorSeat.Index())
	}
	p := s.players[actorSeat.Index()]
	if err := p.Discard(ev.Tile(), ev.Tsumogiri()); err != nil {
		return err
	}
	s.pendingDiscard = nil
	s.nextDraw = *seat.MustSeat((actorSeat.Index() + 1) % common.NumPlayers)
	return nil
}

func (s *State) applyChii(ev *event.Chii) error {
	if !ev.Actor().IsShimochaOf(ev.Target()) {
		return fmt.Errorf("cannot Chii: actor %d is not shimocha of target %d", ev.Actor().Index(), ev.Target().Index())
	}

	chii, err := meld.NewChii(ev.Taken(), ev.Consumed(), ev.Target())
	if err != nil {
		return err
	}
	return s.applyOpenCall(ev.Actor(), ev.Target(), ev.Taken(), func() error {
		return s.players[ev.Actor().Index()].Chii(*chii)
	})
}

func (s *State) applyPon(ev *event.Pon) error {
	pon, err := meld.NewPon(ev.Taken(), ev.Consumed(), ev.Target())
	if err != nil {
		return err
	}
	return s.applyOpenCall(ev.Actor(), ev.Target(), ev.Taken(), func() error {
		return s.players[ev.Actor().Index()].Pon(*pon)
	})
}

func (s *State) applyCalledKan(ev *event.CalledKan) error {
	kan, err := meld.NewCalledKan(ev.Taken(), ev.Consumed(), ev.Target())
	if err != nil {
		return err
	}
	if err := s.applyOpenCall(ev.Actor(), ev.Target(), ev.Taken(), func() error {
		return s.players[ev.Actor().Index()].CalledKan(*kan)
	}); err != nil {
		return err
	}
	s.pendingDoraReveal = true
	return nil
}

func (s *State) applyConcealedKan(ev *event.ConcealedKan) error {
	kan, err := meld.NewConcealedKan(ev.Consumed())
	if err != nil {
		return err
	}
	if err := s.players[ev.Actor().Index()].ConcealedKan(*kan); err != nil {
		return err
	}
	s.pendingDoraReveal = true
	return nil
}

func (s *State) applyPromotedKan(ev *event.PromotedKan) error {
	actor := s.players[ev.Actor().Index()]
	added := ev.Added()
	ponIndex := slices.IndexFunc(actor.Melds(), func(m meld.Meld) bool {
		pon, ok := m.(*meld.Pon)
		return ok && pon.Taken().HasSameSymbol(&added)
	})
	if ponIndex == -1 {
		return fmt.Errorf("cannot PromotedKan: failed to find pon for added tile %s", ev.Added())
	}

	pon := actor.Melds()[ponIndex].(*meld.Pon)
	kan, err := meld.NewPromotedKan(*pon.Taken(), [2]tile.Tile(pon.Consumed()), ev.Added(), *pon.Target())
	if err != nil {
		return err
	}
	if err := actor.PromotedKan(*kan); err != nil {
		return err
	}
	s.pendingDoraReveal = true
	return nil
}

func (s *State) applyDora(ev *event.Dora) error {
	if !s.pendingDoraReveal {
		return fmt.Errorf("cannot reveal dora indicator: not after kan")
	}
	if ev.Indicator().IsUnknown() {
		return fmt.Errorf("cannot reveal unknown dora indicator")
	}
	if len(s.doraIndicators) >= MaxNumDoraIndicators {
		return fmt.Errorf("cannot reveal dora indicator: already have %d indicators", len(s.doraIndicators))
	}
	s.doraIndicators = append(s.doraIndicators, ev.Indicator())
	s.pendingDoraReveal = false
	return nil
}

func (s *State) applyRiichi(ev *event.Riichi) error {
	if s.numLeftTiles < common.NumPlayers {
		return fmt.Errorf("cannot Riichi: no next draw turn remains")
	}
	return s.players[ev.Actor().Index()].Riichi()
}

func (s *State) applyRiichiAccepted(ev *event.RiichiAccepted) error {
	if err := s.players[ev.Actor().Index()].RiichiAccepted(); err != nil {
		return err
	}
	s.applyRiichiAcceptedScoreUpdate(ev)
	s.riichiDeposit++
	return nil
}

func (s *State) applyRiichiAcceptedScoreUpdate(ev *event.RiichiAccepted) {
	if ev.Scores() != nil || ev.Deltas() != nil {
		s.applyScoreUpdate(ev.Scores(), ev.Deltas())
		return
	}
	s.scores[ev.Actor().Index()] -= 1000
}

func (s *State) applyWin(ev *event.Win) error {
	if !s.canApplyWin(ev) {
		return fmt.Errorf("cannot Win: invalid timing")
	}
	s.applyScoreUpdate(ev.Scores(), ev.Deltas())
	return nil
}

func (s *State) canApplyWin(ev *event.Win) bool {
	if ev.Actor() == ev.Target() {
		drawnTile := s.players[ev.Actor().Index()].DrawnTile()
		return drawnTile != nil && isTileMatchKnownEnough(drawnTile, ev.WinningTile())
	}

	targetRiver := s.players[ev.Target().Index()].River()
	return len(targetRiver) > 0 && isTileMatchKnownEnough(&targetRiver[len(targetRiver)-1], ev.WinningTile())
}

func isTileMatchKnownEnough(stateTile *tile.Tile, eventTile *tile.Tile) bool {
	if eventTile == nil {
		return true
	}
	// Invisible players draw "?", while visible players cannot draw "?" by invariant.
	// Therefore an unknown state tile means the event tile cannot be verified here,
	// but a visible player's known tile is still compared strictly.
	if stateTile.IsUnknown() {
		return true
	}
	return *stateTile == *eventTile
}

func (s *State) applyDrawRound(ev *event.DrawRound) error {
	s.applyScoreUpdate(ev.Scores(), ev.Deltas())
	return nil
}

func (s *State) applyScoreUpdate(scores *[common.NumPlayers]int, deltas *[common.NumPlayers]int) {
	if scores != nil {
		s.scores = *scores
		return
	}
	if deltas != nil {
		for i, delta := range deltas {
			s.scores[i] += delta
		}
	}
}

func (s *State) applyOpenCall(actorSeat, targetSeat seat.Seat, taken tile.Tile, applyActor func() error) error {
	if actorSeat == targetSeat {
		return fmt.Errorf("cannot call %s from self", taken)
	}
	if s.numLeftTiles <= 0 {
		return fmt.Errorf("cannot call %s: no tiles left", taken)
	}

	target := s.players[targetSeat.Index()]
	river := target.River()
	if len(river) == 0 {
		return fmt.Errorf("cannot call %s from player %d: target river is empty", taken, targetSeat.Index())
	}
	if river[len(river)-1] != taken {
		return fmt.Errorf("cannot call %s from player %d: last river tile is %s", taken, targetSeat.Index(), river[len(river)-1])
	}

	if err := applyActor(); err != nil {
		return err
	}
	if err := target.TakeFromRiver(taken); err != nil {
		return err
	}
	s.pendingDiscard = &actorSeat
	return nil
}
