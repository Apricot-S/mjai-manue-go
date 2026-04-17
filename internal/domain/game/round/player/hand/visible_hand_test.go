package hand_test

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
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
		want  *hand.TileCounts34
	}{
		{
			name:  "empty hand",
			tiles: nil,
			want:  &hand.TileCounts34{},
		},
		{
			name:  "1m 1m",
			tiles: []tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
			want:  &hand.TileCounts34{0: 2},
		},
		{
			name:  "C",
			tiles: []tile.Tile{*tile.MustTileFromCode("C")},
			want:  &hand.TileCounts34{33: 1},
		},
		{
			name:  "5mr 5pr 5sr",
			tiles: []tile.Tile{*tile.MustTileFromCode("5mr"), *tile.MustTileFromCode("5pr"), *tile.MustTileFromCode("5sr")},
			want:  &hand.TileCounts34{4: 1, 13: 1, 22: 1},
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

func TestVisibleHand_Count(t *testing.T) {
	tests := []struct {
		name  string
		tiles []tile.Tile
		t     *tile.Tile
		want  int
	}{
		{
			name:  "empty hand",
			tiles: nil,
			t:     tile.MustTileFromCode("5m"),
			want:  0,
		},
		{
			name:  "has 1 tile",
			tiles: []tile.Tile{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5mr")},
			t:     tile.MustTileFromCode("5m"),
			want:  1,
		},
		{
			name:  "has 2 tile",
			tiles: []tile.Tile{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5mr")},
			t:     tile.MustTileFromCode("5m"),
			want:  2,
		},
		{
			name:  "red tile",
			tiles: []tile.Tile{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5mr")},
			t:     tile.MustTileFromCode("5mr"),
			want:  1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := hand.NewVisibleHand(tt.tiles)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got := h.Count(tt.t)
			if got != tt.want {
				t.Errorf("Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVisibleHand_Count_Panic(t *testing.T) {
	tests := []struct {
		name  string
		tiles []tile.Tile
	}{
		{
			name:  "empty hand",
			tiles: nil,
		},
		{
			name:  "non empty hand",
			tiles: []tile.Tile{*tile.MustTileFromCode("C")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Expected panic for unknown tile, but did not panic")
				}
			}()

			h, err := hand.NewVisibleHand(tt.tiles)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			h.Count(tile.MustTileFromCode("?"))
		})
	}
}

func TestVisibleHand_Draw(t *testing.T) {
	tests := []struct {
		name      string
		tiles     []tile.Tile
		tile      *tile.Tile
		wantTiles []tile.Tile
		wantErr   bool
	}{
		{
			name:      "visible hand cannot draw an unknown tile",
			tiles:     []tile.Tile{},
			tile:      tile.MustTileFromCode("?"),
			wantTiles: nil,
			wantErr:   true,
		},
		{
			name:      "visible tile",
			tiles:     []tile.Tile{*tile.MustTileFromCode("1m")},
			tile:      tile.MustTileFromCode("1m"),
			wantTiles: []tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
			wantErr:   false,
		},
		{
			name:      "hand can draw 4th identical tiles",
			tiles:     []tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
			tile:      tile.MustTileFromCode("1m"),
			wantTiles: []tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
			wantErr:   false,
		},
		{
			name:      "hand cannot draw 5th identical tiles",
			tiles:     []tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
			tile:      tile.MustTileFromCode("1m"),
			wantTiles: nil,
			wantErr:   true,
		},
		{
			name:      "hand cannot draw 2nd red fives",
			tiles:     []tile.Tile{*tile.MustTileFromCode("5mr")},
			tile:      tile.MustTileFromCode("5mr"),
			wantTiles: nil,
			wantErr:   true,
		},
		{
			name:      "hand can draw 4th normal fives",
			tiles:     []tile.Tile{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m")},
			tile:      tile.MustTileFromCode("5m"),
			wantTiles: []tile.Tile{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m")},
			wantErr:   false,
		},
		{
			name: "hand can draw 14th tile",
			tiles: []tile.Tile{
				*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"),
				*tile.MustTileFromCode("2m"), *tile.MustTileFromCode("2m"), *tile.MustTileFromCode("2m"),
				*tile.MustTileFromCode("3m"), *tile.MustTileFromCode("3m"), *tile.MustTileFromCode("3m"),
				*tile.MustTileFromCode("4m"), *tile.MustTileFromCode("4m"), *tile.MustTileFromCode("4m"),
				*tile.MustTileFromCode("5m"),
			},
			tile: tile.MustTileFromCode("5m"),
			wantTiles: []tile.Tile{
				*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"),
				*tile.MustTileFromCode("2m"), *tile.MustTileFromCode("2m"), *tile.MustTileFromCode("2m"),
				*tile.MustTileFromCode("3m"), *tile.MustTileFromCode("3m"), *tile.MustTileFromCode("3m"),
				*tile.MustTileFromCode("4m"), *tile.MustTileFromCode("4m"), *tile.MustTileFromCode("4m"),
				*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"),
			},
			wantErr: false,
		},
		{
			name: "hand cannot draw 15th tile",
			tiles: []tile.Tile{
				*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"),
				*tile.MustTileFromCode("2m"), *tile.MustTileFromCode("2m"), *tile.MustTileFromCode("2m"),
				*tile.MustTileFromCode("3m"), *tile.MustTileFromCode("3m"), *tile.MustTileFromCode("3m"),
				*tile.MustTileFromCode("4m"), *tile.MustTileFromCode("4m"), *tile.MustTileFromCode("4m"),
				*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"),
			},
			tile:      tile.MustTileFromCode("5m"),
			wantTiles: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := hand.NewVisibleHand(tt.tiles)
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

func TestVisibleHand_Discard(t *testing.T) {
	tests := []struct {
		name      string
		tiles     []tile.Tile
		tile      *tile.Tile
		wantTiles []tile.Tile
		wantErr   bool
	}{
		{
			name:      "visible hand cannot discard an unknown tile",
			tiles:     []tile.Tile{*tile.MustTileFromCode("1m")},
			tile:      tile.MustTileFromCode("?"),
			wantTiles: nil,
			wantErr:   true,
		},
		{
			name:      "visible tile",
			tiles:     []tile.Tile{*tile.MustTileFromCode("1m")},
			tile:      tile.MustTileFromCode("1m"),
			wantTiles: []tile.Tile{},
			wantErr:   false,
		},
		{
			name:      "cannot discard a tile that are not in the hand",
			tiles:     []tile.Tile{*tile.MustTileFromCode("5m")},
			tile:      tile.MustTileFromCode("5mr"),
			wantTiles: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := hand.NewVisibleHand(tt.tiles)
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

func TestVisibleHand_Call(t *testing.T) {
	tests := []struct {
		name      string
		tiles     []tile.Tile
		meld      meld.Meld
		wantTiles []tile.Tile
		wantErr   bool
	}{
		{
			name:  "chii",
			tiles: []tile.Tile{*tile.MustTileFromCode("2m"), *tile.MustTileFromCode("3m")},
			meld: meld.MustChii(
				*tile.MustTileFromCode("1m"),
				[2]tile.Tile{*tile.MustTileFromCode("2m"), *tile.MustTileFromCode("3m")},
				*seat.MustSeat(0),
			),
			wantTiles: []tile.Tile{},
			wantErr:   false,
		},
		{
			name:  "pon",
			tiles: []tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
			meld: meld.MustPon(
				*tile.MustTileFromCode("1m"),
				[2]tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
				*seat.MustSeat(0),
			),
			wantTiles: []tile.Tile{},
			wantErr:   false,
		},
		{
			name:  "called kan",
			tiles: []tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
			meld: meld.MustCalledKan(
				*tile.MustTileFromCode("1m"),
				[3]tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
				*seat.MustSeat(0),
			),
			wantTiles: []tile.Tile{},
			wantErr:   false,
		},
		{
			name:  "concealed kan",
			tiles: []tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
			meld: meld.MustConcealedKan(
				[4]tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
			),
			wantTiles: []tile.Tile{},
			wantErr:   false,
		},
		{
			name:  "promoted kan",
			tiles: []tile.Tile{*tile.MustTileFromCode("5mr")},
			meld: meld.MustPromotedKan(
				*tile.MustTileFromCode("5m"),
				[2]tile.Tile{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m")},
				*tile.MustTileFromCode("5mr"),
				*seat.MustSeat(0),
			),
			wantTiles: []tile.Tile{},
			wantErr:   false,
		},
		{
			name:  "call fails when meld contains red/normal mismatch",
			tiles: []tile.Tile{*tile.MustTileFromCode("4m"), *tile.MustTileFromCode("5mr")},
			meld: meld.MustChii(
				*tile.MustTileFromCode("6m"),
				[2]tile.Tile{*tile.MustTileFromCode("4m"), *tile.MustTileFromCode("5m")},
				*seat.MustSeat(0),
			),
			wantTiles: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := hand.NewVisibleHand(tt.tiles)
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
