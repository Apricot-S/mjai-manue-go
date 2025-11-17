package service_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/service"
)

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
			hand := codesToHand(tt.codes)
			got := service.IsWinningFormKokushimuso(hand)
			if got != tt.want {
				t.Errorf("IsWinningFormKokushimuso() = %v, want %v", got, tt.want)
			}
		})
	}
}
