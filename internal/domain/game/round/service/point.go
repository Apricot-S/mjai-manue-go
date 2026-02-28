package service

// RonPoints calculates the winning point in case of Ron.
func RonPoints(fu, han int, isDealer bool) int {
	var basePoints int
	if han >= 13 {
		basePoints = 8000
	} else if han >= 11 {
		basePoints = 6000
	} else if han >= 8 {
		basePoints = 4000
	} else if han >= 6 {
		basePoints = 3000
	} else if han >= 5 || (han >= 4 && fu >= 40) || (han >= 3 && fu >= 70) {
		basePoints = 2000
	} else if han >= 1 {
		basePoints = fu * (1 << (han + 2))
	} else {
		basePoints = 0
	}

	var multiplier int
	if isDealer {
		multiplier = 6
	} else {
		multiplier = 4
	}

	return (basePoints*multiplier + 99) / 100 * 100
}

func RyukyokuPoints(tenpais [4]bool) [4]int {
	numTenpais := 0
	for _, t := range tenpais {
		if t {
			numTenpais++
		}
	}

	if numTenpais == 0 || numTenpais == 4 {
		return [4]int{0, 0, 0, 0}
	}

	plusPoints := 3000 / numTenpais
	minusPoints := -3000 / (4 - numTenpais)
	var ryukyokuPoints [4]int
	for i, tenpai := range tenpais {
		if tenpai {
			ryukyokuPoints[i] = plusPoints
		} else {
			ryukyokuPoints[i] = minusPoints
		}
	}
	return ryukyokuPoints
}
