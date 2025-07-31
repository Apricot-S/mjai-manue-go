package game

import (
	"cmp"
	"fmt"
	"os"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

func (s *StateImpl) Players() *[NumPlayers]base.Player {
	return &s.players
}

func (s *StateImpl) Bakaze() *base.Pai {
	return &s.bakaze
}

func (s *StateImpl) KyokuNum() int {
	return s.kyokuNum
}

func (s *StateImpl) Honba() int {
	return s.honba
}

func (s *StateImpl) Oya() *base.Player {
	return s.oya
}

func (s *StateImpl) Chicha() *base.Player {
	return s.chicha
}

func (s *StateImpl) DoraMarkers() []base.Pai {
	return s.doraMarkers
}

func (s *StateImpl) NumPipais() int {
	return s.numPipais
}

func (s *StateImpl) Anpais(player *base.Player) []base.Pai {
	return slices.Concat(player.Sutehais(), player.ExtraAnpais())
}

func (s *StateImpl) VisiblePais(player *base.Player) []base.Pai {
	visiblePais := []base.Pai{}

	for _, p := range s.players {
		visiblePais = slices.Concat(visiblePais, p.Ho())
		for _, furo := range p.Furos() {
			visiblePais = slices.Concat(visiblePais, furo.Pais())
		}
	}

	return slices.Concat(visiblePais, s.doraMarkers, player.Tehais())
}

func (s *StateImpl) Doras() []base.Pai {
	doras := make([]base.Pai, len(s.doraMarkers))
	for i, d := range s.doraMarkers {
		doras[i] = *d.NextForDora()
	}
	return doras
}

func (s *StateImpl) Jikaze(player *base.Player) *base.Pai {
	j := 1 + (4+player.ID()-s.oya.ID())%4
	p, _ := base.NewPaiWithDetail(base.TsupaiType, uint8(j), false)
	return p
}

func (s *StateImpl) YakuhaiFan(pai *base.Pai, player *base.Player) int {
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

func (s *StateImpl) NextKyoku() (*base.Pai, int) {
	return getNextKyoku(&s.bakaze, s.kyokuNum)
}

func (s *StateImpl) Turn() float64 {
	return float64(NumInitPipais-s.numPipais) / float64(NumPlayers)
}

func (s *StateImpl) RankedPlayers() [NumPlayers]base.Player {
	players := s.players
	slices.SortStableFunc(players[:], func(p1, p2 base.Player) int {
		if c := cmp.Compare(p1.Score(), p2.Score()); c != 0 {
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
			p.ID(),
			base.PaisToStr(p.Tehais()),
			base.PaisToStr(p.Ho()))
	}
}
