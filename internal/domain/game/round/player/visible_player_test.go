package player_test

import (
	"reflect"
	"slices"
	"sort"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestNewVisiblePlayer(t *testing.T) {
	tests := []struct {
		name      string
		handTiles []tile.Tile
		wantHand  *hand.VisibleHand
		wantErr   bool
	}{
		{
			name: "valid",
			handTiles: []tile.Tile{
				*tile.MustTileFromCode("C"), *tile.MustTileFromCode("9s"), *tile.MustTileFromCode("4m"),
				*tile.MustTileFromCode("2p"), *tile.MustTileFromCode("S"), *tile.MustTileFromCode("4p"),
				*tile.MustTileFromCode("8s"), *tile.MustTileFromCode("6p"), *tile.MustTileFromCode("6s"),
				*tile.MustTileFromCode("7m"), *tile.MustTileFromCode("9s"), *tile.MustTileFromCode("5pr"),
				*tile.MustTileFromCode("5p"),
			},
			wantHand: hand.CodesToHand([]string{"4m", "7m", "2p", "4p", "5p", "5pr", "6p", "6s", "8s", "9s", "9s", "S", "C"}),
			wantErr:  false,
		},
		{
			name: "invalid: 12 tiles",
			handTiles: []tile.Tile{
				*tile.MustTileFromCode("C"), *tile.MustTileFromCode("9s"), *tile.MustTileFromCode("4m"),
				*tile.MustTileFromCode("2p"), *tile.MustTileFromCode("S"), *tile.MustTileFromCode("4p"),
				*tile.MustTileFromCode("8s"), *tile.MustTileFromCode("6p"), *tile.MustTileFromCode("6s"),
				*tile.MustTileFromCode("7m"), *tile.MustTileFromCode("9s"), *tile.MustTileFromCode("5pr"),
			},
			wantHand: nil,
			wantErr:  true,
		},
		{
			name: "invalid: 14 tiles",
			handTiles: []tile.Tile{
				*tile.MustTileFromCode("C"), *tile.MustTileFromCode("9s"), *tile.MustTileFromCode("4m"),
				*tile.MustTileFromCode("2p"), *tile.MustTileFromCode("S"), *tile.MustTileFromCode("4p"),
				*tile.MustTileFromCode("8s"), *tile.MustTileFromCode("6p"), *tile.MustTileFromCode("6s"),
				*tile.MustTileFromCode("7m"), *tile.MustTileFromCode("9s"), *tile.MustTileFromCode("5pr"),
				*tile.MustTileFromCode("5p"), *tile.MustTileFromCode("1m"),
			},
			wantHand: nil,
			wantErr:  true,
		},
		{
			name: "invalid: unknown tiles",
			handTiles: []tile.Tile{
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"),
			},
			wantHand: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := player.NewVisiblePlayer(tt.handTiles)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewVisiblePlayer() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewVisiblePlayer() succeeded unexpectedly")
			}

			h, ok := got.Hand()
			if !ok {
				t.Errorf("NewVisiblePlayer() Hand() returned not ok")
			}
			if *h != *tt.wantHand {
				t.Errorf("NewVisiblePlayer().Hand() = %v, want %v", h, tt.wantHand)
			}

			ts := tile.Tiles(tt.wantHand.ToTiles())
			sort.Sort(ts)
			if !reflect.DeepEqual(got.HandTiles(), []tile.Tile(ts)) {
				t.Errorf("NewVisiblePlayer().HandTiles() = %v, want %v", got.HandTiles(), ts)
			}

			if got.DrawnTile() != nil {
				t.Errorf("NewVisiblePlayer().DrawnTile() = %v, want %v", got.DrawnTile(), nil)
			}

			melds := make([]meld.Meld, 0, 4)
			if !reflect.DeepEqual(got.Melds(), melds) {
				t.Errorf("NewVisiblePlayer().Melds() = %v, want %v", got.Melds(), melds)
			}

			river := make([]tile.Tile, 0, 24)
			if !reflect.DeepEqual(got.River(), river) {
				t.Errorf("NewVisiblePlayer().River() = %v, want %v", got.River(), river)
			}

			discardedTiles := make([]tile.Tile, 0, 27)
			if !reflect.DeepEqual(got.DiscardedTiles(), discardedTiles) {
				t.Errorf("NewVisiblePlayer().DiscardedTiles() = %v, want %v", got.DiscardedTiles(), discardedTiles)
			}

			extraSafeTiles := make([]tile.Tile, 0, 3)
			if !reflect.DeepEqual(got.ExtraSafeTiles(), extraSafeTiles) {
				t.Errorf("NewVisiblePlayer().ExtraSafeTiles() = %v, want %v", got.ExtraSafeTiles(), extraSafeTiles)
			}

			if got.RiichiState() != player.NotRiichi {
				t.Errorf("NewVisiblePlayer().RiichiState() = %v, want %v", got.RiichiState(), player.NotRiichi)
			}
			if got.RiichiRiverIndex() != -1 {
				t.Errorf("NewVisiblePlayer().RiichiRiverIndex() = %v, want %v", got.RiichiRiverIndex(), -1)
			}
			if got.RiichiDiscardedTilesIndex() != -1 {
				t.Errorf("NewVisiblePlayer().RiichiDiscardedTilesIndex() = %v, want %v", got.RiichiDiscardedTilesIndex(), -1)
			}

			if got.CanDiscard() {
				t.Errorf("NewVisiblePlayer().CanDiscard() = %v, want %v", got.CanDiscard(), false)
			}
			if !got.IsConcealed() {
				t.Errorf("NewVisiblePlayer().IsConcealed() = %v, want %v", got.IsConcealed(), true)
			}
		})
	}
}

func TestVisiblePlayer_Draw_AddsTileToHand(t *testing.T) {
	handTiles := []tile.Tile{
		*tile.MustTileFromCode("C"), *tile.MustTileFromCode("9s"), *tile.MustTileFromCode("4m"),
		*tile.MustTileFromCode("2p"), *tile.MustTileFromCode("S"), *tile.MustTileFromCode("4p"),
		*tile.MustTileFromCode("8s"), *tile.MustTileFromCode("6p"), *tile.MustTileFromCode("6s"),
		*tile.MustTileFromCode("7m"), *tile.MustTileFromCode("9s"), *tile.MustTileFromCode("5pr"),
		*tile.MustTileFromCode("5p"),
	}
	p, err := player.NewVisiblePlayer(handTiles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	drawnTile := tile.MustTileFromCode("1m")
	if err := p.Draw(*drawnTile); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if h, _ := p.Hand(); !slices.Contains(h.ToTiles(), *drawnTile) {
		t.Errorf("expected tile %v to be in hand, but it was not", drawnTile)
	}
}
