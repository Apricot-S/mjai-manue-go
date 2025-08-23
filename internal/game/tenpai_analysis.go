package game

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

func IsTenpaiGeneral(ps *base.PaiSet) (bool, error) {
	shanten, _, err := AnalyzeShantenWithOption(ps, 0, 0)
	return shanten <= 0, err
}

func IsTenpaiAll(ps *base.PaiSet) (bool, error) {
	shanten, _, err := AnalyzeShantenWithOption(ps, 0, 0)
	if err != nil {
		return false, err
	}
	if shanten <= 0 {
		return true, nil
	}

	shanten, err = AnalyzeShantenChitoitsu(ps)
	if err != nil {
		return false, err
	}
	if shanten <= 0 {
		return true, nil
	}

	shanten, err = AnalyzeShantenKokushimuso(ps)
	if err != nil {
		return false, err
	}
	if shanten <= 0 {
		return true, nil
	}

	return false, nil
}

func GetWaitedPaisAll(ps *base.PaiSet) (*base.PaiSet, error) {
	numPais, err := countPais(ps)
	if err != nil {
		return nil, err
	}
	if (numPais % 3) != 1 {
		return nil, fmt.Errorf("the waited tiles cannot be calculated if there is a tsumo tile")
	}

	isTenpai, err := IsTenpaiAll(ps)
	if err != nil {
		return nil, err
	}
	if !isTenpai {
		return nil, nil
	}

	waited, _ := base.NewPaiSet(nil)
	for i, c := range ps {
		if c >= 4 {
			continue
		}

		ps[i]++
		isHora, err := IsHoraForm(ps)
		ps[i]--
		if err != nil {
			return nil, err
		}

		if isHora {
			waited[i] = 1
		}
	}

	return waited, nil
}
