package hand_test

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
)

func TestNewInvisibleHand(t *testing.T) {
	tests := []struct {
		name    string
		tiles   []tile.Tile
		want    []tile.Tile
		wantErr bool
	}{
		{
			name:    "empty hand",
			tiles:   []tile.Tile{},
			want:    []tile.Tile{},
			wantErr: false,
		},
		{
			name:    "unknown hand",
			tiles:   []tile.Tile{*tile.MustTileFromCode("?")},
			want:    []tile.Tile{*tile.MustTileFromCode("?")},
			wantErr: false,
		},
		{
			name:    "visible hand",
			tiles:   []tile.Tile{*tile.MustTileFromCode("1m")},
			want:    []tile.Tile{*tile.MustTileFromCode("?")},
			wantErr: false,
		},
		{
			name:    "invisible hand does not validate the number of identical tiles",
			tiles:   []tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
			want:    []tile.Tile{*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?")},
			wantErr: false,
		},
		{
			name: "hand can contain 14 tiles",
			tiles: []tile.Tile{
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
			},
			want: []tile.Tile{
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
			},
			wantErr: false,
		},
		{
			name: "hand cannot contain 15 tiles",
			tiles: []tile.Tile{
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := hand.NewInvisibleHand(tt.tiles)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewInvisibleHand() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewInvisibleHand() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got.ToTiles(), tt.want) {
				t.Errorf("NewInvisibleHand() = %v, want %v", got.ToTiles(), tt.want)
			}
		})
	}
}
