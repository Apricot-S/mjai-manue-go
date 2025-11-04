package tilecount_test

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tilecount"
)

func TestTileCounts34_ToTiles(t *testing.T) {
	tests := []struct {
		name string
		tc34 tilecount.TileCounts34
		want []tile.Tile
	}{
		{
			name: "empty tile counts 34",
			tc34: tilecount.TileCounts34{},
			want: []tile.Tile{},
		},
		{
			name: "tile counts 34 can contain five identical tiles",
			tc34: tilecount.TileCounts34{0: 5},
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
