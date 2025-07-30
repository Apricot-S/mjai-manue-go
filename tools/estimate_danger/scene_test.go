package main

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/game"
)

func getSceneForTest() *Scene {
	return &Scene{evaluators: defaultEvaluators}
}

func mustPai(name string) *game.Pai {
	p, err := game.NewPaiWithName(name)
	if err != nil {
		panic(err)
	}
	return p
}

func mustPaiSet(pais []game.Pai) *game.PaiSet {
	ps, err := game.NewPaiSet(pais)
	if err != nil {
		panic(err)
	}
	return ps
}

type args struct {
	name string
	pai  *game.Pai
}

type testCase struct {
	name    string
	scene   *Scene
	args    args
	want    bool
	wantErr bool
}

func TestScene_Evaluate_Tsupai(t *testing.T) {
	tests := []testCase{}

	{
		tests = append(tests, testCase{
			name:    "E is tsupai",
			scene:   getSceneForTest(),
			args:    args{name: "tsupai", pai: mustPai("E")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "1p is not tsupai",
			scene:   getSceneForTest(),
			args:    args{name: "tsupai", pai: mustPai("1p")},
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

func TestScene_Evaluate_Suji(t *testing.T) {
	tests := []testCase{}

	scene := getSceneForTest()
	scene.anpaiSet = mustPaiSet([]game.Pai{*mustPai("4p")})

	{
		tests = append(tests, testCase{
			name:    "1p is suji of 4p",
			scene:   scene,
			args:    args{name: "suji", pai: mustPai("1p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "1p is weak suji of 4p",
			scene:   scene,
			args:    args{name: "weak_suji", pai: mustPai("1p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "7p is suji of 4p",
			scene:   scene,
			args:    args{name: "suji", pai: mustPai("7p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "7p is weak suji of 4p",
			scene:   scene,
			args:    args{name: "weak_suji", pai: mustPai("7p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "2p is not suji of 4p",
			scene:   scene,
			args:    args{name: "suji", pai: mustPai("2p")},
			want:    false,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "2p is not weak suji of 4p",
			scene:   scene,
			args:    args{name: "weak_suji", pai: mustPai("2p")},
			want:    false,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "1m is not suji of 4p",
			scene:   scene,
			args:    args{name: "suji", pai: mustPai("1m")},
			want:    false,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "1m is not weak suji of 4p",
			scene:   scene,
			args:    args{name: "weak_suji", pai: mustPai("1m")},
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
