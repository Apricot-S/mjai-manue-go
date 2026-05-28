package ai

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestBuildSelfTurnCandidates_BuildsDiscardCandidate(t *testing.T) {
	self := seat.MustSeat(0)
	discard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}

	got, err := buildSelfTurnCandidates([]action.Action{discard}, stubPlayerViewer{
		hand: hand.CodesToHand([]string{
			"1m", "2m", "3m", "4m", "5m", "6m", "7m",
			"1p", "2p", "3p", "E", "E", "S", "S",
		}),
		riichiState: player.NotRiichi,
		drawnTile:   nil,
	})
	if err != nil {
		t.Fatalf("buildSelfTurnCandidates() failed: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(buildSelfTurnCandidates()) = %d, want 1", len(got))
	}
	if got[0].traceKey != "-1.5m" {
		t.Errorf("traceKey = %q, want %q", got[0].traceKey, "-1.5m")
	}
	if got[0].riichi {
		t.Errorf("riichi = true, want false")
	}
	if got[0].action != discard {
		t.Errorf("action = %v, want original discard action", got[0].action)
	}
	if got[0].discardTile != discard.Tile() {
		t.Errorf("discardTile = %v, want %v", got[0].discardTile, discard.Tile())
	}
	if got[0].turnHand == nil {
		t.Fatal("turnHand = nil, want self-turn hand")
	}
	if got[0].turnHand.Count(tile.MustTileFromCode("5m")) != 1 {
		t.Errorf("turnHand 5m count = %d, want 1", got[0].turnHand.Count(tile.MustTileFromCode("5m")))
	}
	if got[0].afterDiscardHand == nil {
		t.Fatal("afterDiscardHand = nil, want hand after discard")
	}
	if got[0].afterDiscardHand.Count(tile.MustTileFromCode("5m")) != 0 {
		t.Errorf("afterDiscardHand 5m count = %d, want 0", got[0].afterDiscardHand.Count(tile.MustTileFromCode("5m")))
	}
}

func TestBuildSelfTurnCandidates_IgnoresConcealedAndPromotedKan(t *testing.T) {
	self := seat.MustSeat(0)
	discard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}
	concealedKan, err := action.NewConcealedKan(self, [4]tile.Tile{
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("5m"),
	})
	if err != nil {
		t.Fatalf("NewConcealedKan() failed: %v", err)
	}
	promotedKan, err := action.NewPromotedKan(self, tile.MustTileFromCode("7p"), [3]tile.Tile{
		tile.MustTileFromCode("7p"),
		tile.MustTileFromCode("7p"),
		tile.MustTileFromCode("7p"),
	})
	if err != nil {
		t.Fatalf("NewPromotedKan() failed: %v", err)
	}

	got, err := buildSelfTurnCandidates([]action.Action{discard, concealedKan, promotedKan}, stubPlayerViewer{
		hand: hand.CodesToHand([]string{
			"1m", "2m", "3m", "5m", "5m", "5m", "5m",
			"1p", "2p", "3p", "7p", "7p", "7p", "E",
		}),
		riichiState: player.NotRiichi,
	})
	if err != nil {
		t.Fatalf("buildSelfTurnCandidates() failed: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(buildSelfTurnCandidates()) = %d, want 1", len(got))
	}
	if got[0].action != discard {
		t.Errorf("action = %v, want discard action only", got[0].action)
	}
}

func TestNormalizedSelfTurnDiscards_PrefersTsumogiriForSameTile(t *testing.T) {
	self := seat.MustSeat(0)
	handDiscard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard(hand) failed: %v", err)
	}
	tsumogiriDiscard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), true)
	if err != nil {
		t.Fatalf("NewDiscard(tsumogiri) failed: %v", err)
	}
	redDiscard, err := action.NewDiscard(self, tile.MustTileFromCode("5mr"), false)
	if err != nil {
		t.Fatalf("NewDiscard(red) failed: %v", err)
	}

	got := normalizedSelfTurnDiscards([]action.Action{
		handDiscard,
		action.NewPass(self),
		redDiscard,
		tsumogiriDiscard,
	})
	if len(got) != 2 {
		t.Fatalf("len(normalizedSelfTurnDiscards()) = %d, want 2", len(got))
	}
	if got[0] != tsumogiriDiscard {
		t.Errorf("normalizedSelfTurnDiscards()[0] = %v, want tsumogiri discard", got[0])
	}
	if got[1] != redDiscard {
		t.Errorf("normalizedSelfTurnDiscards()[1] = %v, want red discard kept separately", got[1])
	}
}

