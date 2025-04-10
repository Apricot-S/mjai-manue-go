package game_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/game"
)

func sum(arr [game.NumIDs]int) int {
	sum := 0
	for _, c := range arr {
		sum += c
	}
	return sum
}

func verifyShantenAndGoals(t *testing.T, paiSet *game.PaiSet, expectedShanten int, expectedGoalsSize int) {
	shanten, goals, err := game.AnalyzeShanten(paiSet)
	if err != nil {
		t.Errorf("AnalyzeShanten() error = %v", err)
		return
	}
	if shanten != expectedShanten {
		t.Errorf("AnalyzeShanten() shanten = %v, want %v", shanten, expectedShanten)
	}
	if len(goals) != expectedGoalsSize {
		t.Errorf("AnalyzeShanten() len(goals) = %v, want %v", len(goals), expectedGoalsSize)
	}

	numRequiredBlock := sum(*paiSet)/3 + 1
	for _, goal := range goals {
		if len(goal.Mentsus) != numRequiredBlock {
			t.Errorf("AnalyzeShanten() len(goal.Mentsus) = %v, want %v", len(goals), expectedGoalsSize)
		}
	}
}

func verifyShantenWithUpperBounds(t *testing.T, paiSet *game.PaiSet, expectedShanten int, expectedGoalsSize int) {
	for i := -1; i <= 8; i++ {
		shanten, goals, err := game.AnalyzeShantenWithOption(paiSet, 0, i)
		if err != nil {
			t.Errorf("i = %v, AnalyzeShantenWithOption() error = %v", i, err)
			return
		}

		expectedShantenWithUpperBound := expectedShanten
		if expectedShanten > i {
			expectedShantenWithUpperBound = game.InfinityShanten
		}
		expectedGoalsSizeWithUpperBound := expectedGoalsSize
		if expectedShanten > i {
			expectedGoalsSizeWithUpperBound = 0
		}
		if shanten != expectedShantenWithUpperBound {
			t.Errorf("i = %v, AnalyzeShantenWithOption() shanten = %v, want %v", i, shanten, expectedShantenWithUpperBound)
		}
		if len(goals) != expectedGoalsSizeWithUpperBound {
			t.Errorf("i = %v, AnalyzeShantenWithOption() len(goals) = %v, want %v", i, len(goals), expectedGoalsSizeWithUpperBound)
		}

		numRequiredBlock := sum(*paiSet)/3 + 1
		for _, goal := range goals {
			if len(goal.Mentsus) != numRequiredBlock {
				t.Errorf("i = %v, AnalyzeShantenWithOption() len(goal.Mentsus) = %v, want %v", i, len(goals), expectedGoalsSizeWithUpperBound)
			}
		}
	}
}

func testAnalyzeShantenInternal(t *testing.T, paiStr string, expectedShanten int, expectedGoalsSize int) {
	pais, _ := game.StrToPais(paiStr)
	paiSet, _ := game.NewPaiSetWithPais(pais)

	verifyShantenAndGoals(t, paiSet, expectedShanten, expectedGoalsSize)
	verifyShantenWithUpperBounds(t, paiSet, expectedShanten, expectedGoalsSize)
}

