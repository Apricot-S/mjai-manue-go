package game_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/game"
)

func sum(arr [base.NumIDs]int) int {
	sum := 0
	for _, c := range arr {
		sum += c
	}
	return sum
}

func verifyShantenAndGoals(t *testing.T, paiSet *base.PaiSet, wantShanten int, wantGoalsCount int) {
	t.Helper()

	shanten, goals, err := game.AnalyzeShanten(paiSet)
	if err != nil {
		t.Errorf("AnalyzeShanten() error = %v", err)
		return
	}
	if shanten != wantShanten {
		t.Errorf("AnalyzeShanten() shanten = %v, want %v", shanten, wantShanten)
	}
	if len(goals) != wantGoalsCount {
		t.Errorf("AnalyzeShanten() len(goals) = %v, want %v", len(goals), wantGoalsCount)
	}

	numRequiredBlock := sum(*paiSet)/3 + 1
	for _, goal := range goals {
		if len(goal.Mentsus) != numRequiredBlock {
			t.Errorf("AnalyzeShanten() len(goal.Mentsus) = %v, want %v", len(goals), wantGoalsCount)
		}
	}
}

func verifyShantenWithUpperBounds(t *testing.T, paiSet *base.PaiSet, wantShanten int, wantGoalsCount int) {
	t.Helper()

	for i := -1; i <= 8; i++ {
		shanten, goals, err := game.AnalyzeShantenWithOption(paiSet, 0, i)
		if err != nil {
			t.Errorf("i = %v, AnalyzeShantenWithOption() error = %v", i, err)
			return
		}

		wantShantenWithUpperBound := wantShanten
		if wantShanten > i {
			wantShantenWithUpperBound = game.InfinityShanten
		}
		wantGoalsCountWithUpperBound := wantGoalsCount
		if wantShanten > i {
			wantGoalsCountWithUpperBound = 0
		}
		if shanten != wantShantenWithUpperBound {
			t.Errorf("i = %v, AnalyzeShantenWithOption() shanten = %v, want %v", i, shanten, wantShantenWithUpperBound)
		}
		if len(goals) != wantGoalsCountWithUpperBound {
			t.Errorf("i = %v, AnalyzeShantenWithOption() len(goals) = %v, want %v", i, len(goals), wantGoalsCountWithUpperBound)
		}

		numRequiredBlock := sum(*paiSet)/3 + 1
		for _, goal := range goals {
			if len(goal.Mentsus) != numRequiredBlock {
				t.Errorf("i = %v, AnalyzeShantenWithOption() len(goal.Mentsus) = %v, want %v", i, len(goals), wantGoalsCountWithUpperBound)
			}
		}
	}
}

