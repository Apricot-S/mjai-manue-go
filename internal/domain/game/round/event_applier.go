package round

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
)

type EventApplier interface {
	Apply(ev event.Event) error
}

func (s *State) Apply(ev event.Event) error {
	switch ev.(type) {
	case *event.StartRound:
		return fmt.Errorf("unimplemented event: %T", ev)
	case *event.EndRound:
		return fmt.Errorf("unimplemented event: %T", ev)
	}
	return fmt.Errorf("unknown event: %T", ev)
}
