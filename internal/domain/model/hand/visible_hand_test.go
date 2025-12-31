package hand_test

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tilecount"
)

func TestNewVisibleHand(t *testing.T) {
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
			name:    "visible hand cannot contain an unknown tile",
			tiles:   []tile.Tile{*tile.MustTileFromCode("?")},
			want:    nil,
			wantErr: true,
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
			name:    "hand cannot contain two red fives",
			tiles:   []tile.Tile{*tile.MustTileFromCode("5mr"), *tile.MustTileFromCode("5mr")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "hand can contain four normal fives",
			tiles:   []tile.Tile{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m")},
			want:    []tile.Tile{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m")},
			wantErr: false,
		},
		{
			name: "hand can contain 14 tiles",
			tiles: []tile.Tile{
				*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"),
				*tile.MustTileFromCode("2m"), *tile.MustTileFromCode("2m"), *tile.MustTileFromCode("2m"),
				*tile.MustTileFromCode("3m"), *tile.MustTileFromCode("3m"), *tile.MustTileFromCode("3m"),
				*tile.MustTileFromCode("4m"), *tile.MustTileFromCode("4m"), *tile.MustTileFromCode("4m"),
				*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"),
			},
			want: []tile.Tile{
				*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"),
				*tile.MustTileFromCode("2m"), *tile.MustTileFromCode("2m"), *tile.MustTileFromCode("2m"),
				*tile.MustTileFromCode("3m"), *tile.MustTileFromCode("3m"), *tile.MustTileFromCode("3m"),
				*tile.MustTileFromCode("4m"), *tile.MustTileFromCode("4m"), *tile.MustTileFromCode("4m"),
				*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"),
			},
			wantErr: false,
		},
		{
			name: "hand cannot contain 15 tiles",
			tiles: []tile.Tile{
				*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"),
				*tile.MustTileFromCode("2m"), *tile.MustTileFromCode("2m"), *tile.MustTileFromCode("2m"),
				*tile.MustTileFromCode("3m"), *tile.MustTileFromCode("3m"), *tile.MustTileFromCode("3m"),
				*tile.MustTileFromCode("4m"), *tile.MustTileFromCode("4m"), *tile.MustTileFromCode("4m"),
				*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := hand.NewVisibleHand(tt.tiles)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewVisibleHand() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewVisibleHand() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got.ToTiles(), tt.want) {
				t.Errorf("NewVisibleHand().ToTiles() = %v, want %v", got.ToTiles(), tt.want)
			}
		})
	}
}

func TestVisibleHand_ToTileCounts34(t *testing.T) {
	tests := []struct {
		name  string
		tiles []tile.Tile
		want  *tilecount.TileCounts34
	}{
		{
			name:  "empty hand",
			tiles: nil,
			want:  &tilecount.TileCounts34{},
		},
		{
			name:  "1m 1m",
			tiles: []tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
			want:  &tilecount.TileCounts34{0: 2},
		},
		{
			name:  "C",
			tiles: []tile.Tile{*tile.MustTileFromCode("C")},
			want:  &tilecount.TileCounts34{33: 1},
		},
		{
			name:  "5mr 5pr 5sr",
			tiles: []tile.Tile{*tile.MustTileFromCode("5mr"), *tile.MustTileFromCode("5pr"), *tile.MustTileFromCode("5sr")},
			want:  &tilecount.TileCounts34{4: 1, 13: 1, 22: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := hand.NewVisibleHand(tt.tiles)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got := h.ToTileCounts34()
			if *got != *tt.want {
				t.Errorf("ToTileCounts34() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVisibleHand_ToTileCounts34_HandAndTileCountsAreIndependent(t *testing.T) {
	tiles := []tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")}
	hand, err := hand.NewVisibleHand(tiles)
	if err != nil {
		t.Fatalf("could not construct receiver type: %v", err)
	}

	counts1 := hand.ToTileCounts34()
	counts1[0]++

	counts2 := hand.ToTileCounts34()
	if counts2[0] != 2 {
		t.Errorf("expected counts2[0] to be 2, but got %v", counts2)
	}
}
