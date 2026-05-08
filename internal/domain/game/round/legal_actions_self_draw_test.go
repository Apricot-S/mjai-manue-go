package round

import (
	"fmt"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestState_LegalActions_PendingDiscard(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	actor := seat.MustSeat(0)
	drawnTile := tile.MustTileFromCode("E")
	if err := s.Apply(event.NewDraw(actor, drawnTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}

	want := map[string]bool{
		"1m:false": false,
		"1p:false": false,
		"1s:false": false,
		"2m:false": false,
		"2p:false": false,
		"2s:false": false,
		"3m:false": false,
		"3p:false": false,
		"3s:false": false,
		"4m:false": false,
		"4p:false": false,
		"4s:false": false,
		"5m:false": false,
		"E:true":   false,
	}
	if len(got) != len(want) {
		t.Fatalf("LegalActions() length = %d, want %d: %v", len(got), len(want), got)
	}
	for _, a := range got {
		discard, ok := a.(*action.Discard)
		if !ok {
			t.Fatalf("LegalActions() contains %T, want only *action.Discard", a)
		}
		if discard.Actor() != actor {
			t.Errorf("Discard.Actor() = %v, want %v", discard.Actor(), actor)
		}
		key := fmt.Sprintf("%s:%t", discard.Tile(), discard.Tsumogiri())
		if _, ok := want[key]; !ok {
			t.Errorf("unexpected discard action: %s", key)
			continue
		}
		want[key] = true
	}
	for key, found := range want {
		if !found {
			t.Errorf("missing discard action: %s", key)
		}
	}
}

func TestState_LegalActions_IncludesRiichi(t *testing.T) {
	hands := newValidHands()
	hands[0] = riichiReadyHandForTest()
	s := mustNewRoundStateForTest(t, hands)
	actor := seat.MustSeat(0)
	if err := s.Apply(event.NewDraw(actor, tile.MustTileFromCode("S"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if !containsRiichi(got, actor) {
		t.Error("LegalActions() does not contain Riichi, want Riichi for concealed tenpai hand")
	}
}

func TestState_LegalActions_ExcludesRiichi(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T) (*State, seat.Seat)
	}{
		{
			name: "not tenpai",
			setup: func(t *testing.T) (*State, seat.Seat) {
				t.Helper()
				s := mustNewRoundStateForTest(t, newValidHands())
				actor := seat.MustSeat(0)
				if err := s.Apply(event.NewDraw(actor, tile.MustTileFromCode("E"))); err != nil {
					t.Fatalf("Apply(Draw) failed: %v", err)
				}
				return s, actor
			},
		},
		{
			name: "no next draw turn remains",
			setup: func(t *testing.T) (*State, seat.Seat) {
				t.Helper()
				hands := newValidHands()
				hands[0] = riichiReadyHandForTest()
				s := mustNewRoundStateForTest(t, hands)
				actor := seat.MustSeat(0)
				if err := s.Apply(event.NewDraw(actor, tile.MustTileFromCode("S"))); err != nil {
					t.Fatalf("Apply(Draw) failed: %v", err)
				}
				s.numLeftTiles = common.NumPlayers - 1
				return s, actor
			},
		},
		{
			name: "open hand",
			setup: func(t *testing.T) (*State, seat.Seat) {
				t.Helper()
				hands := newValidHands()
				hands[1] = openChiiHandForLegalActionsTest()
				s := mustNewRoundStateForTest(t, hands)
				return stateAfterChiiForLegalActionsTest(t, s), seat.MustSeat(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, actor := tt.setup(t)
			got, err := s.LegalActions(actor)
			if err != nil {
				t.Fatalf("LegalActions() failed: %v", err)
			}
			if containsRiichi(got, actor) {
				t.Error("LegalActions() contains Riichi unexpectedly")
			}
		})
	}
}

func TestState_LegalActions_IncludesKyushukyuhai(t *testing.T) {
	hands := newValidHands()
	hands[0] = kyushukyuhaiHandForTest()
	s := mustNewRoundStateForTest(t, hands)
	actor := seat.MustSeat(0)
	if err := s.Apply(event.NewDraw(actor, tile.MustTileFromCode("W"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if !containsKyushukyuhai(got, actor) {
		t.Error("LegalActions() does not contain Kyushukyuhai, want abortive draw with 9 yaochu types")
	}
}

func TestState_LegalActions_ExcludesKyushukyuhaiAfterFirstDiscard(t *testing.T) {
	hands := newValidHands()
	hands[0] = kyushukyuhaiHandForTest()
	s := mustNewRoundStateForTest(t, hands)
	actor := seat.MustSeat(0)
	firstDraw := tile.MustTileFromCode("W")
	if err := s.Apply(event.NewDraw(actor, firstDraw)); err != nil {
		t.Fatalf("Apply(first Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(actor, firstDraw, true)); err != nil {
		t.Fatalf("Apply(first Discard) failed: %v", err)
	}
	for i := 1; i < 4; i++ {
		other := seat.MustSeat(i)
		drawnTile := tile.MustTileFromCode("6m")
		if err := s.Apply(event.NewDraw(other, drawnTile)); err != nil {
			t.Fatalf("Apply(other Draw %d) failed: %v", i, err)
		}
		if err := s.Apply(event.NewDiscard(other, drawnTile, true)); err != nil {
			t.Fatalf("Apply(other Discard %d) failed: %v", i, err)
		}
	}
	if err := s.Apply(event.NewDraw(actor, tile.MustTileFromCode("N"))); err != nil {
		t.Fatalf("Apply(second Draw) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if containsKyushukyuhai(got, actor) {
		t.Error("LegalActions() contains Kyushukyuhai after first discard")
	}
}

func TestState_LegalActions_ExcludesKyushukyuhaiForOtherPlayersAfterConcealedKan(t *testing.T) {
	hands := newValidHands()
	hands[0] = concealedKanHandForTest()
	for i := 1; i < common.NumPlayers; i++ {
		hands[i] = kyushukyuhaiHandForTest()
	}
	s := mustNewRoundStateForTest(t, hands)
	kanActor := seat.MustSeat(0)
	kanTile := tile.MustTileFromCode("E")
	if err := s.Apply(event.NewDraw(kanActor, kanTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewConcealedKan(kanActor, [4]tile.Tile{kanTile, kanTile, kanTile, kanTile})); err != nil {
		t.Fatalf("Apply(ConcealedKan) failed: %v", err)
	}
	if err := s.Apply(event.NewDora(tile.MustTileFromCode("6p"))); err != nil {
		t.Fatalf("Apply(Dora) failed: %v", err)
	}
	if err := s.Apply(event.NewDraw(kanActor, tile.MustTileFromCode("N"))); err != nil {
		t.Fatalf("Apply(replacement Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(kanActor, tile.MustTileFromCode("N"), true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	for i := 1; i < common.NumPlayers; i++ {
		actor := seat.MustSeat(i)
		if err := s.Apply(event.NewDraw(actor, tile.MustTileFromCode("W"))); err != nil {
			t.Fatalf("Apply(Draw %d) failed: %v", i, err)
		}
		got, err := s.LegalActions(actor)
		if err != nil {
			t.Fatalf("LegalActions(%d) failed: %v", i, err)
		}
		if containsKyushukyuhai(got, actor) {
			t.Errorf("LegalActions(%d) contains Kyushukyuhai after earlier concealed kan", i)
		}
		if err := s.Apply(event.NewDiscard(actor, tile.MustTileFromCode("W"), true)); err != nil {
			t.Fatalf("Apply(Discard %d) failed: %v", i, err)
		}
	}
}

func TestState_LegalActions_IncludesTsumoWin(t *testing.T) {
	hands := newValidHands()
	hands[0] = menzenTsumoHandForLegalActionsTest()
	s := mustNewRoundStateForTest(t, hands)
	actor := seat.MustSeat(0)
	winningTile := tile.MustTileFromCode("9s")
	if err := s.Apply(event.NewDraw(actor, winningTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if !containsWin(got, actor, actor, "9s") {
		t.Error("LegalActions() does not contain Win, want menzen tsumo win")
	}
}

func TestState_LegalActions_ExcludesTsumoWinWithoutYaku(t *testing.T) {
	actor := seat.MustSeat(0)
	s := newStateWithOpenNoYakuTsumoForLegalActionsTest(t, 10, false)

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if containsWin(got, actor, actor, "9s") {
		t.Error("LegalActions() contains Win, want open tsumo without yaku excluded")
	}
}

func TestState_LegalActions_IncludesTsumoWinAfterAKan(t *testing.T) {
	actor := seat.MustSeat(0)
	s := newStateWithOpenNoYakuTsumoForLegalActionsTest(t, 10, true)

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if !containsWin(got, actor, actor, "9s") {
		t.Error("LegalActions() does not contain Win, want rinshan tsumo win")
	}
}

func TestState_LegalActions_IncludesTsumoWinLastTile(t *testing.T) {
	actor := seat.MustSeat(0)
	s := newStateWithOpenNoYakuTsumoForLegalActionsTest(t, 0, false)

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if !containsWin(got, actor, actor, "9s") {
		t.Error("LegalActions() does not contain Win, want haitei tsumo win")
	}
}

func TestState_LegalActions_IncludesPromotedKan(t *testing.T) {
	s := newStateBeforePromotedKanForTest(t, 10, 0)
	actor := seat.MustSeat(3)

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if !containsPromotedKan(got, actor, "E", [3]string{"E", "E", "E"}) {
		t.Error("LegalActions() does not contain PromotedKan, want kakan for existing pon")
	}
}

func TestState_LegalActions_ExcludesPromotedKan(t *testing.T) {
	tests := []struct {
		name string
		s    *State
	}{
		{name: "no replacement tile left", s: newStateBeforePromotedKanForTest(t, 0, 0)},
		{name: "fifth kan", s: newStateBeforePromotedKanForTest(t, 10, maxNumKan)},
	}

	actor := seat.MustSeat(3)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.LegalActions(actor)
			if err != nil {
				t.Fatalf("LegalActions() failed: %v", err)
			}
			if containsPromotedKan(got, actor, "E", [3]string{"E", "E", "E"}) {
				t.Error("LegalActions() contains PromotedKan unexpectedly")
			}
		})
	}
}

func TestState_LegalActions_IncludesConcealedKan(t *testing.T) {
	hands := newValidHands()
	hands[0] = concealedKanHandForTest()
	s := mustNewRoundStateForTest(t, hands)
	actor := seat.MustSeat(0)
	kanTile := tile.MustTileFromCode("E")
	if err := s.Apply(event.NewDraw(actor, kanTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if !containsConcealedKan(got, actor, [4]string{"E", "E", "E", "E"}) {
		t.Error("LegalActions() does not contain ConcealedKan, want ankan for four identical tiles")
	}
}

func TestState_LegalActions_ExcludesConcealedKan(t *testing.T) {
	tests := []struct {
		name string
		s    *State
	}{
		{name: "no replacement tile left", s: newStateBeforeConcealedKanForLegalActionsTest(t, 0, 0)},
		{name: "fifth kan", s: newStateBeforeConcealedKanForLegalActionsTest(t, 10, maxNumKan)},
	}

	actor := seat.MustSeat(0)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.LegalActions(actor)
			if err != nil {
				t.Fatalf("LegalActions() failed: %v", err)
			}
			if containsConcealedKan(got, actor, [4]string{"E", "E", "E", "E"}) {
				t.Error("LegalActions() contains ConcealedKan unexpectedly")
			}
		})
	}
}

func TestState_LegalActions_AfterRiichiAcceptedIncludesConcealedKanWhenWaitsDoNotChange(t *testing.T) {
	s := newRiichiAcceptedStateBeforeConcealedKanForTest(t)
	actor := seat.MustSeat(0)
	kanTile := tile.MustTileFromCode("1m")
	if err := s.Apply(event.NewDraw(actor, kanTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if !containsConcealedKan(got, actor, [4]string{"1m", "1m", "1m", "1m"}) {
		t.Error("LegalActions() does not contain ConcealedKan, want ankan after riichi when waits do not change")
	}
}

func TestState_LegalActions_AfterRiichiAcceptedIncludesConcealedKanWhenOnlyWinningFormChanges(t *testing.T) {
	s := newRiichiAcceptedStateBeforeConcealedKanChangingOnlyWinningFormForTest(t)
	actor := seat.MustSeat(0)
	kanTile := tile.MustTileFromCode("2m")
	if err := s.Apply(event.NewDraw(actor, kanTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if !containsConcealedKan(got, actor, [4]string{"2m", "2m", "2m", "2m"}) {
		t.Error("LegalActions() does not contain ConcealedKan, want ankan after riichi when waits stay the same")
	}
}

func TestCanConcealedKanAfterRiichi_ReturnsFalseForFourTilesInHand(t *testing.T) {
	handBeforeKan := hand.CodesToHand([]string{
		"1m", "1m", "1m", "1m",
		"2p", "3p",
		"4s", "5s", "6s",
		"7s", "8s", "9s",
		"E",
	})
	kanTile := tile.MustTileFromCode("1m")

	if canConcealedKanAfterRiichi(handBeforeKan, kanTile, [4]tile.Tile{kanTile, kanTile, kanTile, kanTile}) {
		t.Error("canConcealedKanAfterRiichi() = true, want false for kan made from four tiles already in hand")
	}
}

func TestCanConcealedKanAfterRiichi_HandlesRedFiveInHand(t *testing.T) {
	handBeforeKan := hand.CodesToHand([]string{
		"5m", "5m", "5mr",
		"2p", "3p",
		"4s", "5s", "6s",
		"7s", "8s", "9s",
		"E", "E",
	})
	drawnTile := tile.MustTileFromCode("5m")
	consumed := [4]tile.Tile{
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("5mr"),
	}

	if !canConcealedKanAfterRiichi(handBeforeKan, drawnTile, consumed) {
		t.Error("canConcealedKanAfterRiichi() = false, want true for hand containing red five")
	}
}

func TestCanConcealedKanAfterRiichi_ReturnsFalseWhenWaitsChange(t *testing.T) {
	handBeforeKan := hand.CodesToHand([]string{
		"3m", "4m", "4m", "4m",
		"1p", "2p", "3p",
		"4s", "5s", "6s",
		"7s", "8s", "9s",
	})
	drawnTile := tile.MustTileFromCode("4m")
	consumed := [4]tile.Tile{drawnTile, drawnTile, drawnTile, drawnTile}

	if canConcealedKanAfterRiichi(handBeforeKan, drawnTile, consumed) {
		t.Error("canConcealedKanAfterRiichi() = true, want false when waits change from 2m/3m/5m to 3m")
	}
}

func TestState_LegalActions_AfterRiichiAcceptedAllowsOnlyTsumogiri(t *testing.T) {
	hands := newValidHands()
	hands[0] = riichiReadyHandForTest()
	s := mustNewRoundStateForTest(t, hands)
	actor := seat.MustSeat(0)
	firstDraw := tile.MustTileFromCode("S")
	firstDiscard := tile.MustTileFromCode("W")
	if err := s.Apply(event.NewDraw(actor, firstDraw)); err != nil {
		t.Fatalf("Apply(first Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewRiichi(actor)); err != nil {
		t.Fatalf("Apply(Riichi) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(actor, firstDiscard, false)); err != nil {
		t.Fatalf("Apply(first Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewRiichiAccepted(actor, nil, nil)); err != nil {
		t.Fatalf("Apply(RiichiAccepted) failed: %v", err)
	}

	for i := 1; i < 4; i++ {
		other := seat.MustSeat(i)
		drawnTile := tile.MustTileFromCode("6m")
		if err := s.Apply(event.NewDraw(other, drawnTile)); err != nil {
			t.Fatalf("Apply(other Draw %d) failed: %v", i, err)
		}
		if err := s.Apply(event.NewDiscard(other, drawnTile, true)); err != nil {
			t.Fatalf("Apply(other Discard %d) failed: %v", i, err)
		}
	}

	secondDraw := tile.MustTileFromCode("7m")
	if err := s.Apply(event.NewDraw(actor, secondDraw)); err != nil {
		t.Fatalf("Apply(second Draw) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("LegalActions() length = %d, want 1: %v", len(got), got)
	}
	discard, ok := got[0].(*action.Discard)
	if !ok {
		t.Fatalf("LegalActions()[0] = %T, want *action.Discard", got[0])
	}
	if discard.Tile() != secondDraw {
		t.Errorf("Discard.Tile() = %v, want %v", discard.Tile(), secondDraw)
	}
	if !discard.Tsumogiri() {
		t.Error("Discard.Tsumogiri() = false, want true")
	}
}

func TestState_LegalActions_RiichiDeclarationTileKeepsTenpai(t *testing.T) {
	hands := newValidHands()
	hands[0] = riichiReadyHandForTest()
	s := mustNewRoundStateForTest(t, hands)
	actor := seat.MustSeat(0)
	if err := s.Apply(event.NewDraw(actor, tile.MustTileFromCode("S"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewRiichi(actor)); err != nil {
		t.Fatalf("Apply(Riichi) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if !containsDiscard(got, "W", false) {
		t.Error("LegalActions() does not contain W hand discard, want riichi declaration tile")
	}
	if containsDiscard(got, "1m", false) {
		t.Error("LegalActions() contains 1m hand discard, want discard that breaks tenpai excluded")
	}
}

func TestState_LegalActions_AfterChiiExcludesSwapCallTiles(t *testing.T) {
	hands := newValidHands()
	hands[1] = openChiiHandForLegalActionsTest()
	s := mustNewRoundStateForTest(t, hands)
	actor := seat.MustSeat(1)
	s = stateAfterChiiForLegalActionsTest(t, s)

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	for _, code := range []string{"2m", "5m", "5mr"} {
		if containsDiscard(got, code, false) {
			t.Errorf("LegalActions() contains %s hand discard, want swap-call tile excluded", code)
		}
	}
	if !containsDiscard(got, "1p", false) {
		t.Error("LegalActions() does not contain 1p hand discard, want non-swap-call tile")
	}
}
