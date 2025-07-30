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
			name:    "is tsupai",
			scene:   getSceneForTest(),
			args:    args{name: "tsupai", pai: mustPai("E")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "is not tsupai",
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
