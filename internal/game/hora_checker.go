package game

import (
	"fmt"
)

func IsHoraForm(ps *PaiSet) (bool, error) {
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

	numMentsus := min(sum/3, 4)
	ret := isHoraFormGeneral(ps, numMentsus)
	if sum == 14 {
		ret = ret || isHoraFormChitoitsu(ps)
		ret = ret || isHoraFormKokushimuso(ps)
	}

	return ret, nil
}

func isHoraFormGeneral(ps *PaiSet, numMentsus int) bool {
	return false
}

func isHoraFormChitoitsu(ps *PaiSet) bool {
	numPairs := 0
	for _, c := range ps {
		if c == 2 {
			numPairs++
		}
	}
	return numPairs == 7
}

var yaochuhaiIndices = [13]int{0, 8, 9, 17, 18, 26, 27, 28, 29, 30, 31, 32, 33}

func isHoraFormKokushimuso(ps *PaiSet) bool {
	ret := 1
	for _, i := range yaochuhaiIndices {
		ret *= ps[i]
	}
	return ret == 2
}
