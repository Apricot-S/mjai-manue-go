package block_test

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service/block"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestNewTriplet(t *testing.T) {
	tests := []struct {
		name    string
		tile    tile.Tile
		want    []tile.Tile
		wantErr bool
	}{
		{
			name:    "C triplet",
			tile:    tile.MustTileFromCode("C"),
			want:    []tile.Tile{tile.MustTileFromCode("C"), tile.MustTileFromCode("C"), tile.MustTileFromCode("C")},
			wantErr: false,
		},
		{
			name:    "cannot create triplet from unknown tile",
			tile:    tile.MustTileFromCode("?"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "cannot create triplet from red five",
			tile:    tile.MustTileFromCode("5mr"),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := block.NewTriplet(tt.tile)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewTriplet() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewTriplet() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got.ToTiles(), tt.want) {
				t.Errorf("NewTriplet().ToTiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
