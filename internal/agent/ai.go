package agent

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/ai"
	"github.com/Apricot-S/mjai-manue-go/internal/game"
	"github.com/Apricot-S/mjai-manue-go/internal/message"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type AIAgent struct {
	baseAgent
	ai    ai.AI
	state game.State
}

func NewAIAgent(name string, room string, ai ai.AI) *AIAgent {
	return NewAIAgentWithState(name, room, ai, &game.StateImpl{})
}

func NewAIAgentWithState(name string, room string, ai ai.AI, state game.State) *AIAgent {
	return &AIAgent{
		baseAgent: newBaseAgent(name, room),
		ai:        ai,
		state:     state,
	}
}

func (a *AIAgent) Respond(msgs []jsontext.Value) (jsontext.Value, error) {
	var msg message.Message

	firstMsg := msgs[0]
	if err := json.Unmarshal(firstMsg, &msg); err != nil {
		return nil, err
	}

	// Process messages before and after the game
	switch msg.Type {
	case message.TypeHello:
		return a.makeJoinResponse()
	case message.TypeStartGame:
		if err := a.onStartGame(firstMsg); err != nil {
			return nil, err
		}
		if err := a.state.OnStartGame(firstMsg); err != nil {
			return nil, fmt.Errorf("failed to update state: %w", err)
		}
		return a.makeNoneResponse()
	case message.TypeEndKyoku:
		// Message during the game, but does not affect the game, so process it here
		return a.makeNoneResponse()
	case message.TypeEndGame, message.TypeError:
		a.onEndGame()
		return a.makeNoneResponse()
	}

	if !a.inGame {
		return nil, fmt.Errorf("received message while not in game: %v", msgs)
	}

	// Update state for all messages
	for _, m := range msgs {
		if err := a.state.Update(m); err != nil {
			return nil, fmt.Errorf("failed to update state: %w", err)
		}
	}
	a.state.Print()

	// Ask AI for decision
	action, err := a.ai.DecideAction(a.state, a.playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to decide action: %w", err)
	}
	return action, nil
}
