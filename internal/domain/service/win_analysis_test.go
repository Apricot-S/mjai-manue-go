package service

import (
	"strings"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/hand"
)

func Test_isWinningFormGeneral(t *testing.T) {
	tests := []struct {
		codes []string
		want  bool
	}{
		{
			codes: []string{"3m", "3m", "1p", "2p", "3p", "1s", "2s", "3s", "4s", "5s", "6s", "7s", "8s", "9s"},
			want:  true,
		},
		{
			codes: []string{"1p", "1p", "1p", "2p", "3p", "1s", "2s", "3s", "4s", "5s", "6s", "7s", "8s", "9s"},
			want:  true,
		},
		{
			codes: []string{"1s", "2s", "3s", "4s", "5s", "6s", "7s", "8s", "9s", "E", "E", "C", "C", "C"},
			want:  true,
		},
		{
			codes: []string{"1s", "2s", "3s", "4s", "5s", "5s", "5s", "6s", "7s", "7s", "8s", "8s", "9s", "9s"},
			want:  true,
		},
		{
			codes: []string{"1s", "1s", "1s", "2s", "3s", "4s", "5s", "6s", "7s", "8s", "8s", "9s", "9s", "9s"},
			want:  true,
		},
		{
			codes: []string{"3m", "4m", "5m", "7p", "8p", "9p", "2s", "3s", "3s", "3s", "3s", "4s", "P", "P"},
			want:  true,
		},
		{
			codes: []string{"1m", "1m", "1m", "2m", "3m", "3m", "3m", "4m", "4m", "4m", "5m", "5m", "9m", "9m"},
			want:  true,
		},
		{
			codes: []string{"1m", "1m"},
			want:  true,
		},
		{
			codes: []string{},
			want:  true,
		},
		{
			codes: []string{"1m", "1m", "E", "E", "E"},
			want:  true,
		},

		{
			codes: []string{"1p", "2p", "3p", "4p", "5p", "1s", "2s", "3s", "4s", "5s", "6s", "7s", "8s", "9s"},
			want:  false,
		},
		{
			codes: []string{"1p", "1p", "1p", "4p", "5p", "1s", "1s", "1s", "2s", "2s", "2s", "4s", "4s", "4s"},
			want:  false,
		},
		{
			codes: []string{"1s", "1s", "1s", "2s", "2s", "2s", "3s", "3s", "3s", "5s", "6s", "8s", "8s", "8s"},
			want:  false,
		},
		{
			codes: []string{"2p", "2p", "2p", "2s", "3s", "4s", "6s", "7s"},
			want:  false,
		},
	}
	for _, tt := range tests {
		name := strings.Join(tt.codes, " ")
		t.Run(name, func(t *testing.T) {
			h := hand.CodesToHand(tt.codes)
			tc34 := h.ToTileCounts34()
			got := isWinningFormGeneral(tc34)
			if got != tt.want {
				t.Errorf("isWinningFormGeneral() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isSingleColorWinningFormWithoutPair(t *testing.T) {
	tests := []struct {
		name  string
		codes []string
		want  bool
	}{
		{
			name:  "winning form",
			codes: []string{"1m", "1m", "1m", "1m", "2m", "2m", "3m", "3m", "4m"},
			want:  true,
		},
		{
			name:  "not winning form: with pair",
			codes: []string{"1m", "1m", "1m", "1m", "2m", "2m", "3m", "3m"},
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := hand.CodesToHand(tt.codes)
			tc34 := h.ToTileCounts34()
			singleColorHand := tc34[:9]
			got := isSingleColorWinningFormWithoutPair(singleColorHand)
			if got != tt.want {
				t.Errorf("isSingleColorWinningFormWithoutPair() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isSingleColorWinningFormWithPair(t *testing.T) {
	tests := []struct {
		name  string
		codes []string
		want  bool
	}{
		{
			name:  "winning form",
			codes: []string{"1m", "1m", "1m", "1m", "2m", "2m", "3m", "3m"},
			want:  true,
		},
		{
			name:  "not winning form: without pair",
			codes: []string{"1m", "1m", "1m", "1m", "2m", "2m", "3m", "3m", "4m"},
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := hand.CodesToHand(tt.codes)
			tc34 := h.ToTileCounts34()
			singleColorHand := tc34[:9]
			got := isSingleColorWinningFormWithPair(singleColorHand)
			if got != tt.want {
				t.Errorf("isSingleColorWinningFormWithPair() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isWinningFormChitoitsu(t *testing.T) {
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
			h := hand.CodesToHand(tt.codes)
			tc34 := h.ToTileCounts34()
			got := isWinningFormChitoitsu(tc34)
			if got != tt.want {
				t.Errorf("isWinningFormChitoitsu() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isWinningFormKokushimuso(t *testing.T) {
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
			h := hand.CodesToHand(tt.codes)
			tc34 := h.ToTileCounts34()
			got := isWinningFormKokushimuso(tc34)
			if got != tt.want {
				t.Errorf("isWinningFormKokushimuso() = %v, want %v", got, tt.want)
			}
		})
	}
}
