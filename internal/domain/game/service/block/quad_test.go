package block_test

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/service/block"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestNewQuad(t *testing.T) {
	tests := []struct {
		name    string
		tile    tile.Tile
		want    []tile.Tile
		wantErr bool
	}{
		{
			name:    "C quad",
			tile:    *tile.MustTileFromCode("C"),
			want:    []tile.Tile{*tile.MustTileFromCode("C"), *tile.MustTileFromCode("C"), *tile.MustTileFromCode("C"), *tile.MustTileFromCode("C")},
			wantErr: false,
		},
		{
			name:    "cannot create quad from unknown tile",
			tile:    *tile.MustTileFromCode("?"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "cannot create quad from red five",
			tile:    *tile.MustTileFromCode("5mr"),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := block.NewQuad(tt.tile)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewQuad() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewQuad() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got.ToTiles(), tt.want) {
				t.Errorf("NewQuad().ToTiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
