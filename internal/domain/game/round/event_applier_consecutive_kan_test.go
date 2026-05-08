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

func TestState_Apply_ConsecutiveKan_AllowsOpenThenOpenWithMahjongSoulOrder(t *testing.T) {
	s := newStateBeforeOpenThenOpenKanForTest(t)
	actor := seat.MustSeat(3)
	firstReplacementTile := tile.MustTileFromCode("S")
	secondReplacementTile := tile.MustTileFromCode("W")

	// MahjongSoul order:
	// open kan -> replacement tile -> dora -> open kan -> replacement tile -> dora -> discard.
	applyCalledKanForTest("E")(t, s, actor)
	applyDrawForTest(t, s, actor, firstReplacementTile)
	applyDoraForTest(t, s, "7p")
	applyPromotedKanForTest(firstReplacementTile.String())(t, s, actor)
	applyDrawForTest(t, s, actor, secondReplacementTile)
	applyDoraForTest(t, s, "8p")
	applyDiscardForTest(t, s, actor, secondReplacementTile)
	assertDoraIndicatorCountForTest(t, s, 3)
}

func TestState_Apply_ConsecutiveKan_AllowsOpenThenOpenWithTenhouOrder(t *testing.T) {
	s := newStateBeforeOpenThenOpenKanForTest(t)
	actor := seat.MustSeat(3)
	firstReplacementTile := tile.MustTileFromCode("S")
	secondReplacementTile := tile.MustTileFromCode("W")

	// Tenhou order:
	// open kan -> replacement tile -> open kan -> dora -> replacement tile -> dora -> discard.
	applyCalledKanForTest("E")(t, s, actor)
	applyDrawForTest(t, s, actor, firstReplacementTile)
	applyPromotedKanForTest(firstReplacementTile.String())(t, s, actor)
	applyDoraForTest(t, s, "6p")
	applyDrawForTest(t, s, actor, secondReplacementTile)
	applyDoraForTest(t, s, "8p")
	applyDiscardForTest(t, s, actor, secondReplacementTile)
	assertDoraIndicatorCountForTest(t, s, 3)
}

func TestState_Apply_ConsecutiveKan_AllowsOpenThenConcealedWithMahjongSoulOrder(t *testing.T) {
	s := newStateBeforeOpenThenConcealedKanForTest(t)
	actor := seat.MustSeat(3)
	firstKanTile := tile.MustTileFromCode("E")
	firstReplacementTile := tile.MustTileFromCode("S")
	secondReplacementTile := tile.MustTileFromCode("W")

	// MahjongSoul order:
	// open kan -> replacement tile -> dora -> concealed kan -> dora -> replacement tile -> discard.
	applyCalledKanForTest(firstKanTile.String())(t, s, actor)
	applyDrawForTest(t, s, actor, firstReplacementTile)
	applyDoraForTest(t, s, "6p")
	applyConcealedKanForTest(firstReplacementTile.String())(t, s, actor)
	applyDoraForTest(t, s, "7p")
	applyDrawForTest(t, s, actor, secondReplacementTile)
	applyDiscardForTest(t, s, actor, secondReplacementTile)
	assertDoraIndicatorCountForTest(t, s, 3)
}

func TestState_Apply_ConsecutiveKan_AllowsOpenThenConcealedWithTenhouOrder(t *testing.T) {
	s := newStateBeforeOpenThenConcealedKanForTest(t)
	actor := seat.MustSeat(3)
	firstReplacementTile := tile.MustTileFromCode("S")
	secondReplacementTile := tile.MustTileFromCode("W")

	// Tenhou order:
	// open kan -> replacement tile -> concealed kan -> dora -> dora -> replacement tile -> discard.
	applyCalledKanForTest("E")(t, s, actor)
	applyDrawForTest(t, s, actor, firstReplacementTile)
	applyConcealedKanForTest(firstReplacementTile.String())(t, s, actor)
	applyDoraForTest(t, s, "6p")
	applyDoraForTest(t, s, "7p")
	applyDrawForTest(t, s, actor, secondReplacementTile)
	applyDiscardForTest(t, s, actor, secondReplacementTile)
	assertDoraIndicatorCountForTest(t, s, 3)
}

