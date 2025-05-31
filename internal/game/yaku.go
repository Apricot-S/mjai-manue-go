package game

import "slices"

// Note:
// This function calculates points based on the original specification,
// not the exact score.
func CalculateFan(
	state StateViewer,
	playerID int,
	tehais Pais,
	mentsus []Mentsu,
	furos []Furo,
	reach bool,
) (fu int, fan int, points int, yakus map[string]int) {
	furoMentsus := make([]Mentsu, len(furos))
	for i, f := range furos {
		furoMentsus[i] = f.ToMentsu()
	}
	allMentsus := slices.Concat(mentsus, furoMentsus)
	allPais := make(Pais, 0, 14)
	for _, m := range allMentsus {
		allPais = append(allPais, m.Pais()...)
	}

	fan = 0
	yakus = make(map[string]int)

	if reach {
		fan += 1
		yakus["reach"] = 1
	}

	// TODO: Implement yaku calculation

	isPinfu := false
	if isPinfu || len(furoMentsus) > 0 {
		fu = 30
	} else {
		fu = 40
	}

	isOya := playerID == state.Oya().ID()
	points = GetPoints(fu, fan, isOya)

	return fu, fan, points, yakus
}
