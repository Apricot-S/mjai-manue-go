package ai

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func selfTurnHand(self player.PlayerViewer) (*hand.VisibleHand, error) {
	h, ok := self.Hand()
	if !ok {
		return nil, fmt.Errorf("self hand is not visible")
	}
	drawnTile := self.DrawnTile()
	if drawnTile == nil {
		return h, nil
	}
	withDrawnTile, err := h.Draw(*drawnTile)
	if err != nil {
		return nil, fmt.Errorf("cannot add drawn tile %s to self hand: %w", *drawnTile, err)
	}
	return withDrawnTile, nil
}

func scoreDiscardCandidate(discardTile tile.Tile, shanten int) candidateScore {
	return candidateScore{
		averageRank:    0,
		expectedPoints: 0,
		shanten:        shanten,
		red:            discardTile.IsRed(),
	}
}
