package game

import (
	"fmt"
	"slices"
	"sort"
)

type ReachState int

const (
	None ReachState = iota + 1
	Declared
	Accepted
)

const (
	minPlayerID    = 0
	maxPlayerID    = 3
	initTehaisSize = 13
	maxNumFuro     = 4
	kyotakuPoint   = 1_000
)

type Player struct {
	id                int
	name              string
	tehais            Pais
	furos             []Furo
	ho                []Pai
	sutehais          []Pai
	extraAnpais       []Pai
	reachState        ReachState
	reachHoIndex      *int
	reachSutehaiIndex *int
	score             int
	canDahai          bool
	isMenzen          bool
}

func NewPlayer(id int, name string, initScore int) (*Player, error) {
	if id < minPlayerID || maxPlayerID < id {
		return nil, fmt.Errorf("Player ID is invalid: %d", id)
	}

	return &Player{
		id:                id,
		name:              name,
		tehais:            make(Pais, 0, initTehaisSize+1), // +1 for tsumo
		furos:             make([]Furo, 0, maxNumFuro),
		ho:                nil,
		sutehais:          nil,
		extraAnpais:       nil,
		reachState:        None,
		reachHoIndex:      nil,
		reachSutehaiIndex: nil,
		score:             initScore,
		canDahai:          false,
		isMenzen:          true,
	}, nil
}

func (p *Player) ID() int {
	return p.id
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) Tehais() []Pai {
	return p.tehais
}

func (p *Player) Furos() []Furo {
	return p.furos
}

func (p *Player) Ho() []Pai {
	return p.ho
}

func (p *Player) Sutehais() []Pai {
	return p.sutehais
}

func (p *Player) ExtraAnpais() []Pai {
	return p.extraAnpais
}

func (p *Player) ReachState() ReachState {
	return p.reachState
}

func (p *Player) ReachHoIndex() *int {
	return p.reachHoIndex
}

func (p *Player) ReachSutehaiIndex() *int {
	return p.reachSutehaiIndex
}

func (p *Player) Score() int {
	return p.score
}

func (p *Player) CanDahai() bool {
	return p.canDahai
}

func (p *Player) IsMenzen() bool {
	return p.isMenzen
}

func (p *Player) AddExtraAnpais(pai Pai) {
	p.extraAnpais = append(p.extraAnpais, pai)
}

func (p *Player) OnStartKyoku(tehais []Pai, score *int) error {
	if len(tehais) != initTehaisSize {
		return fmt.Errorf("the length of haipai is not 13: %d", len(tehais))
	}

	p.tehais = p.tehais[:initTehaisSize]
	copy(p.tehais, tehais)
	p.furos = make([]Furo, 0, maxNumFuro)
	p.ho = nil
	p.sutehais = nil
	p.extraAnpais = nil
	p.reachState = None
	p.reachHoIndex = nil
	p.reachSutehaiIndex = nil
	p.canDahai = false
	p.isMenzen = true

	if score != nil {
		p.score = *score
	}

	return nil
}

func (p *Player) OnTsumo(pai Pai) error {
	if p.canDahai {
		return fmt.Errorf("it is not in a state to be tsumo")
	}

	p.tehais = append(p.tehais, pai)
	p.canDahai = true
	return nil
}

func (p *Player) OnDahai(pai Pai) error {
	if !p.canDahai {
		return fmt.Errorf("it is not in a state to be dahai")
	}

	err := p.deleteTehai(&pai)
	if err != nil {
		return fmt.Errorf("failed to delete tehais on dahai: %w", err)
	}

	sort.Sort(p.tehais)
	p.ho = append(p.ho, pai)
	p.sutehais = append(p.sutehais, pai)

	if p.reachState != Accepted {
		p.extraAnpais = nil
	}

	p.canDahai = false
	return nil
}

func (p *Player) OnChiPonKan(furo Furo) error {
	if p.canDahai {
		return fmt.Errorf("it is not in a state to be chi/pon/kan")
	}

	if p.reachState != None {
		return fmt.Errorf("chi/pon/kan are not possible during reach")
	}

	numFuro := len(p.furos)
	if numFuro >= maxNumFuro {
		return fmt.Errorf("a 5th furo is not possible")
	}

	switch furo.typ {
	case Chi, Pon, Daiminkan:
	default:
		return fmt.Errorf("invalid furo for `onChiPonKan`: %v", furo.typ)
	}

	for _, pai := range furo.consumed {
		err := p.deleteTehai(&pai)
		if err != nil {
			return fmt.Errorf("failed to delete tehais on chi/pon/kan: %w", err)
		}
	}

	p.furos = append(p.furos, furo)
	p.canDahai = furo.typ != Daiminkan
	p.isMenzen = false
	return nil
}

