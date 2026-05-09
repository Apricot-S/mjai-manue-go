package ai_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/ai"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestManueAgent_Decide_SelectsRiichiBeforeDiscard(t *testing.T) {
	self := seat.MustSeat(0)
	hands := newValidHands()
	hands[0] = [13]tile.Tile{
		tile.MustTileFromCode("1m"), tile.MustTileFromCode("2m"), tile.MustTileFromCode("3m"),
		tile.MustTileFromCode("4p"), tile.MustTileFromCode("5p"), tile.MustTileFromCode("6p"),
		tile.MustTileFromCode("7s"), tile.MustTileFromCode("8s"), tile.MustTileFromCode("9s"),
		tile.MustTileFromCode("E"), tile.MustTileFromCode("E"), tile.MustTileFromCode("S"),
		tile.MustTileFromCode("W"),
	}
	state := mustNewRoundStateForTest(t, hands)
	drawnTile := tile.MustTileFromCode("S")
	if err := state.Apply(event.NewDraw(self, drawnTile)); err != nil {
		t.Fatalf("Apply() failed: %v", err)
	}
	legalActions, err := state.LegalActions(self)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if len(legalActions) == 0 {
		t.Fatal("LegalActions() returned no actions")
	}

	got, err := ai.NewManueAgent(0).Decide(ai.Request{
		Self:  self,
		Round: state,
	})
	if err != nil {
		t.Fatalf("Decide() failed: %v", err)
	}
	riichi, ok := got.Action.(*action.Riichi)
	if !ok {
		t.Fatalf("Action = %T, want *action.Riichi", got.Action)
	}
	if riichi.Actor() != self {
		t.Errorf("Actor() = %v, want %v", riichi.Actor(), self)
	}
	if !containsRiichiForTest(legalActions) {
		t.Fatal("test setup error: legal actions contain no Riichi")
	}
}

func TestManueAgent_Decide_NoLegalActions(t *testing.T) {
	self := seat.MustSeat(0)
	state := mustNewRoundStateForTest(t, newValidHands())

	if _, err := ai.NewManueAgent(0).Decide(ai.Request{
		Self:  self,
		Round: state,
	}); err == nil {
		t.Fatal("Decide() succeeded unexpectedly")
	}
}

func containsRiichiForTest(actions []action.Action) bool {
	for _, a := range actions {
		if _, ok := a.(*action.Riichi); ok {
			return true
		}
	}
	return false
}
