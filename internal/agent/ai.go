package agent

import (
	"fmt"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/ai"
	"github.com/Apricot-S/mjai-manue-go/internal/game"
	"github.com/Apricot-S/mjai-manue-go/internal/message"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

var (
	selfTurnTypes = []message.Type{
		message.TypeTsumo,
		message.TypeChi,
		message.TypePon,
		message.TypeReach,
	}
	otherTurnTypes = []message.Type{
		message.TypeDahai,
		message.TypeKakan,
	}
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

// isMyTurn checks if the message requires an action from the player.
func isMyTurn(t message.Type, actor, playerID int) bool {
	return actor == playerID && slices.Contains(selfTurnTypes, t)
}

// needsResponse checks if the message requires a response from the player.
func needsResponse(t message.Type, actor, playerID int) bool {
	return actor != playerID && slices.Contains(otherTurnTypes, t)
}

// shouldDecideAction checks if AI needs to make a decision.
func shouldDecideAction(msgs []jsontext.Value, playerID int) (bool, error) {
	for _, m := range slices.Backward(msgs) {
		var action message.Action
		if err := json.Unmarshal(m, &action); err != nil {
			// Skip messages that cannot be parsed as Action
			continue
		}
		if isMyTurn(action.Type, action.Actor, playerID) || needsResponse(action.Type, action.Actor, playerID) {
			return true, nil
		}
	}
	return false, nil
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
	case message.TypeEndGame:
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

	// Check if AI needs to make a decision
	needsDecision, err := shouldDecideAction(msgs, a.playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to check if decision is needed: %w", err)
	}

	// No action needed
	if !needsDecision {
		return a.makeNoneResponse()
	}

	// Ask AI for decision
	action, err := a.ai.DecideAction(a.state, a.playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to decide action: %w", err)
	}
	if action == nil {
		return a.makeNoneResponse()
	}
	res, err := json.Marshal(action)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal action: %w", err)
	}
	return res, nil
}