func TestAnalyzeShanten(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		wantShanten    int
		wantGoalsCount int
	}
	tests := []testCase{
		{
			name:           "case 1",
			input:          "1m 2m 3m 7m 8m 9m 2s 3s 4s S S S W",
			wantShanten:    0,
			wantGoalsCount: 1,
		},
		{
			name:           "case 2",
			input:          "1m 2m 3m 7m 8m 9m 2s 3s S S S W N",
			wantShanten:    1,
			wantGoalsCount: 4,
		},
		{
			name:           "empty : An empty hand is one step away from being a pair wait -> Shanten number: 1",
			input:          "",
			wantShanten:    1,
			wantGoalsCount: 34,
		},
		{
			name:           "thirteen orphans",
			input:          "1m 9m 1p 9p 1s 9s E S W N P F C",
			wantShanten:    8,
			wantGoalsCount: 27675,
		},
		{
			name:           "tenpai",
			input:          "1m 2m 3m 4p 5pr 6p 7s 8s 9s E E S S",
			wantShanten:    0,
			wantGoalsCount: 2,
		},
		{
			name:           "win",
			input:          "1m 2m 3m 4p 5pr 6p 7s 8s 9s E E S S S",
			wantShanten:    -1,
			wantGoalsCount: 1,
		},
		{
			name:           "with meld",
			input:          "1m 2m 3m 4p 5pr 6p 7s 8s 9s E",
			wantShanten:    0,
			wantGoalsCount: 1,
		},
		{
			name:           "without pair",
			input:          "1m 2m 3m 8m 9m 4p 5p 6p 1s 2s 7s 8s 9s E",
			wantShanten:    1,
			wantGoalsCount: 6,
		},
		{
			name:           "too many meld candidates",
			input:          "1m 2m 3m 8m 9m 4p 5p 6p 1s 2s 8s 9s E E",
			wantShanten:    1,
			wantGoalsCount: 3,
		},
		{
			name:           "not enough meld candidates",
			input:          "1m 3m 3m 3m 4m 5m 5m 6m 8m S W F C C",
			wantShanten:    2,
			wantGoalsCount: 1,
		},
		{
			name:           "incomplete hand 4 melds without a pair",
			input:          "2p 3p 4p 5s 6s 7s",
			wantShanten:    1,
			wantGoalsCount: 38,
		},
		{
			name:           "triplet sequence",
			input:          "2p 2p 2p 3p 4p 5p E S W N P F C",
			wantShanten:    4,
			wantGoalsCount: 105,
		},
		{
			name:           "sequence isolated sequence",
			input:          "2p 3p 4p 4p 4p 5p 6p E S W N P F",
			wantShanten:    4,
			wantGoalsCount: 285,
		},
		{
			name:           "pair triplet sequence",
			input:          "1p 1p 2p 2p 2p 3p 4p 5p E S W N P",
			wantShanten:    3,
			wantGoalsCount: 30,
		},
		{
			name:           "pair sequence sequence pair",
			input:          "2p 2p 3p 4p 5p 5p 6p 7p 8p 8p E S W",
			wantShanten:    2,
			wantGoalsCount: 9,
		},
		{
			name:           "waiting for the 5th tile 1",
			input:          "1m 1m 1m 1m 1p 2p 3p 1s 1s 2s 2s 3s 3s",
			wantShanten:    1,
			wantGoalsCount: 40,
		},
		{
			name:           "waiting for the 5th tile 2",
			input:          "1m 1m 1m 1m 2m 3m 4m 4m 4m 4m 1p 1p 1p 1p",
			wantShanten:    1,
			wantGoalsCount: 82,
		},
		{
			name:           "waiting for the 5th tile 3",
			input:          "E E E E S S S S W W W N N N",
			wantShanten:    1,
			wantGoalsCount: 30,
		},
		{
			name:           "2 isolated 4 tiles 1",
			input:          "1m 1m 1m 1m 2m 4m 7m 7m 7m 7m",
			wantShanten:    1,
			wantGoalsCount: 1,
		},
		{
			name:           "2 isolated 4 tiles 2",
			input:          "1m 1m 1m 1m 2m 4m 7m 7m 7m 7m E E E S",
			wantShanten:    1,
			wantGoalsCount: 3,
		},
		{
			name:           "2 isolated 4 tiles 3",
			input:          "1m 1m 1m 1m 4m 4m 4m 4m",
			wantShanten:    1,
			wantGoalsCount: 40,
		},
		{
			name:           "2 isolated 4 tiles 4",
			input:          "1m 1m 1m 1m 2m 4m E E E E",
			wantShanten:    1,
			wantGoalsCount: 1,
		},
		{
			name:           "2 isolated 4 tiles 5",
			input:          "1m 1m 1m 1m 4m 4m 4m 4m 7m 8m",
			wantShanten:    2,
			wantGoalsCount: 89,
		},
		{
			name:           "3 isolated 4 tiles",
			input:          "1m 1m 1m 1m 2m 4m 7m 7m 7m 7m E E E E",
			wantShanten:    1,
			wantGoalsCount: 1,
		},
		{
			name:           "4 honors 1",
			input:          "E E E E",
			wantShanten:    1,
			wantGoalsCount: 33,
		},
		{
			name:           "4 honors 2",
			input:          "1m 2m 3m E E E E",
			wantShanten:    1,
			wantGoalsCount: 34,
		},
		{
			name:           "4 honors 3",
			input:          "E E E E S S S S",
			wantShanten:    1,
			wantGoalsCount: 32,
		},
		{
			name:           "4 honors 4",
			input:          "1m 2m 3m 1p 1p E E E E S S S S",
			wantShanten:    2,
			wantGoalsCount: 37,
		},
		{
			name:           "can be interpreted in multiple set decompositions",
			input:          "1m 1m 1m 2m 2m 2m 3m 3m 3m 7p 8p 9p 9p 9p",
			wantShanten:    -1,
			wantGoalsCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pais, err := base.StrToPais(tt.input)
			if err != nil {
				t.Errorf("StrToPais() error = %v", err)
				return
			}
			paiSet, err := base.NewPaiSet(pais)
			if err != nil {
				t.Errorf("NewPaiSet() error = %v", err)
				return
			}

			verifyShantenAndGoals(t, paiSet, tt.wantShanten, tt.wantGoalsCount)
			verifyShantenWithUpperBounds(t, paiSet, tt.wantShanten, tt.wantGoalsCount)
		})
	}
}

