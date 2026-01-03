package service_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/service"
)

func TestAnalyzeShanten(t *testing.T) {
	tests := []struct {
		name           string
		codes          []string
		wantShanten    int
		wantGoalsCount int
	}{
		{
			name:           "case 1",
			codes:          []string{"1m", "2m", "3m", "7m", "8m", "9m", "2s", "3s", "4s", "S", "S", "S", "W"},
			wantShanten:    0,
			wantGoalsCount: 1,
		},
		{
			name:           "case 2",
			codes:          []string{"1m", "2m", "3m", "7m", "8m", "9m", "2s", "3s", "S", "S", "S", "W", "N"},
			wantShanten:    1,
			wantGoalsCount: 4,
		},
		{
			name:           "empty : An empty hand is one step away from being a pair wait -> Shanten number: 1",
			codes:          nil,
			wantShanten:    1,
			wantGoalsCount: 34,
		},
		{
			name:           "thirteen orphans",
			codes:          []string{"1m", "9m", "1p", "9p", "1s", "9s", "E", "S", "W", "N", "P", "F", "C"},
			wantShanten:    8,
			wantGoalsCount: 27675,
		},
		{
			name:           "tenpai",
			codes:          []string{"1m", "2m", "3m", "4p", "5pr", "6p", "7s", "8s", "9s", "E", "E", "S", "S"},
			wantShanten:    0,
			wantGoalsCount: 2,
		},
		{
			name:           "win",
			codes:          []string{"1m", "2m", "3m", "4p", "5pr", "6p", "7s", "8s", "9s", "E", "E", "S", "S", "S"},
			wantShanten:    -1,
			wantGoalsCount: 1,
		},
		{
			name:           "with meld",
			codes:          []string{"1m", "2m", "3m", "4p", "5pr", "6p", "7s", "8s", "9s", "E"},
			wantShanten:    0,
			wantGoalsCount: 1,
		},
		{
			name:           "without pair",
			codes:          []string{"1m", "2m", "3m", "8m", "9m", "4p", "5p", "6p", "1s", "2s", "7s", "8s", "9s", "E"},
			wantShanten:    1,
			wantGoalsCount: 6,
		},
		{
			name:           "too many meld candidates",
			codes:          []string{"1m", "2m", "3m", "8m", "9m", "4p", "5p", "6p", "1s", "2s", "8s", "9s", "E", "E"},
			wantShanten:    1,
			wantGoalsCount: 3,
		},
		{
			name:           "not enough meld candidates",
			codes:          []string{"1m", "3m", "3m", "3m", "4m", "5m", "5m", "6m", "8m", "S", "W", "F", "C", "C"},
			wantShanten:    2,
			wantGoalsCount: 1,
		},
		{
			name:           "incomplete hand 4 melds without a pair",
			codes:          []string{"2p", "3p", "4p", "5s", "6s", "7s"},
			wantShanten:    1,
			wantGoalsCount: 38,
		},
		{
			name:           "triplet sequence",
			codes:          []string{"2p", "2p", "2p", "3p", "4p", "5p", "E", "S", "W", "N", "P", "F", "C"},
			wantShanten:    4,
			wantGoalsCount: 105,
		},
		{
			name:           "sequence isolated sequence",
			codes:          []string{"2p", "3p", "4p", "4p", "4p", "5p", "6p", "E", "S", "W", "N", "P", "F"},
			wantShanten:    4,
			wantGoalsCount: 285,
		},
		{
			name:           "pair triplet sequence",
			codes:          []string{"1p", "1p", "2p", "2p", "2p", "3p", "4p", "5p", "E", "S", "W", "N", "P"},
			wantShanten:    3,
			wantGoalsCount: 30,
		},
		{
			name:           "pair sequence sequence pair",
			codes:          []string{"2p", "2p", "3p", "4p", "5p", "5p", "6p", "7p", "8p", "8p", "E", "S", "W"},
			wantShanten:    2,
			wantGoalsCount: 9,
		},
		{
			name:           "waiting for the 5th tile 1",
			codes:          []string{"1m", "1m", "1m", "1m", "1p", "2p", "3p", "1s", "1s", "2s", "2s", "3s", "3s"},
			wantShanten:    1,
			wantGoalsCount: 40,
		},
		{
			name:           "waiting for the 5th tile 2",
			codes:          []string{"1m", "1m", "1m", "1m", "2m", "3m", "4m", "4m", "4m", "4m", "1p", "1p", "1p", "1p"},
			wantShanten:    1,
			wantGoalsCount: 82,
		},
		{
			name:           "waiting for the 5th tile 3",
			codes:          []string{"E", "E", "E", "E", "S", "S", "S", "S", "W", "W", "W", "N", "N", "N"},
			wantShanten:    1,
			wantGoalsCount: 30,
		},
		{
			name:           "2 isolated 4 tiles 1",
			codes:          []string{"1m", "1m", "1m", "1m", "2m", "4m", "7m", "7m", "7m", "7m"},
			wantShanten:    1,
			wantGoalsCount: 1,
		},
		{
			name:           "2 isolated 4 tiles 2",
			codes:          []string{"1m", "1m", "1m", "1m", "2m", "4m", "7m", "7m", "7m", "7m", "E", "E", "E", "S"},
			wantShanten:    1,
			wantGoalsCount: 3,
		},
		{
			name:           "2 isolated 4 tiles 3",
			codes:          []string{"1m", "1m", "1m", "1m", "4m", "4m", "4m", "4m"},
			wantShanten:    1,
			wantGoalsCount: 40,
		},
		{
			name:           "2 isolated 4 tiles 4",
			codes:          []string{"1m", "1m", "1m", "1m", "2m", "4m", "E", "E", "E", "E"},
			wantShanten:    1,
			wantGoalsCount: 1,
		},
		{
			name:           "2 isolated 4 tiles 5",
			codes:          []string{"1m", "1m", "1m", "1m", "4m", "4m", "4m", "4m", "7m", "8m"},
			wantShanten:    2,
			wantGoalsCount: 89,
		},
		{
			name:           "3 isolated 4 tiles",
			codes:          []string{"1m", "1m", "1m", "1m", "2m", "4m", "7m", "7m", "7m", "7m", "E", "E", "E", "E"},
			wantShanten:    1,
			wantGoalsCount: 1,
		},
		{
			name:           "4 honors 1",
			codes:          []string{"E", "E", "E", "E"},
			wantShanten:    1,
			wantGoalsCount: 33,
		},
		{
			name:           "4 honors 2",
			codes:          []string{"1m", "2m", "3m", "E", "E", "E", "E"},
			wantShanten:    1,
			wantGoalsCount: 34,
		},
		{
			name:           "4 honors 3",
			codes:          []string{"E", "E", "E", "E", "S", "S", "S", "S"},
			wantShanten:    1,
			wantGoalsCount: 32,
		},
		{
			name:           "4 honors 4",
			codes:          []string{"1m", "2m", "3m", "1p", "1p", "E", "E", "E", "E", "S", "S", "S", "S"},
			wantShanten:    2,
			wantGoalsCount: 37,
		},
		{
			name:           "can be interpreted in 2 set decompositions",
			codes:          []string{"1m", "1m", "1m", "2m", "2m", "2m", "3m", "3m", "3m", "7p", "8p", "9p", "9p", "9p"},
			wantShanten:    -1,
			wantGoalsCount: 2,
		},
		{
			name:           "can be interpreted in 4 set decompositions",
			codes:          []string{"3m", "3m", "3m", "4m", "4m", "4m", "5m", "5m", "5m", "6m", "6m", "6m", "7m", "7m"},
			wantShanten:    -1,
			wantGoalsCount: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand := hand.CodesToHand(tt.codes)
			got, got2 := service.AnalyzeShanten(hand)
			if got != tt.wantShanten {
				t.Errorf("AnalyzeShanten() = %v, want %v", got, tt.wantShanten)
			}
			if len(got2) != tt.wantGoalsCount {
				t.Errorf("AnalyzeShanten() = %v, want %v", len(got2), tt.wantGoalsCount)
			}
		})
	}
}

