package ai_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/ai"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

func TestTsumogiriAgent_Decide(t *testing.T) {
	self := *seat.MustSeat(0)
	drawnTile := *tile.MustTileFromCode("6m")
	roundState := mustNewRoundStateForTest(t, newValidHands())
	if err := roundState.Apply(event.NewDraw(self, drawnTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}

	got, err := ai.NewTsumogiriAgent().Decide(ai.Request{
		Self:  self,
		Round: roundState,
	})
	if err != nil {
		t.Fatalf("Decide() failed: %v", err)
	}
	discard, ok := got.Action.(*action.Discard)
	if !ok {
		t.Fatalf("Action = %T, want *action.Discard", got.Action)
	}
	if discard.Actor() != self {
		t.Errorf("Actor() = %v, want %v", discard.Actor(), self)
	}
	if discard.Tile().ID() != drawnTile.ID() {
		t.Errorf("Tile() = %v, want %v", discard.Tile(), drawnTile)
	}
	if !discard.Tsumogiri() {
		t.Error("Tsumogiri() = false, want true")
	}
}

func TestTsumogiriAgent_Decide_NoDrawnTile(t *testing.T) {
	self := *seat.MustSeat(0)
	roundState := mustNewRoundStateForTest(t, newValidHands())

	got, err := ai.NewTsumogiriAgent().Decide(ai.Request{
		Self:  self,
		Round: roundState,
	})
	if err != nil {
		t.Fatalf("Decide() failed: %v", err)
	}
	pass, ok := got.Action.(*action.Pass)
	if !ok {
		t.Fatalf("Action = %T, want *action.Pass", got.Action)
	}
	if pass.Actor() != self {
		t.Errorf("Actor() = %v, want %v", pass.Actor(), self)
	}
}

func mustNewRoundStateForTest(t *testing.T, hands [common.NumPlayers][common.InitHandSize]tile.Tile) *round.State {
	t.Helper()

	ev := event.NewStartRound(
		wind.East,
		1,
		0,
		0,
		*seat.MustSeat(0),
		*tile.MustTileFromCode("E"),
		&[common.NumPlayers]int{25000, 25000, 25000, 25000},
		hands,
	)

	s, err := round.NewState(ev, [common.NumPlayers]int{25000, 25000, 25000, 25000})
	if err != nil {
		t.Fatalf("round.NewState() failed: %v", err)
	}
	return s
}

func newValidHands() [common.NumPlayers][common.InitHandSize]tile.Tile {
	return [common.NumPlayers][common.InitHandSize]tile.Tile{
		{
			*tile.MustTileFromCode("1m"),
			*tile.MustTileFromCode("2m"),
			*tile.MustTileFromCode("3m"),
			*tile.MustTileFromCode("4m"),
			*tile.MustTileFromCode("5m"),
			*tile.MustTileFromCode("6m"),
			*tile.MustTileFromCode("7m"),
			*tile.MustTileFromCode("8m"),
			*tile.MustTileFromCode("9m"),
			*tile.MustTileFromCode("1p"),
			*tile.MustTileFromCode("2p"),
			*tile.MustTileFromCode("3p"),
			*tile.MustTileFromCode("4p"),
		},
		unknownHand(),
		unknownHand(),
		unknownHand(),
	}
}

func unknownHand() [common.InitHandSize]tile.Tile {
	var hand [common.InitHandSize]tile.Tile
	for i := range common.InitHandSize {
		hand[i] = *tile.MustTileFromCode("?")
	}
	return hand
}
