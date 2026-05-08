package player

import (
	"fmt"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type commonPlayerState struct {
	drawnTile                 *tile.Tile
	melds                     []meld.Meld
	river                     []tile.Tile
	discardedTiles            []tile.Tile
	extraSafeTiles            []tile.Tile
	riichiState               RiichiState
	riichiRiverIndex          int
	riichiDiscardedTilesIndex int
	isConcealed               bool
	swapCallTiles             []tile.Tile
	needsDeadWallDraw         bool
}

func newCommonPlayerState() commonPlayerState {
	return commonPlayerState{
		drawnTile:                 nil,
		melds:                     make([]meld.Meld, 0, maxNumMelds),
		river:                     make([]tile.Tile, 0, maxNumRiver),
		discardedTiles:            make([]tile.Tile, 0, maxNumDiscardedTiles),
		extraSafeTiles:            make([]tile.Tile, 0, 3),
		riichiState:               NotRiichi,
		riichiRiverIndex:          -1,
		riichiDiscardedTilesIndex: -1,
		isConcealed:               true,
		swapCallTiles:             nil,
		needsDeadWallDraw:         false,
	}
}

func (s *commonPlayerState) DrawnTile() *tile.Tile {
	return s.drawnTile
}

func (s *commonPlayerState) Melds() []meld.Meld {
	return slices.Clone(s.melds)
}

func (s *commonPlayerState) River() []tile.Tile {
	return slices.Clone(s.river)
}

func (s *commonPlayerState) DiscardedTiles() []tile.Tile {
	return slices.Clone(s.discardedTiles)
}

func (s *commonPlayerState) ExtraSafeTiles() []tile.Tile {
	return slices.Clone(s.extraSafeTiles)
}

func (s *commonPlayerState) RiichiState() RiichiState {
	return s.riichiState
}

func (s *commonPlayerState) RiichiRiverIndex() int {
	return s.riichiRiverIndex
}

func (s *commonPlayerState) RiichiDiscardedTilesIndex() int {
	return s.riichiDiscardedTilesIndex
}

func (s *commonPlayerState) CanDiscard() bool {
	return !s.needsDeadWallDraw && s.drawnTile != nil || s.swapCallTiles != nil
}

func (s *commonPlayerState) CanChiiPonKan() bool {
	return s.riichiState == NotRiichi && !s.needsDeadWallDraw && !s.CanDiscard() && len(s.melds) < maxNumMelds
}

func (s *commonPlayerState) IsConcealed() bool {
	return s.isConcealed
}

func (s *commonPlayerState) SwapCallTiles() []tile.Tile {
	return slices.Clone(s.swapCallTiles)
}

func (s *commonPlayerState) TakeFromRiver(t tile.Tile) error {
	numRiver := len(s.river)

	if t != s.river[numRiver-1] {
		return fmt.Errorf("cannot take tile %s; last river tile is %s", t, s.river[numRiver-1])
	}

	s.river = slices.Delete(s.river, numRiver-1, numRiver)
	return nil
}
