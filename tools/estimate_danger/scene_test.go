package main

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/testutil"
)

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

	scene, _ := NewSceneFromParams(nil, nil, nil, nil, nil, nil, nil)

	{
		tests = append(tests, testCase{
			name:    "E is tsupai",
			scene:   scene,
			args:    args{name: "tsupai", pai: testutil.MustPai("E")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "1p is not tsupai",
			scene:   scene,
			args:    args{name: "tsupai", pai: testutil.MustPai("1p")},
			want:    false,
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.scene.evaluate(tt.args.name, tt.args.pai)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scene.evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Scene.evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScene_Evaluate_Suji(t *testing.T) {
	tests := []testCase{}

	scene, _ := NewSceneFromParams(nil, testutil.MustPais("4p"), nil, nil, nil, nil, nil)

	{
		tests = append(tests, testCase{
			name:    "1p is suji of 4p",
			scene:   scene,
			args:    args{name: "suji", pai: testutil.MustPai("1p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "1p is weak suji of 4p",
			scene:   scene,
			args:    args{name: "weak_suji", pai: testutil.MustPai("1p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "7p is suji of 4p",
			scene:   scene,
			args:    args{name: "suji", pai: testutil.MustPai("7p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "7p is weak suji of 4p",
			scene:   scene,
			args:    args{name: "weak_suji", pai: testutil.MustPai("7p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "2p is not suji of 4p",
			scene:   scene,
			args:    args{name: "suji", pai: testutil.MustPai("2p")},
			want:    false,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "2p is not weak suji of 4p",
			scene:   scene,
			args:    args{name: "weak_suji", pai: testutil.MustPai("2p")},
			want:    false,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "1m is not suji of 4p",
			scene:   scene,
			args:    args{name: "suji", pai: testutil.MustPai("1m")},
			want:    false,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "1m is not weak suji of 4p",
			scene:   scene,
			args:    args{name: "weak_suji", pai: testutil.MustPai("1m")},
			want:    false,
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.scene.evaluate(tt.args.name, tt.args.pai)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scene.evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Scene.evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScene_Evaluate_NakaSuji(t *testing.T) {
	tests := []testCase{}

	scene, _ := NewSceneFromParams(nil, testutil.MustPais("1p", "7p"), nil, nil, nil, nil, nil)

	{
		tests = append(tests, testCase{
			name:    "4p is suji of 17p",
			scene:   scene,
			args:    args{name: "suji", pai: testutil.MustPai("4p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "4p is weak suji of 17p",
			scene:   scene,
			args:    args{name: "weak_suji", pai: testutil.MustPai("4p")},
			want:    true,
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.scene.evaluate(tt.args.name, tt.args.pai)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scene.evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Scene.evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScene_Evaluate_KataSuji(t *testing.T) {
	tests := []testCase{}

	scene, _ := NewSceneFromParams(nil, testutil.MustPais("1p"), nil, nil, nil, nil, nil)

	{
		tests = append(tests, testCase{
			name:    "4p is not suji of 1p",
			scene:   scene,
			args:    args{name: "suji", pai: testutil.MustPai("4p")},
			want:    false,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "4p is weak suji of 1p",
			scene:   scene,
			args:    args{name: "weak_suji", pai: testutil.MustPai("4p")},
			want:    true,
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.scene.evaluate(tt.args.name, tt.args.pai)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scene.evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Scene.evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScene_Evaluate_ReachSuji(t *testing.T) {
	tests := []testCase{}

	scene, _ := NewSceneFromParams(nil, testutil.MustPais("5p", "4p"), nil, nil, testutil.MustPais("5p", "4p"), nil, nil)

	{
		tests = append(tests, testCase{
			name:    "1p is reach suji of 54p",
			scene:   scene,
			args:    args{name: "reach_suji", pai: testutil.MustPai("1p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "2p is not reach suji of 54p",
			scene:   scene,
			args:    args{name: "reach_suji", pai: testutil.MustPai("2p")},
			want:    false,
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.scene.evaluate(tt.args.name, tt.args.pai)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scene.evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Scene.evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScene_Evaluate_ReachKataSuji(t *testing.T) {
	tests := []testCase{}

	scene, _ := NewSceneFromParams(nil, testutil.MustPais("1p"), nil, nil, testutil.MustPais("1p"), nil, nil)

	{
		tests = append(tests, testCase{
			name:    "4p is reach suji of 1p",
			scene:   scene,
			args:    args{name: "reach_suji", pai: testutil.MustPai("4p")},
			want:    true,
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.scene.evaluate(tt.args.name, tt.args.pai)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scene.evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Scene.evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScene_Evaluate_PrereachSuji(t *testing.T) {
	tests := []testCase{}

	scene, _ := NewSceneFromParams(nil, testutil.MustPais("4p", "E", "4s"), nil, nil, testutil.MustPais("4p", "E"), nil, nil)

	{
		tests = append(tests, testCase{
			name:    "1p is prereach suji of 4pE4s",
			scene:   scene,
			args:    args{name: "prereach_suji", pai: testutil.MustPai("1p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "1s is not prereach suji of 4pE4s",
			scene:   scene,
			args:    args{name: "prereach_suji", pai: testutil.MustPai("1s")},
			want:    false,
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.scene.evaluate(tt.args.name, tt.args.pai)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scene.evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Scene.evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScene_Evaluate_UraSuji(t *testing.T) {
	tests := []testCase{}

	scene, _ := NewSceneFromParams(nil, testutil.MustPais("1p"), nil, nil, testutil.MustPais("1p"), nil, nil)

	{
		tests = append(tests, testCase{
			name:    "2p is urasuji of 1p",
			scene:   scene,
			args:    args{name: "urasuji", pai: testutil.MustPai("2p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "5p is urasuji of 1p",
			scene:   scene,
			args:    args{name: "urasuji", pai: testutil.MustPai("5p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "3p is not urasuji of 1p",
			scene:   scene,
			args:    args{name: "urasuji", pai: testutil.MustPai("3p")},
			want:    false,
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.scene.evaluate(tt.args.name, tt.args.pai)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scene.evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Scene.evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScene_Evaluate_UraSujiOf5(t *testing.T) {
	tests := []testCase{}

	scene, _ := NewSceneFromParams(nil, testutil.MustPais("5p"), nil, nil, testutil.MustPais("5p"), nil, nil)

	{
		tests = append(tests, testCase{
			name:    "1p is urasuji of 5p",
			scene:   scene,
			args:    args{name: "urasuji", pai: testutil.MustPai("1p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "4p is urasuji of 5p",
			scene:   scene,
			args:    args{name: "urasuji", pai: testutil.MustPai("4p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "6p is urasuji of 5p",
			scene:   scene,
			args:    args{name: "urasuji", pai: testutil.MustPai("6p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "9p is urasuji of 5p",
			scene:   scene,
			args:    args{name: "urasuji", pai: testutil.MustPai("9p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "2p is not urasuji of 5p",
			scene:   scene,
			args:    args{name: "urasuji", pai: testutil.MustPai("2p")},
			want:    false,
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.scene.evaluate(tt.args.name, tt.args.pai)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scene.evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Scene.evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScene_Evaluate_UraSuji_ReachPai(t *testing.T) {
	tests := []testCase{}

	scene, _ := NewSceneFromParams(nil, testutil.MustPais("1p", "5p"), nil, nil, testutil.MustPais("1p"), nil, nil)

	{
		tests = append(tests, testCase{
			name:    "2p is not urasuji of reach declaration pai 5p",
			scene:   scene,
			args:    args{name: "urasuji", pai: testutil.MustPai("2p")},
			want:    false,
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.scene.evaluate(tt.args.name, tt.args.pai)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scene.evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Scene.evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScene_Evaluate_EarlyUraSuji_ReachUraSuji(t *testing.T) {
	tests := []testCase{}

	scene, _ := NewSceneFromParams(nil, testutil.MustPais("1p", "E", "S", "W", "1s"), nil, nil, testutil.MustPais("1p", "E", "S", "W", "1s"), nil, nil)

	{
		tests = append(tests, testCase{
			name:    "5p is early urasuji of 1pESW1s",
			scene:   scene,
			args:    args{name: "early_urasuji", pai: testutil.MustPai("5p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "5s is not early urasuji of 1pESW1s",
			scene:   scene,
			args:    args{name: "early_urasuji", pai: testutil.MustPai("5s")},
			want:    false,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "5s is reach urasuji of 1pESW1s",
			scene:   scene,
			args:    args{name: "reach_urasuji", pai: testutil.MustPai("5s")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "5p is not reach urasuji of 1pESW1s",
			scene:   scene,
			args:    args{name: "reach_urasuji", pai: testutil.MustPai("5p")},
			want:    false,
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.scene.evaluate(tt.args.name, tt.args.pai)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scene.evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Scene.evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScene_Evaluate_Only_UraSujiOf5(t *testing.T) {
	tests := []testCase{}

	scene, _ := NewSceneFromParams(nil, testutil.MustPais("1p", "5s"), nil, nil, testutil.MustPais("1p", "5s"), nil, nil)

	{
		tests = append(tests, testCase{
			name:    "1s is urasuji of 5 of 1p5s",
			scene:   scene,
			args:    args{name: "urasuji_of_5", pai: testutil.MustPai("1s")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "2p is not urasuji of 5 of 1p5s",
			scene:   scene,
			args:    args{name: "urasuji_of_5", pai: testutil.MustPai("2p")},
			want:    false,
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.scene.evaluate(tt.args.name, tt.args.pai)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scene.evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Scene.evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScene_Evaluate_Aida4ken(t *testing.T) {
	tests := []testCase{}

	scene, _ := NewSceneFromParams(nil, testutil.MustPais("1p", "6p"), nil, nil, testutil.MustPais("1p", "6p"), nil, nil)

	{
		tests = append(tests, testCase{
			name:    "2p is aida4ken of 1p6p",
			scene:   scene,
			args:    args{name: "aida4ken", pai: testutil.MustPai("2p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "5p is aida4ken of 1p6p",
			scene:   scene,
			args:    args{name: "aida4ken", pai: testutil.MustPai("5p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "3p is not aida4ken of 1p6p",
			scene:   scene,
			args:    args{name: "aida4ken", pai: testutil.MustPai("3p")},
			want:    false,
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.scene.evaluate(tt.args.name, tt.args.pai)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scene.evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Scene.evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScene_Evaluate_MatagiSuji(t *testing.T) {
	tests := []testCase{}

	scene1, _ := NewSceneFromParams(nil, testutil.MustPais("3p"), nil, nil, testutil.MustPais("3p"), nil, nil)

	{
		tests = append(tests, testCase{
			name:    "1p is matagisuji of 3p",
			scene:   scene1,
			args:    args{name: "matagisuji", pai: testutil.MustPai("1p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "2p is matagisuji of 3p",
			scene:   scene1,
			args:    args{name: "matagisuji", pai: testutil.MustPai("2p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "4p is matagisuji of 3p",
			scene:   scene1,
			args:    args{name: "matagisuji", pai: testutil.MustPai("4p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "5p is matagisuji of 3p",
			scene:   scene1,
			args:    args{name: "matagisuji", pai: testutil.MustPai("5p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "6p is not matagisuji of 3p",
			scene:   scene1,
			args:    args{name: "matagisuji", pai: testutil.MustPai("6p")},
			want:    false,
			wantErr: false,
		})
	}

	scene2, _ := NewSceneFromParams(nil, testutil.MustPais("2p"), nil, nil, testutil.MustPais("2p"), nil, nil)

	{
		tests = append(tests, testCase{
			name:    "1p is matagisuji of 2p",
			scene:   scene2,
			args:    args{name: "matagisuji", pai: testutil.MustPai("1p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "4p is matagisuji of 2p",
			scene:   scene2,
			args:    args{name: "matagisuji", pai: testutil.MustPai("4p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "3p is not matagisuji of 2p",
			scene:   scene2,
			args:    args{name: "matagisuji", pai: testutil.MustPai("3p")},
			want:    false,
			wantErr: false,
		})
	}

	scene3, _ := NewSceneFromParams(nil, testutil.MustPais("3p", "4p"), nil, nil, testutil.MustPais("3p"), nil, nil)

	{
		tests = append(tests, testCase{
			name:    "1p is not matagisuji of 3p4p",
			scene:   scene3,
			args:    args{name: "matagisuji", pai: testutil.MustPai("1p")},
			want:    false,
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.scene.evaluate(tt.args.name, tt.args.pai)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scene.evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Scene.evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScene_Evaluate_EarlyMatagiSuji(t *testing.T) {
	tests := []testCase{}

	scene, _ := NewSceneFromParams(nil, testutil.MustPais("3p", "E", "S", "7p", "W"), nil, nil, testutil.MustPais("3p", "E", "S", "7p", "W"), nil, nil)

	{
		tests = append(tests, testCase{
			name:    "1p is early matagisuji of 3pES7pW",
			scene:   scene,
			args:    args{name: "early_matagisuji", pai: testutil.MustPai("1p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "9p is not early matagisuji of 3pES7pW",
			scene:   scene,
			args:    args{name: "early_matagisuji", pai: testutil.MustPai("9p")},
			want:    false,
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.scene.evaluate(tt.args.name, tt.args.pai)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scene.evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Scene.evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScene_Evaluate_LateMatagiSuji(t *testing.T) {
	tests := []testCase{}

	scene, _ := NewSceneFromParams(nil, testutil.MustPais("3p", "E", "S", "7p", "W"), nil, nil, testutil.MustPais("3p", "E", "S", "7p", "W"), nil, nil)

	{
		tests = append(tests, testCase{
			name:    "9p is late matagisuji of 3pES7pW",
			scene:   scene,
			args:    args{name: "late_matagisuji", pai: testutil.MustPai("9p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "1p is not late matagisuji of 3pES7pW",
			scene:   scene,
			args:    args{name: "late_matagisuji", pai: testutil.MustPai("1p")},
			want:    false,
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.scene.evaluate(tt.args.name, tt.args.pai)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scene.evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Scene.evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScene_Evaluate_ReachMatagiSuji(t *testing.T) {
	tests := []testCase{}

	scene, _ := NewSceneFromParams(nil, testutil.MustPais("3p", "E", "S", "7p"), nil, nil, testutil.MustPais("3p", "E", "S", "7p"), nil, nil)

	{
		tests = append(tests, testCase{
			name:    "9p is reach matagisuji of 3pES7p",
			scene:   scene,
			args:    args{name: "reach_matagisuji", pai: testutil.MustPai("9p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "1p is not reach matagisuji of 3pES7p",
			scene:   scene,
			args:    args{name: "reach_matagisuji", pai: testutil.MustPai("1p")},
			want:    false,
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.scene.evaluate(tt.args.name, tt.args.pai)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scene.evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Scene.evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScene_Evaluate_SenkiSuji(t *testing.T) {
	tests := []testCase{}

	scene, _ := NewSceneFromParams(nil, testutil.MustPais("1p"), nil, nil, testutil.MustPais("1p"), nil, nil)

	{
		tests = append(tests, testCase{
			name:    "3p is senkisuji of 1p",
			scene:   scene,
			args:    args{name: "senkisuji", pai: testutil.MustPai("3p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "6p is senkisuji of 1p",
			scene:   scene,
			args:    args{name: "senkisuji", pai: testutil.MustPai("6p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "2p is not senkisuji of 1p",
			scene:   scene,
			args:    args{name: "senkisuji", pai: testutil.MustPai("2p")},
			want:    false,
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.scene.evaluate(tt.args.name, tt.args.pai)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scene.evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Scene.evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Exclude the tile being discarded from the count.
func TestScene_Evaluate_VisibleNOrMore(t *testing.T) {
	tests := []testCase{}

	scene1, _ := NewSceneFromParams(nil, nil, testutil.MustPais("1p", "1p"), nil, nil, nil, nil)

	{
		tests = append(tests, testCase{
			name:    "1 or more 1p are visible",
			scene:   scene1,
			args:    args{name: "visible>=1", pai: testutil.MustPai("1p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "2 or more 1p are not visible",
			scene:   scene1,
			args:    args{name: "visible>=2", pai: testutil.MustPai("1p")},
			want:    false,
			wantErr: false,
		})
	}

	scene2, _ := NewSceneFromParams(nil, nil, testutil.MustPais("1p", "1p", "1p"), nil, nil, nil, nil)

	{
		tests = append(tests, testCase{
			name:    "2 or more 1p are visible",
			scene:   scene2,
			args:    args{name: "visible>=2", pai: testutil.MustPai("1p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "3 or more 1p are not visible",
			scene:   scene2,
			args:    args{name: "visible>=3", pai: testutil.MustPai("1p")},
			want:    false,
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.scene.evaluate(tt.args.name, tt.args.pai)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scene.evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Scene.evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScene_Evaluate_SujiVisible(t *testing.T) {
	tests := []testCase{}

	scene1, _ := NewSceneFromParams(nil, nil, testutil.MustPais("4p"), nil, nil, nil, nil)

	{
		tests = append(tests, testCase{
			name:    "1 or less suji of 1p are visible",
			scene:   scene1,
			args:    args{name: "suji_visible<=1", pai: testutil.MustPai("1p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "0 or less suji of 1p are not visible",
			scene:   scene1,
			args:    args{name: "suji_visible<=0", pai: testutil.MustPai("1p")},
			want:    false,
			wantErr: false,
		})
	}

	scene2, _ := NewSceneFromParams(nil, nil, testutil.MustPais("4p", "4p"), nil, nil, nil, nil)

	{
		tests = append(tests, testCase{
			name:    "2 or less suji of 1p are visible",
			scene:   scene2,
			args:    args{name: "suji_visible<=2", pai: testutil.MustPai("1p")},
			want:    true,
			wantErr: false,
		})
	}
	{
		tests = append(tests, testCase{
			name:    "1 or less suji of 1p are not visible",
			scene:   scene2,
			args:    args{name: "suji_visible<=1", pai: testutil.MustPai("1p")},
			want:    false,
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.scene.evaluate(tt.args.name, tt.args.pai)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scene.evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Scene.evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO Add test for rest of features.
