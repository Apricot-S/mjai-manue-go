package ai

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type actionCandidate struct {
	traceKey         string
	action           action.Action
	riichi           bool
	scoreAsRiichi    bool
	discardTile      tile.Tile
	melds            []meld.Meld
	turnHand         *hand.VisibleHand
	afterDiscardHand *hand.VisibleHand
	baseShanten      int
	shantenGoals     []service.Goal
	score            candidateScore
}
