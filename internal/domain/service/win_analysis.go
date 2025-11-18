package service

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tilecount"
)

func IsWinningForm(hand *hand.Hand) bool {
	tc34 := hand.ToTileCounts34()

	numTiles := tc34.NumTiles()
	if numTiles%3 != 2 {
		return false
	}

	ret := isWinningFormGeneral(tc34)
	if numTiles == 14 {
		ret = ret || isWinningFormChitoitsu(tc34) || isWinningFormKokushimuso(tc34)
	}

	return ret
}

// Reference: https://qiita.com/tomohxx/items/20d886d1991ab89f5522
func isWinningFormGeneral(tc34 *tilecount.TileCounts34) bool {
	colorWithPair := -1

	for i := range 3 {
		sum := 0
		for _, c := range tc34[9*i : 9*i+9] {
			sum += c
		}

		switch sum % 3 {
		case 1:
			return false
		case 2:
			if colorWithPair == -1 {
				colorWithPair = i
			} else {
				return false
			}
		}
	}

	if !isHonorsWinningForm(tc34[27:34], colorWithPair != -1) {
		return false
	}

	for i := range 3 {
		if i == colorWithPair {
			if !isSingleColorWinningFormWithPair(tc34[9*i : 9*i+9]) {
				return false
			}
		} else {
			if !isSingleColorWinningFormWithoutPair(tc34[9*i : 9*i+9]) {
				return false
			}
		}
	}

	return true
}

func isHonorsWinningForm(honorsHand []int, hasPair bool) bool {
	for _, c := range honorsHand {
		switch c % 3 {
		case 1:
			return false
		case 2:
			if hasPair {
				return false
			} else {
				hasPair = true
			}
		}
	}

	return true
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
	sum := 0
	for i := range 9 {
		sum += i * singleColorHand[i]
	}

	for i := sum * 2 % 3; i < 9; i += 3 {
		singleColorHand[i] -= 2
		if singleColorHand[i] >= 0 && isSingleColorWinningFormWithoutPair(singleColorHand) {
			singleColorHand[i] += 2
			return true
		}
		singleColorHand[i] += 2
	}

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
