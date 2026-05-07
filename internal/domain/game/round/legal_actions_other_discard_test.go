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

func TestState_LegalActions_OnOtherDiscardNoAction(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	target := seat.MustSeat(0)
	actor := seat.MustSeat(2)
	discardedTile := tile.MustTileFromCode("6m")
	if err := s.Apply(event.NewDraw(target, discardedTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, discardedTile, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("LegalActions() = %v, want empty", got)
	}
}

func TestState_LegalActions_OnOtherDiscardIncludesRon(t *testing.T) {
	hands := newValidHands()
	hands[1] = ronWithTanyaoHandForLegalActionsTest()
	s := mustNewRoundStateForTest(t, hands)
	target := seat.MustSeat(0)
	actor := seat.MustSeat(1)
	winningTile := tile.MustTileFromCode("3p")
	if err := s.Apply(event.NewDraw(target, winningTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, winningTile, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if !containsWin(got, actor, target, "3p") {
		t.Error("LegalActions() does not contain Win, want ron with tanyao")
	}
	if !containsPass(got, actor) {
		t.Error("LegalActions() does not contain Pass, want pass when ron is available")
	}
}

func TestState_LegalActions_OnOtherDiscardExcludesRonWithoutYaku(t *testing.T) {
	hands := newValidHands()
	hands[1] = ronWithoutYakuHandForLegalActionsTest()
	s := mustNewRoundStateForTest(t, hands)
	target := seat.MustSeat(0)
	actor := seat.MustSeat(1)
	winningTile := tile.MustTileFromCode("9s")
	if err := s.Apply(event.NewDraw(target, winningTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, winningTile, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if containsWin(got, actor, target, "9s") {
		t.Error("LegalActions() contains Win, want no-yaku ron excluded")
	}
}

func TestState_LegalActions_OnOtherDiscardIncludesRonLastTile(t *testing.T) {
	hands := newValidHands()
	hands[1] = ronWithoutYakuHandForLegalActionsTest()
	s := mustNewRoundStateForTest(t, hands)
	target := seat.MustSeat(0)
	actor := seat.MustSeat(1)
	winningTile := tile.MustTileFromCode("9s")
	if err := s.Apply(event.NewDraw(target, winningTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	s.numLeftTiles = 0
	if err := s.Apply(event.NewDiscard(target, winningTile, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if !containsWin(got, actor, target, "9s") {
		t.Error("LegalActions() does not contain Win, want houtei ron")
	}
}

func TestState_LegalActions_OnRobbingKanIncludesRon(t *testing.T) {
	s := newStateBeforeRobbingKanForLegalActionsTest(t)
	target := seat.MustSeat(3)
	actor := seat.MustSeat(1)
	added := tile.MustTileFromCode("E")
	if err := s.Apply(event.NewPromotedKan(target, added, [3]tile.Tile{added, added, added})); err != nil {
		t.Fatalf("Apply(PromotedKan) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if !containsWin(got, actor, target, "E") {
		t.Error("LegalActions() does not contain Win, want robbing-a-kan ron")
	}
	if !containsPass(got, actor) {
		t.Error("LegalActions() does not contain Pass, want pass when robbing-a-kan ron is available")
	}
}

func ronWithTanyaoHandForLegalActionsTest() [common.InitHandSize]tile.Tile {
	return [common.InitHandSize]tile.Tile{
		tile.MustTileFromCode("2m"),
		tile.MustTileFromCode("3m"),
		tile.MustTileFromCode("4m"),
		tile.MustTileFromCode("2p"),
		tile.MustTileFromCode("4p"),
		tile.MustTileFromCode("3s"),
		tile.MustTileFromCode("4s"),
		tile.MustTileFromCode("5s"),
		tile.MustTileFromCode("6s"),
		tile.MustTileFromCode("7s"),
		tile.MustTileFromCode("8s"),
		tile.MustTileFromCode("6m"),
		tile.MustTileFromCode("6m"),
	}
}

func newStateBeforeRobbingKanForLegalActionsTest(t *testing.T) *State {
	t.Helper()

	actor := seat.MustSeat(1)
	target := seat.MustSeat(3)
	players := [common.NumPlayers]player.Player{
		player.NewInvisiblePlayer(),
		robbingKanPlayerForLegalActionsTest(t),
		player.NewInvisiblePlayer(),
		playerBeforePromotedKanForTest(t),
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
	s.pendingDiscard = &target
	if actor == target {
		t.Fatal("invalid test setup: actor and target must differ")
	}
	return &s
}

func robbingKanPlayerForLegalActionsTest(t *testing.T) player.Player {
	t.Helper()

	handTiles := [common.InitHandSize]tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("2p"),
		tile.MustTileFromCode("3p"),
		tile.MustTileFromCode("4p"),
		tile.MustTileFromCode("3s"),
		tile.MustTileFromCode("4s"),
		tile.MustTileFromCode("5s"),
		tile.MustTileFromCode("6s"),
		tile.MustTileFromCode("6s"),
		tile.MustTileFromCode("6s"),
		tile.MustTileFromCode("E"),
		tile.MustTileFromCode("W"),
	}
	p, err := player.NewVisiblePlayer(handTiles)
	if err != nil {
		t.Fatalf("player.NewVisiblePlayer() failed: %v", err)
	}
	pon := meld.MustPon(
		tile.MustTileFromCode("1m"),
		[2]tile.Tile{tile.MustTileFromCode("1m"), tile.MustTileFromCode("1m")},
		seat.MustSeat(0),
	)
	if err := p.Pon(*pon); err != nil {
		t.Fatalf("Pon() failed: %v", err)
	}
	if err := p.Discard(tile.MustTileFromCode("W"), false); err != nil {
		t.Fatalf("Discard() failed: %v", err)
	}
	return p
}

func ronWithoutYakuHandForLegalActionsTest() [common.InitHandSize]tile.Tile {
	return [common.InitHandSize]tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("2p"),
		tile.MustTileFromCode("3p"),
		tile.MustTileFromCode("4p"),
		tile.MustTileFromCode("3s"),
		tile.MustTileFromCode("4s"),
		tile.MustTileFromCode("5s"),
		tile.MustTileFromCode("6s"),
		tile.MustTileFromCode("6s"),
		tile.MustTileFromCode("6s"),
		tile.MustTileFromCode("9s"),
	}
}
