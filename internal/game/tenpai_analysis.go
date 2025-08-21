package game

import "github.com/Apricot-S/mjai-manue-go/internal/base"

func IsTenpaiGeneral(ps *base.PaiSet) (bool, error) {
	shanten, _, err := AnalyzeShantenWithOption(ps, 0, 0)
	return shanten <= 0, err
}
