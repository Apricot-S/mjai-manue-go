package ai_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/ai"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestManueAgent_Decide(t *testing.T) {
	self := seat.MustSeat(0)
	state := mustNewRoundStateForTest(t, newValidHands())
	drawnTile := tile.MustTileFromCode("1m")
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
	if got.Action != legalActions[0] {
		t.Errorf("Action = %v, want first legal action %v", got.Action, legalActions[0])
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
