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

	if score != nil {
		p.Score = *score
	}

	return nil
}

func (p *Player) OnTsumo(pai Pai) error {
	numTehai := len(p.Tehais)
	numFuro := len(p.Furos)
	if (numTehai + numFuro*3) != initTehaisSize {
		return fmt.Errorf("it is not in a state to be tsumo")
	}

	p.Tehais = append(p.Tehais, pai)
	return nil
}

func (p *Player) OnDahai(pai Pai) error {
	numTehai := len(p.Tehais)
	numFuro := len(p.Furos)
	if (numTehai + numFuro*3) != (initTehaisSize + 1) {
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
		p.ExtraAnpais = make([]Pai, 0)
	}

	return nil
}

func (p *Player) OnChiPonKan(furo Furo) error {
	numTehai := len(p.Tehais)
	numFuro := len(p.Furos)
	if numFuro >= maxNumFuro {
		return fmt.Errorf("a 5th furo is not possible")
	}
	if (numTehai + numFuro*3) != initTehaisSize {
		return fmt.Errorf("it is not in a state to be chi/pon/kan")
	}

	switch furo.Type {
	case Chi, Pon, Daiminkan:
	default:
		return fmt.Errorf("invalid furo for `onChiPonKan`: %v", furo.Type)
	}

	for _, pai := range furo.Consumed {
		err := p.deleteTehai(&pai)
		if err != nil {
			return fmt.Errorf("failed to delete tehais on chi/pon/kan: %w", err)
		}
	}

	p.Furos = append(p.Furos, furo)
	return nil
}

func (p *Player) OnAnkan(furo Furo) error {
	numTehai := len(p.Tehais)
	numFuro := len(p.Furos)
	if numFuro >= maxNumFuro {
		return fmt.Errorf("a 5th furo is not possible")
	}
	if (numTehai + numFuro*3) != (initTehaisSize + 1) {
		return fmt.Errorf("it is not in a state to be ankan")
	}

	if furo.Type != Ankan {
		return fmt.Errorf("invalid furo for `onAnkan`: %v", furo.Type)
	}

	for _, pai := range furo.Consumed {
		err := p.deleteTehai(&pai)
		if err != nil {
			return fmt.Errorf("failed to delete tehais on ankan: %w", err)
		}
	}

	p.Furos = append(p.Furos, furo)
	return nil
}

func (p *Player) OnKakan(furo Furo) error {
	numTehai := len(p.Tehais)
	numFuro := len(p.Furos)
	if (numTehai + numFuro*3) != (initTehaisSize + 1) {
		return fmt.Errorf("it is not in a state to be kakan")
	}

	if furo.Type != Kakan {
		return fmt.Errorf("invalid furo for `onKakan`: %v", furo.Type)
	}

	err := p.deleteTehai(&furo.Taken)
	if err != nil {
		return fmt.Errorf("failed to delete tehais on kakan: %w", err)
	}

	ponIndex := -1
	for i, f := range p.Furos {
		if slices.Contains(f.Pais, furo.Taken) {
			ponIndex = i
			break
		}
	}
	if ponIndex == -1 {
		return fmt.Errorf("failed to find pon mentsu for kakan: %v", furo)
	}

	ponMentsu := p.Furos[ponIndex]
	ponMentsu.Consumed = append(ponMentsu.Consumed, furo.Taken)
	kanMentsu, err := NewFuro(Kakan, &ponMentsu.Taken, ponMentsu.Consumed, ponMentsu.Target)
	if err != nil {
		return fmt.Errorf("failed to create kakan mentsu: %w", err)
	}

	p.Furos[ponIndex] = *kanMentsu
	return nil
}

func (p *Player) OnReach() error {
	if p.ReachState != None {
		return fmt.Errorf("it is not in a state to be reach declaration")
	}

	p.ReachState = Declared
	return nil
}

func (p *Player) OnReachAccepted(score *int) error {
	if p.ReachState != Declared {
		return fmt.Errorf("it is not in a state to be reach acception")
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
	switch furo.Type {
	case Chi, Pon, Daiminkan:
		pai := p.Ho[len(p.Ho)-1]
		if pai != furo.Taken {
			return fmt.Errorf("pai %v is not equal to taken %v", pai, furo.Taken)
		}
		p.Ho = p.Ho[:len(p.Ho)-1]
	}

	return nil
}

func (p *Player) deleteTehai(pai *Pai) error {
	paiIndex := -1
	for i, v := range p.Tehais {
		if v == *pai {
			paiIndex = i
			break
		}
	}
	if paiIndex == -1 {
		for i, v := range p.Tehais {
			if v == *Unknown {
				paiIndex = i
				break
			}
		}
	}

	if paiIndex == -1 {
		return fmt.Errorf("trying to delete %s which is not in tehais: %v", pai.ToString(), p.Tehais)
	}

	p.Tehais = p.Tehais[:paiIndex+copy(p.Tehais[paiIndex:], p.Tehais[paiIndex+1:])]
	return nil
}
