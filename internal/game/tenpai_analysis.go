package game

import "github.com/Apricot-S/mjai-manue-go/internal/base"

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
