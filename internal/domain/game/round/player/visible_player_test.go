package player_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/id"
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
				t.Errorf("NewVisiblePlayer().Hand() returned not ok")
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

func TestVisiblePlayer_Draw_Success(t *testing.T) {
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

	if h, _ := p.Hand(); *h != *hand.MustVisibleHand(handTiles) {
		t.Errorf("Hand() must remain unchanged after Draw(); got %+v", h)
	}

	sortedHandTiles := tile.Tiles(handTiles)
	sort.Sort(sortedHandTiles)
	if !reflect.DeepEqual(p.HandTiles(), []tile.Tile(sortedHandTiles)) {
		t.Errorf("HandTiles() must remain unchanged after Draw(); got %+v", p.HandTiles())
	}

	if *p.DrawnTile() != *drawnTile {
		t.Errorf("DrawnTile() mismatch: expected %v but got %v", drawnTile, p.DrawnTile())
	}

	if !p.CanDiscard() {
		t.Errorf("player must be able to discard after Draw; CanDiscard() returned false")
	}
}

func TestVisiblePlayer_Draw_CannotDrawUnknown(t *testing.T) {
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

	drawnTile := tile.MustTileFromCode("?")
	if err := p.Draw(*drawnTile); err == nil {
		t.Errorf("Draw() succeeded unexpectedly")
	}
}

func TestVisiblePlayer_Draw_CannotDrawTwice(t *testing.T) {
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

	firstTile := tile.MustTileFromCode("1m")
	if err := p.Draw(*firstTile); err != nil {
		t.Fatalf("unexpected error on first Draw: %v", err)
	}

	secondTile := tile.MustTileFromCode("2m")
	if err := p.Draw(*secondTile); err == nil {
		t.Errorf("Draw should fail when called twice without a discard; expected error but got nil")
	}
}

func TestVisiblePlayer_Draw_CannotDrawBeforeRiichiAccepted(t *testing.T) {
	handTiles := []tile.Tile{
		*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m"), *tile.MustTileFromCode("3m"),
		*tile.MustTileFromCode("4p"), *tile.MustTileFromCode("5p"), *tile.MustTileFromCode("6p"),
		*tile.MustTileFromCode("7s"), *tile.MustTileFromCode("8s"), *tile.MustTileFromCode("9s"),
		*tile.MustTileFromCode("E"), *tile.MustTileFromCode("E"), *tile.MustTileFromCode("S"),
		*tile.MustTileFromCode("W"),
	}

	p, err := player.NewVisiblePlayer(handTiles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	firstTile := tile.MustTileFromCode("S")
	if err := p.Draw(*firstTile); err != nil {
		t.Fatalf("unexpected error on first Draw: %v", err)
	}

	if err := p.Riichi(); err != nil {
		t.Fatalf("unexpected error on Riichi: %v", err)
	}

	discardedTile := tile.MustTileFromCode("W")
	if err := p.Discard(*discardedTile, false); err != nil {
		t.Fatalf("unexpected error on Discard: %v", err)
	}

	secondTile := tile.MustTileFromCode("2m")
	if err := p.Draw(*secondTile); err == nil {
		t.Errorf("cannot Draw: riichi has been declared but the discard has not yet been accepted")
	}
}

func TestVisiblePlayer_Discard_TileInHand(t *testing.T) {
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
		t.Fatalf("unexpected error on Draw: %v", err)
	}

	discardedTile := tile.MustTileFromCode("C")
	if err := p.Discard(*discardedTile, false); err != nil {
		t.Errorf("unexpected error on Discard: %v", err)
	}

	afterHandTiles := []tile.Tile{
		*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("4m"), *tile.MustTileFromCode("7m"),
		*tile.MustTileFromCode("2p"), *tile.MustTileFromCode("4p"), *tile.MustTileFromCode("5p"),
		*tile.MustTileFromCode("5pr"), *tile.MustTileFromCode("6p"), *tile.MustTileFromCode("6s"),
		*tile.MustTileFromCode("8s"), *tile.MustTileFromCode("9s"), *tile.MustTileFromCode("9s"),
		*tile.MustTileFromCode("S"),
	}
	h := hand.MustVisibleHand(afterHandTiles)

	if gotHand, _ := p.Hand(); *gotHand != *h {
		t.Errorf("Hand() mismatch after Discard: got %v, want %v", gotHand, h)
	}

	if !reflect.DeepEqual(p.HandTiles(), []tile.Tile(afterHandTiles)) {
		t.Errorf("HandTiles() mismatch after Discard: got %v, want %v", p.HandTiles(), afterHandTiles)
	}

	if p.DrawnTile() != nil {
		t.Errorf("DrawnTile() should be nil after Discard; got %v", p.DrawnTile())
	}

	river := []tile.Tile{*discardedTile}

	if !reflect.DeepEqual(p.River(), river) {
		t.Errorf("River() mismatch after Discard: got %v, want %v", p.River(), river)
	}

	if !reflect.DeepEqual(p.DiscardedTiles(), river) {
		t.Errorf("DiscardedTiles() mismatch after Discard: got %v, want %v", p.DiscardedTiles(), river)
	}

	if p.CanDiscard() {
		t.Errorf("CanDiscard() should be false after Discard; got true")
	}
}