func TestState_Apply_ConsecutiveKan_AllowsConcealedThenOpen(t *testing.T) {
	s := newStateBeforeConcealedThenOpenKanForTest(t)
	actor := seat.MustSeat(3)
	firstReplacementTile := tile.MustTileFromCode("5p")
	secondReplacementTile := tile.MustTileFromCode("W")

	// concealed kan -> dora -> replacement tile -> open kan -> replacement tile -> dora -> discard.
	applyConcealedKanForTest("E")(t, s, actor)
	applyDoraForTest(t, s, "6p")
	applyDrawForTest(t, s, actor, firstReplacementTile)
	applyPromotedKanForTest(firstReplacementTile.String())(t, s, actor)
	applyDrawForTest(t, s, actor, secondReplacementTile)
	applyDoraForTest(t, s, "7p")
	applyDiscardForTest(t, s, actor, secondReplacementTile)
	assertDoraIndicatorCountForTest(t, s, 3)
}

func TestState_Apply_ConsecutiveKan_AllowsConcealedThenConcealed(t *testing.T) {
	s := newStateBeforeConcealedThenConcealedKanForTest(t)
	actor := seat.MustSeat(3)
	firstReplacementTile := tile.MustTileFromCode("S")
	secondReplacementTile := tile.MustTileFromCode("W")

	applyConcealedKanForTest("E")(t, s, actor)
	applyDoraForTest(t, s, "6p")
	applyDrawForTest(t, s, actor, firstReplacementTile)
	applyConcealedKanForTest(firstReplacementTile.String())(t, s, actor)
	applyDoraForTest(t, s, "7p")
	applyDrawForTest(t, s, actor, secondReplacementTile)
	applyDiscardForTest(t, s, actor, secondReplacementTile)
	assertDoraIndicatorCountForTest(t, s, 3)
}

func TestState_Apply_ConsecutiveKan_ReturnsErrorWhenConcealedKanReplacementTilePrecedesDora(t *testing.T) {
	s := newStateBeforeConcealedThenOpenKanForTest(t)
	actor := seat.MustSeat(3)

	applyConcealedKanForTest("E")(t, s, actor)

	if err := s.Apply(event.NewDraw(actor, tile.MustTileFromCode("5p"))); err == nil {
		t.Fatal("Apply(Draw) succeeded before concealed-kan dora reveal")
	}
}

func TestState_Apply_ConsecutiveKan_ReturnsErrorWhenOpenThenOpenDiscardPrecedesSecondDoraReveal(t *testing.T) {
	s := newStateBeforeOpenThenOpenKanForTest(t)
	actor := seat.MustSeat(3)
	firstReplacementTile := tile.MustTileFromCode("S")
	secondReplacementTile := tile.MustTileFromCode("W")

	applyCalledKanForTest("E")(t, s, actor)
	applyDrawForTest(t, s, actor, firstReplacementTile)
	applyDoraForTest(t, s, "6p")
	applyPromotedKanForTest(firstReplacementTile.String())(t, s, actor)
	applyDrawForTest(t, s, actor, secondReplacementTile)

	if err := s.Apply(event.NewDiscard(actor, secondReplacementTile, true)); err == nil {
		t.Fatal("Apply(Discard) succeeded before revealing pending dora indicators")
	}
}

func TestState_Apply_ConsecutiveKan_ReturnsErrorWhenTenhouOrderDrawHasOnlyOneDoraReveal(t *testing.T) {
	s := newStateBeforeOpenThenConcealedKanForTest(t)
	actor := seat.MustSeat(3)
	firstReplacementTile := tile.MustTileFromCode("S")
	secondReplacementTile := tile.MustTileFromCode("W")

	applyCalledKanForTest("E")(t, s, actor)
	applyDrawForTest(t, s, actor, firstReplacementTile)
	applyConcealedKanForTest(firstReplacementTile.String())(t, s, actor)
	applyDoraForTest(t, s, "6p")

	if err := s.Apply(event.NewDraw(actor, secondReplacementTile)); err == nil {
		t.Fatal("Apply(Draw) succeeded before revealing all pending dora indicators")
	}
}

