package round

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
)

type StateUpdater interface {
	Update(ev event.Event) error
}

func (s *State) Update(ev event.Event) error {
	switch ev.(type) {
	case event.EndRound:
	case event.EndGame:
	case event.RequestAction:
	}
	return fmt.Errorf("unknown event type: %s", ev.EventType())
}
