package service

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tilecount"
)

func IsWinningForm(hand *hand.Hand) bool {
	return false
}

// Reference: https://qiita.com/tomohxx/items/20d886d1991ab89f5522
func isWinningFormGeneral(tc34 *tilecount.TileCounts34) bool {
	return false
}

func isSingleColorWinningFormWithoutPair(singleColorHand []int) bool {
	var r int
	a := singleColorHand[0]
	b := singleColorHand[1]

	for i := range 7 {
		r = a % 3
		c := singleColorHand[i+2]
		if b < r || c < r {
			return false
		}
		a = b - r
		b = c - r
	}

	return a%3 == 0 && b%3 == 0
}

func isSingleColorWinningFormWithPair(singleColorHand []int) bool {
	return false
}

func isWinningFormChitoitsu(tc34 *tilecount.TileCounts34) bool {
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

func isWinningFormKokushimuso(tc34 *tilecount.TileCounts34) bool {
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
