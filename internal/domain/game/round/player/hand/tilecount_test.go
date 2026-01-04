package hand_test

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestTileCounts34_ToTiles(t *testing.T) {
	tests := []struct {
		name string
		tc34 hand.TileCounts34
		want []tile.Tile
	}{
		{
			name: "empty tile counts 34",
			tc34: hand.TileCounts34{},
			want: []tile.Tile{},
		},
		{
			name: "tile counts 34 can contain five identical tiles",
			tc34: hand.TileCounts34{0: 5},
			want: []tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.tc34.ToTiles()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToTiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTileCounts34_NumTiles(t *testing.T) {
	tests := []struct {
		name string
		tc34 hand.TileCounts34
		want int
	}{
		{
			name: "empty tile counts 34",
			tc34: hand.TileCounts34{},
			want: 0,
		},
		{
			name: "tile counts 34 can contain five identical tiles",
			tc34: hand.TileCounts34{0: 5, 33: 1},
			want: 6,
		},
		{
			name: "tile counts 34 can contain negative count",
			tc34: hand.TileCounts34{0: -1, 33: -1},
			want: -2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.tc34.NumTiles()
			if got != tt.want {
				t.Errorf("NumTiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
