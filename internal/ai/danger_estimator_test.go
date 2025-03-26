package ai_test

import (
	"slices"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/ai"
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

func TestScene_Evaluate(t *testing.T) {
	type args struct {
		name string
		pai  *game.Pai
	}
	type testCase struct {
		name    string
		scene   *ai.Scene
		args    args
		want    bool
		wantErr bool
	}
	tests := []testCase{}

	{
		state := NewMockState(nil, nil, nil, nil, nil, nil, nil)
		scene, _ := ai.NewScene(state, &state.players[0], &state.players[1])
		east, _ := game.NewPaiWithName("E")
		man1, _ := game.NewPaiWithName("1m")

		tests = append(tests, testCase{
			name:    "tsupai true",
			scene:   scene,
			args:    args{name: "tsupai", pai: east},
			want:    true,
			wantErr: false,
		})
		tests = append(tests, testCase{
			name:    "tsupai false",
			scene:   scene,
			args:    args{name: "tsupai", pai: man1},
			want:    false,
			wantErr: false,
		})
	}

	{
		pin4, _ := game.NewPaiWithName("4p")
		anpais := []game.Pai{*pin4}
		state := NewMockState(nil, nil, nil, nil, anpais, nil, nil)
		scene, _ := ai.NewScene(state, &state.players[0], &state.players[1])
		pin1, _ := game.NewPaiWithName("1p")
		pin7, _ := game.NewPaiWithName("7p")
		pin2, _ := game.NewPaiWithName("2p")
		man1, _ := game.NewPaiWithName("1m")

		tests = append(tests, testCase{
			name:    "suji true",
			scene:   scene,
			args:    args{name: "suji", pai: pin1},
			want:    true,
			wantErr: false,
		})
		tests = append(tests, testCase{
			name:    "weak_suji true",
			scene:   scene,
			args:    args{name: "weak_suji", pai: pin1},
			want:    true,
			wantErr: false,
		})
		tests = append(tests, testCase{
			name:    "suji true 2",
			scene:   scene,
			args:    args{name: "suji", pai: pin7},
			want:    true,
			wantErr: false,
		})
		tests = append(tests, testCase{
			name:    "weak_suji true 2",
			scene:   scene,
			args:    args{name: "weak_suji", pai: pin7},
			want:    true,
			wantErr: false,
		})
		tests = append(tests, testCase{
			name:    "suji false",
			scene:   scene,
			args:    args{name: "suji", pai: pin2},
			want:    false,
			wantErr: false,
		})
		tests = append(tests, testCase{
			name:    "weak_suji false",
			scene:   scene,
			args:    args{name: "weak_suji", pai: pin2},
			want:    false,
			wantErr: false,
		})
		tests = append(tests, testCase{
			name:    "suji false 2",
			scene:   scene,
			args:    args{name: "suji", pai: man1},
			want:    false,
			wantErr: false,
		})
		tests = append(tests, testCase{
			name:    "weak_suji false 2",
			scene:   scene,
			args:    args{name: "weak_suji", pai: man1},
			want:    false,
			wantErr: false,
		})
	}

	{
		pin1, _ := game.NewPaiWithName("1p")
		pin7, _ := game.NewPaiWithName("7p")
		anpais := []game.Pai{*pin1, *pin7}
		state := NewMockState(nil, nil, nil, nil, anpais, nil, nil)
		scene, _ := ai.NewScene(state, &state.players[0], &state.players[1])
		pin4, _ := game.NewPaiWithName("4p")

		tests = append(tests, testCase{
			name:    "suji true",
			scene:   scene,
			args:    args{name: "suji", pai: pin4},
			want:    true,
			wantErr: false,
		})
		tests = append(tests, testCase{
			name:    "weak_suji true 2",
			scene:   scene,
			args:    args{name: "weak_suji", pai: pin4},
			want:    true,
			wantErr: false,
		})
	}

	{
		pin5, _ := game.NewPaiWithName("5p")
		pin4, _ := game.NewPaiWithName("4p")
		anpais := []game.Pai{*pin5, *pin4}
		state := NewMockState(nil, anpais, nil, nil, anpais, nil, nil)
		scene, _ := ai.NewScene(state, &state.players[0], &state.players[1])
		pin1, _ := game.NewPaiWithName("1p")
		pin2, _ := game.NewPaiWithName("2p")

		tests = append(tests, testCase{
			name:    "reach_suji true",
			scene:   scene,
			args:    args{name: "reach_suji", pai: pin1},
			want:    true,
			wantErr: false,
		})
		tests = append(tests, testCase{
			name:    "reach_suji false",
			scene:   scene,
			args:    args{name: "reach_suji", pai: pin2},
			want:    false,
			wantErr: false,
		})
	}

	{
		pin1, _ := game.NewPaiWithName("1p")
		anpais := []game.Pai{*pin1}
		state := NewMockState(nil, anpais, nil, nil, anpais, nil, nil)
		scene, _ := ai.NewScene(state, &state.players[0], &state.players[1])
		pin4, _ := game.NewPaiWithName("4p")

		tests = append(tests, testCase{
			name:    "suji true",
			scene:   scene,
			args:    args{name: "reach_suji", pai: pin4},
			want:    true,
			wantErr: false,
		})
	}

	{
		pin4, _ := game.NewPaiWithName("4p")
		east, _ := game.NewPaiWithName("E")
		sou4, _ := game.NewPaiWithName("4s")
		anpais := []game.Pai{*pin4, *east, *sou4}
		prereachSutehais := []game.Pai{*pin4, *east}
		state := NewMockState(nil, prereachSutehais, nil, nil, anpais, nil, nil)
		scene, _ := ai.NewScene(state, &state.players[0], &state.players[1])
		pin1, _ := game.NewPaiWithName("1p")
		sou1, _ := game.NewPaiWithName("1s")

		tests = append(tests, testCase{
			name:    "prereach_suji true",
			scene:   scene,
			args:    args{name: "prereach_suji", pai: pin1},
			want:    true,
			wantErr: false,
		})
		tests = append(tests, testCase{
			name:    "prereach_suji false",
			scene:   scene,
			args:    args{name: "prereach_suji", pai: sou1},
			want:    false,
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.scene.Evaluate(tt.args.name, tt.args.pai)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scene.Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Scene.Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}
