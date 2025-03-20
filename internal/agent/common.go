package agent

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/message"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

func makeNoneResponse() (jsontext.Value, error) {
	none := message.NewNone()
	res, err := json.Marshal(&none)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal none message: %w", err)
	}
	return res, nil
}

func makeJoinResponse(name string, room string) (jsontext.Value, error) {
	join := message.NewJoin(name, room)
	res, err := json.Marshal(&join)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal join message: %w", err)
	}
	return res, nil
}

type gameStateHandler interface {
	setPlayerID(id int)
	setInGame(inGame bool)
}

func onStartGame(handler gameStateHandler, rawMsg jsontext.Value) (jsontext.Value, error) {
	var startGame message.StartGame
	if err := json.Unmarshal(rawMsg, &startGame); err != nil {
		return nil, fmt.Errorf("failed to unmarshal start_game message: %w", err)
	}
	handler.setPlayerID(startGame.ID)
	handler.setInGame(true)
	return makeNoneResponse()
}

func onEndGame(handler gameStateHandler) (jsontext.Value, error) {
	handler.setPlayerID(-1)
	handler.setInGame(false)
	return makeNoneResponse()
}
