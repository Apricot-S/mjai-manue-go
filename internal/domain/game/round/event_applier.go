package round

import (
	"fmt"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type EventApplier interface {
	Apply(ev event.Event) error
}

func (s *State) Apply(ev event.Event) error {
	if s.roundEnded {
		if _, ok := ev.(*event.Win); !ok || !s.roundEndedByWin {
			return fmt.Errorf("cannot apply %T: round already ended", ev)
		}
	}

	var err error
	switch ev := ev.(type) {
	case *event.Draw:
		err = s.applyDraw(ev)
	case *event.Discard:
		err = s.applyDiscard(ev)
	case *event.Chii:
		err = s.applyChii(ev)
	case *event.Pon:
		err = s.applyPon(ev)
	case *event.CalledKan:
		err = s.applyCalledKan(ev)
	case *event.ConcealedKan:
		err = s.applyConcealedKan(ev)
	case *event.PromotedKan:
		err = s.applyPromotedKan(ev)
	case *event.Dora:
		err = s.applyDora(ev)
	case *event.Riichi:
		err = s.applyRiichi(ev)
	case *event.RiichiAccepted:
		err = s.applyRiichiAccepted(ev)
	case *event.Win:
		err = s.applyWin(ev)
	case *event.DrawRound:
		err = s.applyDrawRound(ev)
	default:
		err = fmt.Errorf("unknown event: %T", ev)
	}

	if err != nil {
		return err
	}

	s.legalActionsCache = nil
	return nil
}

func (s *State) applyDraw(ev *event.Draw) error {
	if s.pendingDoraReveals > 0 && (s.kanProgress == noKanProgress || s.kanProgress == waitingReplacementAfterDora) {
		return fmt.Errorf("cannot Draw: dora indicator must be revealed first")
	}
	if ev.Actor() != s.nextDraw {
		return fmt.Errorf("cannot Draw: actor %d is not next draw player %d", ev.Actor().Index(), s.nextDraw.Index())
	}
	isReplacementTileDraw := s.isWaitingReplacementTileDraw(ev.Actor())
	if s.numLeftTiles <= 0 {
		return fmt.Errorf("cannot Draw: no tiles left")
	}

	actorSeat := ev.Actor()
	p := s.players[actorSeat.Index()]
	if err := p.Draw(ev.Tile()); err != nil {
		return err
	}

	s.numLeftTiles--
	if isReplacementTileDraw {
		s.pendingRobbedKanTile = nil
		s.kanProgress = noKanProgress
		if s.pendingDoraReveals == 0 {
			s.pendingDiscard = &actorSeat
			s.pendingKanActor = nil
		}
	} else {
		s.pendingDiscard = &actorSeat
	}
	s.lastActor = &actorSeat
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
	if p.RiichiState() == player.RiichiDeclared {
		s.pendingRiichiAcceptance = &actorSeat
	}
	s.canKyushukyuhai[actorSeat.Index()] = false
	s.pendingDiscard = nil
	s.nextDraw = *seat.MustSeat((actorSeat.Index() + 1) % common.NumPlayers)
	s.lastActor = &actorSeat
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
	actorSeat := ev.Actor()
	if err := s.applyOpenCall(actorSeat, ev.Target(), ev.Taken(), func() error {
		return s.players[actorSeat.Index()].Chii(*chii)
	}); err != nil {
		return err
	}
	s.disableKyushukyuhaiForAll()
	s.pendingDiscard = &actorSeat
	return nil
}

func (s *State) applyPon(ev *event.Pon) error {
	pon, err := meld.NewPon(ev.Taken(), ev.Consumed(), ev.Target())
	if err != nil {
		return err
	}
	actorSeat := ev.Actor()
	if err := s.applyOpenCall(actorSeat, ev.Target(), ev.Taken(), func() error {
		return s.players[actorSeat.Index()].Pon(*pon)
	}); err != nil {
		return err
	}
	s.disableKyushukyuhaiForAll()
	s.pendingDiscard = &actorSeat
	return nil
}

func (s *State) applyCalledKan(ev *event.CalledKan) error {
	if s.numKans >= maxNumKan {
		return fmt.Errorf("cannot CalledKan: already have %d kans", s.numKans)
	}
	kan, err := meld.NewCalledKan(ev.Taken(), ev.Consumed(), ev.Target())
	if err != nil {
		return err
	}
	actorSeat := ev.Actor()
	if err := s.applyOpenCall(actorSeat, ev.Target(), ev.Taken(), func() error {
		return s.players[actorSeat.Index()].CalledKan(*kan)
	}); err != nil {
		return err
	}
	s.disableKyushukyuhaiForAll()
	s.numKans++
	s.pendingKanActor = &actorSeat
	s.pendingDoraReveals++
	s.kanProgress = waitingReplacementBeforeDora
	s.nextDraw = actorSeat
	s.pendingDiscard = nil
	return nil
}

func (s *State) applyConcealedKan(ev *event.ConcealedKan) error {
	if s.numKans >= maxNumKan {
		return fmt.Errorf("cannot ConcealedKan: already have %d kans", s.numKans)
	}
	if s.numLeftTiles <= 0 {
		return fmt.Errorf("cannot ConcealedKan: no replacement tile left")
	}
	kan, err := meld.NewConcealedKan(ev.Consumed())
	if err != nil {
		return err
	}
	if err := s.players[ev.Actor().Index()].ConcealedKan(*kan); err != nil {
		return err
	}
	s.numKans++
	actorSeat := ev.Actor()
	s.disableKyushukyuhaiForAll()
	s.pendingKanActor = &actorSeat
	s.pendingDoraReveals++
	s.kanProgress = waitingReplacementAfterDora
	s.nextDraw = actorSeat
	s.pendingDiscard = nil
	s.lastActor = &actorSeat
	return nil
}

func (s *State) applyPromotedKan(ev *event.PromotedKan) error {
	if s.numKans >= maxNumKan {
		return fmt.Errorf("cannot PromotedKan: already have %d kans", s.numKans)
	}
	if s.numLeftTiles <= 0 {
		return fmt.Errorf("cannot PromotedKan: no replacement tile left")
	}
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
	s.numKans++
	actorSeat := ev.Actor()
	s.pendingKanActor = &actorSeat
	s.pendingDoraReveals++
	s.kanProgress = waitingReplacementBeforeDora
	s.pendingRobbedKanTile = &added
	s.nextDraw = actorSeat
	s.pendingDiscard = nil
	s.lastActor = &actorSeat
	return nil
}

func (s *State) applyDora(ev *event.Dora) error {
	if s.pendingDoraReveals <= 0 {
		return fmt.Errorf("cannot reveal dora indicator: not after kan")
	}
	if ev.Indicator().IsUnknown() {
		return fmt.Errorf("cannot reveal unknown dora indicator")
	}
	if len(s.doraIndicators) >= MaxNumDoraIndicators {
		return fmt.Errorf("cannot reveal dora indicator: already have %d indicators", len(s.doraIndicators))
	}
	s.doraIndicators = append(s.doraIndicators, ev.Indicator())
	s.pendingDoraReveals--
	if s.pendingKanActor != nil && s.pendingDoraReveals == 0 {
		switch s.kanProgress {
		case noKanProgress:
			s.pendingDiscard = s.pendingKanActor
			s.pendingKanActor = nil
		case waitingReplacementAfterDora:
			s.kanProgress = waitingReplacementBeforeDora
			s.nextDraw = *s.pendingKanActor
		}
	}
	return nil
}

func (s *State) isWaitingReplacementTileDraw(actor seat.Seat) bool {
	if s.pendingKanActor == nil || *s.pendingKanActor != actor {
		return false
	}
	return s.kanProgress == waitingReplacementBeforeDora
}

func (s *State) applyRiichi(ev *event.Riichi) error {
	if s.pendingDiscard == nil || *s.pendingDiscard != ev.Actor() {
		return fmt.Errorf("cannot Riichi: actor %d is not pending discard player", ev.Actor().Index())
	}
	if s.numLeftTiles < common.NumPlayers {
		return fmt.Errorf("cannot Riichi: no next draw turn remains")
	}
	if err := s.players[ev.Actor().Index()].Riichi(); err != nil {
		return err
	}
	actorSeat := ev.Actor()
	s.lastActor = &actorSeat
	return nil
}

func (s *State) applyRiichiAccepted(ev *event.RiichiAccepted) error {
	if s.pendingRiichiAcceptance == nil || *s.pendingRiichiAcceptance != ev.Actor() {
		return fmt.Errorf("cannot accept Riichi: actor %d is not pending riichi acceptance player", ev.Actor().Index())
	}
	if err := s.players[ev.Actor().Index()].RiichiAccepted(); err != nil {
		return err
	}
	s.applyRiichiAcceptedScoreUpdate(ev)
	s.riichiDeposit++
	s.pendingRiichiAcceptance = nil
	actorSeat := ev.Actor()
	s.lastActor = &actorSeat
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
	if s.roundEnded {
		if !s.canApplyAdditionalWin(ev) {
			return fmt.Errorf("cannot Win: round already ended")
		}
	} else if !s.canApplyWin(ev) {
		return fmt.Errorf("cannot Win: invalid timing")
	}
	s.applyScoreUpdate(ev.Scores(), ev.Deltas())
	s.roundEnded = true
	s.roundEndedByWin = true
	actorSeat := ev.Actor()
	targetSeat := ev.Target()
	s.winTarget = &targetSeat
	s.winActors[actorSeat.Index()] = true
	s.lastActor = &actorSeat
	return nil
}

func (s *State) canApplyAdditionalWin(ev *event.Win) bool {
	if !s.roundEndedByWin || s.winTarget == nil {
		// Additional wins are only possible after a previous win in this round.
		return false
	}
	if ev.Actor() == ev.Target() {
		// Tsumo immediately ends the round, so it cannot have double/triple wins.
		return false
	}
	if ev.Target() != *s.winTarget {
		// Double/triple ron must be against the same discarding player.
		return false
	}
	if s.winActors[ev.Actor().Index()] {
		// The same player cannot win twice from the same discard.
		return false
	}
	return s.canApplyWin(ev)
}

func (s *State) canApplyWin(ev *event.Win) bool {
	if s.pendingRobbedKanTile != nil {
		return s.canApplyRobbingKan(ev)
	}

	if ev.Actor() == ev.Target() {
		drawnTile := s.players[ev.Actor().Index()].DrawnTile()
		return drawnTile != nil && isTileMatchKnownEnough(drawnTile, ev.WinningTile())
	}

	targetRiver := s.players[ev.Target().Index()].River()
	return len(targetRiver) > 0 && isTileMatchKnownEnough(&targetRiver[len(targetRiver)-1], ev.WinningTile())
}

func (s *State) canApplyRobbingKan(ev *event.Win) bool {
	if s.pendingKanActor == nil || *s.pendingKanActor != ev.Target() {
		return false
	}
	// This currently covers robbing a promoted kan. MahjongSoul's Kokushi Musou
	// robbing-a-concealed-kan rule is not supported because concealed kan does
	// not expose a pendingRobbedKanTile.
	if s.pendingRobbedKanTile == nil {
		return false
	}
	return isTileMatchKnownEnough(s.pendingRobbedKanTile, ev.WinningTile())
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
	s.roundEnded = true
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

func (s *State) disableKyushukyuhaiForAll() {
	for i := range s.canKyushukyuhai {
		s.canKyushukyuhai[i] = false
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
	s.lastActor = &actorSeat
	return nil
}
