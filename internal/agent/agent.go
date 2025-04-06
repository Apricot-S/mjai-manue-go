package agent

import (
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

func (b *baseAgent) setPlayerID(id int) {
	b.playerID = id
}

func (b *baseAgent) setInGame(inGame bool) {
	b.inGame = inGame
}
