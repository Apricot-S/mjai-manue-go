package agent

import (
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/outbound"
)

type Agent interface {
	Respond(events []inbound.Event) (outbound.Event, error)
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

func (a *baseAgent) makeNoneResponse() outbound.Event {
	return outbound.NewNone()
}

func (a *baseAgent) makeJoinResponse() outbound.Event {
	return outbound.NewJoin(a.name, a.room)
}

func (a *baseAgent) onStartGame(event inbound.StartGame) {
	a.playerID = event.ID
	a.inGame = true
}

func (a *baseAgent) onEndGame() {
	a.playerID = -1
	a.inGame = false
}