func (p *Player) OnAnkan(furo Furo) error {
	if !p.canDahai {
		return fmt.Errorf("it is not in a state to be ankan")
	}

	numFuro := len(p.furos)
	if numFuro >= maxNumFuro {
		return fmt.Errorf("a 5th furo is not possible")
	}

	if furo.typ != Ankan {
		return fmt.Errorf("invalid furo for `onAnkan`: %v", furo.typ)
	}

	for _, pai := range furo.consumed {
		err := p.deleteTehai(&pai)
		if err != nil {
			return fmt.Errorf("failed to delete tehais on ankan: %w", err)
		}
	}

	p.furos = append(p.furos, furo)
	p.canDahai = false
	return nil
}

func (p *Player) OnKakan(furo Furo) error {
	if !p.canDahai {
		return fmt.Errorf("it is not in a state to be kakan")
	}

	if p.reachState != None {
		return fmt.Errorf("kakan is not possible during reach")
	}

	if furo.typ != Kakan {
		return fmt.Errorf("invalid furo for `onKakan`: %v", furo.typ)
	}

	ponIndex := slices.IndexFunc(p.furos, func(f Furo) bool {
		return slices.Contains(f.pais, furo.taken)
	})
	if ponIndex == -1 {
		return fmt.Errorf("failed to find pon mentsu for kakan: %v", furo)
	}

	err := p.deleteTehai(&furo.taken)
	if err != nil {
		return fmt.Errorf("failed to delete tehais on kakan: %w", err)
	}

	ponMentsu := p.furos[ponIndex]
	consumed := append(ponMentsu.consumed, furo.taken)
	kanMentsu, err := NewFuro(Kakan, &ponMentsu.taken, consumed, ponMentsu.target)
	if err != nil {
		return fmt.Errorf("failed to create kakan mentsu: %w", err)
	}

	p.furos[ponIndex] = *kanMentsu
	p.canDahai = false
	return nil
}

func (p *Player) OnReach() error {
	if !p.canDahai {
		return fmt.Errorf("it is not in a state to be reach declaration")
	}

	if p.reachState != None {
		return fmt.Errorf("reach again is not possible during a reach")
	}

	if !p.isMenzen {
		return fmt.Errorf("reach is not possible after furo")
	}

	p.reachState = Declared
	return nil
}

func (p *Player) OnReachAccepted(score *int) error {
	if p.canDahai {
		return fmt.Errorf("it is not in a state to be reach acception")
	}

	if p.reachState != Declared {
		return fmt.Errorf("reach acception cannot be made except after reach declaration")
	}

	if !p.isMenzen {
		return fmt.Errorf("reach acception is not possible after furo")
	}

	p.reachState = Accepted
	p.reachHoIndex = new(int)
	*p.reachHoIndex = len(p.ho) - 1
	p.reachSutehaiIndex = new(int)
	*p.reachSutehaiIndex = len(p.sutehais) - 1

	if score != nil {
		p.score = *score
	} else {
		p.score -= kyotakuPoint
	}

	return nil
}

func (p *Player) OnTargeted(furo Furo) error {
	switch furo.typ {
	case Ankan, Kakan:
		return fmt.Errorf("invalid furo for `onTargeted`: %v", furo.typ)
	}

	if *furo.target != p.id {
		return fmt.Errorf("furo target is not me: %d", *furo.target)
	}

	numHo := len(p.ho)
	pai := p.ho[numHo-1]
	if pai != furo.taken {
		return fmt.Errorf("pai %v is not equal to taken %v", pai, furo.taken)
	}
	p.ho = slices.Delete(p.ho, numHo-1, numHo)

	return nil
}

func (p *Player) deleteTehai(pai *Pai) error {
	paiIndex := -1
	for i, v := range slices.Backward(p.tehais) {
		if v == *pai {
			paiIndex = i
			break
		}
	}
	if paiIndex == -1 {
		for i, v := range slices.Backward(p.tehais) {
			if v == *Unknown {
				paiIndex = i
				break
			}
		}
	}

	if paiIndex == -1 {
		return fmt.Errorf("trying to delete %s which is not in tehais: %v", pai.ToString(), p.tehais)
	}

	p.tehais = slices.Delete(p.tehais, paiIndex, paiIndex+1)
	return nil
}