func TestVisiblePlayer_Discard_DrawnTile(t *testing.T) {
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
		t.Fatalf("unexpected error on Draw: %v", err)
	}

	discardedTile := tile.MustTileFromCode("1m")
	if err := p.Discard(*discardedTile, true); err != nil {
		t.Errorf("unexpected error on Discard: %v", err)
	}

	afterHandTiles := tile.Tiles(handTiles)
	sort.Sort(afterHandTiles)
	h := hand.MustVisibleHand(afterHandTiles)

	if gotHand, _ := p.Hand(); *gotHand != *h {
		t.Errorf("Hand() mismatch after Discard: got %v, want %v", gotHand, h)
	}

	if !reflect.DeepEqual(p.HandTiles(), []tile.Tile(afterHandTiles)) {
		t.Errorf("HandTiles() mismatch after Discard: got %v, want %v", p.HandTiles(), afterHandTiles)
	}

	if p.DrawnTile() != nil {
		t.Errorf("DrawnTile() should be nil after Discard; got %v", p.DrawnTile())
	}

	river := []tile.Tile{*discardedTile}

	if !reflect.DeepEqual(p.River(), river) {
		t.Errorf("River() mismatch after Discard: got %v, want %v", p.River(), river)
	}

	if !reflect.DeepEqual(p.DiscardedTiles(), river) {
		t.Errorf("DiscardedTiles() mismatch after Discard: got %v, want %v", p.DiscardedTiles(), river)
	}

	if p.CanDiscard() {
		t.Errorf("CanDiscard() should be false after Discard; got true")
	}
}

func TestVisiblePlayer_Discard_ClearExtraSafeTiles(t *testing.T) {
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

	p.AddExtraSafeTiles(*tile.MustTileFromCode("1m"))
	p.AddExtraSafeTiles(*tile.MustTileFromCode("2m"))
	p.AddExtraSafeTiles(*tile.MustTileFromCode("3m"))
	extraSafeTiles := []tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m"), *tile.MustTileFromCode("3m")}
	if !reflect.DeepEqual(p.ExtraSafeTiles(), extraSafeTiles) {
		t.Fatalf("ExtraSafeTiles() failed")
	}

	drawnTile := tile.MustTileFromCode("1m")
	if err := p.Draw(*drawnTile); err != nil {
		t.Fatalf("unexpected error on Draw: %v", err)
	}

	discardedTile := tile.MustTileFromCode("C")
	if err := p.Discard(*discardedTile, false); err != nil {
		t.Fatalf("unexpected error on Discard: %v", err)
	}

	got := p.ExtraSafeTiles()
	want := []tile.Tile{}
	if !reflect.DeepEqual(p.ExtraSafeTiles(), want) {
		t.Errorf("ExtraSafeTiles() = %v, want %v", got, want)
	}
}

