package agent

import (
	"github.com/Apricot-S/mjai-manue-go/internal/ai"
	"github.com/Apricot-S/mjai-manue-go/internal/game"
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
	panic("unimplemented!")
}
