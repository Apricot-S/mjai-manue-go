package meld

import (
	"fmt"
	"strings"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/block"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/playerid"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
)

type Meld interface {
	Consumed() []tile.Tile
	ToTiles() []tile.Tile
	ToBlock() block.Block
	ToString() string
}

type OpenMeld interface {
	Meld
	Taken() *tile.Tile
	Target() *playerid.PlayerID
}

type ChiiPon interface {
	OpenMeld
	SwapCallTiles() []tile.Tile
}

func countRed(tiles []tile.Tile) int {
	numRed := 0
	for _, t := range tiles {
		if t.IsRed() {
			numRed++
		}
	}
	return numRed
}

func meldToString(m OpenMeld) string {
	consumedStrs := make([]string, len(m.Consumed()))
	for i, t := range m.Consumed() {
		consumedStrs[i] = t.Code()
	}

	return fmt.Sprintf("[%s(%d)/%s]", m.Taken().Code(), m.Target().Index(), strings.Join(consumedStrs, " "))
}
