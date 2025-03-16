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
	typ      FuroType
	taken    Pai
	consumed []Pai
	target   *int
	pais     []Pai
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
		typ:      t,
		taken:    tk,
		consumed: c,
		target:   tg,
		pais:     pais,
	}, nil
}

func (f *Furo) Type() FuroType {
	return f.typ
}

func (f *Furo) Taken() Pai {
	return f.taken
}

func (f *Furo) Consumed() []Pai {
	return f.consumed
}

func (f *Furo) Target() *int {
	return f.target
}

func (f *Furo) Pais() []Pai {
	return f.pais
}

type FuroBase interface {
	Taken() *Pai
	Consumed() []Pai
	Target() *int
	Pais() []Pai
}

type FuroChi struct {
	taken    Pai
	consumed [2]Pai
	target   int
	pais     []Pai
}

func NewChi(taken Pai, consumed [2]Pai, target int) (*FuroChi, error) {
	if target < 0 || target > 3 {
		return nil, fmt.Errorf("chi: invalid target player index (must be 0-3, got: %d)", target)
	}

	var pais Pais = []Pai{taken, consumed[0], consumed[1]}
	sort.Sort(pais)

	return &FuroChi{
		taken:    taken,
		consumed: consumed,
		target:   target,
		pais:     pais,
	}, nil
}

func (c *FuroChi) Taken() *Pai {
	return &c.taken
}

func (c *FuroChi) Consumed() []Pai {
	return c.consumed[:]
}

func (c *FuroChi) Target() *int {
	return &c.target
}

func (c *FuroChi) Pais() []Pai {
	return c.pais
}

type FuroPon struct {
	taken    Pai
	consumed [2]Pai
	target   int
	pais     []Pai
}

func NewPon(taken Pai, consumed [2]Pai, target int) (*FuroPon, error) {
	if target < 0 || target > 3 {
		return nil, fmt.Errorf("pon: invalid target player index (must be 0-3, got: %d)", target)
	}

	var pais Pais = []Pai{taken, consumed[0], consumed[1]}
	sort.Sort(pais)

	return &FuroPon{
		taken:    taken,
		consumed: consumed,
		target:   target,
		pais:     pais,
	}, nil
}

func (p *FuroPon) Taken() *Pai {
	return &p.taken
}

func (p *FuroPon) Consumed() []Pai {
	return p.consumed[:]
}

func (p *FuroPon) Target() *int {
	return &p.target
}

func (p *FuroPon) Pais() []Pai {
	return p.pais
}

type FuroDaiminkan struct {
	taken    Pai
	consumed [3]Pai
	target   int
	pais     []Pai
}

func NewDaiminkan(taken Pai, consumed [3]Pai, target int) (*FuroDaiminkan, error) {
	if target < 0 || target > 3 {
		return nil, fmt.Errorf("daiminkan: invalid target player index (must be 0-3, got: %d)", target)
	}

	var pais Pais = []Pai{taken, consumed[0], consumed[1], consumed[2]}
	sort.Sort(pais)

	return &FuroDaiminkan{
		taken:    taken,
		consumed: consumed,
		target:   target,
		pais:     pais,
	}, nil
}

func (d *FuroDaiminkan) Taken() *Pai {
	return &d.taken
}

func (d *FuroDaiminkan) Consumed() []Pai {
	return d.consumed[:]
}

func (d *FuroDaiminkan) Target() *int {
	return &d.target
}

func (d *FuroDaiminkan) Pais() []Pai {
	return d.pais
}

type FuroAnkan struct {
	consumed [4]Pai
	pais     []Pai
}

func NewAnkan(consumed [4]Pai) (*FuroAnkan, error) {
	var pais Pais = []Pai{consumed[0], consumed[1], consumed[2], consumed[3]}
	sort.Sort(pais)

	return &FuroAnkan{
		consumed: consumed,
		pais:     pais,
	}, nil
}

func (a *FuroAnkan) Taken() *Pai {
	return nil
}

func (a *FuroAnkan) Consumed() []Pai {
	return a.consumed[:]
}

func (a *FuroAnkan) Target() *int {
	return nil
}

func (a *FuroAnkan) Pais() []Pai {
	return a.pais
}

type FuroKakan struct {
	taken    Pai
	consumed [3]Pai
	target   *int
	pais     []Pai
}

func NewKakan(taken Pai, consumed [3]Pai, target *int) (*FuroKakan, error) {
	if (target != nil) && (*target < 0 || *target > 3) {
		return nil, fmt.Errorf("kakan: invalid target player index (must be 0-3, got: %d)", target)
	}

	var pais Pais = []Pai{taken, consumed[0], consumed[1], consumed[2]}
	sort.Sort(pais)

	return &FuroKakan{
		taken:    taken,
		consumed: consumed,
		target:   target,
		pais:     pais,
	}, nil
}

func (k *FuroKakan) Taken() *Pai {
	return &k.taken
}

func (k *FuroKakan) Consumed() []Pai {
	return k.consumed[:]
}

func (k *FuroKakan) Target() *int {
	return k.target
}

func (k *FuroKakan) Pais() []Pai {
	return k.pais
}
