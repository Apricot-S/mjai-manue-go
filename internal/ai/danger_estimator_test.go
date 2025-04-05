package ai_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/ai"
	"github.com/Apricot-S/mjai-manue-go/internal/game"
	"github.com/go-json-experiment/json/jsontext"
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

	s.players[0] = *game.NewPlayerForTest(0, s.tehais, nil, nil, nil, game.None, -1)

	allSutehais := slices.Concat(s.prereachSutehais, s.postreachSutehais)
	reachSutehaiIndex := -1
	if len(s.prereachSutehais) > 0 {
		reachSutehaiIndex = len(s.prereachSutehais) - 1
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

	s.players[2] = *game.NewPlayerForTest(2, nil, nil, nil, nil, game.None, -1)
	s.players[3] = *game.NewPlayerForTest(3, nil, nil, nil, nil, game.None, -1)

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

func (s *MockState) KyokuNum() int                                       { panic("not implemented") }
func (s *MockState) Honba() int                                          { panic("not implemented") }
func (s *MockState) Oya() *game.Player                                   { panic("not implemented") }
func (s *MockState) Chicha() *game.Player                                { panic("not implemented") }
func (s *MockState) DoraMarkers() []game.Pai                             { panic("not implemented") }
func (s *MockState) NumPipais() int                                      { panic("not implemented") }
func (s *MockState) Turn() int                                           { panic("not implemented") }
func (s *MockState) RankedPlayers() [4]game.Player                       { panic("not implemented") }
func (s *MockState) OnStartGame(event jsontext.Value) error              { panic("not implemented") }
func (s *MockState) Update(event jsontext.Value) error                   { panic("not implemented") }
func (s *MockState) Print()                                              { panic("not implemented") }
func (s *MockState) DahaiCandidates(player *game.Player) []game.Pai      { panic("not implemented") }
func (s *MockState) ReachDahaiCandidates(player *game.Player) []game.Pai { panic("not implemented") }
func (s *MockState) ChiCandidates(player *game.Player) []game.Pai        { panic("not implemented") }
func (s *MockState) PonCandidates(player *game.Player) []game.Pai        { panic("not implemented") }
func (s *MockState) CanHora(player *game.Player) bool                    { panic("not implemented") }

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
		anpais, _ := game.StrToPais("4p")
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
		anpais, _ := game.StrToPais("1p 7p")
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
		anpais, _ := game.StrToPais("5p 4p")
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
		anpais, _ := game.StrToPais("1p")
		state := NewMockState(nil, anpais, nil, nil, anpais, nil, nil)
		scene, _ := ai.NewScene(state, &state.players[0], &state.players[1])
		pin4, _ := game.NewPaiWithName("4p")

		tests = append(tests, testCase{
			name:    "reach_suji true",
			scene:   scene,
			args:    args{name: "reach_suji", pai: pin4},
			want:    true,
			wantErr: false,
		})
	}

	{
		anpais, _ := game.StrToPais("4p E 4s")
		prereachSutehais, _ := game.StrToPais("4p E")
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

	{
		anpais, _ := game.StrToPais("1p")
		state := NewMockState(nil, anpais, nil, nil, anpais, nil, nil)
		scene, _ := ai.NewScene(state, &state.players[0], &state.players[1])
		man2, _ := game.NewPaiWithName("2m")

		for i, v := range [9]bool{false, true, false, false, true, false, false, false, false} {
			pin, _ := game.NewPaiWithName(fmt.Sprintf("%dp", i+1))
			tests = append(tests, testCase{
				name:    fmt.Sprintf("urasuji %v", v),
				scene:   scene,
				args:    args{name: "urasuji", pai: pin},
				want:    v,
				wantErr: false,
			})
		}

		tests = append(tests, testCase{
			name:    "urasuji false",
			scene:   scene,
			args:    args{name: "urasuji", pai: man2},
			want:    false,
			wantErr: false,
		})
	}

	{
		anpais, _ := game.StrToPais("5p")
		state := NewMockState(nil, anpais, nil, nil, anpais, nil, nil)
		scene, _ := ai.NewScene(state, &state.players[0], &state.players[1])
		man1, _ := game.NewPaiWithName("1m")

		for i, v := range [9]bool{true, false, false, true, false, true, false, false, true} {
			pin, _ := game.NewPaiWithName(fmt.Sprintf("%dp", i+1))
			tests = append(tests, testCase{
				name:    fmt.Sprintf("urasuji %v", v),
				scene:   scene,
				args:    args{name: "urasuji", pai: pin},
				want:    v,
				wantErr: false,
			})
		}

		tests = append(tests, testCase{
			name:    "urasuji false",
			scene:   scene,
			args:    args{name: "urasuji", pai: man1},
			want:    false,
			wantErr: false,
		})
	}

	{
		anpais, _ := game.StrToPais("1p 5p")
		prereachSutehais, _ := game.StrToPais("1p")
		state := NewMockState(nil, prereachSutehais, nil, nil, anpais, nil, nil)
		scene, _ := ai.NewScene(state, &state.players[0], &state.players[1])
		pin2, _ := game.NewPaiWithName("2p")

		tests = append(tests, testCase{
			name:    "urasuji false",
			scene:   scene,
			args:    args{name: "urasuji", pai: pin2},
			want:    false,
			wantErr: false,
		})
	}

	{
		anpais, _ := game.StrToPais("1p E S W 1s")
		prereachSutehais, _ := game.StrToPais("1p E S W 1s")
		state := NewMockState(nil, prereachSutehais, nil, nil, anpais, nil, nil)
		scene, _ := ai.NewScene(state, &state.players[0], &state.players[1])
		pin5, _ := game.NewPaiWithName("5p")
		sou5, _ := game.NewPaiWithName("5s")

		tests = append(tests, testCase{
			name:    "early_urasuji true",
			scene:   scene,
			args:    args{name: "early_urasuji", pai: pin5},
			want:    true,
			wantErr: false,
		})
		tests = append(tests, testCase{
			name:    "early_urasuji false",
			scene:   scene,
			args:    args{name: "early_urasuji", pai: sou5},
			want:    false,
			wantErr: false,
		})
		tests = append(tests, testCase{
			name:    "reach_urasuji true",
			scene:   scene,
			args:    args{name: "reach_urasuji", pai: sou5},
			want:    true,
			wantErr: false,
		})
		tests = append(tests, testCase{
			name:    "reach_urasuji false",
			scene:   scene,
			args:    args{name: "reach_urasuji", pai: pin5},
			want:    false,
			wantErr: false,
		})
	}

	{
		anpais, _ := game.StrToPais("1p 6p")
		prereachSutehais, _ := game.StrToPais("1p 6p")
		state := NewMockState(nil, prereachSutehais, nil, nil, anpais, nil, nil)
		scene, _ := ai.NewScene(state, &state.players[0], &state.players[1])
		man2, _ := game.NewPaiWithName("2m")

		for i, v := range [9]bool{false, true, false, false, true, false, false, false, false} {
			pin, _ := game.NewPaiWithName(fmt.Sprintf("%dp", i+1))
			tests = append(tests, testCase{
				name:    fmt.Sprintf("aida4ken %v", v),
				scene:   scene,
				args:    args{name: "aida4ken", pai: pin},
				want:    v,
				wantErr: false,
			})
		}

		tests = append(tests, testCase{
			name:    "aida4ken false",
			scene:   scene,
			args:    args{name: "aida4ken", pai: man2},
			want:    false,
			wantErr: false,
		})
	}

	{
		anpais, _ := game.StrToPais("3p")
		prereachSutehais, _ := game.StrToPais("3p")
		state := NewMockState(nil, prereachSutehais, nil, nil, anpais, nil, nil)
		scene, _ := ai.NewScene(state, &state.players[0], &state.players[1])
		man1, _ := game.NewPaiWithName("1m")

		for i, v := range [9]bool{true, true, false, true, true, false, false, false, false} {
			pin, _ := game.NewPaiWithName(fmt.Sprintf("%dp", i+1))
			tests = append(tests, testCase{
				name:    fmt.Sprintf("matagisuji %v", v),
				scene:   scene,
				args:    args{name: "matagisuji", pai: pin},
				want:    v,
				wantErr: false,
			})
		}

		tests = append(tests, testCase{
			name:    "matagisuji false for different suit",
			scene:   scene,
			args:    args{name: "matagisuji", pai: man1},
			want:    false,
			wantErr: false,
		})
	}

	{
		anpais, _ := game.StrToPais("2p")
		prereachSutehais, _ := game.StrToPais("2p")
		state := NewMockState(nil, prereachSutehais, nil, nil, anpais, nil, nil)
		scene, _ := ai.NewScene(state, &state.players[0], &state.players[1])

		for i, v := range [9]bool{true, false, false, true, false, false, false, false, false} {
			pin, _ := game.NewPaiWithName(fmt.Sprintf("%dp", i+1))
			tests = append(tests, testCase{
				name:    fmt.Sprintf("matagisuji %v", v),
				scene:   scene,
				args:    args{name: "matagisuji", pai: pin},
				want:    v,
				wantErr: false,
			})
		}
	}

	{
		anpais, _ := game.StrToPais("3p 4p")
		prereachSutehais, _ := game.StrToPais("3p")
		state := NewMockState(nil, prereachSutehais, nil, nil, anpais, nil, nil)
		scene, _ := ai.NewScene(state, &state.players[0], &state.players[1])
		pin1, _ := game.NewPaiWithName("1p")

		tests = append(tests, testCase{
			name:    "matagisuji false",
			scene:   scene,
			args:    args{name: "matagisuji", pai: pin1},
			want:    false,
			wantErr: false,
		})
	}

	{
		anpais, _ := game.StrToPais("3p E S 7p W")
		prereachSutehais, _ := game.StrToPais("3p E S 7p W")
		state := NewMockState(nil, prereachSutehais, nil, nil, anpais, nil, nil)
		scene, _ := ai.NewScene(state, &state.players[0], &state.players[1])
		pin1, _ := game.NewPaiWithName("1p")
		pin9, _ := game.NewPaiWithName("9p")

		tests = append(tests, []testCase{
			{
				name:    "late_matagisuji true for 9p",
				scene:   scene,
				args:    args{name: "late_matagisuji", pai: pin9},
				want:    true,
				wantErr: false,
			},
			{
				name:    "early_matagisuji false for 9p",
				scene:   scene,
				args:    args{name: "early_matagisuji", pai: pin9},
				want:    false,
				wantErr: false,
			},
			{
				name:    "early_matagisuji true for 1p",
				scene:   scene,
				args:    args{name: "early_matagisuji", pai: pin1},
				want:    true,
				wantErr: false,
			},
			{
				name:    "late_matagisuji false for 1p",
				scene:   scene,
				args:    args{name: "late_matagisuji", pai: pin1},
				want:    false,
				wantErr: false,
			},
		}...)
	}

	{
		anpais, _ := game.StrToPais("3p E S 7p")
		prereachSutehais, _ := game.StrToPais("3p E S 7p")
		state := NewMockState(nil, prereachSutehais, nil, nil, anpais, nil, nil)
		scene, _ := ai.NewScene(state, &state.players[0], &state.players[1])
		pin1, _ := game.NewPaiWithName("1p")
		pin9, _ := game.NewPaiWithName("9p")

		tests = append(tests, []testCase{
			{
				name:    "reach_matagisuji true",
				scene:   scene,
				args:    args{name: "reach_matagisuji", pai: pin9},
				want:    true,
				wantErr: false,
			},
			{
				name:    "reach_matagisuji false",
				scene:   scene,
				args:    args{name: "reach_matagisuji", pai: pin1},
				want:    false,
				wantErr: false,
			},
		}...)
	}

	{
		anpais, _ := game.StrToPais("1p")
		prereachSutehais, _ := game.StrToPais("1p")
		state := NewMockState(nil, prereachSutehais, nil, nil, anpais, nil, nil)
		scene, _ := ai.NewScene(state, &state.players[0], &state.players[1])
		man3, _ := game.NewPaiWithName("3m")

		for i, v := range [9]bool{false, false, true, false, false, true, false, false, false} {
			pin, _ := game.NewPaiWithName(fmt.Sprintf("%dp", i+1))
			tests = append(tests, testCase{
				name:    fmt.Sprintf("senkisuji %v", v),
				scene:   scene,
				args:    args{name: "senkisuji", pai: pin},
				want:    v,
				wantErr: false,
			})
		}

		tests = append(tests, testCase{
			name:    "senkisuji false for different suit",
			scene:   scene,
			args:    args{name: "senkisuji", pai: man3},
			want:    false,
			wantErr: false,
		})
	}

	// Doesn't count the pai which I'm going to discard.
	{
		visiblePais, _ := game.StrToPais("1p 1p")
		state := NewMockState(nil, nil, nil, nil, nil, visiblePais, nil)
		scene, _ := ai.NewScene(state, &state.players[0], &state.players[1])
		pin1, _ := game.NewPaiWithName("1p")

		tests = append(tests, []testCase{
			{
				name:    "visible>=1 true",
				scene:   scene,
				args:    args{name: "visible>=1", pai: pin1},
				want:    true,
				wantErr: false,
			},
			{
				name:    "visible>=2 false",
				scene:   scene,
				args:    args{name: "visible>=2", pai: pin1},
				want:    false,
				wantErr: false,
			},
		}...)
	}

	{
		visiblePais, _ := game.StrToPais("1p 1p 1p")
		state := NewMockState(nil, nil, nil, nil, nil, visiblePais, nil)
		scene, _ := ai.NewScene(state, &state.players[0], &state.players[1])
		pin1, _ := game.NewPaiWithName("1p")

		tests = append(tests, []testCase{
			{
				name:    "visible>=2 true",
				scene:   scene,
				args:    args{name: "visible>=2", pai: pin1},
				want:    true,
				wantErr: false,
			},
			{
				name:    "visible>=3 false",
				scene:   scene,
				args:    args{name: "visible>=3", pai: pin1},
				want:    false,
				wantErr: false,
			},
		}...)
	}

	{
		visiblePais, _ := game.StrToPais("4p")
		state := NewMockState(nil, nil, nil, nil, nil, visiblePais, nil)
		scene, _ := ai.NewScene(state, &state.players[0], &state.players[1])
		pin1, _ := game.NewPaiWithName("1p")

		tests = append(tests, []testCase{
			{
				name:    "suji_visible<=1 true",
				scene:   scene,
				args:    args{name: "suji_visible<=1", pai: pin1},
				want:    true,
				wantErr: false,
			},
			{
				name:    "suji_visible<=0 false",
				scene:   scene,
				args:    args{name: "suji_visible<=0", pai: pin1},
				want:    false,
				wantErr: false,
			},
		}...)
	}

	{
		visiblePais, _ := game.StrToPais("4p 4p")
		state := NewMockState(nil, nil, nil, nil, nil, visiblePais, nil)
		scene, _ := ai.NewScene(state, &state.players[0], &state.players[1])
		pin1, _ := game.NewPaiWithName("1p")

		tests = append(tests, []testCase{
			{
				name:    "suji_visible<=2 true",
				scene:   scene,
				args:    args{name: "suji_visible<=2", pai: pin1},
				want:    true,
				wantErr: false,
			},
			{
				name:    "suji_visible<=1 false",
				scene:   scene,
				args:    args{name: "suji_visible<=1", pai: pin1},
				want:    false,
				wantErr: false,
			},
		}...)
	}

	// TODO Add test for rest of features.

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