func assertDoraIndicatorCountForTest(t *testing.T, s *State, want int) {
	t.Helper()

	if got := len(s.DoraIndicators()); got != want {
		t.Fatalf("DoraIndicators() length = %d, want %d", got, want)
	}
}

func applyCalledKanForTest(code string) func(*testing.T, *State, seat.Seat) {
	return func(t *testing.T, s *State, actor seat.Seat) {
		t.Helper()

		target := seat.MustSeat(0)
		taken := tile.MustTileFromCode(code)
		if err := s.Apply(event.NewDraw(target, taken)); err != nil {
			t.Fatalf("Apply(Draw target %s) failed: %v", code, err)
		}
		if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
			t.Fatalf("Apply(Discard target %s) failed: %v", code, err)
		}
		if err := s.Apply(event.NewCalledKan(actor, target, taken, [3]tile.Tile{taken, taken, taken})); err != nil {
			t.Fatalf("Apply(CalledKan %s) failed: %v", code, err)
		}
	}
}

func applyPromotedKanForTest(code string) func(*testing.T, *State, seat.Seat) {
	return func(t *testing.T, s *State, actor seat.Seat) {
		t.Helper()

		added := tile.MustTileFromCode(code)
		if err := s.Apply(event.NewPromotedKan(actor, added, [3]tile.Tile{added, added, added})); err != nil {
			t.Fatalf("Apply(PromotedKan %s) failed: %v", code, err)
		}
	}
}

func applyConcealedKanForTest(code string) func(*testing.T, *State, seat.Seat) {
	return func(t *testing.T, s *State, actor seat.Seat) {
		t.Helper()

		kanTile := tile.MustTileFromCode(code)
		if err := s.Apply(event.NewConcealedKan(actor, [4]tile.Tile{kanTile, kanTile, kanTile, kanTile})); err != nil {
			t.Fatalf("Apply(ConcealedKan %s) failed: %v", code, err)
		}
	}
}

func applyDrawForTest(t *testing.T, s *State, actor seat.Seat, drawn tile.Tile) {
	t.Helper()

	if err := s.Apply(event.NewDraw(actor, drawn)); err != nil {
		t.Fatalf("Apply(Draw %s) failed: %v", drawn, err)
	}
}

func applyDoraForTest(t *testing.T, s *State, code string) {
	t.Helper()

	if err := s.Apply(event.NewDora(tile.MustTileFromCode(code))); err != nil {
		t.Fatalf("Apply(Dora %s) failed: %v", code, err)
	}
}

func applyDiscardForTest(t *testing.T, s *State, actor seat.Seat, discarded tile.Tile) {
	t.Helper()

	if err := s.Apply(event.NewDiscard(actor, discarded, true)); err != nil {
		t.Fatalf("Apply(Discard %s) failed: %v", discarded, err)
	}
}

func newStateBeforeOpenThenOpenKanForTest(t *testing.T) *State {
	t.Helper()
	return newStateWithActorForConsecutiveKanTest(t, playerBeforeOpenThenOpenKanForTest(t))
}

func newStateBeforeOpenThenConcealedKanForTest(t *testing.T) *State {
	t.Helper()
	return newStateWithActorForConsecutiveKanTest(t, playerBeforeOpenThenConcealedKanForTest(t))
}

func newStateBeforeConcealedThenOpenKanForTest(t *testing.T) *State {
	t.Helper()
	return newStateWithActorForConsecutiveKanTest(t, playerBeforeConcealedThenOpenKanForTest(t))
}

