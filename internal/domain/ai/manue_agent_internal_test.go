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

func TestManueAgent_DecideActionSkeleton(t *testing.T) {
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
		hand        *hand.VisibleHand
		riichiState player.RiichiState
		drawnTile   *tile.Tile
		want        action.Action
		decide      func(*ManueAgent, []action.Action, player.PlayerViewer) (Decision, error)
	}{
		{
			name:        "win first",
			actions:     []action.Action{handDiscard, win, pass},
			riichiState: player.NotRiichi,
			drawnTile:   &drawnTile,
			want:        win,
			decide: func(agent *ManueAgent, actions []action.Action, self player.PlayerViewer) (Decision, error) {
				if win := firstActionOfType[*action.Win](actions); win != nil {
					return Decision{Action: win}, nil
				}
				t.Fatal("firstActionOfType[*action.Win]() returned nil")
				return Decision{}, nil
			},
		},
		{
			name:        "riichi accepted tsumogiri",
			actions:     []action.Action{handDiscard, tsumogiriDiscard},
			riichiState: player.RiichiAccepted,
			drawnTile:   &drawnTile,
			want:        tsumogiriDiscard,
			decide: func(agent *ManueAgent, actions []action.Action, self player.PlayerViewer) (Decision, error) {
				return agent.decideSelfTurn(actions, self)
			},
		},
		{
			name:        "riichi before discard",
			actions:     []action.Action{handDiscard, riichi},
			riichiState: player.NotRiichi,
			drawnTile:   &drawnTile,
			want:        riichi,
			decide: func(agent *ManueAgent, actions []action.Action, self player.PlayerViewer) (Decision, error) {
				return agent.decideSelfTurn(actions, self)
			},
		},
		{
			name:        "discard before call",
			actions:     []action.Action{chii, handDiscard, pass},
			hand:        hand.CodesToHand([]string{"1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m", "1p", "1p", "E", "S"}),
			riichiState: player.NotRiichi,
			drawnTile:   &drawnTile,
			want:        handDiscard,
			decide: func(agent *ManueAgent, actions []action.Action, self player.PlayerViewer) (Decision, error) {
				return agent.decideSelfTurn(actions, self)
			},
		},
		{
			name:        "call before pass",
			actions:     []action.Action{pass, chii},
			riichiState: player.NotRiichi,
			drawnTile:   nil,
			want:        chii,
			decide: func(agent *ManueAgent, actions []action.Action, self player.PlayerViewer) (Decision, error) {
				return agent.decideOtherDiscardReaction(actions)
			},
		},
		{
			name:        "pass only",
			actions:     []action.Action{pass},
			riichiState: player.NotRiichi,
			drawnTile:   nil,
			want:        pass,
			decide: func(agent *ManueAgent, actions []action.Action, self player.PlayerViewer) (Decision, error) {
				return agent.decideOtherDiscardReaction(actions)
			},
		},
	}

	agent := NewManueAgent(0)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decision, err := tt.decide(agent, tt.actions, stubPlayerViewer{
				hand:        tt.hand,
				riichiState: tt.riichiState,
				drawnTile:   tt.drawnTile,
			})
			if err != nil {
				t.Fatalf("decide failed: %v", err)
			}
			if decision.Action != tt.want {
				t.Errorf("Action = %T %[1]v, want %T %[2]v", decision.Action, tt.want)
			}
		})
	}
}

func TestManueAgent_decideSelfTurn_ReturnsErrorWithoutTsumogiriAfterRiichiAccepted(t *testing.T) {
	self := seat.MustSeat(0)
	handDiscard, err := action.NewDiscard(self, tile.MustTileFromCode("1m"), false)
	if err != nil {
		t.Fatalf("NewDiscard(hand) failed: %v", err)
	}
	drawnTile := tile.MustTileFromCode("5p")

	_, err = NewManueAgent(0).decideSelfTurn([]action.Action{handDiscard}, stubPlayerViewer{
		riichiState: player.RiichiAccepted,
		drawnTile:   &drawnTile,
	})
	if err == nil {
		t.Fatal("selectAction() succeeded unexpectedly")
	}
}

type stubPlayerViewer struct {
	hand        *hand.VisibleHand
	riichiState player.RiichiState
	drawnTile   *tile.Tile
}

func (p stubPlayerViewer) Hand() (*hand.VisibleHand, bool) {
	if p.hand == nil {
		return nil, false
	}
	return p.hand, true
}
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
