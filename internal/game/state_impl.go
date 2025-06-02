package game

import (
	"cmp"
	"fmt"
	"os"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/message"
)

type StateImpl struct {
	players     [NumPlayers]Player
	bakaze      Pai
	kyokuNum    int
	honba       int
	oya         *Player
	chicha      *Player
	doraMarkers []Pai
	numPipais   int

	prevEventType message.Type
	// -1 if prev action is not dahai
	prevDahaiActor   int
	prevDahaiPai     *Pai
	currentEventType message.Type

	playerID int
	// -1 if there is no action
	lastActor      int
	lastActionType message.Type

	// The tiles that cannot be discarded because they would result in swap calling (喰い替え)
	kuikaePais     []Pai
	isRinshanTsumo bool
}

func (s *StateImpl) Players() *[NumPlayers]Player {
	return &s.players
}

func (s *StateImpl) Bakaze() *Pai {
	return &s.bakaze
}

func (s *StateImpl) KyokuNum() int {
	return s.kyokuNum
}

func (s *StateImpl) Honba() int {
	return s.honba
}

func (s *StateImpl) Oya() *Player {
	return s.oya
}

func (s *StateImpl) Chicha() *Player {
	return s.chicha
}

func (s *StateImpl) DoraMarkers() []Pai {
	return s.doraMarkers
}

func (s *StateImpl) NumPipais() int {
	return s.numPipais
}

func (s *StateImpl) Anpais(player *Player) []Pai {
	return slices.Concat(player.sutehais, player.extraAnpais)
}

func (s *StateImpl) VisiblePais(player *Player) []Pai {
	visiblePais := []Pai{}

	for _, p := range s.players {
		visiblePais = slices.Concat(visiblePais, p.ho)
		for _, furo := range p.furos {
			visiblePais = slices.Concat(visiblePais, furo.Pais())
		}
	}

	return slices.Concat(visiblePais, s.doraMarkers, player.tehais)
}

func (s *StateImpl) Doras() []Pai {
	doras := make([]Pai, len(s.doraMarkers))
	for i, d := range s.doraMarkers {
		doras[i] = *d.NextForDora()
	}
	return doras
}

func (s *StateImpl) Jikaze(player *Player) *Pai {
	j := 1 + (4+player.id-s.oya.id)%4
	p, _ := NewPaiWithDetail(tsupaiType, uint8(j), false)
	return p
}

func (s *StateImpl) YakuhaiFan(pai *Pai, player *Player) int {
	if !pai.IsTsupai() {
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

func (s *StateImpl) NextKyoku() (*Pai, int) {
	return getNextKyoku(&s.bakaze, s.kyokuNum)
}

func (s *StateImpl) Turn() int {
	return (NumInitPipais - s.numPipais) / NumPlayers
}

func (s *StateImpl) RankedPlayers() [NumPlayers]Player {
	players := s.players
	slices.SortStableFunc(players[:], func(p1, p2 Player) int {
		if c := cmp.Compare(p1.score, p2.score); c != 0 {
			// Sort descending.
			return -c
		}
		// In case of a tie, sort by closest to chicha.
		return GetPlayerDistance(&p1, s.chicha) - GetPlayerDistance(&p2, s.chicha)
	})

	return players
}

func (s *StateImpl) Print() {
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
