package base

import (
	"fmt"
	"sort"
)

type Furo interface {
	Taken() *Pai
	Consumed() []Pai
	Target() *int
	Pais() []Pai
	ToMentsu() Mentsu
}

type Chi struct {
	taken    Pai
	consumed [2]Pai
	target   int
	pais     []Pai
}

func NewChi(taken Pai, consumed [2]Pai, target int) (*Chi, error) {
	if target < 0 || target > 3 {
		return nil, fmt.Errorf("chi: invalid target player index (must be 0-3, got: %d)", target)
	}

	var pais Pais = []Pai{taken, consumed[0], consumed[1]}
	sort.Sort(pais)

	return &Chi{
		taken:    taken,
		consumed: consumed,
		target:   target,
		pais:     pais,
	}, nil
}

func (c *Chi) Taken() *Pai {
	return &c.taken
}

func (c *Chi) Consumed() []Pai {
	return c.consumed[:]
}

func (c *Chi) Target() *int {
	return &c.target
}

func (c *Chi) Pais() []Pai {
	return c.pais
}

func (c *Chi) ToMentsu() Mentsu {
	return NewShuntsu(c.pais[0], c.pais[1], c.pais[2])
}

type Pon struct {
	taken    Pai
	consumed [2]Pai
	target   int
	pais     []Pai
}

func NewPon(taken Pai, consumed [2]Pai, target int) (*Pon, error) {
	if target < 0 || target > 3 {
		return nil, fmt.Errorf("pon: invalid target player index (must be 0-3, got: %d)", target)
	}

	var pais Pais = []Pai{taken, consumed[0], consumed[1]}
	sort.Sort(pais)

	return &Pon{
		taken:    taken,
		consumed: consumed,
		target:   target,
		pais:     pais,
	}, nil
}

func (p *Pon) Taken() *Pai {
	return &p.taken
}

func (p *Pon) Consumed() []Pai {
	return p.consumed[:]
}

func (p *Pon) Target() *int {
	return &p.target
}

func (p *Pon) Pais() []Pai {
	return p.pais
}

func (p *Pon) ToMentsu() Mentsu {
	return NewKotsu(p.pais[0], p.pais[1], p.pais[2])
}

type Daiminkan struct {
	taken    Pai
	consumed [3]Pai
	target   int
	pais     []Pai
}

func NewDaiminkan(taken Pai, consumed [3]Pai, target int) (*Daiminkan, error) {
	if target < 0 || target > 3 {
		return nil, fmt.Errorf("daiminkan: invalid target player index (must be 0-3, got: %d)", target)
	}

	var pais Pais = []Pai{taken, consumed[0], consumed[1], consumed[2]}
	sort.Sort(pais)

	return &Daiminkan{
		taken:    taken,
		consumed: consumed,
		target:   target,
		pais:     pais,
	}, nil
}

func (d *Daiminkan) Taken() *Pai {
	return &d.taken
}

func (d *Daiminkan) Consumed() []Pai {
	return d.consumed[:]
}

func (d *Daiminkan) Target() *int {
	return &d.target
}

func (d *Daiminkan) Pais() []Pai {
	return d.pais
}

func (d *Daiminkan) ToMentsu() Mentsu {
	return NewKantsu(d.pais[0], d.pais[1], d.pais[2], d.pais[3])
}

type Ankan struct {
	consumed [4]Pai
	pais     []Pai
}

func NewAnkan(consumed [4]Pai) (*Ankan, error) {
	var pais Pais = []Pai{consumed[0], consumed[1], consumed[2], consumed[3]}
	sort.Sort(pais)

	return &Ankan{
		consumed: consumed,
		pais:     pais,
	}, nil
}

func (a *Ankan) Taken() *Pai {
	return nil
}

func (a *Ankan) Consumed() []Pai {
	return a.consumed[:]
}

func (a *Ankan) Target() *int {
	return nil
}

func (a *Ankan) Pais() []Pai {
	return a.pais
}

func (a *Ankan) ToMentsu() Mentsu {
	return NewKantsu(a.pais[0], a.pais[1], a.pais[2], a.pais[3])
}

type Kakan struct {
	// Taken tile for the original Pon (nil in event data)
	taken *Pai
	// 2 tiles from the player's hand used for the original Pon
	consumed [2]Pai
	// Tile added from hand to form Kakan
	added Pai
	// Source player index for original Pon (nil in event data)
	target *int
	pais   []Pai
}

func NewKakan(taken Pai, consumed [2]Pai, added Pai, target int) (*Kakan, error) {
	if target < 0 || target > 3 {
		return nil, fmt.Errorf("kakan: invalid target player index (must be 0-3, got: %d)", target)
	}

	var pais Pais = []Pai{taken, consumed[0], consumed[1], added}
	sort.Sort(pais)

	return &Kakan{
		taken:    &taken,
		consumed: consumed,
		added:    added,
		target:   &target,
		pais:     pais,
	}, nil
}

// NewKakanFromEvent constructs a Kakan from event data.
// The event provides `added` and `consumed` (3 tiles from the Pon meld, including `taken`).
// Internally, one tile from consumed is heuristically assigned as `taken`.
// Note: The exact identity of `taken` is not guaranteed.
func NewKakanFromEvent(added Pai, consumed [3]Pai) (*Kakan, error) {
	var pais Pais = []Pai{added, consumed[0], consumed[1], consumed[2]}
	sort.Sort(pais)

	// Heuristic: The last tile in the consumed is considered `taken`.
	// This is a simplification and may not reflect the actual game state.
	taken := &consumed[2]

	return &Kakan{
		taken:    taken,
		consumed: [2]Pai{consumed[0], consumed[1]},
		added:    added,
		target:   nil, // Target is not provided in event data
		pais:     pais,
	}, nil
}

func (k *Kakan) Taken() *Pai {
	return k.taken
}

func (k *Kakan) Consumed() []Pai {
	return k.consumed[:]
}

func (k *Kakan) Added() *Pai {
	return &k.added
}

func (k *Kakan) Target() *int {
	return k.target
}

func (k *Kakan) Pais() []Pai {
	return k.pais
}

func (k *Kakan) ToMentsu() Mentsu {
	return NewKantsu(k.pais[0], k.pais[1], k.pais[2], k.pais[3])
}

func IsKuikae(furo Furo, dahai *Pai) bool {
	taken := furo.Taken()
	if dahai.HasSameSymbol(taken) {
		return true
	}

	chi, isChi := furo.(*Chi)
	if !isChi {
		// There is no suji swap calling for pon or daiminkan
		return false
	}

	pais := chi.Pais()
	if dahai.Type() != pais[0].Type() {
		return false
	}

	if taken.Number() == pais[1].Number() {
		// There is no suji swap calling for kanchan chi
		return false
	}

	number := dahai.Number()
	if number > 3 && number-3 == pais[0].Number() {
		return true
	}
	if number < 7 && number+3 == pais[2].Number() {
		return true
	}
	return false
}
