package agent

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/outbound"
)

type TsumogiriAgent struct {
	baseAgent
}

func NewTsumogiriAgent(name string, room string) *TsumogiriAgent {
	return &TsumogiriAgent{
		baseAgent: newBaseAgent(name, room),
	}
}

func (a *TsumogiriAgent) Respond(events []inbound.Event) (outbound.Event, error) {
	// Process messages before and after the game
	switch firstEvent := events[0].(type) {
	case *inbound.Hello:
		return a.makeJoinResponse()
	case *inbound.StartGame:
		if err := a.onStartGame(*firstEvent); err != nil {
			return nil, err
		}
		return a.makeNoneResponse()
	case *inbound.EndKyoku:
		// Message during the game, but does not affect the game, so process it here
		return a.makeNoneResponse()
	case *inbound.EndGame, *inbound.Error:
		a.onEndGame()
		return a.makeNoneResponse()
	}

	if !a.inGame {
		return nil, fmt.Errorf("received message while not in game: %v", events)
	}

	if lastEvent, ok := events[len(events)-1].(*inbound.Tsumo); ok {
		if lastEvent.Actor != a.playerID {
			// Not self tsumo
			return a.makeNoneResponse()
		}

		// Self tsumo
		dahai, err := outbound.NewDahai(a.playerID, lastEvent.Pai, true, "")
		if err != nil {
			return nil, fmt.Errorf("failed to make dahai: %w", err)
		}

		return dahai, nil
	}

	// Not tsumo
	return a.makeNoneResponse()
}
