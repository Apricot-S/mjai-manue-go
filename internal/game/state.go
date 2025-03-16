package game

import (
	"cmp"
	"fmt"
	"os"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/message"
)

const (
	numPlayers        = 4
	initScore         = 25_000
	maxNumDoraMarkers = 5
	numInitPipais     = NumIDs*4 - 13*numPlayers - 14
	finalTurn         = numInitPipais / numPlayers
)

func getDistance(p1 *Player, p2 *Player) int {
	return (numPlayers + p1.ID() - p2.ID()) % numPlayers
}

func getNextKyoku(bakaze *Pai, kyokuNum int) (*Pai, int) {
	if kyokuNum == 4 {
		return bakaze.NextForDora(), 1
	}
	return bakaze, kyokuNum + 1
}

type State struct {
	players     [numPlayers]Player
	bakaze      Pai
	kyokuNum    int
	honba       int
	oya         *Player
	chicha      *Player
	doraMarkers []Pai
	numPipais   int

	prevActionType message.Type
	// -1 if prev action is not dahai
	prevDahaiActor    int
	prevDahaiPai      *Pai
	currentActionType message.Type
	tenpais           [numPlayers]bool
}

func (s *State) Players() *[numPlayers]Player {
	return &s.players
}

func (s *State) Bakaze() *Pai {
	return &s.bakaze
}

func (s *State) KyokuNum() int {
	return s.kyokuNum
}

func (s *State) Honba() int {
	return s.honba
}

func (s *State) Oya() *Player {
	return s.oya
}

func (s *State) Chicha() *Player {
	return s.chicha
}

func (s *State) DoraMarkers() []Pai {
	return s.doraMarkers
}

func (s *State) NumPipais() int {
	return s.numPipais
}

func (s *State) Update(event any) error {
	s.prevActionType = s.currentActionType

	panic("unimplemented!")
}

func (s *State) Print() {
	for _, p := range s.players {
		fmt.Fprintf(
			os.Stderr,
			`[%d] tehai: %s
       ho: %s

`,
			p.id,
			PaisToStr(p.tehais),
			PaisToStr(p.ho))
	}
}

func (s *State) Anpais(player *Player) []Pai {
	return slices.Concat(player.sutehais, player.extraAnpais)
}

func (s *State) VisiblePais(player *Player) []Pai {
	visiblePais := []Pai{}

	for _, p := range s.players {
		visiblePais = slices.Concat(visiblePais, p.ho)
		for _, furo := range p.furos {
			visiblePais = slices.Concat(visiblePais, furo.Pais())
		}
	}

	return slices.Concat(visiblePais, s.doraMarkers, player.tehais)
}

func (s *State) Doras() []Pai {
	doras := make([]Pai, len(s.doraMarkers))
	for i, d := range s.doraMarkers {
		doras[i] = *d.NextForDora()
	}
	return doras
}

func (s *State) Jikaze(player *Player) *Pai {
	j := 1 + (4+player.id-s.oya.id)%4
	p, _ := NewPaiWithDetail(tsupaiType, uint8(j), false)
	return p
}

func (s *State) YakuhaiFan(pai *Pai, player *Player) int {
	if pai.Type() != tsupaiType {
		// Suhai
		return 0
	}

	// Jihai
	n := pai.Number()
	if 5 <= n && n <= 7 {
		// Sangenpai
		return 1
	}

	// Kazehai
	fan := 0
	if pai.HasSameSymbol(&s.bakaze) {
		fan++
	}
	if pai.HasSameSymbol(s.Jikaze(player)) {
		fan++
	}
	return fan
}

func (s *State) Turn() int {
	return (numInitPipais - s.numPipais) / numPlayers
}

func (s *State) RankedPlayers() [numPlayers]Player {
	players := s.players
	slices.SortStableFunc(players[:], func(p1, p2 Player) int {
		if c := cmp.Compare(p1.score, p2.score); c != 0 {
			// Sort descending.
			return -c
		}
		// In case of a tie, sort by closest to chicha.
		return getDistance(&p1, s.chicha) - getDistance(&p2, s.chicha)
	})

	return players
}
