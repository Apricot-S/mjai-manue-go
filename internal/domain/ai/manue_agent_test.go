package ai_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/configs"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/ai"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestManueAgent_Decide_SelectsLegalSelfTurnAction(t *testing.T) {
	self := seat.MustSeat(0)
	hands := newValidHands()
	hands[0] = [13]tile.Tile{
		tile.MustTileFromCode("1m"), tile.MustTileFromCode("2m"), tile.MustTileFromCode("3m"),
		tile.MustTileFromCode("4m"), tile.MustTileFromCode("5m"), tile.MustTileFromCode("6m"),
		tile.MustTileFromCode("7m"), tile.MustTileFromCode("8m"), tile.MustTileFromCode("9m"),
		tile.MustTileFromCode("1p"), tile.MustTileFromCode("1p"), tile.MustTileFromCode("E"),
		tile.MustTileFromCode("E"),
	}
	state := mustNewRoundStateForTest(t, hands)
	drawnTile := tile.MustTileFromCode("5m")
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

	agent := newManueAgentForTest(t)
	got, err := agent.Decide(ai.Request{
		Self:  self,
		Round: state,
	})
	if err != nil {
		t.Fatalf("Decide() failed: %v", err)
	}
	if got.Action == nil {
		t.Fatal("Action = nil, want legal action")
	}
	if actorAction, ok := got.Action.(interface{ Actor() seat.Seat }); ok && actorAction.Actor() != self {
		t.Errorf("Actor() = %v, want %v", actorAction.Actor(), self)
	}
	if !containsActionForTest(legalActions, got.Action) {
		t.Errorf("Action = %T %[1]v, want one of legal actions", got.Action)
	}
}

func TestManueAgent_Decide_NoLegalActions(t *testing.T) {
	self := seat.MustSeat(0)
	state := mustNewRoundStateForTest(t, newValidHands())

	agent := newManueAgentForTest(t)
	if _, err := agent.Decide(ai.Request{
		Self:  self,
		Round: state,
	}); err == nil {
		t.Fatal("Decide() succeeded unexpectedly")
	}
}

func newManueAgentForTest(t *testing.T) *ai.ManueAgent {
	t.Helper()
	stats, err := configs.LoadGameStats()
	if err != nil {
		t.Fatalf("LoadGameStats() failed: %v", err)
	}
	dangerTree, err := configs.LoadDangerTree()
	if err != nil {
		t.Fatalf("LoadDangerTree() failed: %v", err)
	}
	agent, err := ai.NewManueAgent(0, ai.ManueAgentDeps{
		Stats:  stats,
		Danger: ai.NewDangerEstimator(dangerTree),
	})
	if err != nil {
		t.Fatalf("NewManueAgent() failed: %v", err)
	}
	return agent
}

func containsActionForTest(actions []action.Action, target action.Action) bool {
	for _, a := range actions {
		if a == target {
			return true
		}
	}
	return false
}