func TestBuildSelfTurnCandidates_BuildsRiichiCandidate(t *testing.T) {
	self := seat.MustSeat(0)
	riichi := action.NewRiichi(self)
	discard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}

	got, err := buildSelfTurnCandidates([]action.Action{discard, riichi}, stubPlayerViewer{
		hand: hand.CodesToHand([]string{
			"1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m",
			"1p", "1p", "E", "E", "5m",
		}),
		riichiState: player.NotRiichi,
		drawnTile:   nil,
	})
	if err != nil {
		t.Fatalf("buildSelfTurnCandidates() failed: %v", err)
	}
	var gotRiichi *actionCandidate
	for i := range got {
		if got[i].riichi {
			gotRiichi = &got[i]
			break
		}
	}
	if gotRiichi == nil {
		t.Fatalf("buildSelfTurnCandidates() contains no riichi candidate")
	}
	if gotRiichi.traceKey != "0.5m" {
		t.Errorf("traceKey = %q, want %q", gotRiichi.traceKey, "0.5m")
	}
	if gotRiichi.action != riichi {
		t.Errorf("action = %v, want riichi action", gotRiichi.action)
	}
	if gotRiichi.discardTile != discard.Tile() {
		t.Errorf("discardTile = %v, want %v", gotRiichi.discardTile, discard.Tile())
	}
}

func TestBuildSelfTurnCandidates_FiltersNonTenpaiRiichiCandidate(t *testing.T) {
	self := seat.MustSeat(0)
	riichi := action.NewRiichi(self)
	discard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}

	got, err := buildSelfTurnCandidates([]action.Action{discard, riichi}, stubPlayerViewer{
		hand: hand.CodesToHand([]string{
			"1m", "2m", "3m", "4m", "5m", "6m", "7m",
			"1p", "2p", "3p", "E", "E", "S", "S",
		}),
		riichiState: player.NotRiichi,
		drawnTile:   nil,
	})
	if err != nil {
		t.Fatalf("buildSelfTurnCandidates() failed: %v", err)
	}
	for _, candidate := range got {
		if candidate.riichi {
			t.Fatalf("buildSelfTurnCandidates() contains riichi candidate unexpectedly: %+v", candidate)
		}
	}
}

func TestBuildSelfTurnCandidates_RiichiDeclaredScoresLegalActionsAsRiichi(t *testing.T) {
	self := seat.MustSeat(0)
	tenpaiDiscard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard(tenpai) failed: %v", err)
	}
	nonTenpaiDiscard, err := action.NewDiscard(self, tile.MustTileFromCode("1p"), false)
	if err != nil {
		t.Fatalf("NewDiscard(non-tenpai) failed: %v", err)
	}

	got, err := buildSelfTurnCandidates([]action.Action{tenpaiDiscard, nonTenpaiDiscard}, stubPlayerViewer{
		hand: hand.CodesToHand([]string{
			"1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m",
			"1p", "1p", "E", "E", "5m",
		}),
		riichiState: player.RiichiDeclared,
		drawnTile:   nil,
	})
	if err != nil {
		t.Fatalf("buildSelfTurnCandidates() failed: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(buildSelfTurnCandidates()) = %d, want 2", len(got))
	}
	wantActions := map[action.Action]bool{
		tenpaiDiscard:    false,
		nonTenpaiDiscard: false,
	}
	for _, candidate := range got {
		if candidate.riichi {
			t.Errorf("riichi = true for %q, want false because declaration action is already done", candidate.traceKey)
		}
		if !candidate.scoreAsRiichi {
			t.Errorf("scoreAsRiichi = false for %q, want true", candidate.traceKey)
		}
		if _, ok := wantActions[candidate.action]; !ok {
			t.Errorf("action = %v, want one of legal-action inputs", candidate.action)
			continue
		}
		wantActions[candidate.action] = true
	}
	for action, found := range wantActions {
		if !found {
			t.Errorf("buildSelfTurnCandidates() did not include %v", action)
		}
	}
}

func TestBuildSelfTurnCandidates_IncludesRiichiAndDiscardCandidates(t *testing.T) {
	self := seat.MustSeat(0)
	riichi := action.NewRiichi(self)
	discard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}

	got, err := buildSelfTurnCandidates([]action.Action{discard, riichi}, stubPlayerViewer{
		hand: hand.CodesToHand([]string{
			"1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m",
			"1p", "1p", "E", "E", "5m",
		}),
		riichiState: player.NotRiichi,
		drawnTile:   nil,
	})
	if err != nil {
		t.Fatalf("buildSelfTurnCandidates() failed: %v", err)
	}
	wantTraceKeys := map[string]bool{"0.5m": false, "-1.5m": false}
	for _, candidate := range got {
		if _, ok := wantTraceKeys[candidate.traceKey]; ok {
			wantTraceKeys[candidate.traceKey] = true
		}
	}
	for traceKey, found := range wantTraceKeys {
		if !found {
			t.Errorf("buildSelfTurnCandidates() does not contain %q", traceKey)
		}
	}
}
