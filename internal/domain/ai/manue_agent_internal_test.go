package ai

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestManueAgent_selectAction_Priorities(t *testing.T) {
	self := seat.MustSeat(0)
	target := seat.MustSeat(1)
	drawnTile := tile.MustTileFromCode("5p")
	handDiscard, err := action.NewDiscard(self, tile.MustTileFromCode("1m"), false)
	if err != nil {
		t.Fatalf("NewDiscard(hand) failed: %v", err)
	}
	tsumogiriDiscard, err := action.NewDiscard(self, drawnTile, true)
	if err != nil {
		t.Fatalf("NewDiscard(tsumogiri) failed: %v", err)
	}
	win, err := action.NewWin(self, self, drawnTile)
	if err != nil {
		t.Fatalf("NewWin() failed: %v", err)
	}
	riichi := action.NewRiichi(self)
	chii, err := action.NewChii(self, target, tile.MustTileFromCode("3m"), [2]tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("2m"),
	})
	if err != nil {
		t.Fatalf("NewChii() failed: %v", err)
	}
	pass := action.NewPass(self)

	tests := []struct {
		name        string
		actions     []action.Action
		riichiState player.RiichiState
		drawnTile   *tile.Tile
		want        action.Action
	}{
		{
			name:        "win first",
			actions:     []action.Action{handDiscard, win, pass},
			riichiState: player.NotRiichi,
			drawnTile:   &drawnTile,
			want:        win,
		},
		{
			name:        "riichi accepted tsumogiri",
			actions:     []action.Action{handDiscard, tsumogiriDiscard},
			riichiState: player.RiichiAccepted,
			drawnTile:   &drawnTile,
			want:        tsumogiriDiscard,
		},
		{
			name:        "riichi before discard",
			actions:     []action.Action{handDiscard, riichi},
			riichiState: player.NotRiichi,
			drawnTile:   &drawnTile,
			want:        riichi,
		},
		{
			name:        "discard before call",
			actions:     []action.Action{chii, handDiscard, pass},
			riichiState: player.NotRiichi,
			drawnTile:   &drawnTile,
			want:        handDiscard,
		},
		{
			name:        "call before pass",
			actions:     []action.Action{pass, chii},
			riichiState: player.NotRiichi,
			drawnTile:   nil,
			want:        chii,
		},
		{
			name:        "pass only",
			actions:     []action.Action{pass},
			riichiState: player.NotRiichi,
			drawnTile:   nil,
			want:        pass,
		},
	}

	agent := NewManueAgent(0)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := agent.selectAction(tt.actions, stubPlayerViewer{
				riichiState: tt.riichiState,
				drawnTile:   tt.drawnTile,
			})
			if err != nil {
				t.Fatalf("selectAction() failed: %v", err)
			}
			if got != tt.want {
				t.Errorf("selectAction() = %T %[1]v, want %T %[2]v", got, tt.want)
			}
		})
	}
}

func TestManueAgent_selectAction_ReturnsErrorWithoutTsumogiriAfterRiichiAccepted(t *testing.T) {
	self := seat.MustSeat(0)
	handDiscard, err := action.NewDiscard(self, tile.MustTileFromCode("1m"), false)
	if err != nil {
		t.Fatalf("NewDiscard(hand) failed: %v", err)
	}
	drawnTile := tile.MustTileFromCode("5p")

	_, err = NewManueAgent(0).selectAction([]action.Action{handDiscard}, stubPlayerViewer{
		riichiState: player.RiichiAccepted,
		drawnTile:   &drawnTile,
	})
	if err == nil {
		t.Fatal("selectAction() succeeded unexpectedly")
	}
}

type stubPlayerViewer struct {
	riichiState player.RiichiState
	drawnTile   *tile.Tile
}

func (p stubPlayerViewer) Hand() (*hand.VisibleHand, bool) { return nil, false }
func (p stubPlayerViewer) HandTiles() []tile.Tile          { return nil }
func (p stubPlayerViewer) DrawnTile() *tile.Tile           { return p.drawnTile }
func (p stubPlayerViewer) Melds() []meld.Meld              { return nil }
func (p stubPlayerViewer) River() []tile.Tile              { return nil }
func (p stubPlayerViewer) DiscardedTiles() []tile.Tile     { return nil }
func (p stubPlayerViewer) ExtraSafeTiles() []tile.Tile     { return nil }
func (p stubPlayerViewer) IsFuriten() bool                 { return false }
func (p stubPlayerViewer) CanRonBy(*tile.Tile) bool        { return true }
func (p stubPlayerViewer) RiichiState() player.RiichiState { return p.riichiState }
func (p stubPlayerViewer) RiichiRiverIndex() int           { return -1 }
func (p stubPlayerViewer) RiichiDiscardedTilesIndex() int  { return -1 }
func (p stubPlayerViewer) CanDiscard() bool                { return p.drawnTile != nil }
func (p stubPlayerViewer) CanChiiPonKan() bool             { return p.drawnTile == nil }
func (p stubPlayerViewer) IsConcealed() bool               { return true }
func (p stubPlayerViewer) SwapCallTiles() []tile.Tile      { return nil }