func TestVisiblePlayer_Discard_NotClearExtraSafeTilesAfterRiichi(t *testing.T) {
	handTiles := []tile.Tile{
		*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m"), *tile.MustTileFromCode("3m"),
		*tile.MustTileFromCode("4p"), *tile.MustTileFromCode("5p"), *tile.MustTileFromCode("6p"),
		*tile.MustTileFromCode("7s"), *tile.MustTileFromCode("8s"), *tile.MustTileFromCode("9s"),
		*tile.MustTileFromCode("E"), *tile.MustTileFromCode("E"), *tile.MustTileFromCode("S"),
		*tile.MustTileFromCode("W"),
	}

	p, err := player.NewVisiblePlayer(handTiles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	drawnTile := tile.MustTileFromCode("S")
	if err := p.Draw(*drawnTile); err != nil {
		t.Fatalf("unexpected error on Draw: %v", err)
	}

	if err := p.Riichi(); err != nil {
		t.Fatalf("unexpected error on Riichi: %v", err)
	}

	discardedTile := tile.MustTileFromCode("W")
	if err := p.Discard(*discardedTile, false); err != nil {
		t.Fatalf("unexpected error on Discard: %v", err)
	}

	if err := p.RiichiAccepted(); err != nil {
		t.Fatalf("unexpected error on RiichiAccepted: %v", err)
	}

	p.AddExtraSafeTiles(*tile.MustTileFromCode("1m"))
	p.AddExtraSafeTiles(*tile.MustTileFromCode("2m"))
	p.AddExtraSafeTiles(*tile.MustTileFromCode("3m"))
	extraSafeTiles := []tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m"), *tile.MustTileFromCode("3m")}
	if !reflect.DeepEqual(p.ExtraSafeTiles(), extraSafeTiles) {
		t.Fatalf("ExtraSafeTiles() failed")
	}

	drawnTile2 := tile.MustTileFromCode("1s")
	if err := p.Draw(*drawnTile2); err != nil {
		t.Fatalf("unexpected error on Draw: %v", err)
	}

	if err := p.Discard(*drawnTile2, true); err != nil {
		t.Fatalf("unexpected error on Discard: %v", err)
	}

	got := p.ExtraSafeTiles()
	if !reflect.DeepEqual(p.ExtraSafeTiles(), extraSafeTiles) {
		t.Errorf("ExtraSafeTiles() = %v, want %v", got, extraSafeTiles)
	}
}

func TestVisiblePlayer_Discard_CannotDiscardBeforeDraw(t *testing.T) {
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

	discardedTile := tile.MustTileFromCode("C")
	if err := p.Discard(*discardedTile, false); err == nil {
		t.Errorf("Discard should fail before any Draw; expected error but got nil")
	}
}

func TestVisiblePlayer_Discard_CannotDiscardUnknown(t *testing.T) {
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
		t.Fatalf("unexpected error on Draw: %v", err)
	}

	discardedTile := tile.MustTileFromCode("?")
	if err := p.Discard(*discardedTile, false); err == nil {
		t.Errorf("Discard should fail to discard an unknown tile; expected error but got nil")
	}
}

func TestVisiblePlayer_Discard_CannotDiscardTileNotInHand(t *testing.T) {
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
		t.Fatalf("unexpected error on Draw: %v", err)
	}

	discardedTile := tile.MustTileFromCode("1m")
	if err := p.Discard(*discardedTile, false); err == nil {
		t.Errorf("Discard should fail when tile is not in hand; expected error but got nil")
	}
}

func TestVisiblePlayer_Discard_CannotDiscardTileNotDrawnTile(t *testing.T) {
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
		t.Fatalf("unexpected error on Draw: %v", err)
	}

	discardedTile := tile.MustTileFromCode("C")
	if err := p.Discard(*discardedTile, true); err == nil {
		t.Errorf("Discard should fail when tsumogiri=true but a hand tile was specified; expected error but got nil")
	}
}

