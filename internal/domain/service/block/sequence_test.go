package block_test

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/service/block"
)

func TestNewSequence(t *testing.T) {
	tests := []struct {
		name    string
		tile    tile.Tile
		want    []tile.Tile
		wantErr bool
	}{
		{
			name:    "can create sequence starting with 7",
			tile:    *tile.MustTileFromCode("7s"),
			want:    []tile.Tile{*tile.MustTileFromCode("7s"), *tile.MustTileFromCode("8s"), *tile.MustTileFromCode("9s")},
			wantErr: false,
		},
		{
			name:    "cannot create sequence starting with 8",
			tile:    *tile.MustTileFromCode("8m"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "cannot create sequence from honors",
			tile:    *tile.MustTileFromCode("C"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "cannot create sequence from unknown tile",
			tile:    *tile.MustTileFromCode("?"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "cannot create sequence from red five",
			tile:    *tile.MustTileFromCode("5mr"),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := block.NewSequence(tt.tile)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewSequence() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewSequence() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got.ToTiles(), tt.want) {
				t.Errorf("NewSequence().ToTiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