func TestAnalyzeShanten(t *testing.T) {
	type args struct {
		ps string
	}
	type testCase struct {
		args  args
		want  int
		want1 int
	}
	tests := []testCase{
		// case 1
		{args{"1m 2m 3m 7m 8m 9m 2s 3s 4s S S S W"}, 0, 1},
		// case 2
		{args{"1m 2m 3m 7m 8m 9m 2s 3s S S S W N"}, 1, 4},
		// empty : An empty hand is one step away from being a pair wait -> Shanten number: 1
		{args{""}, 1, 34},
		// thirteen orphans
		{args{"1m 9m 1p 9p 1s 9s E S W N P F C"}, 8, 27675},
		// tenpai
		{args{"1m 2m 3m 4p 5pr 6p 7s 8s 9s E E S S"}, 0, 2},
		// win
		{args{"1m 2m 3m 4p 5pr 6p 7s 8s 9s E E S S S"}, -1, 1},
		// with meld
		{args{"1m 2m 3m 4p 5pr 6p 7s 8s 9s E"}, 0, 1},
		// without pair
		{args{"1m 2m 3m 8m 9m 4p 5p 6p 1s 2s 7s 8s 9s E"}, 1, 6},
		// too many meld candidates
		{args{"1m 2m 3m 8m 9m 4p 5p 6p 1s 2s 8s 9s E E"}, 1, 3},
		// not enough meld candidates
		{args{"1m 3m 3m 3m 4m 5m 5m 6m 8m S W F C C"}, 2, 1},
		// incomplete hand 4 melds without a pair
		{args{"2p 3p 4p 5s 6s 7s"}, 1, 38},
		// triplet sequence
		{args{"2p 2p 2p 3p 4p 5p E S W N P F C"}, 4, 105},
		// sequence isolated sequence
		{args{"2p 3p 4p 4p 4p 5p 6p E S W N P F"}, 4, 285},
		// pair triplet sequence
		{args{"1p 1p 2p 2p 2p 3p 4p 5p E S W N P"}, 3, 30},
		// pair sequence sequence pair
		{args{"2p 2p 3p 4p 5p 5p 6p 7p 8p 8p E S W"}, 2, 9},
		// waiting for the 5th tile 1
		{args{"1m 1m 1m 1m 1p 2p 3p 1s 1s 2s 2s 3s 3s"}, 1, 40},
		// waiting for the 5th tile 2
		{args{"1m 1m 1m 1m 2m 3m 4m 4m 4m 4m 1p 1p 1p 1p"}, 1, 82},
		// waiting for the 5th tile 3
		{args{"E E E E S S S S W W W N N N"}, 1, 30},
		// 2 isolated 4 tiles 1
		{args{"1m 1m 1m 1m 2m 4m 7m 7m 7m 7m"}, 1, 1},
		// 2 isolated 4 tiles 2
		{args{"1m 1m 1m 1m 2m 4m 7m 7m 7m 7m E E E S"}, 1, 3},
		// 2 isolated 4 tiles 3
		{args{"1m 1m 1m 1m 4m 4m 4m 4m"}, 1, 40},
		// 2 isolated 4 tiles 4
		{args{"1m 1m 1m 1m 2m 4m E E E E"}, 1, 1},
		// 2 isolated 4 tiles 5
		{args{"1m 1m 1m 1m 4m 4m 4m 4m 7m 8m"}, 2, 89},
		// 3 isolated 4 tiles
		{args{"1m 1m 1m 1m 2m 4m 7m 7m 7m 7m E E E E"}, 1, 1},
		// 4 honors 1
		{args{"E E E E"}, 1, 33},
		// 4 honors 2
		{args{"1m 2m 3m E E E E"}, 1, 34},
		// 4 honors 3
		{args{"E E E E S S S S"}, 1, 32},
		// 4 honors 4
		{args{"1m 2m 3m 1p 1p E E E E S S S S"}, 2, 37},
		// can be interpreted in multiple set decompositions
		{args{"1m 1m 1m 2m 2m 2m 3m 3m 3m 7p 8p 9p 9p 9p"}, -1, 2},
	}

	for _, tt := range tests {
		t.Run(tt.args.ps, func(t *testing.T) {
			testAnalyzeShantenInternal(t, tt.args.ps, tt.want, tt.want1)
		})
	}
}

func TestAnalyzeShanten_Invalid(t *testing.T) {
	type args struct {
		ps string
	}
	type testCase struct {
		args    args
		shanten int
		goals   []game.Goal
	}
	tests := []testCase{
		// 5 identical tiles
		{args{"1m 1m 1m 1m 1m"}, math.MaxInt, nil},
	}

	for _, tt := range tests {
		t.Run(tt.args.ps, func(t *testing.T) {
			pais, _ := game.StrToPais(tt.args.ps)
			paiSet, _ := game.NewPaiSetWithPais(pais)

			shanten, goals, err := game.AnalyzeShanten(paiSet)
			if err == nil {
				t.Errorf("AnalyzeShanten() error = %v", err)
				return
			}
			if shanten != tt.shanten {
				t.Errorf("AnalyzeShanten() shanten = %v, want %v", shanten, tt.shanten)
			}
			if !reflect.DeepEqual(goals, tt.goals) {
				t.Errorf("AnalyzeShanten() goals = %v, want %v", goals, tt.goals)
			}
		})
	}
}
