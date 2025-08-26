package agent

import (
	"fmt"
	"os"

	"github.com/Apricot-S/mjai-manue-go/internal/ai"
	"github.com/Apricot-S/mjai-manue-go/internal/game"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/outbound"
)

type AIAgent struct {
	baseAgent
	ai    ai.AI
	state game.State
}

func NewAIAgent(name string, room string, ai ai.AI, state game.State) *AIAgent {
	return &AIAgent{
		baseAgent: newBaseAgent(name, room),
		ai:        ai,
		state:     state,
	}
}

func NewAIAgentDefault(name string, room string, ai ai.AI) *AIAgent {
	return NewAIAgent(name, room, ai, &game.StateImpl{})
}

func (a *AIAgent) Respond(events []inbound.Event) (outbound.Event, error) {
	// Process messages before and after the game
	switch firstEvent := events[0].(type) {
	case *inbound.Hello:
		return a.makeJoinResponse(), nil
	case *inbound.StartGame:
		a.onStartGame(*firstEvent)
		if err := a.state.Update(firstEvent); err != nil {
			return nil, fmt.Errorf("failed to update state: %w", err)
		}
		a.ai.Initialize()
		return a.makeNoneResponse(), nil
	case *inbound.EndKyoku:
		// Message during the game, but does not affect the game, so process it here
		return a.makeNoneResponse(), nil
	case *inbound.EndGame, *inbound.Error:
		a.onEndGame()
		return a.makeNoneResponse(), nil
	}

	if !a.inGame {
		return nil, fmt.Errorf("received message while not in game: %v", events)
	}

	// Update state for all messages
	for _, ev := range events {
		if err := a.state.Update(ev); err != nil {
			return nil, fmt.Errorf("failed to update state: %w", err)
		}
	}
	fmt.Fprint(os.Stderr, a.state.RenderBoard())

	// Ask AI for decision
	action, err := a.ai.DecideAction(a.state, a.playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to decide action: %w", err)
	}
	return action, nil
}
