package service_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/service"
)

func TestIsWinningFormChitoitsu(t *testing.T) {
	tests := []struct {
		name  string
		codes []string
		want  bool
	}{
		{
			name:  "winning form",
			codes: []string{"1p", "1p", "9p", "9p", "1s", "1s", "3s", "3s", "5s", "5s", "7s", "7s", "9s", "9s"},
			want:  true,
		},
		{
			name:  "winning form",
			codes: []string{"1m", "1m", "1p", "1p", "9p", "9p", "2s", "2s", "4s", "4s", "S", "S", "C", "C"},
			want:  true,
		},
		{
			name:  "winning form: Ryanpeko",
			codes: []string{"1m", "1m", "2m", "2m", "3m", "3m", "4m", "4m", "5m", "5m", "6m", "6m", "7m", "7m"},
			want:  true,
		},
		{
			name:  "not winning form: has 4 1p",
			codes: []string{"1p", "1p", "1p", "1p", "3p", "3p", "4p", "4p", "5p", "5p", "6p", "6p", "7p", "7p"},
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand := hand.CodesToHand(tt.codes)
			got := service.IsWinningFormChitoitsu(hand)
			if got != tt.want {
				t.Errorf("IsWinningFormChitoitsu() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsWinningFormKokushimuso(t *testing.T) {
	tests := []struct {
		name  string
		codes []string
		want  bool
	}{
		{
			name:  "winning form: 9m pair",
			codes: []string{"1m", "9m", "9m", "1p", "9p", "1s", "9s", "E", "S", "W", "N", "P", "F", "C"},
			want:  true,
		},
		{
			name:  "winning form: E pair",
			codes: []string{"1m", "9m", "1p", "9p", "1s", "9s", "E", "E", "S", "W", "N", "P", "F", "C"},
			want:  true,
		},
		{
			name:  "winning form: C pair",
			codes: []string{"1m", "9m", "1p", "9p", "1s", "9s", "E", "S", "W", "N", "P", "F", "C", "C"},
			want:  true,
		},
		{
			name:  "not winning form: has 2m",
			codes: []string{"1m", "2m", "9m", "1p", "9p", "1s", "9s", "E", "S", "W", "N", "P", "F", "C"},
			want:  false,
		},
		{
			name:  "not winning form: has 3 E",
			codes: []string{"1m", "9m", "1p", "9p", "1s", "9s", "E", "E", "E", "W", "N", "P", "F", "C"},
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand := hand.CodesToHand(tt.codes)
			got := service.IsWinningFormKokushimuso(hand)
			if got != tt.want {
				t.Errorf("IsWinningFormKokushimuso() = %v, want %v", got, tt.want)
			}
		})
	}
}
