package round

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

func TestState_Apply_PromotedKan(t *testing.T) {
	s := newStateBeforePromotedKanForTest(t, 10, 0)
	actor := *seat.MustSeat(3)
	added := tile.MustTileFromCode("E")

	if err := s.Apply(event.NewPromotedKan(actor, added, [3]tile.Tile{added, added, added})); err != nil {
		t.Fatalf("Apply(PromotedKan) failed: %v", err)
	}

	if got := s.Player(actor).Melds(); len(got) != 1 {
		t.Fatalf("Melds() length = %d, want 1", len(got))
	}
	if s.Player(actor).CanDiscard() {
		t.Error("CanDiscard() = true, want false after PromotedKan")
	}
	if s.Player(actor).IsConcealed() {
		t.Error("IsConcealed() = true, want false after PromotedKan")
	}
}

func TestState_Apply_PromotedKan_ReturnsErrorWhenNoReplacementTileLeft(t *testing.T) {
	s := newStateBeforePromotedKanForTest(t, 0, 0)
	actor := *seat.MustSeat(3)
	added := tile.MustTileFromCode("E")

	if err := s.Apply(event.NewPromotedKan(actor, added, [3]tile.Tile{added, added, added})); err == nil {
		t.Fatal("Apply(PromotedKan) succeeded unexpectedly")
	}

	if got := s.Player(actor).DrawnTile(); got == nil || *got != added {
		t.Fatalf("DrawnTile() = %v, want %v", got, added)
	}
}

func TestState_Apply_PromotedKan_ReturnsErrorOnFifthKan(t *testing.T) {
	s := newStateBeforePromotedKanForTest(t, 10, maxNumKan)
	actor := *seat.MustSeat(3)
	added := tile.MustTileFromCode("E")

	if err := s.Apply(event.NewPromotedKan(actor, added, [3]tile.Tile{added, added, added})); err == nil {
		t.Fatal("Apply(PromotedKan) succeeded unexpectedly")
	}

	if got := s.Player(actor).DrawnTile(); got == nil || *got != added {
		t.Fatalf("DrawnTile() = %v, want %v", got, added)
	}
}

func TestState_Apply_PromotedKan_ReturnsErrorWhenDiscardFollowsPromotedKan(t *testing.T) {
	s := newStateBeforePromotedKanForTest(t, 10, 0)
	actor := *seat.MustSeat(3)
	added := tile.MustTileFromCode("E")

	if err := s.Apply(event.NewPromotedKan(actor, added, [3]tile.Tile{added, added, added})); err != nil {
		t.Fatalf("Apply(PromotedKan) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(actor, tile.MustTileFromCode("W"), false)); err == nil {
		t.Fatal("Apply(Discard) succeeded unexpectedly")
	}
}

func TestState_Apply_PromotedKan_AllowsDoraAfterReplacementTileDraw(t *testing.T) {
	s := newStateBeforePromotedKanForTest(t, 10, 0)
	actor := *seat.MustSeat(3)
	added := tile.MustTileFromCode("E")
	replacementTile := tile.MustTileFromCode("W")
	doraIndicator := tile.MustTileFromCode("6p")

	if err := s.Apply(event.NewPromotedKan(actor, added, [3]tile.Tile{added, added, added})); err != nil {
		t.Fatalf("Apply(PromotedKan) failed: %v", err)
	}
	if err := s.Apply(event.NewDraw(actor, replacementTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDora(doraIndicator)); err != nil {
		t.Fatalf("Apply(Dora) failed: %v", err)
	}

	if got := s.DoraIndicators(); len(got) != 2 || got[1] != doraIndicator {
		t.Fatalf("DoraIndicators() = %v, want appended %v", got, doraIndicator)
	}
}

func TestState_Apply_PromotedKan_AllowsDiscardAfterDoraReveal(t *testing.T) {
	s := newStateBeforePromotedKanForTest(t, 10, 0)
	actor := *seat.MustSeat(3)
	added := tile.MustTileFromCode("E")
	replacementTile := tile.MustTileFromCode("W")

	if err := s.Apply(event.NewPromotedKan(actor, added, [3]tile.Tile{added, added, added})); err != nil {
		t.Fatalf("Apply(PromotedKan) failed: %v", err)
	}
	if err := s.Apply(event.NewDraw(actor, replacementTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDora(tile.MustTileFromCode("6p"))); err != nil {
		t.Fatalf("Apply(Dora) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(actor, replacementTile, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
}

func newStateBeforePromotedKanForTest(t *testing.T, numLeftTiles int, numKans int) *State {
	t.Helper()

	actor := *seat.MustSeat(3)
	p := playerBeforePromotedKanForTest(t)
	players := [common.NumPlayers]player.Player{
		player.NewInvisiblePlayer(),
		player.NewInvisiblePlayer(),
		player.NewInvisiblePlayer(),
		p,
	}
	s := NewStateForTest(
		wind.East,
		1,
		0,
		0,
		[common.NumPlayers]int{25000, 25000, 25000, 25000},
		*seat.MustSeat(0),
		*seat.MustSeat(0),
		tile.Tiles{tile.MustTileFromCode("E")},
		numLeftTiles,
		players,
	)
	s.pendingDiscard = &actor
	s.nextDraw = actor
	s.numKans = numKans
	return &s
}

func playerBeforePromotedKanForTest(t *testing.T) player.Player {
	t.Helper()

	handTiles := [common.InitHandSize]tile.Tile{
		tile.MustTileFromCode("1m"), tile.MustTileFromCode("2m"), tile.MustTileFromCode("3m"),
		tile.MustTileFromCode("4p"), tile.MustTileFromCode("5p"), tile.MustTileFromCode("6p"),
		tile.MustTileFromCode("7s"), tile.MustTileFromCode("8s"), tile.MustTileFromCode("9s"),
		tile.MustTileFromCode("E"), tile.MustTileFromCode("E"), tile.MustTileFromCode("S"),
		tile.MustTileFromCode("W"),
	}
	p, err := player.NewVisiblePlayer(handTiles)
	if err != nil {
		t.Fatalf("player.NewVisiblePlayer() failed: %v", err)
	}
	pon := meld.MustPon(
		tile.MustTileFromCode("E"),
		[2]tile.Tile{tile.MustTileFromCode("E"), tile.MustTileFromCode("E")},
		*seat.MustSeat(0),
	)
	if err := p.Pon(*pon); err != nil {
		t.Fatalf("Pon() failed: %v", err)
	}
	if err := p.Discard(tile.MustTileFromCode("S"), false); err != nil {
		t.Fatalf("Discard() failed: %v", err)
	}
	if err := p.Draw(tile.MustTileFromCode("E")); err != nil {
		t.Fatalf("Draw() failed: %v", err)
	}
	return p
}
