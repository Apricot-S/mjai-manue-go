package ai_test

import (
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/game"
)

type MockState struct {
	tehais            []game.Pai
	prereachSutehais  []game.Pai
	postreachSutehais []game.Pai
	doras             []game.Pai
	anpais            []game.Pai
	visiblePais       []game.Pai

	players [4]game.Player
	bakaze  *game.Pai
}

func NewMockState(
	tehais []game.Pai,
	prereachSutehais []game.Pai,
	postreachSutehais []game.Pai,
	doras []game.Pai,
	anpais []game.Pai,
	visiblePais []game.Pai,
	bakaze *game.Pai,
) *MockState {
	s := &MockState{
		tehais:            tehais,
		prereachSutehais:  prereachSutehais,
		postreachSutehais: postreachSutehais,
		doras:             doras,
		anpais:            anpais,
		visiblePais:       visiblePais,
		bakaze:            bakaze,
	}

	s.players[0] = *game.NewPlayerForTest(0, s.tehais, nil, nil, nil, game.None, nil)

	allSutehais := slices.Concat(s.prereachSutehais, s.postreachSutehais)
	var reachSutehaiIndex *int
	if len(s.prereachSutehais) > 0 {
		i := len(s.prereachSutehais) - 1
		reachSutehaiIndex = &i
	}
	s.players[1] = *game.NewPlayerForTest(
		1,
		nil,
		nil,
		allSutehais,
		allSutehais,
		game.None,
		reachSutehaiIndex,
	)

	s.players[2] = *game.NewPlayerForTest(2, nil, nil, nil, nil, game.None, nil)
	s.players[3] = *game.NewPlayerForTest(3, nil, nil, nil, nil, game.None, nil)

	return s
}

func (s *MockState) Players() *[4]game.Player {
	return &s.players
}

func (s *MockState) Bakaze() *game.Pai {
	return s.bakaze
}

func (s *MockState) Jikaze(player *game.Player) *game.Pai {
	p, _ := game.NewPaiWithDetail('t', uint8(1+player.ID()), false)
	return p
}

func (s *MockState) Anpais(player *game.Player) []game.Pai {
	if player.ID() == 1 {
		return s.anpais
	}
	panic("not implemented")
}

func (s *MockState) VisiblePais(player *game.Player) []game.Pai {
	if player.ID() == 0 {
		return s.visiblePais
	}
	panic("not implemented")
}

func (s *MockState) Doras() []game.Pai {
	return s.doras
}

func (s *MockState) YakuhaiFan(pai *game.Pai, player *game.Player) int {
	if pai.IsTsupai() && pai.Number() >= 5 {
		return 1
	}

	fan := 0
	if pai.HasSameSymbol(s.Bakaze()) {
		fan++
	}

	if pai.HasSameSymbol(s.Jikaze(player)) {
		fan++
	}

	return fan
}

func (s *MockState) KyokuNum() int                 { panic("not implemented") }
func (s *MockState) Honba() int                    { panic("not implemented") }
func (s *MockState) Oya() *game.Player             { panic("not implemented") }
func (s *MockState) Chicha() *game.Player          { panic("not implemented") }
func (s *MockState) DoraMarkers() []game.Pai       { panic("not implemented") }
func (s *MockState) NumPipais() int                { panic("not implemented") }
func (s *MockState) Turn() int                     { panic("not implemented") }
func (s *MockState) RankedPlayers() [4]game.Player { panic("not implemented") }
func (s *MockState) Update(event any) error        { panic("not implemented") }
func (s *MockState) Print()                        { panic("not implemented") }
