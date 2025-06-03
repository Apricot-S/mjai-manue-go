package game

import "slices"

// Note:
// This function calculates points based on the original specification,
// not the exact points.
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

		if yakuFan > 0 {
			yakus[name] = yakuFan
			fan += yakuFan
		}
	}

	addYaku(reach, "reach", 1, 0)
	addYaku(isTanyaochu(allPais), "tyc", 1, 1)
	addYaku(isChantaiyao(allMentsus), "cty", 2, 1)
	pinfu := isPinfu(state, playerID, allMentsus)
	addYaku(pinfu, "pf", 1, 0)
	yakuhaiFan := sumYakuhaiFan(state, playerID, allMentsus)
	addYaku(yakuhaiFan > 0, "ykh", yakuhaiFan, yakuhaiFan)
	addYaku(isIpeko(allMentsus), "ipk", 1, 0)
	addYaku(isSanshokuDojun(allMentsus), "ssj", 2, 1)
	addYaku(isIkkiTsukan(allMentsus), "ikt", 2, 1)
	addYaku(isToitoiho(allMentsus), "tth", 2, 2)
	if isChiniso(allMentsus) {
		addYaku(true, "cis", 6, 5)
	} else if isHoniso(allMentsus) {
		addYaku(true, "his", 3, 2)
	}

	if fan > 0 {
		doras := state.Doras()
		numDoras := 0
		for _, p := range allPais {
			for _, d := range doras {
				if p.HasSameSymbol(&d) {
					numDoras++
				}
			}
		}
		addYaku(numDoras > 0, "dr", numDoras, numDoras)

		var currentPais []Pai = slices.Clone(tehais)
		for _, f := range furos {
			currentPais = slices.Concat(currentPais, f.Pais())
		}
		numAkadoras := 0
		for _, cp := range currentPais {
			if !cp.IsRed() {
				continue
			}
			for _, p := range allPais {
				if p.HasSameSymbol(&cp) {
					numAkadoras++
					break
				}
			}
		}
		addYaku(numAkadoras > 0, "adr", numAkadoras, numAkadoras)
	}

	// TODO Calculate fu more accurately
	if pinfu || len(furoMentsus) > 0 {
		fu = 30
	} else {
		fu = 40
	}

	isOya := playerID == state.Oya().ID()
	points = GetPoints(fu, fan, isOya)

	return fu, fan, points, yakus
}