func newStateBeforeConcealedThenConcealedKanForTest(t *testing.T) *State {
	t.Helper()
	return newStateWithActorForConsecutiveKanTest(t, playerBeforeConcealedThenConcealedKanForTest(t))
}

func newStateWithActorForConsecutiveKanTest(t *testing.T, actorPlayer player.Player) *State {
	t.Helper()

	players := [common.NumPlayers]player.Player{
		player.NewInvisiblePlayer(),
		player.NewInvisiblePlayer(),
		player.NewInvisiblePlayer(),
		actorPlayer,
	}
	s := NewStateForTest(
		wind.East,
		1,
		0,
		0,
		[common.NumPlayers]int{25000, 25000, 25000, 25000},
		seat.MustSeat(0),
		seat.MustSeat(0),
		tile.Tiles{tile.MustTileFromCode("E")},
		10,
		players,
	)
	return &s
}

func playerBeforeOpenThenOpenKanForTest(t *testing.T) player.Player {
	t.Helper()

	p := visiblePlayerForConsecutiveKanTest(t, "E", "E", "E", "S", "S")
	applyPonToPlayerForTest(t, p, "S")
	prepareAfterCallForKanTest(t, p, "7s")
	return p
}

func playerBeforeOpenThenConcealedKanForTest(t *testing.T) player.Player {
	t.Helper()

	p := visiblePlayerForConsecutiveKanTest(t, "E", "E", "E", "S", "S", "S")
	return p
}

func playerBeforeConcealedThenOpenKanForTest(t *testing.T) player.Player {
	t.Helper()

	p := visiblePlayerForConsecutiveKanTest(t, "E", "E", "E", "5p", "5pr")
	applyPonToPlayerForTest(t, p, "5p")
	prepareDrawnTileForKanTest(t, p, "7s", "E")
	return p
}

func playerBeforeConcealedThenConcealedKanForTest(t *testing.T) player.Player {
	t.Helper()

	p := visiblePlayerForConsecutiveKanTest(t, "E", "E", "E", "S", "S", "S")
	if err := p.Draw(tile.MustTileFromCode("E")); err != nil {
		t.Fatalf("Draw(E) failed: %v", err)
	}
	return p
}

func visiblePlayerForConsecutiveKanTest(t *testing.T, codes ...string) player.Player {
	t.Helper()

	baseCodes := []string{"1m", "2m", "3m", "4p", "6p", "7s", "8s", "9s", "W", "N", "P", "F", "C"}
	copy(baseCodes, codes)

	var handTiles [common.InitHandSize]tile.Tile
	for i, code := range baseCodes {
		handTiles[i] = tile.MustTileFromCode(code)
	}
	p, err := player.NewVisiblePlayer(handTiles)
	if err != nil {
		t.Fatalf("player.NewVisiblePlayer() failed: %v", err)
	}
	return p
}

func applyPonToPlayerForTest(t *testing.T, p player.Player, code string) {
	t.Helper()

	ponTile := tile.MustTileFromCode(code)
	consumed := [2]tile.Tile{ponTile, ponTile}
	if code == "5p" {
		consumed[1] = tile.MustTileFromCode("5pr")
	}
	pon := meld.MustPon(ponTile, consumed, seat.MustSeat(0))
	if err := p.Pon(*pon); err != nil {
		t.Fatalf("Pon(%s) failed: %v", code, err)
	}
}

func prepareAfterCallForKanTest(t *testing.T, p player.Player, discardCode string) {
	t.Helper()

	if err := p.Discard(tile.MustTileFromCode(discardCode), false); err != nil {
		t.Fatalf("Discard(%s) failed: %v", discardCode, err)
	}
}

func prepareDrawnTileForKanTest(t *testing.T, p player.Player, discardCode string, drawCode string) {
	t.Helper()

	prepareAfterCallForKanTest(t, p, discardCode)
	if err := p.Draw(tile.MustTileFromCode(drawCode)); err != nil {
		t.Fatalf("Draw(%s) failed: %v", drawCode, err)
	}
}
