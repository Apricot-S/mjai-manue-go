package game

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

func IsHoraForm(ps *base.PaiSet) (bool, error) {
	sum := 0
	for _, c := range ps {
		if c < 0 {
			return false, fmt.Errorf("negative number of tiles in the PaiSet")
		}
		if c > 4 {
			return false, fmt.Errorf("more than 4 tiles of the same type in the PaiSet")
		}
		sum += c
	}
	if sum > 14 {
		return false, fmt.Errorf("too many tiles in hand %d", sum)
	}
	if sum%3 != 2 {
		return false, fmt.Errorf("invalid hand length %d", sum)
	}

	ret := isHoraFormGeneral(ps)
	if sum == 14 {
		ret = ret || isHoraFormChitoitsu(ps)
		ret = ret || isHoraFormKokushimuso(ps)
	}

	return ret, nil
}

// Reference: https://qiita.com/tomohxx/items/20d886d1991ab89f5522
func isHoraFormGeneral(ps *base.PaiSet) bool {
	colorWithPair := -1

	for i := range 3 {
		sum := 0
		for _, c := range ps[9*i : 9*i+9] {
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

	for i := 27; i < 34; i++ {
		switch ps[i] % 3 {
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

	for i := range 3 {
		if i == colorWithPair {
			if !isSingleColorHoraFormWithPair(ps[9*i : 9*i+9]) {
				return false
			}
		} else {
			if !isSingleColorHoraFormWithoutPair(ps[9*i : 9*i+9]) {
				return false
			}
		}
	}

	return true
}

func isSingleColorHoraFormWithoutPair(ps []int) bool {
	var r int
	a := ps[0]
	b := ps[1]

	for i := range 7 {
		r = a % 3
		c := ps[i+2]
		if b < r || c < r {
			return false
		}
		a = b - r
		b = c - r
	}

	return a%3 == 0 && b%3 == 0
}

func isSingleColorHoraFormWithPair(ps []int) bool {
	p := 0
	for i := range 9 {
		p += i * ps[i]
	}

	for i := p * 2 % 3; i < 9; i += 3 {
		ps[i] -= 2
		if ps[i] >= 0 && isSingleColorHoraFormWithoutPair(ps) {
			ps[i] += 2
			return true
		}
		ps[i] += 2
	}
	return false
}

func isHoraFormChitoitsu(ps *base.PaiSet) bool {
	numPairs := 0
	for _, c := range ps {
		if c == 2 {
			numPairs++
		}
	}
	return numPairs == 7
}

var yaochuhaiIndices = [13]int{0, 8, 9, 17, 18, 26, 27, 28, 29, 30, 31, 32, 33}

func isHoraFormKokushimuso(ps *base.PaiSet) bool {
	ret := 1
	for _, i := range yaochuhaiIndices {
		ret *= ps[i]
	}
	return ret == 2
}
