package hand_test

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
)

func TestNewHand(t *testing.T) {
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
			name:    "hand can contain four identical tiles",
			tiles:   []tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
			want:    []tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
			wantErr: false,
		},
		{
			name:    "hand cannot contain five identical tiles",
			tiles:   []tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "hand can contain five unknown tiles",
			tiles:   []tile.Tile{*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?")},
			want:    []tile.Tile{*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?")},
			wantErr: false,
		},
		{
			name:    "hand can contain two red fives",
			tiles:   []tile.Tile{*tile.MustTileFromCode("5mr"), *tile.MustTileFromCode("5mr")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "hand cannot contain four normal fives",
			tiles:   []tile.Tile{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m")},
			want:    nil,
			wantErr: true,
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
			got, gotErr := hand.NewHand(tt.tiles)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewHand() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewHand() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got.ToTiles(), tt.tiles) {
				t.Errorf("NewHand() = %v, want %v", got, tt.want)
			}
		})
	}
}
