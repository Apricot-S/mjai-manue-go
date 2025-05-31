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

	addYaku := func(existsYaku bool, name string, menzenFan int, kuiFan int) {
		if !existsYaku {
			return
		}

		var yakuFan int
		if len(furoMentsus) == 0 {
			yakuFan = menzenFan
		} else {
			yakuFan = kuiFan
		}

		yakus[name] = yakuFan
		fan += yakuFan
	}

	// TODO: Implement yaku calculation
	addYaku(reach, "reach", 1, 0)
	addYaku(isTanyaochu(allPais), "tyc", 1, 1)
	addYaku(isChantaiyao(allMentsus), "cty", 2, 1)

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

func isTanyaochu(allPais Pais) bool {
	for _, p := range allPais {
		if p.IsYaochu() {
			return false
		}
	}
	return true
}

func isChantaiyao(allMentsus []Mentsu) bool {
	for _, m := range allMentsus {
		isYaochuMentsu := false
		for _, p := range m.Pais() {
			if p.IsYaochu() {
				isYaochuMentsu = true
				break
			}
		}
		if !isYaochuMentsu {
			return false
		}
	}
	return true
}
