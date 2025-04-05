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
	name     string
	room     string
	ai       ai.AI
	playerID int
	inGame   bool
	state    game.State
}

func NewAIAgent(name string, room string, ai ai.AI) *AIAgent {
	return NewAIAgentWithState(name, room, ai, &game.StateImpl{})
}

func NewAIAgentWithState(name string, room string, ai ai.AI, state game.State) *AIAgent {
	return &AIAgent{
		name:     name,
		room:     room,
		ai:       ai,
		playerID: -1,
		inGame:   false,
		state:    state,
	}
}

func (a *AIAgent) setPlayerID(id int) {
	a.playerID = id
}

func (a *AIAgent) setInGame(inGame bool) {
	a.inGame = inGame
}

func (a *AIAgent) Respond(msgs []jsontext.Value) (jsontext.Value, error) {
	var msg message.Message

	firstMsg := msgs[0]
	if err := json.Unmarshal(firstMsg, &msg); err != nil {
		return nil, err
	}

	switch msg.Type {
	case message.TypeHello:
		return makeJoinResponse(a.name, a.room)
	case message.TypeStartGame:
		return onStartGame(a, firstMsg)
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

	// Get last message to determine response
	lastMsg := msgs[len(msgs)-1]
	if err := json.Unmarshal(lastMsg, &msg); err != nil {
		return nil, err
	}

	switch msg.Type {
	case message.TypeTsumo:
		var tsumo message.Tsumo
		if err := json.Unmarshal(lastMsg, &tsumo); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tsumo message: %w", err)
		}

		if tsumo.Actor != a.playerID {
			// Not self tsumo
			return makeNoneResponse()
		}

		// Ask AI for decision
		action, err := a.ai.DecideAction(a.state, a.playerID)
		if err != nil {
			return nil, fmt.Errorf("failed to decide action: %w", err)
		}

		if action == nil {
			// No action needed
			return makeNoneResponse()
		}

		res, err := json.Marshal(action)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal action: %w", err)
		}
		return res, nil

	case message.TypeDahai:
		var dahai message.Dahai
		if err := json.Unmarshal(lastMsg, &dahai); err != nil {
			return nil, fmt.Errorf("failed to unmarshal dahai message: %w", err)
		}

		if dahai.Actor == a.playerID {
			// Self dahai
			return makeNoneResponse()
		}

		// Ask AI for decision on opponent's dahai
		action, err := a.ai.DecideAction(a.state, a.playerID)
		if err != nil {
			return nil, fmt.Errorf("failed to decide action: %w", err)
		}

		if action == nil {
			// No action needed
			return makeNoneResponse()
		}

		res, err := json.Marshal(action)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal action: %w", err)
		}
		return res, nil

	case message.TypeEndGame:
		return onEndGame(a)
	default:
		return makeNoneResponse()
	}
}
