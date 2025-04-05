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

// isSelfTurnType はそのメッセージが自分の手番を表すかを判定します
func isSelfTurnType(t message.Type) bool {
	return t == message.TypeTsumo || t == message.TypeChi || t == message.TypePon
}

// isOtherTurnType はそのメッセージが他家の手番を表すかを判定します
func isOtherTurnType(t message.Type) bool {
	return t == message.TypeDahai || t == message.TypeKakan
}

// findLastAction は配列の最後から、自分の手番または他家の手番のメッセージを探し、そのメッセージとTypeを返します
func findLastAction(msgs []jsontext.Value) (*jsontext.Value, message.Type, error) {
	for i, m := range slices.Backward(msgs) {
		var msg message.Message
		if err := json.Unmarshal(m, &msg); err != nil {
			return nil, "", err
		}
		if isSelfTurnType(msg.Type) || isOtherTurnType(msg.Type) {
			return &msgs[len(msgs)-1-i], msg.Type, nil
		}
	}
	return nil, "", nil
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

	// Find last tsumo or dahai message
	lastActionMsg, lastActionType, err := findLastAction(msgs)
	if err != nil {
		return nil, fmt.Errorf("failed to find last action: %w", err)
	}

	if lastActionMsg == nil {
		// Not self action
		return makeNoneResponse()
	}

	switch lastActionType {
	case message.TypeTsumo:
		var tsumo message.Tsumo
		if err := json.Unmarshal(*lastActionMsg, &tsumo); err != nil {
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
			return makeNoneResponse()
		}
		res, err := json.Marshal(action)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal action: %w", err)
		}
		return res, nil

	case message.TypeChi:
		var chi message.Chi
		if err := json.Unmarshal(*lastActionMsg, &chi); err != nil {
			return nil, fmt.Errorf("failed to unmarshal chi message: %w", err)
		}

		if chi.Actor != a.playerID {
			// Not self action
			return makeNoneResponse()
		}

		// Self action
		action, err := a.ai.DecideAction(a.state, a.playerID)
		if err != nil {
			return nil, fmt.Errorf("failed to decide action: %w", err)
		}
		if action == nil {
			return makeNoneResponse()
		}
		res, err := json.Marshal(action)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal action: %w", err)
		}
		return res, nil

	case message.TypePon:
		var pon message.Pon
		if err := json.Unmarshal(*lastActionMsg, &pon); err != nil {
			return nil, fmt.Errorf("failed to unmarshal pon message: %w", err)
		}

		if pon.Actor != a.playerID {
			// Not self action
			return makeNoneResponse()
		}

		action, err := a.ai.DecideAction(a.state, a.playerID)
		if err != nil {
			return nil, fmt.Errorf("failed to decide action: %w", err)
		}
		if action == nil {
			return makeNoneResponse()
		}
		res, err := json.Marshal(action)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal action: %w", err)
		}
		return res, nil

	case message.TypeDahai:
		var dahai message.Dahai
		if err := json.Unmarshal(*lastActionMsg, &dahai); err != nil {
			return nil, fmt.Errorf("failed to unmarshal dahai message: %w", err)
		}

		if dahai.Actor == a.playerID {
			return makeNoneResponse()
		}

		// Ask AI for decision on opponent's action
		action, err := a.ai.DecideAction(a.state, a.playerID)
		if err != nil {
			return nil, fmt.Errorf("failed to decide action: %w", err)
		}
		if action == nil {
			return makeNoneResponse()
		}
		res, err := json.Marshal(action)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal action: %w", err)
		}
		return res, nil

	case message.TypeKakan:
		var kakan message.Kakan
		if err := json.Unmarshal(*lastActionMsg, &kakan); err != nil {
			return nil, fmt.Errorf("failed to unmarshal kakan message: %w", err)
		}

		if kakan.Actor == a.playerID {
			return makeNoneResponse()
		}

		// Ask AI for decision on opponent's action
		action, err := a.ai.DecideAction(a.state, a.playerID)
		if err != nil {
			return nil, fmt.Errorf("failed to decide action: %w", err)
		}
		if action == nil {
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
