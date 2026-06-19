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

func TestState_LegalActions_OnOtherDiscardExcludesRonWhenFuriten(t *testing.T) {
	hands := newValidHands()
	hands[2] = ronWithTanyaoHandForLegalActionsTest()
	s := mustNewRoundStateForTest(t, hands)
	actor := seat.MustSeat(2)
	winningTile := tile.MustTileFromCode("3p")
	p := s.players[actor.Index()].(*player.VisiblePlayer)
	p.AddExtraSafeTiles(winningTile)
	if !p.IsFuriten() {
		t.Fatal("test setup failed: actor is not furiten")
	}

	target := seat.MustSeat(0)
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
	if containsWin(got, actor, target, "3p") {
		t.Error("LegalActions() contains Win, want furiten ron excluded")
	}
	if containsPass(got, actor) {
		t.Error("LegalActions() contains Pass, want no pass when furiten ron is excluded")
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

func TestState_LegalActions_OnOtherDiscardIncludesPon(t *testing.T) {
	hands := newValidHands()
	hands[1] = ponHandForLegalActionsTest("E", "E")
	s := mustNewRoundStateForTest(t, hands)
	target := seat.MustSeat(0)
	actor := seat.MustSeat(1)
	taken := tile.MustTileFromCode("E")
	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if !containsPon(got, actor, target, "E", [2]string{"E", "E"}) {
		t.Error("LegalActions() does not contain Pon, want pon with two matching tiles")
	}
	if !containsPass(got, actor) {
		t.Error("LegalActions() does not contain Pass, want pass when pon is available")
	}
}

func TestState_LegalActions_OnOtherDiscardExcludesPonWithOneMatchingTile(t *testing.T) {
	hands := newValidHands()
	hands[1] = ponHandForLegalActionsTest("E", "S")
	s := mustNewRoundStateForTest(t, hands)
	target := seat.MustSeat(0)
	actor := seat.MustSeat(1)
	taken := tile.MustTileFromCode("E")
	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if containsPon(got, actor, target, "E", [2]string{"E", "E"}) {
		t.Error("LegalActions() contains Pon, want pon excluded with only one matching tile")
	}
}

func TestState_LegalActions_OnOtherDiscardIncludesPonRedFiveChoices(t *testing.T) {
	hands := newValidHands()
	hands[1] = ponHandForLegalActionsTest("5m", "5mr")
	s := mustNewRoundStateForTest(t, hands)
	target := seat.MustSeat(0)
	actor := seat.MustSeat(1)
	taken := tile.MustTileFromCode("5m")
	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if !containsPon(got, actor, target, "5m", [2]string{"5m", "5mr"}) {
		t.Error("LegalActions() does not contain Pon, want red-five consumed choice")
	}
}

func TestState_LegalActions_OnOtherDiscardIncludesTwoPonChoicesWithRedFive(t *testing.T) {
	hands := newValidHands()
	hands[1] = ponHandWithThreeMatchingTilesForLegalActionsTest("5m", "5m", "5mr")
	s := mustNewRoundStateForTest(t, hands)
	target := seat.MustSeat(0)
	actor := seat.MustSeat(1)
	taken := tile.MustTileFromCode("5m")
	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if !containsPon(got, actor, target, "5m", [2]string{"5m", "5m"}) {
		t.Error("LegalActions() does not contain Pon with two normal fives")
	}
	if !containsPon(got, actor, target, "5m", [2]string{"5m", "5mr"}) {
		t.Error("LegalActions() does not contain Pon with normal and red five")
	}
}

func TestState_LegalActions_OnOtherDiscardExcludesPonAfterRiichiAccepted(t *testing.T) {
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
	if err := s.Apply(event.NewDiscard(actor, tile.MustTileFromCode("W"), false)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewRiichiAccepted(actor, nil, nil)); err != nil {
		t.Fatalf("Apply(RiichiAccepted) failed: %v", err)
	}

	target := seat.MustSeat(1)
	taken := tile.MustTileFromCode("E")
	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if containsPon(got, actor, target, "E", [2]string{"E", "E"}) {
		t.Error("LegalActions() contains Pon, want pon excluded after riichi accepted")
	}
}

func TestState_LegalActions_AfterRiichiAcceptedReturnsNoActions(t *testing.T) {
	hands := newValidHands()
	hands[0] = riichiReadyHandForTest()
	hands[3] = ponHandForLegalActionsTest("W", "W")
	s := mustNewRoundStateForTest(t, hands)
	riichiActor := seat.MustSeat(0)
	declarationTile := tile.MustTileFromCode("W")

	if err := s.Apply(event.NewDraw(riichiActor, tile.MustTileFromCode("S"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewRiichi(riichiActor)); err != nil {
		t.Fatalf("Apply(Riichi) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(riichiActor, declarationTile, false)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewRiichiAccepted(riichiActor, nil, nil)); err != nil {
		t.Fatalf("Apply(RiichiAccepted) failed: %v", err)
	}

	for i := range common.NumPlayers {
		actor := seat.MustSeat(i)
		got, err := s.LegalActions(actor)
		if err != nil {
			t.Fatalf("LegalActions(%d) failed: %v", i, err)
		}
		if len(got) != 0 {
			t.Errorf("LegalActions(%d) = %v, want empty after RiichiAccepted", i, got)
		}
	}
}

func TestState_LegalActions_OnOtherDiscardIncludesCalledKan(t *testing.T) {
	hands := newValidHands()
	hands[1] = calledKanHandForLegalActionsTest("E", "E", "E")
	s := mustNewRoundStateForTest(t, hands)
	target := seat.MustSeat(0)
	actor := seat.MustSeat(1)
	taken := tile.MustTileFromCode("E")
	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if !containsCalledKan(got, actor, target, "E", [3]string{"E", "E", "E"}) {
		t.Error("LegalActions() does not contain CalledKan, want daiminkan with three matching tiles")
	}
	if !containsPon(got, actor, target, "E", [2]string{"E", "E"}) {
		t.Error("LegalActions() does not contain Pon, want pon when daiminkan is available")
	}
	if !containsPass(got, actor) {
		t.Error("LegalActions() does not contain Pass, want pass when called kan is available")
	}
}

func TestState_LegalActions_OnOtherDiscardExcludesCalledKanWithTwoMatchingTiles(t *testing.T) {
	hands := newValidHands()
	hands[1] = ponHandForLegalActionsTest("E", "E")
	s := mustNewRoundStateForTest(t, hands)
	target := seat.MustSeat(0)
	actor := seat.MustSeat(1)
	taken := tile.MustTileFromCode("E")
	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if containsCalledKan(got, actor, target, "E", [3]string{"E", "E", "E"}) {
		t.Error("LegalActions() contains CalledKan, want called kan excluded with only two matching tiles")
	}
}

func TestState_LegalActions_OnOtherDiscardIncludesCalledKanRedFive(t *testing.T) {
	hands := newValidHands()
	hands[1] = calledKanHandForLegalActionsTest("5m", "5m", "5mr")
	s := mustNewRoundStateForTest(t, hands)
	target := seat.MustSeat(0)
	actor := seat.MustSeat(1)
	taken := tile.MustTileFromCode("5m")
	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if !containsCalledKan(got, actor, target, "5m", [3]string{"5m", "5m", "5mr"}) {
		t.Error("LegalActions() does not contain CalledKan, want daiminkan containing red five")
	}
}

func TestState_LegalActions_OnOtherDiscardExcludesCalledKanOnFifthKan(t *testing.T) {
	hands := newValidHands()
	hands[1] = calledKanHandForLegalActionsTest("E", "E", "E")
	s := mustNewRoundStateForTest(t, hands)
	s.numKans = maxNumKan
	target := seat.MustSeat(0)
	actor := seat.MustSeat(1)
	taken := tile.MustTileFromCode("E")
	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if containsCalledKan(got, actor, target, "E", [3]string{"E", "E", "E"}) {
		t.Error("LegalActions() contains CalledKan, want fifth kan excluded")
	}
}

func TestState_LegalActions_OnOtherDiscardIncludesChii(t *testing.T) {
	hands := newValidHands()
	hands[1] = chiiHandForLegalActionsTest("2m", "3m")
	s := mustNewRoundStateForTest(t, hands)
	target := seat.MustSeat(0)
	actor := seat.MustSeat(1)
	taken := tile.MustTileFromCode("4m")
	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if !containsChii(got, actor, target, "4m", [2]string{"2m", "3m"}) {
		t.Error("LegalActions() does not contain Chii, want chii from kamicha discard")
	}
	if !containsPass(got, actor) {
		t.Error("LegalActions() does not contain Pass, want pass when chii is available")
	}
}

func TestState_LegalActions_OnOtherDiscardExcludesChiiFromNonKamicha(t *testing.T) {
	hands := newValidHands()
	hands[2] = chiiHandForLegalActionsTest("2m", "3m")
	s := mustNewRoundStateForTest(t, hands)
	target := seat.MustSeat(0)
	actor := seat.MustSeat(2)
	taken := tile.MustTileFromCode("4m")
	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if containsChii(got, actor, target, "4m", [2]string{"2m", "3m"}) {
		t.Error("LegalActions() contains Chii, want chii excluded from non-kamicha discard")
	}
}

func TestState_LegalActions_OnOtherDiscardIncludesChiiRedFiveChoices(t *testing.T) {
	hands := newValidHands()
	hands[1] = chiiHandWithRedFiveForLegalActionsTest()
	s := mustNewRoundStateForTest(t, hands)
	target := seat.MustSeat(0)
	actor := seat.MustSeat(1)
	taken := tile.MustTileFromCode("4m")
	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if !containsChii(got, actor, target, "4m", [2]string{"3m", "5m"}) {
		t.Error("LegalActions() does not contain Chii with normal five")
	}
	if !containsChii(got, actor, target, "4m", [2]string{"3m", "5mr"}) {
		t.Error("LegalActions() does not contain Chii with red five")
	}
}

func TestState_LegalActions_OnOtherDiscardExcludesChiiWhenAllRemainingTilesAreSwapCallTilesAfterTwoMelds(t *testing.T) {
	actor := seat.MustSeat(1)
	target := seat.MustSeat(0)
	players := [common.NumPlayers]player.Player{
		player.NewInvisiblePlayer(),
		openPlayerForChiiSwapCallLegalActionsTest(t, []tile.Tile{
			tile.MustTileFromCode("1m"), tile.MustTileFromCode("1m"), tile.MustTileFromCode("1m"),
			tile.MustTileFromCode("2m"), tile.MustTileFromCode("3m"), tile.MustTileFromCode("4m"), tile.MustTileFromCode("4m"),
		}, 2),
		player.NewInvisiblePlayer(),
		player.NewInvisiblePlayer(),
	}
	s := newStateForOtherDiscardLegalActionsTest(players)
	taken := tile.MustTileFromCode("1m")
	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if containsChii(got, actor, target, "1m", [2]string{"2m", "3m"}) {
		t.Error("LegalActions() contains Chii, want chii excluded when all remaining tiles are swap-call tiles")
	}
	if containsChii(got, actor, target, "4m", [2]string{"2m", "3m"}) {
		t.Error("LegalActions() contains Chii, want chii excluded when all remaining tiles are swap-call tiles")
	}
}

func TestState_LegalActions_OnOtherDiscardExcludesChiiWhenAllRemainingTilesAreSwapCallTilesAfterThreeMelds(t *testing.T) {
	actor := seat.MustSeat(1)
	target := seat.MustSeat(0)
	players := [common.NumPlayers]player.Player{
		player.NewInvisiblePlayer(),
		openPlayerForChiiSwapCallLegalActionsTest(t, []tile.Tile{
			tile.MustTileFromCode("2m"), tile.MustTileFromCode("3m"), tile.MustTileFromCode("3m"), tile.MustTileFromCode("4m"),
		}, 3),
		player.NewInvisiblePlayer(),
		player.NewInvisiblePlayer(),
	}
	s := newStateForOtherDiscardLegalActionsTest(players)
	taken := tile.MustTileFromCode("3m")
	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if containsChii(got, actor, target, "3m", [2]string{"2m", "4m"}) {
		t.Error("LegalActions() contains Chii, want chii excluded when all remaining tiles are swap-call tiles")
	}
}

func TestState_LegalActions_OnOtherDiscardIncludesMaxActions(t *testing.T) {
	hands := newValidHands()
	hands[1] = maxOtherDiscardActionsHandForLegalActionsTest()
	s := mustNewRoundStateForTest(t, hands)
	target := seat.MustSeat(0)
	actor := seat.MustSeat(1)
	taken := tile.MustTileFromCode("4m")
	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if len(got) != maxNumActionsOnOtherDiscard {
		t.Fatalf("LegalActions() length = %d, want %d: %v", len(got), maxNumActionsOnOtherDiscard, got)
	}
	if !containsWin(got, actor, target, "4m") {
		t.Error("LegalActions() does not contain Win")
	}
	if !containsPon(got, actor, target, "4m", [2]string{"4m", "4m"}) {
		t.Error("LegalActions() does not contain Pon")
	}
	if !containsCalledKan(got, actor, target, "4m", [3]string{"4m", "4m", "4m"}) {
		t.Error("LegalActions() does not contain CalledKan")
	}
	for _, consumed := range [][2]string{
		{"2m", "3m"},
		{"3m", "5m"},
		{"3m", "5mr"},
		{"5m", "6m"},
		{"5mr", "6m"},
	} {
		if !containsChii(got, actor, target, "4m", consumed) {
			t.Errorf("LegalActions() does not contain Chii consumed=%v", consumed)
		}
	}
	if !containsPass(got, actor) {
		t.Error("LegalActions() does not contain Pass")
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

func TestState_LegalActions_OnRobbingKanExcludesRonWhenFuriten(t *testing.T) {
	s := newStateBeforeRobbingKanForLegalActionsTest(t)
	target := seat.MustSeat(3)
	actor := seat.MustSeat(1)
	added := tile.MustTileFromCode("E")
	p := s.players[actor.Index()].(*player.VisiblePlayer)
	p.AddExtraSafeTiles(added)
	if !p.IsFuriten() {
		t.Fatal("test setup failed: actor is not furiten")
	}
	if err := s.Apply(event.NewPromotedKan(target, added, [3]tile.Tile{added, added, added})); err != nil {
		t.Fatalf("Apply(PromotedKan) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if containsWin(got, actor, target, "E") {
		t.Error("LegalActions() contains Win, want furiten robbing-a-kan ron excluded")
	}
	if containsPass(got, actor) {
		t.Error("LegalActions() contains Pass, want no pass when furiten robbing-a-kan ron is excluded")
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

func ponHandForLegalActionsTest(firstCode, secondCode string) [common.InitHandSize]tile.Tile {
	return [common.InitHandSize]tile.Tile{
		tile.MustTileFromCode(firstCode),
		tile.MustTileFromCode(secondCode),
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("2m"),
		tile.MustTileFromCode("3m"),
		tile.MustTileFromCode("1p"),
		tile.MustTileFromCode("2p"),
		tile.MustTileFromCode("3p"),
		tile.MustTileFromCode("1s"),
		tile.MustTileFromCode("2s"),
		tile.MustTileFromCode("3s"),
		tile.MustTileFromCode("9p"),
		tile.MustTileFromCode("9s"),
	}
}

func ponHandWithThreeMatchingTilesForLegalActionsTest(firstCode, secondCode, thirdCode string) [common.InitHandSize]tile.Tile {
	return [common.InitHandSize]tile.Tile{
		tile.MustTileFromCode(firstCode),
		tile.MustTileFromCode(secondCode),
		tile.MustTileFromCode(thirdCode),
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("2m"),
		tile.MustTileFromCode("3m"),
		tile.MustTileFromCode("1p"),
		tile.MustTileFromCode("2p"),
		tile.MustTileFromCode("3p"),
		tile.MustTileFromCode("1s"),
		tile.MustTileFromCode("2s"),
		tile.MustTileFromCode("3s"),
		tile.MustTileFromCode("9p"),
	}
}

func calledKanHandForLegalActionsTest(firstCode, secondCode, thirdCode string) [common.InitHandSize]tile.Tile {
	return [common.InitHandSize]tile.Tile{
		tile.MustTileFromCode(firstCode),
		tile.MustTileFromCode(secondCode),
		tile.MustTileFromCode(thirdCode),
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("2m"),
		tile.MustTileFromCode("3m"),
		tile.MustTileFromCode("1p"),
		tile.MustTileFromCode("2p"),
		tile.MustTileFromCode("3p"),
		tile.MustTileFromCode("1s"),
		tile.MustTileFromCode("2s"),
		tile.MustTileFromCode("3s"),
		tile.MustTileFromCode("9p"),
	}
}

func chiiHandForLegalActionsTest(firstCode, secondCode string) [common.InitHandSize]tile.Tile {
	return [common.InitHandSize]tile.Tile{
		tile.MustTileFromCode(firstCode),
		tile.MustTileFromCode(secondCode),
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("7m"),
		tile.MustTileFromCode("8m"),
		tile.MustTileFromCode("9m"),
		tile.MustTileFromCode("1p"),
		tile.MustTileFromCode("2p"),
		tile.MustTileFromCode("3p"),
		tile.MustTileFromCode("1s"),
		tile.MustTileFromCode("2s"),
		tile.MustTileFromCode("3s"),
	}
}

func chiiHandWithRedFiveForLegalActionsTest() [common.InitHandSize]tile.Tile {
	return [common.InitHandSize]tile.Tile{
		tile.MustTileFromCode("3m"),
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("5mr"),
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("7m"),
		tile.MustTileFromCode("8m"),
		tile.MustTileFromCode("9m"),
		tile.MustTileFromCode("1p"),
		tile.MustTileFromCode("2p"),
		tile.MustTileFromCode("3p"),
		tile.MustTileFromCode("1s"),
		tile.MustTileFromCode("2s"),
	}
}

func maxOtherDiscardActionsHandForLegalActionsTest() [common.InitHandSize]tile.Tile {
	return [common.InitHandSize]tile.Tile{
		tile.MustTileFromCode("2m"),
		tile.MustTileFromCode("3m"),
		tile.MustTileFromCode("4m"),
		tile.MustTileFromCode("4m"),
		tile.MustTileFromCode("4m"),
		tile.MustTileFromCode("5mr"),
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("6m"),
		tile.MustTileFromCode("P"),
		tile.MustTileFromCode("P"),
		tile.MustTileFromCode("P"),
	}
}

func openPlayerForChiiSwapCallLegalActionsTest(t *testing.T, finalHand []tile.Tile, numMelds int) player.Player {
	t.Helper()

	callData := []struct {
		taken    tile.Tile
		consumed [2]tile.Tile
		discard  tile.Tile
	}{
		{
			taken:    tile.MustTileFromCode("7p"),
			consumed: [2]tile.Tile{tile.MustTileFromCode("8p"), tile.MustTileFromCode("9p")},
			discard:  tile.MustTileFromCode("E"),
		},
		{
			taken:    tile.MustTileFromCode("1s"),
			consumed: [2]tile.Tile{tile.MustTileFromCode("2s"), tile.MustTileFromCode("3s")},
			discard:  tile.MustTileFromCode("S"),
		},
		{
			taken:    tile.MustTileFromCode("7s"),
			consumed: [2]tile.Tile{tile.MustTileFromCode("8s"), tile.MustTileFromCode("9s")},
			discard:  tile.MustTileFromCode("W"),
		},
	}

	tiles := append([]tile.Tile{}, finalHand...)
	for i := range numMelds {
		tiles = append(tiles, callData[i].consumed[:]...)
		tiles = append(tiles, callData[i].discard)
	}
	if len(tiles) != common.InitHandSize {
		t.Fatalf("test setup hand size = %d, want %d", len(tiles), common.InitHandSize)
	}

	var handTiles [common.InitHandSize]tile.Tile
	copy(handTiles[:], tiles)
	p, err := player.NewVisiblePlayer(handTiles)
	if err != nil {
		t.Fatalf("player.NewVisiblePlayer() failed: %v", err)
	}
	for i := range numMelds {
		chii := meld.MustChii(callData[i].taken, callData[i].consumed, seat.MustSeat(0))
		if err := p.Chii(*chii); err != nil {
			t.Fatalf("Chii(%d) failed: %v", i, err)
		}
		if err := p.Discard(callData[i].discard, false); err != nil {
			t.Fatalf("Discard(%d) failed: %v", i, err)
		}
	}
	return p
}

func newStateForOtherDiscardLegalActionsTest(players [common.NumPlayers]player.Player) *State {
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
