package base

import (
	"fmt"
	"slices"
	"sort"
)

type ReachState int

const (
	NotReach ReachState = iota + 1
	ReachDeclared
	ReachAccepted
)

const (
	MinPlayerID    = 0
	MaxPlayerID    = 3
	InitTehaisSize = 13
	MaxNumFuro     = 4
	// Reference: <https://note.com/daku_longyi/n/n51fe08566f1b>
	maxNumHo       = 24
	maxNumSutehais = 27
	KyotakuPoint   = 1_000
)

type Player struct {
	// Player ID
	// 0: the dealer at the start of a game (起家)
	// 1: the right next to the 0th seat (起家の下家)
	// 2: the one across from the 0th seat (起家の対面)
	// 3: the left next to the 0th seat (起家の上家)
	id int
	// Player name
	name string
	// Hand (手牌) excluding the furos (副露)
	// The last element is for tsumo.
	tehais Pais
	// Furos (副露)
	furos []Furo
	// River (河)
	// It does not include the tiles that have been called.
	ho []Pai
	// Discarded tiles (捨て牌)
	// It includes the tiles that have been called.
	sutehais []Pai
	// Extra safe tiles (安全牌)
	// The tiles that are safe in the same turn and the tiles that are safe after reach.
	extraAnpais []Pai
	// Reach state
	reachState ReachState
	// The index of the tile that was declared as reach in the river.
	// It is -1 if the player has not declared reach.
	reachHoIndex int
	// The index of the tile that was declared as reach in the discarded tiles.
	// It is -1 if the player has not declared reach.
	reachSutehaiIndex int
	// Player score
	score int
	// Whether the player can discard a tile (打牌)
	canDahai bool
	// Whether the player hand is concealed (門前)
	isMenzen bool
}

func NewPlayer(id int, name string, initScore int) (*Player, error) {
	if id < MinPlayerID || MaxPlayerID < id {
		return nil, fmt.Errorf("player ID is invalid: %d", id)
	}

	return &Player{
		id:                id,
		name:              name,
		tehais:            make(Pais, 0, InitTehaisSize+1), // +1 for tsumo
		furos:             make([]Furo, 0, MaxNumFuro),
		ho:                make([]Pai, 0, maxNumHo),
		sutehais:          make([]Pai, 0, maxNumSutehais),
		extraAnpais:       nil,
		reachState:        NotReach,
		reachHoIndex:      -1,
		reachSutehaiIndex: -1,
		score:             initScore,
		canDahai:          false,
		isMenzen:          true,
	}, nil
}

