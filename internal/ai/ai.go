package ai

import (
	"github.com/Apricot-S/mjai-manue-go/internal/game"
	"github.com/go-json-experiment/json/jsontext"
)

type AI interface {
	DecideAction(state game.StateAnalyzer, playerID int) (jsontext.Value, error)
}
