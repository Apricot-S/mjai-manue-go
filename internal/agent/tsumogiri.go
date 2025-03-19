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

func (a *TsumogiriAgent) Respond(msgs []jsontext.Value) (jsontext.Value, error) {
	var msg message.Message

	firstMsg := msgs[0]
	if err := json.Unmarshal(firstMsg, &msg); err != nil {
		return nil, err
	}

	switch msg.Type {
	case message.TypeStartGame:
		var startGame message.StartGame
		if err := json.Unmarshal(firstMsg, &startGame); err != nil {
			return nil, fmt.Errorf("failed to unmarshal start_game message: %w", err)
		}
		a.playerID = startGame.ID
		a.inGame = true
		return makeNoneResponse()
	case message.TypeEndGame:
		a.playerID = -1
		a.inGame = false
		return makeNoneResponse()
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
			return nil, fmt.Errorf("failed to unmarshal dahai message: %w", err)
		}

		return res, nil
	case message.TypeHello:
		return makeJoinResponse(a.name, a.room)
	default:
		return makeNoneResponse()
	}
}