func TestAnalyzeShanten_Options(t *testing.T) {
	tests := []struct {
		name              string
		codes             []string
		allowedExtraTiles int
		upperBound        int
		wantShanten       int
		wantGoalsCount    int
	}{
		{
			name:              "case 1, allowedExtraTiles: 0, upperBound: 8",
			codes:             []string{"1m", "2m", "3m", "7m", "8m", "9m", "2s", "3s", "4s", "S", "S", "S", "W"},
			allowedExtraTiles: 0,
			upperBound:        8,
			wantShanten:       0,
			wantGoalsCount:    1,
		},
		{
			name:              "case 1, allowedExtraTiles: 1, upperBound: 8",
			codes:             []string{"1m", "2m", "3m", "7m", "8m", "9m", "2s", "3s", "4s", "S", "S", "S", "W"},
			allowedExtraTiles: 1,
			upperBound:        8,
			wantShanten:       0,
			wantGoalsCount:    42,
		},
		{
			name:              "case 1, allowedExtraTiles: 0, upperBound: -1",
			codes:             []string{"1m", "2m", "3m", "7m", "8m", "9m", "2s", "3s", "4s", "S", "S", "S", "W"},
			allowedExtraTiles: 0,
			upperBound:        -1,
			wantShanten:       service.InfinityShanten,
			wantGoalsCount:    0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand := hand.CodesToHand(tt.codes)
			got, got2 := service.AnalyzeShanten(
				hand,
				service.AllowedExtraTiles(tt.allowedExtraTiles),
				service.UpperBound(tt.upperBound),
			)
			if got != tt.wantShanten {
				t.Errorf("AnalyzeShanten() = %v, want %v", got, tt.wantShanten)
			}
			if len(got2) != tt.wantGoalsCount {
				t.Errorf("AnalyzeShanten() = %v, want %v", len(got2), tt.wantGoalsCount)
			}
		})
	}
}

