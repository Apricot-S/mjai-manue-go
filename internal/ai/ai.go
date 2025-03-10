package ai

import (
	"github.com/Apricot-S/mjai-manue-go/internal/game"
	"github.com/go-json-experiment/json/jsontext"
)

type AI interface {
	Respond(s *game.State) (jsontext.Value, error)
}
