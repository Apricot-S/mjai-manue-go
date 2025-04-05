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
	}
	otherTurnTypes = []message.Type{
		message.TypeDahai,
		message.TypeKakan,
	}
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

// isMyTurn はそのメッセージが自分のアクションを要求するものかを判定します
func isMyTurn(t message.Type, actor, playerID int) bool {
	return actor == playerID && slices.Contains(selfTurnTypes, t)
}

// needsResponse はそのメッセージが自分の応答を要求するものかを判定します
func needsResponse(t message.Type, actor, playerID int) bool {
	return actor != playerID && slices.Contains(otherTurnTypes, t)
}

// findRelevantAction は配列の最後から、自分に関係のあるメッセージのTypeとActorを返します
func findRelevantAction(msgs []jsontext.Value, playerID int) (message.Type, error) {
	for _, m := range slices.Backward(msgs) {
		var action message.Action
		if err := json.Unmarshal(m, &action); err != nil {
			// Action型としてパースできないメッセージはスキップ
			continue
		}
		if isMyTurn(action.Type, action.Actor, playerID) || needsResponse(action.Type, action.Actor, playerID) {
			return action.Type, nil
		}
	}
	return "", nil
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
	case message.TypeEndKyoku:
		return makeNoneResponse()
	case message.TypeEndGame:
		return onEndGame(a)
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

	// Find relevant action
	lastActionType, err := findRelevantAction(msgs, a.playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to find last action: %w", err)
	}

	// No action needed
	if lastActionType == "" {
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
}
