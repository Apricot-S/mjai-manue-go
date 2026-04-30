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
	if s.numLeftTiles <= 0 {
		return fmt.Errorf("cannot Draw: no tiles left")
	}

	actorSeat := ev.Actor()
	p := s.players[actorSeat.Index()]
	if err := p.Draw(ev.Tile()); err != nil {
		return err
	}

	s.numLeftTiles--
	return nil
}

func (s *State) applyDiscard(ev *event.Discard) error {
	actorSeat := ev.Actor()
	p := s.players[actorSeat.Index()]
	return p.Discard(ev.Tile(), ev.Tsumogiri())
}

func (s *State) applyChii(ev *event.Chii) error {
	chii, err := meld.NewChii(ev.Taken(), ev.Consumed(), ev.Target())
	if err != nil {
		return err
	}
	return s.applyOpenCall(ev.Target(), ev.Taken(), func() error {
		return s.players[ev.Actor().Index()].Chii(*chii)
	})
}

func (s *State) applyPon(ev *event.Pon) error {
	pon, err := meld.NewPon(ev.Taken(), ev.Consumed(), ev.Target())
	if err != nil {
		return err
	}
	return s.applyOpenCall(ev.Target(), ev.Taken(), func() error {
		return s.players[ev.Actor().Index()].Pon(*pon)
	})
}

func (s *State) applyCalledKan(ev *event.CalledKan) error {
	kan, err := meld.NewCalledKan(ev.Taken(), ev.Consumed(), ev.Target())
	if err != nil {
		return err
	}
	return s.applyOpenCall(ev.Target(), ev.Taken(), func() error {
		return s.players[ev.Actor().Index()].CalledKan(*kan)
	})
}

func (s *State) applyConcealedKan(ev *event.ConcealedKan) error {
	kan, err := meld.NewConcealedKan(ev.Consumed())
	if err != nil {
		return err
	}
	return s.players[ev.Actor().Index()].ConcealedKan(*kan)
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
	return actor.PromotedKan(*kan)
}

func (s *State) applyDora(ev *event.Dora) error {
	if ev.Indicator().IsUnknown() {
		return fmt.Errorf("cannot add unknown dora indicator")
	}
	if len(s.doraIndicators) >= MaxNumDoraIndicators {
		return fmt.Errorf("cannot add dora indicator: already have %d indicators", len(s.doraIndicators))
	}
	s.doraIndicators = append(s.doraIndicators, ev.Indicator())
	return nil
}

func (s *State) applyRiichi(ev *event.Riichi) error {
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
	s.applyScoreUpdate(ev.Scores(), ev.Deltas())
	return nil
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

func (s *State) applyOpenCall(targetSeat seat.Seat, taken tile.Tile, applyActor func() error) error {
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
	return target.TakeFromRiver(taken)
}
