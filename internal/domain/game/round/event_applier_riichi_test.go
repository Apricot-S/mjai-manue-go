package round

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func riichiReadyHandForTest() [common.InitHandSize]tile.Tile {
	return [common.InitHandSize]tile.Tile{
		*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m"), *tile.MustTileFromCode("3m"),
		*tile.MustTileFromCode("4p"), *tile.MustTileFromCode("5p"), *tile.MustTileFromCode("6p"),
		*tile.MustTileFromCode("7s"), *tile.MustTileFromCode("8s"), *tile.MustTileFromCode("9s"),
		*tile.MustTileFromCode("E"), *tile.MustTileFromCode("E"), *tile.MustTileFromCode("S"),
		*tile.MustTileFromCode("W"),
	}
}

func TestState_Apply_Riichi(t *testing.T) {
	hands := newValidHands()
	hands[0] = riichiReadyHandForTest()
	s := mustNewRoundStateForTest(t, hands)
	actor := *seat.MustSeat(0)

	if err := s.Apply(event.NewDraw(actor, *tile.MustTileFromCode("S"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewRiichi(actor)); err != nil {
		t.Fatalf("Apply(Riichi) failed: %v", err)
	}

	if got := s.Player(actor).RiichiState(); got != player.RiichiDeclared {
		t.Fatalf("RiichiState() = %v, want %v", got, player.RiichiDeclared)
	}
}

func TestState_Apply_RiichiAccepted(t *testing.T) {
	hands := newValidHands()
	hands[0] = riichiReadyHandForTest()
	s := mustNewRoundStateForTest(t, hands)
	actor := *seat.MustSeat(0)

	if err := s.Apply(event.NewDraw(actor, *tile.MustTileFromCode("S"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewRiichi(actor)); err != nil {
		t.Fatalf("Apply(Riichi) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(actor, *tile.MustTileFromCode("W"), false)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	scores := [common.NumPlayers]int{24000, 25000, 25000, 25000}
	deltas := [common.NumPlayers]int{-1000, 0, 0, 0}
	if err := s.Apply(event.NewRiichiAccepted(actor, &deltas, &scores)); err != nil {
		t.Fatalf("Apply(RiichiAccepted) failed: %v", err)
	}

	if got := s.Player(actor).RiichiState(); got != player.RiichiAccepted {
		t.Errorf("RiichiState() = %v, want %v", got, player.RiichiAccepted)
	}
	if got := s.RiichiDeposit(); got != 1 {
		t.Errorf("RiichiDeposit() = %d, want 1", got)
	}
	if got := s.Scores(); got != scores {
		t.Errorf("Scores() = %v, want %v", got, scores)
	}
}