func TestAnalyzeShanten_Invalid(t *testing.T) {
	type testCase struct {
		name        string
		input       string
		wantShanten int
		wantGoals   []game.Goal
	}
	tests := []testCase{
		{
			name:        "5 identical tiles",
			input:       "1m 1m 1m 1m 1m",
			wantShanten: math.MaxInt,
			wantGoals:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pais, _ := base.StrToPais(tt.input)
			paiSet, _ := base.NewPaiSet(pais)

			shanten, goals, err := game.AnalyzeShanten(paiSet)
			if err == nil {
				t.Errorf("AnalyzeShanten() error = %v", err)
				return
			}
			if shanten != tt.wantShanten {
				t.Errorf("AnalyzeShanten() shanten = %v, want %v", shanten, tt.wantShanten)
			}
			if !reflect.DeepEqual(goals, tt.wantGoals) {
				t.Errorf("AnalyzeShanten() goals = %v, want %v", goals, tt.wantGoals)
			}
		})
	}
}

func TestAnalyzeShantenChitoitsu(t *testing.T) {
	type testCase struct {
		name        string
		input       string
		wantShanten int
	}
	tests := []testCase{
		{
			name:        "without pair",
			input:       "1m 9m 1p 9p 1s 9s E S W N P F C",
			wantShanten: 6,
		},
		{
			name:        "with quadruple",
			input:       "1m 1m 8m 8m 2p 8p 8p 5s 5s E E E E",
			wantShanten: 2,
		},
		{
			name:        "with triplet",
			input:       "1m 1m 8m 8m 2p 3p 8p 8p 5s 5s E E E",
			wantShanten: 1,
		},
		{
			name:        "with 2 triplets",
			input:       "1m 1m 8m 8m 2p 8p 8p 5s 5s 5s E E E",
			wantShanten: 2,
		},
		{
			name:        "tenpai",
			input:       "1m 1m 8m 8m 2p 8p 8p 5s 5s E E C C",
			wantShanten: 0,
		},
		{
			name:        "win",
			input:       "1m 1m 8m 8m 2p 2p 8p 8p 5s 5s E E C C",
			wantShanten: -1,
		},
		{
			name:        "incomplete_hand",
			input:       "1m 1m 8m 8m 5s 5s E E S S",
			wantShanten: game.InfinityShanten,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pais, err := base.StrToPais(tt.input)
			if err != nil {
				t.Errorf("StrToPais() error = %v", err)
				return
			}
			paiSet, err := base.NewPaiSet(pais)
			if err != nil {
				t.Errorf("NewPaiSet() error = %v", err)
				return
			}

			shanten, err := game.AnalyzeShantenChitoitsu(paiSet)
			if err != nil {
				t.Errorf("AnalyzeShantenChitoitsu() error = %v", err)
				return
			}
			if shanten != tt.wantShanten {
				t.Errorf("AnalyzeShantenChitoitsu() shanten = %v, want %v", shanten, tt.wantShanten)
			}

		})
	}
}
