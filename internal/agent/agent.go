package agent

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/message"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Agent interface {
	Respond(msgs []jsontext.Value) (jsontext.Value, error)
}

type baseAgent struct {
	name     string
	room     string
	playerID int
	inGame   bool
}

func newBaseAgent(name string, room string) baseAgent {
	return baseAgent{
		name:     name,
		room:     room,
		playerID: -1,
		inGame:   false,
	}
}

func (a *baseAgent) makeNoneResponse() (jsontext.Value, error) {
	none := message.NewNone()
	res, err := json.Marshal(&none)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal none message: %w", err)
	}
	return res, nil
}

func (a *baseAgent) makeJoinResponse() (jsontext.Value, error) {
	join := message.NewJoin(a.name, a.room)
	res, err := json.Marshal(&join)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal join message: %w", err)
	}
	return res, nil
}

func (a *baseAgent) onStartGame(rawMsg jsontext.Value) error {
	var startGame message.StartGame
	if err := json.Unmarshal(rawMsg, &startGame); err != nil {
		return fmt.Errorf("failed to unmarshal start_game message: %w", err)
	}

	a.playerID = startGame.ID
	a.inGame = true

	return nil
}

func (a *baseAgent) onEndGame() {
	a.playerID = -1
	a.inGame = false
}
