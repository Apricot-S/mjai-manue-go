package agent

import (
	"github.com/Apricot-S/mjai-manue-go/internal/ai"
	"github.com/Apricot-S/mjai-manue-go/internal/game"
	"github.com/go-json-experiment/json/jsontext"
)

type Agent struct {
	name     string
	playerID int
	state    *game.State
	ai       ai.AI
}

func (a *Agent) Respond(message *jsontext.Value) (jsontext.Value, error) {
	// Dummy implementation
	return []byte{}, nil
}
