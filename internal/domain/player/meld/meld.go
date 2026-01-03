package meld

import (
	"fmt"
	"strings"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/playerid"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/tile"
)

type Meld interface {
	Consumed() []tile.Tile
	ToTiles() []tile.Tile
	String() string
}

type OpenMeld interface {
	Meld
	Taken() *tile.Tile
	Target() *playerid.PlayerID
}

type ChiiPon interface {
	OpenMeld
	// Red five is included.
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
		consumedStrs[i] = t.String()
	}

	taken := m.Taken().String()
	target := m.Target().Index()
	consumed := strings.Join(consumedStrs, " ")

	return fmt.Sprintf("[%s(%d)/%s]", taken, target, consumed)
}
