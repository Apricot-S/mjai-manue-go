package service_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/service"
)

func codesToHand(codes []string) *hand.Hand {
	tiles := make([]tile.Tile, len(codes))
	for i, code := range codes {
		tiles[i] = *tile.MustTileFromCode(code)
	}

	h, err := hand.NewHand(tiles)
	if err != nil {
		panic(err)
	}

	return h
}

func TestAnalyzeShantenChitoitsu(t *testing.T) {
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
			hand := codesToHand(tt.codes)
			got := service.AnalyzeShantenChitoitsu(hand)
			if got != tt.want {
				t.Errorf("AnalyzeShantenChitoitsu() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnalyzeShantenKokushimuso(t *testing.T) {
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
			hand := codesToHand(tt.codes)
			got := service.AnalyzeShantenKokushimuso(hand)
			if got != tt.want {
				t.Errorf("AnalyzeShantenKokushimuso() = %v, want %v", got, tt.want)
			}
		})
	}
}
