package service

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
)

func IsWinningForm(hand *hand.Hand) bool {
	return false
}

// Reference: https://qiita.com/tomohxx/items/20d886d1991ab89f5522
func isWinningFormGeneral(hand *hand.Hand) bool {
	return false
}

func isWinningFormChitoitsu(hand *hand.Hand) bool {
	tc34 := hand.ToTileCounts34()

	if tc34.NumTiles() != 14 {
		return false
	}

	numPairs := 0
	for _, c := range tc34 {
		if c == 2 {
			numPairs++
		}
	}

	return numPairs == 7
}

func isWinningFormKokushimuso(hand *hand.Hand) bool {
	tc34 := hand.ToTileCounts34()

	if tc34.NumTiles() != 14 {
		return false
	}

	numKinds := 0
	hasPair := false
	for _, i := range tile.YaochuhaiIDs {
		if tc34[i] >= 1 {
			numKinds++
		}
		if tc34[i] >= 2 {
			hasPair = true
		}
	}

	return numKinds == 13 && hasPair
}
