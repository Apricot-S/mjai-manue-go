package ai

func getPoints(fu, fan int, oya bool) int {
	var basePoints int
	if fan >= 13 {
		basePoints = 8000
	} else if fan >= 11 {
		basePoints = 6000
	} else if fan >= 8 {
		basePoints = 4000
	} else if fan >= 6 {
		basePoints = 3000
	} else if fan >= 5 || (fan >= 4 && fu >= 40) || (fan >= 3 && fu >= 70) {
		basePoints = 2000
	} else if fan >= 1 {
		basePoints = fu * (1 << (fan + 2))
	} else {
		basePoints = 0
	}

	var multiplier int
	if oya {
		multiplier = 6
	} else {
		multiplier = 4
	}

	return (basePoints*multiplier + 99) / 100 * 100
}

func tenpaisToRyukyokuPoints(tenpais [4]bool) [4]int {
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
	for i := range tenpais {
		if tenpais[i] {
			ryukyokuPoints[i] = plusPoints
		} else {
			ryukyokuPoints[i] = minusPoints
		}
	}
	return ryukyokuPoints
}
