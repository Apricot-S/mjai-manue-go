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

func (a *baseAgent) makeNoneResponse() (outbound.Event, error) {
	return outbound.NewNone(), nil
}

func (a *baseAgent) makeJoinResponse() (outbound.Event, error) {
	return outbound.NewJoin(a.name, a.room), nil
}

func (a *baseAgent) onStartGame(event inbound.StartGame) error {
	a.playerID = event.ID
	a.inGame = true

	return nil
}

func (a *baseAgent) onEndGame() {
	a.playerID = -1
	a.inGame = false
}