func Has1Fan(
	state StateViewer,
	playerID int,
	tehais Pais,
	furos []Furo,
	horaPai *Pai,
	isTsumo bool,
) (bool, error) {
	horaTehais := slices.Clone(tehais)
	if !isTsumo {
		horaTehais = append(horaTehais, *horaPai)
	}
	horaTehaisCounts, err := NewPaiSet(horaTehais)
	if err != nil {
		return false, err
	}

	isMenzen := state.Players()[playerID].IsMenzen()
	if isMenzen {
		if isChitoitsu := isHoraFormChitoitsu(horaTehaisCounts); isChitoitsu {
			return true, nil
		}
		if isKokushimuso := isHoraFormKokushimuso(horaTehaisCounts); isKokushimuso {
			return true, nil
		}
	}

	_, goals, err := AnalyzeShantenWithOption(horaTehaisCounts, 0, -1)
	if err != nil {
		return false, err
	}

	furoMentsus := make([]Mentsu, len(furos))
	for i, f := range furos {
		furoMentsus[i] = f.ToMentsu()
	}

	for _, goal := range goals {
		allMentsus := slices.Concat(goal.Mentsus, furoMentsus)
		allPais := make(Pais, 0, 14)
		for _, m := range allMentsus {
			allPais = append(allPais, m.Pais()...)
		}

		if sumYakuhaiFan(state, playerID, allMentsus) > 0 {
			return true, nil
		}
		if isTanyaochu(allPais) {
			return true, nil
		}
		if isMenzen && isIpeko(allMentsus) {
			return true, nil
		}
		if isMenzen && isPinfuStrict(state, playerID, allMentsus, horaPai) {
			return true, nil
		}
		if isChantaiyao(allMentsus) {
			return true, nil
		}
		if isIkkiTsukan(allMentsus) {
			return true, nil
		}
		if isSanshokuDojun(allMentsus) {
			return true, nil
		}
		if isSanshokuDoko(allMentsus) {
			return true, nil
		}
		if isSankantsu(allMentsus) {
			return true, nil
		}
		if isToitoiho(allMentsus) {
			return true, nil
		}
		// TODO: sananko
		if isHoniso(allMentsus) {
			return true, nil
		}
	}

	return false, nil
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

// TODO Consider ryanmen criteria
func isPinfu(state StateViewer, playerID int, allMentsus []Mentsu) bool {
	player := &state.Players()[playerID]
	for _, m := range allMentsus {
		switch m.(type) {
		case *Kotsu, *Kantsu:
			return false
		case *Toitsu:
			pai := &m.Pais()[0]
			if state.YakuhaiFan(pai, player) > 0 {
				return false
			}
		}
	}
	return true
}

func isPinfuStrict(state StateViewer, playerID int, allMentsus []Mentsu, horaPai *Pai) bool {
	player := &state.Players()[playerID]
	foundShuntsuWithHoraPai := false

	for _, m := range allMentsus {
		switch m.(type) {
		case *Kotsu, *Kantsu:
			return false
		case *Toitsu:
			pai := &m.Pais()[0]
			if state.YakuhaiFan(pai, player) > 0 {
				return false
			}
		case *Shuntsu:
			pais := m.Pais()
			// Check if the shuntsu contains the horaPai
			if !slices.ContainsFunc(pais, func(p Pai) bool {
				return p.HasSameSymbol(horaPai)
			}) {
				continue
			}
			// The middle (kanchan) wait is not allowed
			if pais[1].HasSameSymbol(horaPai) {
				continue
			}
			// The edge (penchan) wait is not allowed
			if horaPai.Number() == 3 && pais[2].HasSameSymbol(horaPai) {
				continue
			}
			if horaPai.Number() == 7 && pais[0].HasSameSymbol(horaPai) {
				continue
			}
			// If the shuntsu contains the horaPai, it must be a ryanmen wait
			foundShuntsuWithHoraPai = true
		}
	}

	return foundShuntsuWithHoraPai
}

func sumYakuhaiFan(state StateViewer, playerID int, allMentsus []Mentsu) int {
	player := &state.Players()[playerID]
	fan := 0
	for _, m := range allMentsus {
		switch m.(type) {
		case *Kotsu, *Kantsu:
			pai := &m.Pais()[0]
			fan += state.YakuhaiFan(pai, player)
		}
	}
	return fan
}

func isIpeko(allMentsus []Mentsu) bool {
	for i, m1 := range allMentsus {
		if _, ok := m1.(*Shuntsu); !ok {
			continue
		}
		for _, m2 := range allMentsus[i+1:] {
			if _, ok := m2.(*Shuntsu); !ok {
				continue
			}
			if m1.Pais()[0].HasSameSymbol(&m2.Pais()[0]) {
				return true
			}
		}
	}
	return false
}

func isSanshokuDojun(allMentsus []Mentsu) bool {
	typeNumMap := map[rune]map[uint8]bool{
		'm': {},
		'p': {},
		's': {},
	}
	for _, m := range allMentsus {
		shuntsu, ok := m.(*Shuntsu)
		if !ok {
			continue
		}
		pai := shuntsu.Pais()[0]
		t := pai.Type()
		n := pai.Number()
		typeNumMap[t][n] = true
	}
	for n := uint8(1); n <= 7; n++ {
		if typeNumMap['m'][n] && typeNumMap['p'][n] && typeNumMap['s'][n] {
			return true
		}
	}
	return false
}

func isIkkiTsukan(allMentsus []Mentsu) bool {
	typeNumMap := map[rune]map[uint8]bool{
		'm': {},
		'p': {},
		's': {},
	}
	for _, m := range allMentsus {
		shuntsu, ok := m.(*Shuntsu)
		if !ok {
			continue
		}
		pai := shuntsu.Pais()[0]
		t := pai.Type()
		n := pai.Number()
		typeNumMap[t][n] = true
	}
	for _, t := range []rune{'m', 'p', 's'} {
		if typeNumMap[t][1] && typeNumMap[t][4] && typeNumMap[t][7] {
			return true
		}
	}
	return false
}

func isToitoiho(allMentsus []Mentsu) bool {
	for _, m := range allMentsus {
		if _, ok := m.(*Shuntsu); ok {
			return false
		}
	}
	return true
}

func isChiniso(allMentsus []Mentsu) bool {
	var suit rune
	for i, m := range allMentsus {
		t := m.Pais()[0].Type()
		if t == tsupaiType {
			return false
		}
		if i == 0 {
			suit = t
		} else if t != suit {
			return false
		}
	}
	return true
}

func isHoniso(allMentsus []Mentsu) bool {
	var suit rune
	for i, m := range allMentsus {
		t := m.Pais()[0].Type()
		if i == 0 {
			suit = t
		} else if t != suit && t != tsupaiType {
			return false
		}
	}
	return true
}

func isSanshokuDoko(allMentsus []Mentsu) bool {
	typeNumMap := map[rune]map[uint8]bool{
		'm': {},
		'p': {},
		's': {},
	}
	for _, m := range allMentsus {
		_, isKotsu := m.(*Kotsu)
		_, isKantsu := m.(*Kantsu)
		if !isKotsu && !isKantsu {
			continue
		}
		pai := m.Pais()[0]
		t := pai.Type()
		n := pai.Number()
		typeNumMap[t][n] = true
	}
	for n := uint8(1); n <= 9; n++ {
		if typeNumMap['m'][n] && typeNumMap['p'][n] && typeNumMap['s'][n] {
			return true
		}
	}
	return false
}

func isSankantsu(allMentsus []Mentsu) bool {
	numKantsu := 0
	for _, m := range allMentsus {
		if _, ok := m.(*Kantsu); ok {
			numKantsu++
		}
	}
	return numKantsu >= 3
}