// For test only.
func NewPlayerForTest(
	id int,
	name string,
	tehais []Pai,
	furos []Furo,
	ho []Pai,
	sutehais []Pai,
	extraAnpais []Pai,
	reachState ReachState,
	reachHoIndex int,
	reachSutehaiIndex int,
	score int,
	canDahai bool,
	isMenzen bool,
) *Player {
	return &Player{
		id:                id,
		name:              name,
		tehais:            tehais,
		furos:             furos,
		ho:                ho,
		sutehais:          sutehais,
		extraAnpais:       extraAnpais,
		reachState:        reachState,
		reachHoIndex:      reachHoIndex,
		reachSutehaiIndex: reachSutehaiIndex,
		score:             score,
		canDahai:          canDahai,
		isMenzen:          isMenzen,
	}
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

func (p *Player) ReachHoIndex() int {
	return p.reachHoIndex
}

func (p *Player) ReachSutehaiIndex() int {
	return p.reachSutehaiIndex
}

func (p *Player) Score() int {
	return p.score
}

func (p *Player) SetScore(score int) {
	p.score = score
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
	if len(tehais) != InitTehaisSize {
		return fmt.Errorf("the length of haipai is not 13: %d", len(tehais))
	}

	p.tehais = p.tehais[:InitTehaisSize]
	copy(p.tehais, tehais)
	sort.Sort(p.tehais)
	p.furos = make([]Furo, 0, MaxNumFuro)
	p.ho = make([]Pai, 0, maxNumHo)
	p.sutehais = make([]Pai, 0, maxNumSutehais)
	p.extraAnpais = nil
	p.reachState = NotReach
	p.reachHoIndex = -1
	p.reachSutehaiIndex = -1
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

	if p.reachState != ReachAccepted {
		p.extraAnpais = nil
	}

	p.canDahai = false
	return nil
}

func (p *Player) OnChi(furo *Chi) error {
	if p.canDahai {
		return fmt.Errorf("it is not in a state to be chi")
	}

	if p.reachState != NotReach {
		return fmt.Errorf("chi is not possible during reach")
	}

	numFuro := len(p.furos)
	if numFuro >= MaxNumFuro {
		return fmt.Errorf("a 5th furo is not possible")
	}

	for _, pai := range furo.Consumed() {
		err := p.deleteTehai(&pai)
		if err != nil {
			return fmt.Errorf("failed to delete tehais on chi: %w", err)
		}
	}

	p.furos = append(p.furos, furo)
	p.canDahai = true
	p.isMenzen = false
	return nil
}

func (p *Player) OnPon(furo *Pon) error {
	if p.canDahai {
		return fmt.Errorf("it is not in a state to be pon")
	}

	if p.reachState != NotReach {
		return fmt.Errorf("pon is not possible during reach")
	}

	numFuro := len(p.furos)
	if numFuro >= MaxNumFuro {
		return fmt.Errorf("a 5th furo is not possible")
	}

	for _, pai := range furo.Consumed() {
		err := p.deleteTehai(&pai)
		if err != nil {
			return fmt.Errorf("failed to delete tehais on pon: %w", err)
		}
	}

	p.furos = append(p.furos, furo)
	p.canDahai = true
	p.isMenzen = false
	return nil
}

func (p *Player) OnDaiminkan(furo *Daiminkan) error {
	if p.canDahai {
		return fmt.Errorf("it is not in a state to be daiminkan")
	}

	if p.reachState != NotReach {
		return fmt.Errorf("daiminkan is not possible during reach")
	}

	numFuro := len(p.furos)
	if numFuro >= MaxNumFuro {
		return fmt.Errorf("a 5th furo is not possible")
	}

	for _, pai := range furo.Consumed() {
		err := p.deleteTehai(&pai)
		if err != nil {
			return fmt.Errorf("failed to delete tehais on daiminkan: %w", err)
		}
	}

	p.furos = append(p.furos, furo)
	p.canDahai = false
	p.isMenzen = false
	return nil
}

func (p *Player) OnAnkan(furo *Ankan) error {
	if furo == nil {
		return fmt.Errorf("furo is nil")
	}

	if !p.canDahai {
		return fmt.Errorf("it is not in a state to be ankan")
	}

	numFuro := len(p.furos)
	if numFuro >= MaxNumFuro {
		return fmt.Errorf("a 5th furo is not possible")
	}

	for _, pai := range furo.Consumed() {
		err := p.deleteTehai(&pai)
		if err != nil {
			return fmt.Errorf("failed to delete tehais on ankan: %w", err)
		}
	}

	p.furos = append(p.furos, furo)
	p.canDahai = false
	return nil
}

func (p *Player) OnKakan(furo *Kakan) error {
	if furo == nil {
		return fmt.Errorf("furo is nil")
	}

	if !p.canDahai {
		return fmt.Errorf("it is not in a state to be kakan")
	}

	if p.reachState != NotReach {
		return fmt.Errorf("kakan is not possible during reach")
	}

	ponIndex := slices.IndexFunc(p.furos, func(f Furo) bool {
		p, isPon := f.(*Pon)
		if !isPon {
			return false
		}
		return slices.Contains(p.Pais(), *furo.Taken().RemoveRed())
	})
	if ponIndex == -1 {
		return fmt.Errorf("failed to find pon mentsu for kakan: %v", furo)
	}

	err := p.deleteTehai(furo.Taken())
	if err != nil {
		return fmt.Errorf("failed to delete tehais on kakan: %w", err)
	}

	ponMentsu := p.furos[ponIndex]
	consumed := append(ponMentsu.Consumed(), *furo.Taken())
	kanMentsu, err := NewKakan(*ponMentsu.Taken(), [3]Pai(consumed), ponMentsu.Target())
	if err != nil {
		return fmt.Errorf("failed to create kakan mentsu: %w", err)
	}

	p.furos[ponIndex] = kanMentsu
	p.canDahai = false
	return nil
}

func (p *Player) OnReach() error {
	if !p.canDahai {
		return fmt.Errorf("it is not in a state to be reach declaration")
	}

	if p.reachState != NotReach {
		return fmt.Errorf("reach again is not possible during a reach")
	}

	if !p.isMenzen {
		return fmt.Errorf("reach is not possible after furo")
	}

	p.reachState = ReachDeclared
	return nil
}

func (p *Player) OnReachAccepted(score *int) error {
	if p.canDahai {
		return fmt.Errorf("it is not in a state to be reach acception")
	}

	if p.reachState != ReachDeclared {
		return fmt.Errorf("reach acception cannot be made except after reach declaration")
	}

	if !p.isMenzen {
		return fmt.Errorf("reach acception is not possible after furo")
	}

	p.reachState = ReachAccepted
	p.reachHoIndex = len(p.ho) - 1
	p.reachSutehaiIndex = len(p.sutehais) - 1

	if score != nil {
		p.score = *score
	} else {
		p.score -= KyotakuPoint
	}

	return nil
}

func (p *Player) OnTargeted(furo Furo) error {
	switch furo.(type) {
	case *Ankan, *Kakan:
		return fmt.Errorf("invalid furo for `onTargeted`: %v", furo)
	}

	if *furo.Target() != p.id {
		return fmt.Errorf("furo target is not me: %d", *furo.Target())
	}

	numHo := len(p.ho)
	pai := p.ho[numHo-1]
	if pai != *furo.Taken() {
		return fmt.Errorf("pai %v is not equal to taken %v", pai, *furo.Taken())
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

	// If the pai is not found, check if it is an unknown tile.
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
