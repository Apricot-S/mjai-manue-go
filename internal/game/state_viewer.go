package game

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
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

func (s *StateImpl) RenderBoard() string {
	var sb strings.Builder

	if _, ok := s.currentEvent.(*inbound.StartGame); ok {
		sb.WriteString("\n")
		sb.WriteString(strings.Repeat("-", 80))
		sb.WriteString("\n")
		return sb.String()
	}

	sb.WriteString(fmt.Sprintf("%s-%d kyoku %d honba  ", s.Bakaze().ToString(), s.KyokuNum(), s.Honba()))
	sb.WriteString(fmt.Sprintf("pipai: %d  ", s.NumPipais()))

	doraMarkers := slices.Collect(func(yield func(string) bool) {
		for _, d := range s.DoraMarkers() {
			if !yield(d.ToString()) {
				return
			}
		}
	})
	sb.WriteString(fmt.Sprintf("dora_marker: %s  ", strings.Join(doraMarkers, " ")))

	sb.WriteString("\n")

	oyaID := s.Oya().ID()
	for i, player := range s.Players() {
		var actorMark string
		if player.ID() == s.lastActor {
			actorMark = "*"
		} else {
			actorMark = " "
		}

		var playerNum string
		if player.ID() == oyaID {
			playerNum = fmt.Sprintf("{%d}", i)
		} else {
			playerNum = fmt.Sprintf("[%d]", i)
		}

		furoStrs := slices.Collect(func(yield func(string) bool) {
			for _, f := range player.Furos() {
				if !yield(f.ToString()) {
					return
				}
			}
		})

		sb.WriteString(fmt.Sprintf(
			"%s%s tehai: %s %s\n",
			actorMark, playerNum, base.PaisToStr(player.Tehais()), strings.Join(furoStrs, " "),
		))

		var hoStr string
		reachHoIndex := player.ReachHoIndex()
		ho := player.Ho()
		if reachHoIndex >= 0 {
			// If the player has declared Riichi, insert "=" just before the Riichi declaration tile.
			hoStr = base.PaisToStr(ho[:reachHoIndex]) + "=" + base.PaisToStr(ho[reachHoIndex:])
		} else {
			hoStr = base.PaisToStr(ho)
		}
		sb.WriteString(fmt.Sprintf("     ho:    %s\n", hoStr))
	}

	sb.WriteString(strings.Repeat("-", 80))
	sb.WriteString("\n")
	return sb.String()
}
