package player_test

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestNewInvisiblePlayer(t *testing.T) {
	tests := []struct {
		name      string
		handTiles [13]tile.Tile
		wantErr   bool
	}{
		{
			name: "from visible tiles",
			handTiles: [13]tile.Tile{
				*tile.MustTileFromCode("C"), *tile.MustTileFromCode("9s"), *tile.MustTileFromCode("4m"),
				*tile.MustTileFromCode("2p"), *tile.MustTileFromCode("S"), *tile.MustTileFromCode("4p"),
				*tile.MustTileFromCode("8s"), *tile.MustTileFromCode("6p"), *tile.MustTileFromCode("6s"),
				*tile.MustTileFromCode("7m"), *tile.MustTileFromCode("9s"), *tile.MustTileFromCode("5pr"),
				*tile.MustTileFromCode("5p"),
			},
			wantErr: false,
		},
		{
			name: "from unknown tiles",
			handTiles: [13]tile.Tile{
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("?"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := player.NewInvisiblePlayer(tt.handTiles)
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
			if ok {
				t.Errorf("NewVisiblePlayer().Hand() returned ok")
			}
			if h != nil {
				t.Errorf("NewVisiblePlayer().Hand() = %v, want %v", h, nil)
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
			if !got.CanChiiPonKan() {
				t.Errorf("NewVisiblePlayer().CanChiiPonKan() = %v, want %v", got.CanChiiPonKan(), true)
			}
			if !got.IsConcealed() {
				t.Errorf("NewVisiblePlayer().IsConcealed() = %v, want %v", got.IsConcealed(), true)
			}
		})
	}
}

func TestInvisiblePlayer_Draw_Unknown(t *testing.T) {
	handTiles := [13]tile.Tile{
		*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
		*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
		*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
		*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
		*tile.MustTileFromCode("?"),
	}

	p, err := player.NewInvisiblePlayer(handTiles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	drawnTile := tile.MustTileFromCode("?")
	if err := p.Draw(*drawnTile); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if *p.DrawnTile() != *drawnTile {
		t.Errorf("DrawnTile() mismatch: expected %v but got %v", drawnTile, p.DrawnTile())
	}

	if !p.CanDiscard() {
		t.Errorf("player must be able to discard after Draw; CanDiscard() returned false")
	}
	if p.CanChiiPonKan() {
		t.Errorf("player must not be able to call after Draw; CanChiiPonKan() returned true")
	}
}

func TestInvisiblePlayer_Draw_Visible(t *testing.T) {
	handTiles := [13]tile.Tile{
		*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
		*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
		*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
		*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
		*tile.MustTileFromCode("?"),
	}

	p, err := player.NewInvisiblePlayer(handTiles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	drawnTile := tile.MustTileFromCode("1m")
	if err := p.Draw(*drawnTile); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if *p.DrawnTile() != *drawnTile {
		t.Errorf("DrawnTile() mismatch: expected %v but got %v", drawnTile, p.DrawnTile())
	}

	if !p.CanDiscard() {
		t.Errorf("player must be able to discard after Draw; CanDiscard() returned false")
	}
	if p.CanChiiPonKan() {
		t.Errorf("player must not be able to call after Draw; CanChiiPonKan() returned true")
	}
}

func TestInvisiblePlayer_Draw_CannotDrawTwice(t *testing.T) {
	handTiles := [13]tile.Tile{
		*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
		*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
		*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
		*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"),
		*tile.MustTileFromCode("?"),
	}

	p, err := player.NewInvisiblePlayer(handTiles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	firstTile := tile.MustTileFromCode("?")
	if err := p.Draw(*firstTile); err != nil {
		t.Fatalf("unexpected error on first Draw: %v", err)
	}

	secondTile := tile.MustTileFromCode("?")
	if err := p.Draw(*secondTile); err == nil {
		t.Errorf("Draw should fail when called twice without a discard; expected error but got nil")
	}
}