func TestVisiblePlayer_Discard_CannotDiscardNotRiichiDeclarationTile(t *testing.T) {
	handTiles := []tile.Tile{
		*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m"), *tile.MustTileFromCode("3m"),
		*tile.MustTileFromCode("4p"), *tile.MustTileFromCode("5p"), *tile.MustTileFromCode("6p"),
		*tile.MustTileFromCode("7s"), *tile.MustTileFromCode("8s"), *tile.MustTileFromCode("9s"),
		*tile.MustTileFromCode("E"), *tile.MustTileFromCode("E"), *tile.MustTileFromCode("S"),
		*tile.MustTileFromCode("W"),
	}

	p, err := player.NewVisiblePlayer(handTiles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	drawnTile := tile.MustTileFromCode("S")
	if err := p.Draw(*drawnTile); err != nil {
		t.Fatalf("unexpected error on Draw: %v", err)
	}

	if err := p.Riichi(); err != nil {
		t.Fatalf("unexpected error on Riichi: %v", err)
	}

	discardedTile := tile.MustTileFromCode("1m")
	if err := p.Discard(*discardedTile, false); err == nil {
		t.Errorf("Discard should fail when the player is in riichi and attempted to discard a tile other than the riichi declaration tile; expected error but got nil")
	}
}

func TestVisiblePlayer_Discard_CannotDiscardFromHandAfterRiichi(t *testing.T) {
	handTiles := []tile.Tile{
		*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m"), *tile.MustTileFromCode("3m"),
		*tile.MustTileFromCode("4p"), *tile.MustTileFromCode("5p"), *tile.MustTileFromCode("6p"),
		*tile.MustTileFromCode("7s"), *tile.MustTileFromCode("8s"), *tile.MustTileFromCode("9s"),
		*tile.MustTileFromCode("E"), *tile.MustTileFromCode("E"), *tile.MustTileFromCode("S"),
		*tile.MustTileFromCode("W"),
	}

	p, err := player.NewVisiblePlayer(handTiles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	drawnTile := tile.MustTileFromCode("S")
	if err := p.Draw(*drawnTile); err != nil {
		t.Fatalf("unexpected error on Draw: %v", err)
	}

	if err := p.Riichi(); err != nil {
		t.Fatalf("unexpected error on Riichi: %v", err)
	}

	discardedTile := tile.MustTileFromCode("W")
	if err := p.Discard(*discardedTile, false); err != nil {
		t.Fatalf("unexpected error on Discard: %v", err)
	}

	if err := p.RiichiAccepted(); err != nil {
		t.Fatalf("unexpected error on RiichiAccepted: %v", err)
	}

	drawnTile2 := tile.MustTileFromCode("1m")
	if err := p.Draw(*drawnTile2); err != nil {
		t.Fatalf("unexpected error on Draw: %v", err)
	}

	discardedTile2 := tile.MustTileFromCode("1m")
	if err := p.Discard(*discardedTile2, false); err == nil {
		t.Errorf("Discard should fail when the player accepted riichi discarded a tile in the hand")
	}
}

func TestVisiblePlayer_Pon_Success(t *testing.T) {
	handTiles := []tile.Tile{
		*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m"), *tile.MustTileFromCode("3m"),
		*tile.MustTileFromCode("4p"), *tile.MustTileFromCode("5p"), *tile.MustTileFromCode("6p"),
		*tile.MustTileFromCode("7s"), *tile.MustTileFromCode("8s"), *tile.MustTileFromCode("9s"),
		*tile.MustTileFromCode("E"), *tile.MustTileFromCode("E"), *tile.MustTileFromCode("S"),
		*tile.MustTileFromCode("W"),
	}

	p, err := player.NewVisiblePlayer(handTiles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	pon := meld.MustPon(
		*tile.MustTileFromCode("E"),
		[2]tile.Tile{*tile.MustTileFromCode("E"), *tile.MustTileFromCode("E")},
		*id.MustID(0),
	)
	if err := p.Pon(*pon); err != nil {
		t.Errorf("Pon() failed: %v", err)
	}

	wantHand := hand.CodesToHand([]string{"1m", "2m", "3m", "4p", "5p", "6p", "7s", "8s", "9s", "S", "W"})
	if h, _ := p.Hand(); *h != *wantHand {
		t.Errorf("Hand() = %v, want %v", h, wantHand)
	}

	wantMelds := []meld.Meld{pon}
	if !reflect.DeepEqual(p.Melds(), wantMelds) {
		t.Errorf("Hand() = %v, want %v", p.Melds(), wantMelds)
	}

	if !p.CanDiscard() {
		t.Errorf("CanDiscard() = %v, want %v", p.CanDiscard(), true)
	}
	if p.IsConcealed() {
		t.Errorf("IsConcealed() = %v, want %v", p.IsConcealed(), false)
	}
}

func TestVisiblePlayer_Pon_CannotAfterDraw(t *testing.T) {
	handTiles := []tile.Tile{
		*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m"), *tile.MustTileFromCode("3m"),
		*tile.MustTileFromCode("4p"), *tile.MustTileFromCode("5p"), *tile.MustTileFromCode("6p"),
		*tile.MustTileFromCode("7s"), *tile.MustTileFromCode("8s"), *tile.MustTileFromCode("9s"),
		*tile.MustTileFromCode("E"), *tile.MustTileFromCode("E"), *tile.MustTileFromCode("S"),
		*tile.MustTileFromCode("W"),
	}

	p, err := player.NewVisiblePlayer(handTiles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	drawnTile := tile.MustTileFromCode("S")
	if err := p.Draw(*drawnTile); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	pon := meld.MustPon(
		*tile.MustTileFromCode("E"),
		[2]tile.Tile{*tile.MustTileFromCode("E"), *tile.MustTileFromCode("E")},
		*id.MustID(0),
	)
	if err := p.Pon(*pon); err == nil {
		t.Errorf("Pon should fail when the player has a drawn tile")
	}
}

func TestVisiblePlayer_Pon_CannotAfterRiichi(t *testing.T) {
	handTiles := []tile.Tile{
		*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m"), *tile.MustTileFromCode("3m"),
		*tile.MustTileFromCode("4p"), *tile.MustTileFromCode("5p"), *tile.MustTileFromCode("6p"),
		*tile.MustTileFromCode("7s"), *tile.MustTileFromCode("8s"), *tile.MustTileFromCode("9s"),
		*tile.MustTileFromCode("E"), *tile.MustTileFromCode("E"), *tile.MustTileFromCode("S"),
		*tile.MustTileFromCode("W"),
	}

	p, err := player.NewVisiblePlayer(handTiles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	drawnTile := tile.MustTileFromCode("S")
	if err := p.Draw(*drawnTile); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := p.Riichi(); err != nil {
		t.Fatalf("unexpected error on Riichi: %v", err)
	}

	discardedTile := tile.MustTileFromCode("W")
	if err := p.Discard(*discardedTile, false); err != nil {
		t.Fatalf("unexpected error on Discard: %v", err)
	}

	if err := p.RiichiAccepted(); err != nil {
		t.Fatalf("unexpected error on RiichiAccepted: %v", err)
	}

	pon := meld.MustPon(
		*tile.MustTileFromCode("E"),
		[2]tile.Tile{*tile.MustTileFromCode("E"), *tile.MustTileFromCode("E")},
		*id.MustID(0),
	)
	if err := p.Pon(*pon); err == nil {
		t.Errorf("Pon should fail when the player is already in riichi state")
	}
}

func TestVisiblePlayer_Pon_Cannot5thCall(t *testing.T) {
	handTiles := []tile.Tile{
		*tile.MustTileFromCode("E"), *tile.MustTileFromCode("E"),
		*tile.MustTileFromCode("S"), *tile.MustTileFromCode("S"),
		*tile.MustTileFromCode("W"), *tile.MustTileFromCode("W"),
		*tile.MustTileFromCode("N"), *tile.MustTileFromCode("N"),
		*tile.MustTileFromCode("P"), *tile.MustTileFromCode("P"),
		*tile.MustTileFromCode("P"), *tile.MustTileFromCode("P"),
		*tile.MustTileFromCode("C"),
	}

	p, err := player.NewVisiblePlayer(handTiles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, wind := range []string{"E", "S", "W", "N"} {
		pon := meld.MustPon(
			*tile.MustTileFromCode(wind),
			[2]tile.Tile{*tile.MustTileFromCode(wind), *tile.MustTileFromCode(wind)},
			*id.MustID(0),
		)
		if err := p.Pon(*pon); err != nil {
			t.Fatalf("unexpected error on Pon: %v", err)
		}

		discardedTile := tile.MustTileFromCode("P")
		if err := p.Discard(*discardedTile, false); err != nil {
			t.Fatalf("unexpected error on Discard: %v", err)
		}
	}

	pon := meld.MustPon(
		*tile.MustTileFromCode("C"),
		[2]tile.Tile{*tile.MustTileFromCode("C"), *tile.MustTileFromCode("C")},
		*id.MustID(0),
	)
	if err := p.Pon(*pon); err == nil {
		t.Errorf("Pon should fail when the player has four melds")
	}
}

func TestVisiblePlayer_Riichi_Success(t *testing.T) {
	handTiles := []tile.Tile{
		*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m"), *tile.MustTileFromCode("3m"),
		*tile.MustTileFromCode("4p"), *tile.MustTileFromCode("5p"), *tile.MustTileFromCode("6p"),
		*tile.MustTileFromCode("7s"), *tile.MustTileFromCode("8s"), *tile.MustTileFromCode("9s"),
		*tile.MustTileFromCode("E"), *tile.MustTileFromCode("E"), *tile.MustTileFromCode("S"),
		*tile.MustTileFromCode("W"),
	}

	p, err := player.NewVisiblePlayer(handTiles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	drawnTile := tile.MustTileFromCode("S")
	if err := p.Draw(*drawnTile); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := p.Riichi(); err != nil {
		t.Errorf("Riichi() failed: %v", err)
	}

	if p.RiichiState() != player.RiichiDeclared {
		t.Errorf("RiichiState() = %v, want %v", p.RiichiState(), player.RiichiDeclared)
	}
	if p.RiichiRiverIndex() != -1 {
		t.Errorf("RiichiRiverIndex() = %v, want %v", p.RiichiRiverIndex(), -1)
	}
	if p.RiichiDiscardedTilesIndex() != -1 {
		t.Errorf("RiichiDiscardedTilesIndex() = %v, want %v", p.RiichiDiscardedTilesIndex(), -1)
	}

	if !p.CanDiscard() {
		t.Errorf("player must be able to discard after Riichi; CanDiscard() returned false")
	}
}

func TestVisiblePlayer_Riichi_CannotDeclareTwice(t *testing.T) {
	handTiles := []tile.Tile{
		*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m"), *tile.MustTileFromCode("3m"),
		*tile.MustTileFromCode("4p"), *tile.MustTileFromCode("5p"), *tile.MustTileFromCode("6p"),
		*tile.MustTileFromCode("7s"), *tile.MustTileFromCode("8s"), *tile.MustTileFromCode("9s"),
		*tile.MustTileFromCode("E"), *tile.MustTileFromCode("E"), *tile.MustTileFromCode("S"),
		*tile.MustTileFromCode("W"),
	}

	p, err := player.NewVisiblePlayer(handTiles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	drawnTile := tile.MustTileFromCode("S")
	if err := p.Draw(*drawnTile); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := p.Riichi(); err != nil {
		t.Fatalf("unexpected error on first Riichi: %v", err)
	}

	if err := p.Riichi(); err == nil {
		t.Errorf("Riichi should fail when called twice; expected error but got nil")
	}
}

func TestVisiblePlayer_Riichi_CannotDeclareBeforeDraw(t *testing.T) {
	handTiles := []tile.Tile{
		*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m"), *tile.MustTileFromCode("3m"),
		*tile.MustTileFromCode("4p"), *tile.MustTileFromCode("5p"), *tile.MustTileFromCode("6p"),
		*tile.MustTileFromCode("7s"), *tile.MustTileFromCode("8s"), *tile.MustTileFromCode("9s"),
		*tile.MustTileFromCode("E"), *tile.MustTileFromCode("E"), *tile.MustTileFromCode("S"),
		*tile.MustTileFromCode("W"),
	}

	p, err := player.NewVisiblePlayer(handTiles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := p.Riichi(); err == nil {
		t.Errorf("Riichi should fail when called before Draw; expected error but got nil")
	}
}

func TestVisiblePlayer_Riichi_CannotDeclareWithOpenMeld(t *testing.T) {
	handTiles := []tile.Tile{
		*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m"), *tile.MustTileFromCode("3m"),
		*tile.MustTileFromCode("4p"), *tile.MustTileFromCode("5p"), *tile.MustTileFromCode("6p"),
		*tile.MustTileFromCode("7s"), *tile.MustTileFromCode("8s"), *tile.MustTileFromCode("9s"),
		*tile.MustTileFromCode("E"), *tile.MustTileFromCode("E"), *tile.MustTileFromCode("S"),
		*tile.MustTileFromCode("W"),
	}

	p, err := player.NewVisiblePlayer(handTiles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	pon := meld.MustPon(
		*tile.MustTileFromCode("E"),
		[2]tile.Tile{*tile.MustTileFromCode("E"), *tile.MustTileFromCode("E")},
		*id.MustID(0),
	)
	if err := p.Pon(*pon); err != nil {
		t.Fatalf("unexpected error on Pon: %v", err)
	}

	discardedTile := tile.MustTileFromCode("W")
	if err := p.Discard(*discardedTile, false); err != nil {
		t.Fatalf("unexpected error on Discard: %v", err)
	}

	drawnTile := tile.MustTileFromCode("W")
	if err := p.Draw(*drawnTile); err != nil {
		t.Fatalf("unexpected error on Draw: %v", err)
	}

	if err := p.Riichi(); err == nil {
		t.Errorf("Riichi should fail when called with open melds; expected error but got nil")
	}
}

func TestVisiblePlayer_Riichi_CannotDeclareNotTenpai(t *testing.T) {
	handTiles := []tile.Tile{
		*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m"), *tile.MustTileFromCode("3m"),
		*tile.MustTileFromCode("4p"), *tile.MustTileFromCode("5p"), *tile.MustTileFromCode("6p"),
		*tile.MustTileFromCode("7s"), *tile.MustTileFromCode("8s"), *tile.MustTileFromCode("9s"),
		*tile.MustTileFromCode("E"), *tile.MustTileFromCode("E"), *tile.MustTileFromCode("S"),
		*tile.MustTileFromCode("W"),
	}

	p, err := player.NewVisiblePlayer(handTiles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	drawnTile := tile.MustTileFromCode("N")
	if err := p.Draw(*drawnTile); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := p.Riichi(); err == nil {
		t.Errorf("Riichi should fail when the player is not tenpai; expected error but got nil")
	}
}

func TestVisiblePlayer_RiichiAccepted_Success(t *testing.T) {
	handTiles := []tile.Tile{
		*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m"), *tile.MustTileFromCode("3m"),
		*tile.MustTileFromCode("4p"), *tile.MustTileFromCode("5p"), *tile.MustTileFromCode("6p"),
		*tile.MustTileFromCode("7s"), *tile.MustTileFromCode("8s"), *tile.MustTileFromCode("9s"),
		*tile.MustTileFromCode("E"), *tile.MustTileFromCode("E"), *tile.MustTileFromCode("S"),
		*tile.MustTileFromCode("W"),
	}

	p, err := player.NewVisiblePlayer(handTiles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	firstTile := tile.MustTileFromCode("S")
	if err := p.Draw(*firstTile); err != nil {
		t.Fatalf("unexpected error on first Draw: %v", err)
	}

	if err := p.Riichi(); err != nil {
		t.Fatalf("unexpected error on Riichi: %v", err)
	}

	discardedTile := tile.MustTileFromCode("W")
	if err := p.Discard(*discardedTile, false); err != nil {
		t.Fatalf("unexpected error on Discard: %v", err)
	}

	if err := p.RiichiAccepted(); err != nil {
		t.Errorf("RiichiAccepted() failed: %v", err)
	}

	if p.RiichiState() != player.RiichiAccepted {
		t.Errorf("RiichiState() = %v, want %v", p.RiichiState(), player.RiichiAccepted)
	}
	if p.RiichiRiverIndex() != 0 {
		t.Errorf("RiichiRiverIndex() = %v, want %v", p.RiichiRiverIndex(), 0)
	}
	if p.RiichiDiscardedTilesIndex() != 0 {
		t.Errorf("RiichiDiscardedTilesIndex() = %v, want %v", p.RiichiDiscardedTilesIndex(), 0)
	}

	secondTile := tile.MustTileFromCode("2m")
	if err := p.Draw(*secondTile); err != nil {
		t.Errorf("Draw should succeed after riichi has been accepted; got error: %v", err)
	}
}

func TestVisiblePlayer_RiichiAccepted_CannotAcceptTwice(t *testing.T) {
	handTiles := []tile.Tile{
		*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m"), *tile.MustTileFromCode("3m"),
		*tile.MustTileFromCode("4p"), *tile.MustTileFromCode("5p"), *tile.MustTileFromCode("6p"),
		*tile.MustTileFromCode("7s"), *tile.MustTileFromCode("8s"), *tile.MustTileFromCode("9s"),
		*tile.MustTileFromCode("E"), *tile.MustTileFromCode("E"), *tile.MustTileFromCode("S"),
		*tile.MustTileFromCode("W"),
	}

	p, err := player.NewVisiblePlayer(handTiles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	firstTile := tile.MustTileFromCode("S")
	if err := p.Draw(*firstTile); err != nil {
		t.Fatalf("unexpected error on first Draw: %v", err)
	}

	if err := p.Riichi(); err != nil {
		t.Fatalf("unexpected error on Riichi: %v", err)
	}

	discardedTile := tile.MustTileFromCode("W")
	if err := p.Discard(*discardedTile, false); err != nil {
		t.Fatalf("unexpected error on Discard: %v", err)
	}

	if err := p.RiichiAccepted(); err != nil {
		t.Fatalf("unexpected error on RiichiAccepted: %v", err)
	}

	if err := p.RiichiAccepted(); err == nil {
		t.Errorf("RiichiAccepted should fail when called twice; expected error but got nil")
	}
}

func TestVisiblePlayer_RiichiAccepted_NotRiichi(t *testing.T) {
	handTiles := []tile.Tile{
		*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m"), *tile.MustTileFromCode("3m"),
		*tile.MustTileFromCode("4p"), *tile.MustTileFromCode("5p"), *tile.MustTileFromCode("6p"),
		*tile.MustTileFromCode("7s"), *tile.MustTileFromCode("8s"), *tile.MustTileFromCode("9s"),
		*tile.MustTileFromCode("E"), *tile.MustTileFromCode("E"), *tile.MustTileFromCode("S"),
		*tile.MustTileFromCode("W"),
	}

	p, err := player.NewVisiblePlayer(handTiles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	firstTile := tile.MustTileFromCode("S")
	if err := p.Draw(*firstTile); err != nil {
		t.Fatalf("unexpected error on first Draw: %v", err)
	}

	discardedTile := tile.MustTileFromCode("W")
	if err := p.Discard(*discardedTile, false); err != nil {
		t.Fatalf("unexpected error on Discard: %v", err)
	}

	if err := p.RiichiAccepted(); err == nil {
		t.Errorf("RiichiAccepted should fail when called before Riichi; expected error but got nil")
	}
}

func TestVisiblePlayer_AddExtraSafeTiles(t *testing.T) {
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

	p.AddExtraSafeTiles(*tile.MustTileFromCode("5s"))
	got := p.ExtraSafeTiles()
	want := []tile.Tile{*tile.MustTileFromCode("5s")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ExtraSafeTiles() = %v, want %v", got, want)
	}

	p.AddExtraSafeTiles(*tile.MustTileFromCode("5sr"))
	got = p.ExtraSafeTiles()
	want = []tile.Tile{*tile.MustTileFromCode("5s"), *tile.MustTileFromCode("5sr")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ExtraSafeTiles() = %v, want %v", got, want)
	}
}

func TestVisiblePlayer_AddExtraSafeTiles_Panic(t *testing.T) {
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

	t.Run("unknown tile", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for unknown tile, but did not panic")
			}
		}()

		p.AddExtraSafeTiles(*tile.MustTileFromCode("?"))
	})
}

func TestVisiblePlayer_TakeFromRiver(t *testing.T) {
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

	drawnTile := tile.MustTileFromCode("P")
	if err := p.Draw(*drawnTile); err != nil {
		t.Fatalf("unexpected error on Draw: %v", err)
	}

	discardedTile := tile.MustTileFromCode("2p")
	if err := p.Discard(*discardedTile, false); err != nil {
		t.Fatalf("unexpected error on Discard: %v", err)
	}

	if err := p.TakeFromRiver(*discardedTile); err != nil {
		t.Errorf("TakeFromRiver() failed: %v", err)
	}

	gotRiver := p.River()
	wantRiver := []tile.Tile{}
	if !reflect.DeepEqual(gotRiver, wantRiver) {
		t.Errorf("River() mismatch after TakeFromRiver: got %v, want %v", gotRiver, wantRiver)
	}

	gotDiscardedTiles := p.DiscardedTiles()
	wantDiscardedTiles := []tile.Tile{*discardedTile}
	if !reflect.DeepEqual(gotDiscardedTiles, wantDiscardedTiles) {
		t.Errorf("DiscardedTiles() mismatch after TakeFromRiver: got %v, want %v", gotDiscardedTiles, wantDiscardedTiles)
	}
}

func TestVisiblePlayer_TakeFromRiver_Mismatch(t *testing.T) {
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

	drawnTile := tile.MustTileFromCode("P")
	if err := p.Draw(*drawnTile); err != nil {
		t.Fatalf("unexpected error on Draw: %v", err)
	}

	discardedTile := tile.MustTileFromCode("2p")
	if err := p.Discard(*discardedTile, false); err != nil {
		t.Fatalf("unexpected error on Discard: %v", err)
	}

	taken := tile.MustTileFromCode("3p")
	if err := p.TakeFromRiver(*taken); err == nil {
		t.Errorf("TakeFromRiver() succeeded unexpectedly")
	}
}

func TestVisiblePlayer_TakeFromRiver_Unknown(t *testing.T) {
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

	drawnTile := tile.MustTileFromCode("P")
	if err := p.Draw(*drawnTile); err != nil {
		t.Fatalf("unexpected error on Draw: %v", err)
	}

	discardedTile := tile.MustTileFromCode("2p")
	if err := p.Discard(*discardedTile, false); err != nil {
		t.Fatalf("unexpected error on Discard: %v", err)
	}

	taken := tile.MustTileFromCode("?")
	if err := p.TakeFromRiver(*taken); err == nil {
		t.Errorf("TakeFromRiver() succeeded unexpectedly")
	}
}
