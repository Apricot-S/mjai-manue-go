package ai

import (
	"github.com/Apricot-S/mjai-manue-go/internal/game"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/outbound"
)

type AI interface {
	Initialize()
	DecideAction(state game.StateAnalyzer, playerID int) (outbound.Event, error)
}
