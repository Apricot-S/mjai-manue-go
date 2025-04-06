package agent

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/message"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type TsumogiriAgent struct {
	name     string
	room     string
	playerID int
	inGame   bool
}

func NewTsumogiriAgent(name string, room string) *TsumogiriAgent {
	return &TsumogiriAgent{
		name:     name,
		room:     room,
		playerID: -1,
		inGame:   false,
	}
}

func (a *TsumogiriAgent) setPlayerID(id int) {
	a.playerID = id
}

func (a *TsumogiriAgent) setInGame(inGame bool) {
	a.inGame = inGame
}

func (a *TsumogiriAgent) Respond(msgs []jsontext.Value) (jsontext.Value, error) {
	var msg message.Message

	firstMsg := msgs[0]
	if err := json.Unmarshal(firstMsg, &msg); err != nil {
		return nil, err
	}

	// Process messages before and after the game
	switch msg.Type {
	case message.TypeHello:
		return makeJoinResponse(a.name, a.room)
	case message.TypeStartGame:
		if err := onStartGame(a, firstMsg); err != nil {
			return nil, err
		}
		return makeNoneResponse()
	case message.TypeEndKyoku:
		// Message during the game, but does not affect the game, so process it here
		return makeNoneResponse()
	case message.TypeEndGame:
		onEndGame(a)
		return makeNoneResponse()
	}

	if !a.inGame {
		return nil, fmt.Errorf("received message while not in game: %v", msgs)
	}

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

		// Self tsumo
		dahai, err := message.NewDahai(a.playerID, tsumo.Pai, true, "")
		if err != nil {
			return nil, fmt.Errorf("failed to make dahai: %w", err)
		}

		res, err := json.Marshal(&dahai)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal dahai message: %w", err)
		}

		return res, nil
	default:
		return makeNoneResponse()
	}
}
