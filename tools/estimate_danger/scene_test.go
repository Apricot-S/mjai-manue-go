package main

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

func getSceneForTest() *Scene {
	return &Scene{evaluators: defaultEvaluators}
}

func mustPai(name string) *base.Pai {
	p, err := base.NewPaiWithName(name)
	if err != nil {
		panic(err)
	}
	return p
}

func mustPaiSet(names ...string) *base.PaiSet {
	pais := make([]base.Pai, len(names))
	for i, name := range names {
		pais[i] = *mustPai(name)
	}

	ps, err := base.NewPaiSet(pais)
	if err != nil {
		panic(err)
	}
	return ps
}

type args struct {
	name string
	pai  *base.Pai
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
	scene.anpaiSet = mustPaiSet("4p")

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

func TestScene_Evaluate_NakaSuji(t *testing.T) {
	tests := []testCase{}

	scene := getSceneForTest()
	scene.anpaiSet = mustPaiSet("1p", "7p")

	{
		tests = append(tests, testCase{
			name:    "4p is suji of 17p",
			scene:   scene,
			args:    args{name: "suji", pai: mustPai("4p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "4p is weak suji of 17p",
			scene:   scene,
			args:    args{name: "weak_suji", pai: mustPai("4p")},
			want:    true,
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

func TestScene_Evaluate_KataSuji(t *testing.T) {
	tests := []testCase{}

	scene := getSceneForTest()
	scene.anpaiSet = mustPaiSet("1p")

	{
		tests = append(tests, testCase{
			name:    "4p is not suji of 1p",
			scene:   scene,
			args:    args{name: "suji", pai: mustPai("4p")},
			want:    false,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "4p is weak suji of 1p",
			scene:   scene,
			args:    args{name: "weak_suji", pai: mustPai("4p")},
			want:    true,
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

func TestScene_Evaluate_ReachSuji(t *testing.T) {
	tests := []testCase{}

	scene := getSceneForTest()
	scene.anpaiSet = mustPaiSet("5p", "4p")
	scene.prereachSutehaiSet = mustPaiSet("5p", "4p")
	scene.reachPaiSet = mustPaiSet("4p")

	{
		tests = append(tests, testCase{
			name:    "1p is reach suji of 54p",
			scene:   scene,
			args:    args{name: "reach_suji", pai: mustPai("1p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "2p is not reach suji of 54p",
			scene:   scene,
			args:    args{name: "reach_suji", pai: mustPai("2p")},
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

func TestScene_Evaluate_ReachKataSuji(t *testing.T) {
	tests := []testCase{}

	scene := getSceneForTest()
	scene.anpaiSet = mustPaiSet("1p")
	scene.prereachSutehaiSet = mustPaiSet("1p")
	scene.reachPaiSet = mustPaiSet("1p")

	{
		tests = append(tests, testCase{
			name:    "4p is reach suji of 1p",
			scene:   scene,
			args:    args{name: "reach_suji", pai: mustPai("4p")},
			want:    true,
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
