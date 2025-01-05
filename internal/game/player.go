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
	ID                int
	Name              string
	Tehais            Pais
	Furos             []Furo
	Ho                []Pai
	Sutehais          []Pai
	ExtraAnpais       []Pai
	ReachState        ReachState
	ReachHoIndex      *int
	ReachSutehaiIndex *int
	Score             int
	CanDahai          bool
	IsMenzen          bool
}

func NewPlayer(id int, name string, initScore int) (*Player, error) {
	if id < minPlayerID || maxPlayerID < id {
		return nil, fmt.Errorf("Player ID is invalid: %d", id)
	}

	return &Player{
		ID:                id,
		Name:              name,
		Tehais:            make(Pais, 0, initTehaisSize+1), // +1 for tsumo
		Furos:             make([]Furo, 0, maxNumFuro),
		Ho:                nil,
		Sutehais:          nil,
		ExtraAnpais:       nil,
		ReachState:        None,
		ReachHoIndex:      nil,
		ReachSutehaiIndex: nil,
		Score:             initScore,
		CanDahai:          false,
		IsMenzen:          true,
	}, nil
}

func (p *Player) AddExtraAnpais(pai Pai) {
	p.ExtraAnpais = append(p.ExtraAnpais, pai)
}

func (p *Player) OnStartKyoku(tehais []Pai, score *int) error {
	if len(tehais) != initTehaisSize {
		return fmt.Errorf("the length of haipai is not 13: %d", len(tehais))
	}

	p.Tehais = p.Tehais[:len(tehais)]
	copy(p.Tehais, tehais)
	p.Furos = make([]Furo, 0, maxNumFuro)
	p.Ho = nil
	p.Sutehais = nil
	p.ExtraAnpais = nil
	p.ReachState = None
	p.ReachHoIndex = nil
	p.ReachSutehaiIndex = nil
	p.CanDahai = false
	p.IsMenzen = true

	if score != nil {
		p.Score = *score
	}

	return nil
}

func (p *Player) OnTsumo(pai Pai) error {
	if p.CanDahai {
		return fmt.Errorf("it is not in a state to be tsumo")
	}

	p.Tehais = append(p.Tehais, pai)
	p.CanDahai = true
	return nil
}

func (p *Player) OnDahai(pai Pai) error {
	if !p.CanDahai {
		return fmt.Errorf("it is not in a state to be dahai")
	}

	err := p.deleteTehai(&pai)
	if err != nil {
		return fmt.Errorf("failed to delete tehais on dahai: %w", err)
	}

	sort.Sort(p.Tehais)
	p.Ho = append(p.Ho, pai)
	p.Sutehais = append(p.Sutehais, pai)

	if p.ReachState != Accepted {
		p.ExtraAnpais = nil
	}

	p.CanDahai = false
	return nil
}

func (p *Player) OnChiPonKan(furo Furo) error {
	if p.CanDahai {
		return fmt.Errorf("it is not in a state to be chi/pon/kan")
	}

	if p.ReachState != None {
		return fmt.Errorf("chi/pon/kan are not possible during reach")
	}

	numFuro := len(p.Furos)
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

	p.Furos = append(p.Furos, furo)
	p.CanDahai = furo.typ != Daiminkan
	p.IsMenzen = false
	return nil
}

func (p *Player) OnAnkan(furo Furo) error {
	if !p.CanDahai {
		return fmt.Errorf("it is not in a state to be ankan")
	}

	numFuro := len(p.Furos)
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

	p.Furos = append(p.Furos, furo)
	p.CanDahai = false
	return nil
}

func (p *Player) OnKakan(furo Furo) error {
	if !p.CanDahai {
		return fmt.Errorf("it is not in a state to be kakan")
	}

	if p.ReachState != None {
		return fmt.Errorf("kakan is not possible during reach")
	}

	if furo.typ != Kakan {
		return fmt.Errorf("invalid furo for `onKakan`: %v", furo.typ)
	}

	ponIndex := slices.IndexFunc(p.Furos, func(f Furo) bool {
		return slices.Contains(f.pais, furo.taken)
	})
	if ponIndex == -1 {
		return fmt.Errorf("failed to find pon mentsu for kakan: %v", furo)
	}

	err := p.deleteTehai(&furo.taken)
	if err != nil {
		return fmt.Errorf("failed to delete tehais on kakan: %w", err)
	}

	ponMentsu := p.Furos[ponIndex]
	consumed := append(ponMentsu.consumed, furo.taken)
	kanMentsu, err := NewFuro(Kakan, &ponMentsu.taken, consumed, ponMentsu.target)
	if err != nil {
		return fmt.Errorf("failed to create kakan mentsu: %w", err)
	}

	p.Furos[ponIndex] = *kanMentsu
	p.CanDahai = false
	return nil
}

func (p *Player) OnReach() error {
	if !p.CanDahai {
		return fmt.Errorf("it is not in a state to be reach declaration")
	}

	if p.ReachState != None {
		return fmt.Errorf("reach again is not possible during a reach")
	}

	if !p.IsMenzen {
		return fmt.Errorf("reach is not possible after furo")
	}

	p.ReachState = Declared
	return nil
}

func (p *Player) OnReachAccepted(score *int) error {
	if p.CanDahai {
		return fmt.Errorf("it is not in a state to be reach acception")
	}

	if p.ReachState != Declared {
		return fmt.Errorf("reach acception cannot be made except after reach declaration")
	}

	if !p.IsMenzen {
		return fmt.Errorf("reach acception is not possible after furo")
	}

	p.ReachState = Accepted
	p.ReachHoIndex = new(int)
	*p.ReachHoIndex = len(p.Ho) - 1
	p.ReachSutehaiIndex = new(int)
	*p.ReachSutehaiIndex = len(p.Sutehais) - 1

	if score != nil {
		p.Score = *score
	} else {
		p.Score -= kyotakuPoint
	}

	return nil
}

func (p *Player) OnTargeted(furo Furo) error {
	switch furo.typ {
	case Ankan, Kakan:
		return fmt.Errorf("invalid furo for `onTargeted`: %v", furo.typ)
	}

	if *furo.target != p.ID {
		return fmt.Errorf("furo target is not me: %d", *furo.target)
	}

	numHo := len(p.Ho)
	pai := p.Ho[numHo-1]
	if pai != furo.taken {
		return fmt.Errorf("pai %v is not equal to taken %v", pai, furo.taken)
	}
	p.Ho = slices.Delete(p.Ho, numHo-1, numHo)

	return nil
}

func (p *Player) deleteTehai(pai *Pai) error {
	paiIndex := -1
	for i, v := range slices.Backward(p.Tehais) {
		if v == *pai {
			paiIndex = i
			break
		}
	}
	if paiIndex == -1 {
		for i, v := range slices.Backward(p.Tehais) {
			if v == *Unknown {
				paiIndex = i
				break
			}
		}
	}

	if paiIndex == -1 {
		return fmt.Errorf("trying to delete %s which is not in tehais: %v", pai.ToString(), p.Tehais)
	}

	p.Tehais = slices.Delete(p.Tehais, paiIndex, paiIndex+1)
	return nil
}
