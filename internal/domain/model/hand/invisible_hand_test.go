package hand_test

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/playerid"
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

func TestInvisibleHand_Draw(t *testing.T) {
	tests := []struct {
		name      string
		tiles     []tile.Tile
		tile      *tile.Tile
		wantTiles []tile.Tile
		wantErr   bool
	}{
		{
			name:      "unknown tile",
			tiles:     []tile.Tile{},
			tile:      tile.MustTileFromCode("?"),
			wantTiles: []tile.Tile{*tile.MustTileFromCode("?")},
			wantErr:   false,
		},
		{
			name:      "visible tile",
			tiles:     []tile.Tile{*tile.MustTileFromCode("?")},
			tile:      tile.MustTileFromCode("1m"),
			wantTiles: []tile.Tile{*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?")},
			wantErr:   false,
		},
		{
			name: "hand can draw 14th tile",
			tiles: []tile.Tile{
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"),
			},
			tile: tile.MustTileFromCode("?"),
			wantTiles: []tile.Tile{
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
			},
			wantErr: false,
		},
		{
			name: "hand cannot draw 15th tile",
			tiles: []tile.Tile{
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
			},
			tile:      tile.MustTileFromCode("?"),
			wantTiles: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := hand.NewInvisibleHand(tt.tiles)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got, gotErr := h.Draw(tt.tile)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Draw() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Draw() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got.ToTiles(), tt.wantTiles) {
				t.Errorf("Draw() = %v, want %v", got.ToTiles(), tt.wantTiles)
			}
		})
	}
}

func TestInvisibleHand_Discard(t *testing.T) {
	tests := []struct {
		name      string
		tiles     []tile.Tile
		tile      *tile.Tile
		wantTiles []tile.Tile
		wantErr   bool
	}{
		{
			name:      "unknown tile",
			tiles:     []tile.Tile{*tile.MustTileFromCode("?")},
			tile:      tile.MustTileFromCode("?"),
			wantTiles: []tile.Tile{},
			wantErr:   false,
		},
		{
			name:      "visible tile",
			tiles:     []tile.Tile{*tile.MustTileFromCode("?")},
			tile:      tile.MustTileFromCode("1m"),
			wantTiles: []tile.Tile{},
			wantErr:   false,
		},
		{
			name:      "cannot discard from an empty hand",
			tiles:     []tile.Tile{},
			tile:      tile.MustTileFromCode("?"),
			wantTiles: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := hand.NewInvisibleHand(tt.tiles)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got, gotErr := h.Discard(tt.tile)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Discard() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Discard() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got.ToTiles(), tt.wantTiles) {
				t.Errorf("Discard() = %v, want %v", got.ToTiles(), tt.wantTiles)
			}
		})
	}
}

func TestInvisibleHand_Call(t *testing.T) {
	tests := []struct {
		name      string
		tiles     []tile.Tile
		meld      meld.Meld
		wantTiles []tile.Tile
		wantErr   bool
	}{
		{
			name:  "chii",
			tiles: []tile.Tile{*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?")},
			meld: meld.MustChii(
				*tile.MustTileFromCode("1m"),
				[2]tile.Tile{*tile.MustTileFromCode("2m"), *tile.MustTileFromCode("3m")},
				*playerid.MustPlayerID(0),
			),
			wantTiles: []tile.Tile{},
			wantErr:   false,
		},
		{
			name:  "pon",
			tiles: []tile.Tile{*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?")},
			meld: meld.MustPon(
				*tile.MustTileFromCode("1m"),
				[2]tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
				*playerid.MustPlayerID(0),
			),
			wantTiles: []tile.Tile{},
			wantErr:   false,
		},
		{
			name:  "called kan",
			tiles: []tile.Tile{*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?")},
			meld: meld.MustCalledKan(
				*tile.MustTileFromCode("1m"),
				[3]tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
				*playerid.MustPlayerID(0),
			),
			wantTiles: []tile.Tile{},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := hand.NewInvisibleHand(tt.tiles)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got, gotErr := h.Call(tt.meld)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Call() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Call() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got.ToTiles(), tt.wantTiles) {
				t.Errorf("Call() = %v, want %v", got.ToTiles(), tt.wantTiles)
			}
		})
	}
}
