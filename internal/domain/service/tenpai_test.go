package service_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/service"
)

func TestIsTenpaiGeneral(t *testing.T) {
	tests := []struct {
		name  string
		codes []string
		want  bool
	}{
		{
			name:  "empty : An empty hand is one step away from being a pair wait -> noten",
			codes: nil,
			want:  false,
		},
		{
			name:  "chiitoitsu",
			codes: []string{"1m", "1m", "8m", "8m", "2p", "8p", "8p", "5s", "5s", "E", "E", "C", "C"},
			want:  false,
		},
		{
			name:  "kokushimusou",
			codes: []string{"1m", "9m", "1p", "9p", "1s", "9s", "E", "S", "W", "N", "P", "F", "C"},
			want:  false,
		},
		{
			name:  "tenpai",
			codes: []string{"1m", "2m", "3m", "4p", "5pr", "6p", "7s", "8s", "9s", "E", "E", "S", "S"},
			want:  true,
		},
		{
			name:  "win",
			codes: []string{"1m", "2m", "3m", "4p", "5pr", "6p", "7s", "8s", "9s", "E", "E", "S", "S", "S"},
			want:  true,
		},
		{
			name:  "with meld",
			codes: []string{"1m", "2m", "3m", "4p", "5pr", "6p", "7s", "8s", "9s", "E"},
			want:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand := hand.CodesToHand(tt.codes)
			got := service.IsTenpaiGeneral(hand)
			if got != tt.want {
				t.Errorf("IsTenpaiGeneral() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsTenpaiAll(t *testing.T) {
	tests := []struct {
		name  string
		codes []string
		want  bool
	}{
		{
			name:  "empty : An empty hand is one step away from being a pair wait -> noten",
			codes: nil,
			want:  false,
		},
		{
			name:  "chiitoitsu",
			codes: []string{"1m", "1m", "8m", "8m", "2p", "8p", "8p", "5s", "5s", "E", "E", "C", "C"},
			want:  true,
		},
		{
			name:  "kokushimusou",
			codes: []string{"1m", "9m", "1p", "9p", "1s", "9s", "E", "S", "W", "N", "P", "F", "C"},
			want:  true,
		},
		{
			name:  "tenpai",
			codes: []string{"1m", "2m", "3m", "4p", "5pr", "6p", "7s", "8s", "9s", "E", "E", "S", "S"},
			want:  true,
		},
		{
			name:  "win",
			codes: []string{"1m", "2m", "3m", "4p", "5pr", "6p", "7s", "8s", "9s", "E", "E", "S", "S", "S"},
			want:  true,
		},
		{
			name:  "with meld",
			codes: []string{"1m", "2m", "3m", "4p", "5pr", "6p", "7s", "8s", "9s", "E"},
			want:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand := hand.CodesToHand(tt.codes)
			got := service.IsTenpaiAll(hand)
			if got != tt.want {
				t.Errorf("IsTenpaiAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
