package round

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
)

type EventApplier interface {
	Apply(ev event.Event) error
}

func (s *State) Apply(ev event.Event) error {
	switch ev := ev.(type) {
	case *event.Draw:
		return s.applyDraw(ev)
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