func TestAnalyzeShantenChiitoitsu(t *testing.T) {
	tests := []struct {
		name  string
		codes []string
		want  int
	}{
		{
			name:  "without pair",
			codes: []string{"1m", "9m", "1p", "9p", "1s", "9s", "E", "S", "W", "N", "P", "F", "C"},
			want:  6,
		},
		{
			name:  "with quadruple",
			codes: []string{"1m", "1m", "8m", "8m", "2p", "8p", "8p", "5s", "5s", "E", "E", "E", "E"},
			want:  2,
		},
		{
			name:  "with triplet",
			codes: []string{"1m", "1m", "8m", "8m", "2p", "3p", "8p", "8p", "5s", "5s", "E", "E", "E"},
			want:  1,
		},
		{
			name:  "with 2 triplets",
			codes: []string{"1m", "1m", "8m", "8m", "2p", "8p", "8p", "5s", "5s", "5s", "E", "E", "E"},
			want:  2,
		},
		{
			name:  "tenpai",
			codes: []string{"1m", "1m", "8m", "8m", "2p", "8p", "8p", "5s", "5s", "E", "E", "C", "C"},
			want:  0,
		},
		{
			name:  "win",
			codes: []string{"1m", "1m", "8m", "8m", "2p", "2p", "8p", "8p", "5s", "5s", "E", "E", "C", "C"},
			want:  -1,
		},
		{
			name:  "incomplete_hand",
			codes: []string{"1m", "1m", "8m", "8m", "5s", "5s", "E", "E", "S", "S"},
			want:  service.InfinityShanten,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand := hand.CodesToHand(tt.codes)
			got := service.AnalyzeShantenChiitoitsu(hand)
			if got != tt.want {
				t.Errorf("AnalyzeShantenChiitoitsu() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnalyzeShantenKokushimusou(t *testing.T) {
	tests := []struct {
		name  string
		codes []string
		want  int
	}{
		{
			name:  "no terminals and honors",
			codes: []string{"2m", "3m", "4m", "5m", "5m", "3p", "4p", "5p", "4s", "5s", "6s", "7s", "8s"},
			want:  13,
		},
		{
			name:  "without pair",
			codes: []string{"1m", "8m", "9m", "1p", "2p", "2s", "4s", "9s", "E", "S", "W", "N", "P"},
			want:  4,
		},
		{
			name:  "with_pair",
			codes: []string{"1m", "1m", "9m", "1p", "2p", "2s", "9s", "9s", "E", "S", "W", "N", "P"},
			want:  3,
		},
		{
			name:  "tenpai",
			codes: []string{"1m", "1m", "1p", "9p", "1s", "9s", "E", "S", "W", "N", "P", "F", "C"},
			want:  0,
		},
		{
			name:  "tenpai 13 wait",
			codes: []string{"1m", "9m", "1p", "9p", "1s", "9s", "E", "S", "W", "N", "P", "F", "C"},
			want:  0,
		},
		{
			name:  "win",
			codes: []string{"1m", "1m", "9m", "1p", "9p", "1s", "9s", "E", "S", "W", "N", "P", "F", "C"},
			want:  -1,
		},
		{
			name:  "incomplete_hand",
			codes: []string{"9m", "1p", "9p", "1s", "9s", "E", "S", "W", "N", "P", "F", "C"},
			want:  service.InfinityShanten,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand := hand.CodesToHand(tt.codes)
			got := service.AnalyzeShantenKokushimusou(hand)
			if got != tt.want {
				t.Errorf("AnalyzeShantenKokushimusou() = %v, want %v", got, tt.want)
			}
		})
	}
}
