package game

import (
	"fmt"
	"slices"
	"sort"
)

type FuroType int

const (
	Chi FuroType = iota + 1
	Pon
	Daiminkan
	Ankan
	Kakan
)

type Furo struct {
	Type     FuroType
	Taken    Pai
	Consumed []Pai
	Target   *int
	Pais     []Pai
}

func NewFuro(t FuroType, taken *Pai, consumed []Pai, target *int) (*Furo, error) {
	switch t {
	case Ankan:
		if len(consumed) != 4 || taken != nil || target != nil {
			return nil, fmt.Errorf("invalid ankan")
		}
	case Kakan:
		if len(consumed) != 3 || taken == nil {
			return nil, fmt.Errorf("invalid kakan")
		}
	case Daiminkan:
		if len(consumed) != 3 || taken == nil || target == nil {
			return nil, fmt.Errorf("invalid daiminkan")
		}
	case Chi, Pon:
		if len(consumed) != 2 || taken == nil || target == nil {
			return nil, fmt.Errorf("invalid chi or pon")
		}
	default:
		return nil, fmt.Errorf("invalid furo type")
	}

	var tk Pai
	if taken != nil {
		tk = *taken
	}

	c := slices.Clone(consumed)

	var tg *int
	if target != nil {
		tgCopy := *target
		tg = &tgCopy
	}

	var pais Pais = slices.Clone(c)
	if taken != nil {
		pais = append(pais, *taken)
	}
	sort.Sort(pais)

	return &Furo{
		Type:     t,
		Taken:    tk,
		Consumed: c,
		Target:   tg,
		Pais:     pais,
	}, nil
}
