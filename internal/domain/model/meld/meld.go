package meld

import (
	"fmt"
	"strings"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/block"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
)

type Meld interface {
	Taken() *tile.Tile
	Consumed() []tile.Tile
	Target() int
	ToTiles() []tile.Tile
	ToBlock() block.Block
	ToString() string
}

func isValidTarget(target int) bool {
	return 0 <= target && target <= 3
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

func meldToString(m Meld) string {
	consumedStrs := make([]string, len(m.Consumed()))
	for i, t := range m.Consumed() {
		consumedStrs[i] = t.Code()
	}

	return fmt.Sprintf("[%s(%d)/%s]", m.Taken().Code(), m.Target(), strings.Join(consumedStrs, " "))
}
