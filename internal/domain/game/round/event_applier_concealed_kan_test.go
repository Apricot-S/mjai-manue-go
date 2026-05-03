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

func concealedKanHandForTest() [common.InitHandSize]tile.Tile {
	return [common.InitHandSize]tile.Tile{
		*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m"), *tile.MustTileFromCode("3m"),
		*tile.MustTileFromCode("4p"), *tile.MustTileFromCode("5p"), *tile.MustTileFromCode("6p"),
		*tile.MustTileFromCode("7s"), *tile.MustTileFromCode("8s"), *tile.MustTileFromCode("9s"),
		*tile.MustTileFromCode("E"), *tile.MustTileFromCode("E"), *tile.MustTileFromCode("E"),
		*tile.MustTileFromCode("W"),
	}
}

func TestState_Apply_ConcealedKan(t *testing.T) {
	hands := newValidHands()
	hands[0] = concealedKanHandForTest()
	s := mustNewRoundStateForTest(t, hands)
	actor := *seat.MustSeat(0)
	kanTile := *tile.MustTileFromCode("E")
	consumed := [4]tile.Tile{kanTile, kanTile, kanTile, kanTile}

	if err := s.Apply(event.NewDraw(actor, kanTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewConcealedKan(actor, consumed)); err != nil {
		t.Fatalf("Apply(ConcealedKan) failed: %v", err)
	}

	if got := s.Player(actor).Melds(); len(got) != 1 {
		t.Fatalf("Melds() length = %d, want 1", len(got))
	}
	if s.Player(actor).CanDiscard() {
		t.Error("CanDiscard() = true, want false after ConcealedKan")
	}
	if !s.Player(actor).IsConcealed() {
		t.Error("IsConcealed() = false, want true after ConcealedKan from concealed hand")
	}
}

func TestState_Apply_ConcealedKan_KeepsOpenHandOpen(t *testing.T) {
	actor := *seat.MustSeat(0)
	kanTile := *tile.MustTileFromCode("E")
	p := openPlayerWithDrawnKanTileForTest(t, kanTile)
	players := [common.NumPlayers]player.Player{
		p,
		player.NewInvisiblePlayer(),
		player.NewInvisiblePlayer(),
		player.NewInvisiblePlayer(),
	}
	s := NewStateForTest(
		wind.East,
		1,
		0,
		0,
		[common.NumPlayers]int{25000, 25000, 25000, 25000},
		actor,
		actor,
		tile.Tiles{*tile.MustTileFromCode("E")},
		10,
		players,
	)
	s.pendingDiscard = &actor
	consumed := [4]tile.Tile{kanTile, kanTile, kanTile, kanTile}

	if err := s.Apply(event.NewConcealedKan(actor, consumed)); err != nil {
		t.Fatalf("Apply(ConcealedKan) failed: %v", err)
	}

	if s.Player(actor).IsConcealed() {
		t.Error("IsConcealed() = true, want false after ConcealedKan from already open hand")
	}
}

func TestState_Apply_ConcealedKan_RequiresDoraBeforeReplacementTile(t *testing.T) {
	hands := newValidHands()
	hands[0] = concealedKanHandForTest()
	s := mustNewRoundStateForTest(t, hands)
	actor := *seat.MustSeat(0)
	kanTile := *tile.MustTileFromCode("E")
	consumed := [4]tile.Tile{kanTile, kanTile, kanTile, kanTile}

	if err := s.Apply(event.NewDraw(actor, kanTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewConcealedKan(actor, consumed)); err != nil {
		t.Fatalf("Apply(ConcealedKan) failed: %v", err)
	}

	if err := s.Apply(event.NewDraw(actor, *tile.MustTileFromCode("W"))); err == nil {
		t.Fatal("Apply(Draw) succeeded unexpectedly")
	}
}

func TestState_Apply_ConcealedKan_ReturnsErrorWhenNoReplacementTileLeft(t *testing.T) {
	hands := newValidHands()
	hands[0] = concealedKanHandForTest()
	s := NewStateForTest(
		wind.East,
		1,
		0,
		0,
		[common.NumPlayers]int{25000, 25000, 25000, 25000},
		*seat.MustSeat(0),
		*seat.MustSeat(0),
		tile.Tiles{*tile.MustTileFromCode("E")},
		1,
		newVisiblePlayersForTest(t, hands),
	)
	actor := *seat.MustSeat(0)
	kanTile := *tile.MustTileFromCode("E")
	consumed := [4]tile.Tile{kanTile, kanTile, kanTile, kanTile}

	if err := s.Apply(event.NewDraw(actor, kanTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if got := s.NumLeftTiles(); got != 0 {
		t.Fatalf("NumLeftTiles() = %d, want 0", got)
	}

	if err := s.Apply(event.NewConcealedKan(actor, consumed)); err == nil {
		t.Fatal("Apply(ConcealedKan) succeeded unexpectedly")
	}

	if got := s.Player(actor).Melds(); len(got) != 0 {
		t.Fatalf("Melds() length = %d, want 0", len(got))
	}
	if got := s.Player(actor).DrawnTile(); got == nil || *got != kanTile {
		t.Fatalf("DrawnTile() = %v, want %v", got, kanTile)
	}
}

func TestState_Apply_ConcealedKan_ReturnsErrorOnFifthKan(t *testing.T) {
	hands := newValidHands()
	hands[0] = concealedKanHandForTest()
	s := newStateForTestWithNumKans(mustNewRoundStateForTest(t, hands), maxNumKan)
	actor := *seat.MustSeat(0)
	kanTile := *tile.MustTileFromCode("E")
	consumed := [4]tile.Tile{kanTile, kanTile, kanTile, kanTile}

	if err := s.Apply(event.NewDraw(actor, kanTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}

	if err := s.Apply(event.NewConcealedKan(actor, consumed)); err == nil {
		t.Fatal("Apply(ConcealedKan) succeeded unexpectedly")
	}

	if got := s.Player(actor).Melds(); len(got) != 0 {
		t.Fatalf("Melds() length = %d, want 0", len(got))
	}
	if got := s.Player(actor).DrawnTile(); got == nil || *got != kanTile {
		t.Fatalf("DrawnTile() = %v, want %v", got, kanTile)
	}
}

func TestState_Apply_ConcealedKan_AllowsReplacementTileAfterDora(t *testing.T) {
	hands := newValidHands()
	hands[0] = concealedKanHandForTest()
	s := mustNewRoundStateForTest(t, hands)
	actor := *seat.MustSeat(0)
	kanTile := *tile.MustTileFromCode("E")
	replacementTile := *tile.MustTileFromCode("W")
	consumed := [4]tile.Tile{kanTile, kanTile, kanTile, kanTile}

	if err := s.Apply(event.NewDraw(actor, kanTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewConcealedKan(actor, consumed)); err != nil {
		t.Fatalf("Apply(ConcealedKan) failed: %v", err)
	}
	if err := s.Apply(event.NewDora(*tile.MustTileFromCode("6p"))); err != nil {
		t.Fatalf("Apply(Dora) failed: %v", err)
	}

	if err := s.Apply(event.NewDraw(actor, replacementTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if got := s.Player(actor).DrawnTile(); got == nil || *got != replacementTile {
		t.Fatalf("DrawnTile() = %v, want %v", got, replacementTile)
	}
}

func openPlayerWithDrawnKanTileForTest(t *testing.T, kanTile tile.Tile) player.Player {
	t.Helper()

	handTiles := [common.InitHandSize]tile.Tile{
		*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m"), *tile.MustTileFromCode("3m"),
		*tile.MustTileFromCode("4p"), *tile.MustTileFromCode("5p"), *tile.MustTileFromCode("6p"),
		*tile.MustTileFromCode("7s"), *tile.MustTileFromCode("8s"), *tile.MustTileFromCode("9s"),
		kanTile, kanTile, kanTile, *tile.MustTileFromCode("5pr"),
	}
	p, err := player.NewVisiblePlayer(handTiles)
	if err != nil {
		t.Fatalf("player.NewVisiblePlayer() failed: %v", err)
	}
	pon := meld.MustPon(
		*tile.MustTileFromCode("5p"),
		[2]tile.Tile{*tile.MustTileFromCode("5p"), *tile.MustTileFromCode("5pr")},
		*seat.MustSeat(1),
	)
	if err := p.Pon(*pon); err != nil {
		t.Fatalf("Pon() failed: %v", err)
	}
	if err := p.Discard(*tile.MustTileFromCode("1m"), false); err != nil {
		t.Fatalf("Discard() failed: %v", err)
	}
	if err := p.Draw(kanTile); err != nil {
		t.Fatalf("Draw() failed: %v", err)
	}
	return p
}
