package ai

import (
	"fmt"
	"math/rand/v2"
	"reflect"
	"strings"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service/block"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

func TestManueAgent_DecideActionSkeleton(t *testing.T) {
	self := seat.MustSeat(0)
	drawnTile := tile.MustTileFromCode("5p")
	handDiscard, err := action.NewDiscard(self, tile.MustTileFromCode("1m"), false)
	if err != nil {
		t.Fatalf("NewDiscard(hand) failed: %v", err)
	}
	tsumogiriDiscard, err := action.NewDiscard(self, drawnTile, true)
	if err != nil {
		t.Fatalf("NewDiscard(tsumogiri) failed: %v", err)
	}
	win, err := action.NewWin(self, self, drawnTile)
	if err != nil {
		t.Fatalf("NewWin() failed: %v", err)
	}
	riichi := action.NewRiichi(self)
	pass := action.NewPass(self)

	tests := []struct {
		name        string
		actions     []action.Action
		hand        *hand.VisibleHand
		riichiState player.RiichiState
		drawnTile   *tile.Tile
		want        action.Action
		decide      func(*ManueAgent, []action.Action, player.PlayerViewer) (Decision, error)
	}{
		{
			name:        "win first",
			actions:     []action.Action{handDiscard, win, pass},
			riichiState: player.NotRiichi,
			drawnTile:   &drawnTile,
			want:        win,
			decide: func(agent *ManueAgent, actions []action.Action, self player.PlayerViewer) (Decision, error) {
				if win := firstActionOfType[*action.Win](actions); win != nil {
					return Decision{Action: win}, nil
				}
				t.Fatal("firstActionOfType[*action.Win]() returned nil")
				return Decision{}, nil
			},
		},
		{
			name:        "riichi accepted tsumogiri",
			actions:     []action.Action{handDiscard, tsumogiriDiscard},
			riichiState: player.RiichiAccepted,
			drawnTile:   &drawnTile,
			want:        tsumogiriDiscard,
			decide: func(agent *ManueAgent, actions []action.Action, self player.PlayerViewer) (Decision, error) {
				return agent.decideSelfTurn(actions, nil, seat.MustSeat(0), self)
			},
		},
		{
			name:        "riichi before discard",
			actions:     []action.Action{handDiscard, riichi},
			hand:        hand.CodesToHand([]string{"1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m", "1p", "1p", "E", "E", "1m"}),
			riichiState: player.NotRiichi,
			drawnTile:   nil,
			want:        riichi,
			decide: func(agent *ManueAgent, actions []action.Action, self player.PlayerViewer) (Decision, error) {
				return agent.decideSelfTurn(actions, nil, seat.MustSeat(0), self)
			},
		},
	}

	agent := newTestManueAgent(t, 0)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decision, err := tt.decide(agent, tt.actions, stubPlayerViewer{
				hand:        tt.hand,
				riichiState: tt.riichiState,
				drawnTile:   tt.drawnTile,
			})
			if err != nil {
				t.Fatalf("decide failed: %v", err)
			}
			if decision.Action != tt.want {
				t.Errorf("Action = %T %[1]v, want %T %[2]v", decision.Action, tt.want)
			}
		})
	}
}

func TestNewManueAgent(t *testing.T) {
	stats := validStubManueStats()
	agent, err := NewManueAgent(123, ManueAgentDeps{
		Stats:  stats,
		Danger: stubDangerEstimator{},
	})
	if err != nil {
		t.Fatalf("NewManueAgent() failed: %v", err)
	}
	if agent.seed != 123 {
		t.Errorf("seed = %d, want 123", agent.seed)
	}
	if agent.deps.Stats == nil {
		t.Fatalf("deps.Stats = nil, want stats")
	}
	if got := agent.deps.Stats.NumWins(); got != stats.numWins {
		t.Errorf("deps.Stats.NumWins() = %d, want %d", got, stats.numWins)
	}
	if agent.rng == nil {
		t.Errorf("rng = nil, want initialized rng")
	}
}

func TestManueAgent_decideSelfTurn_ReturnsOriginalStyleActionLog(t *testing.T) {
	self := seat.MustSeat(0)
	discard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}

	decision, err := newTestManueAgent(t, 0).decideSelfTurn([]action.Action{discard}, nil, seat.MustSeat(0), stubPlayerViewer{
		hand:        hand.CodesToHand([]string{"1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m", "1p", "1p", "E", "E", "5m"}),
		riichiState: player.NotRiichi,
		drawnTile:   nil,
	})
	if err != nil {
		t.Fatalf("decideSelfTurn() failed: %v", err)
	}
	if !strings.Contains(decision.Log, "| action |") {
		t.Errorf("Log = %q, want metrics table", decision.Log)
	}
	if strings.Contains(decision.Log, "decidedKey") {
		t.Errorf("Log = %q, should not contain decidedKey", decision.Log)
	}
	if !strings.HasSuffix(decision.Log, "\n\n\n") {
		t.Errorf("Log = %q, want original blank-line suffix", decision.Log)
	}
	if !strings.Contains(decision.Trace, "decidedKey -1.5m\n") {
		t.Errorf("Trace = %q, want decidedKey", decision.Trace)
	}
}

func TestManueAgent_decideSelfTurn_ReturnsErrorWithoutTsumogiriAfterRiichiAccepted(t *testing.T) {
	self := seat.MustSeat(0)
	handDiscard, err := action.NewDiscard(self, tile.MustTileFromCode("1m"), false)
	if err != nil {
		t.Fatalf("NewDiscard(hand) failed: %v", err)
	}
	drawnTile := tile.MustTileFromCode("5p")

	_, err = newTestManueAgent(t, 0).decideSelfTurn([]action.Action{handDiscard}, nil, seat.MustSeat(0), stubPlayerViewer{
		riichiState: player.RiichiAccepted,
		drawnTile:   &drawnTile,
	})
	if err == nil {
		t.Fatal("selectAction() succeeded unexpectedly")
	}
}

func TestManueAgent_chooseBestCandidate_PrefersBlackTileOnTie(t *testing.T) {
	self := seat.MustSeat(0)
	redDiscard, err := action.NewDiscard(self, tile.MustTileFromCode("5mr"), false)
	if err != nil {
		t.Fatalf("NewDiscard(red) failed: %v", err)
	}
	blackDiscard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard(black) failed: %v", err)
	}

	candidates, err := getSelfTurnCandidates([]action.Action{redDiscard, blackDiscard}, stubPlayerViewer{
		hand: hand.CodesToHand([]string{
			"1m", "2m", "3m", "4m", "5m", "5mr", "6m", "7m",
			"1p", "2p", "3p", "E", "E", "S",
		}),
		riichiState: player.NotRiichi,
		drawnTile:   nil,
	})
	if err != nil {
		t.Fatalf("getSelfTurnCandidates() failed: %v", err)
	}
	got := chooseBestCandidate(candidates, true)
	if got.action != blackDiscard {
		t.Errorf("chooseBestCandidate() = %v, want black discard", got)
	}
}

func TestManueAgent_getSelfTurnCandidates_BuildsDiscardCandidate(t *testing.T) {
	self := seat.MustSeat(0)
	discard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}

	got, err := getSelfTurnCandidates([]action.Action{discard}, stubPlayerViewer{
		hand: hand.CodesToHand([]string{
			"1m", "2m", "3m", "4m", "5m", "6m", "7m",
			"1p", "2p", "3p", "E", "E", "S", "S",
		}),
		riichiState: player.NotRiichi,
		drawnTile:   nil,
	})
	if err != nil {
		t.Fatalf("getSelfTurnCandidates() failed: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(getSelfTurnCandidates()) = %d, want 1", len(got))
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

func TestCandidateShantenUsesThrowableVector(t *testing.T) {
	got := candidateShanten(
		tile.MustTileFromCode("1m"),
		3,
		[]service.Goal{
			{Shanten: 2, ThrowableVector: hand.TileCounts34{0: 1}},
			{Shanten: 1, ThrowableVector: hand.TileCounts34{0: 1}},
			{Shanten: 0, ThrowableVector: hand.TileCounts34{1: 1}},
		},
	)
	if got != 1 {
		t.Errorf("candidateShanten() = %d, want 1", got)
	}
}

func TestCandidateShantenReturnsBaseForNone(t *testing.T) {
	got := candidateShanten(tile.MustTileFromCode("?"), 3, nil)
	if got != 3 {
		t.Errorf("candidateShanten(?) = %d, want 3", got)
	}
}

func TestCandidateShantenReturnsInfinityWhenTileIsNotThrowable(t *testing.T) {
	got := candidateShanten(
		tile.MustTileFromCode("1m"),
		3,
		[]service.Goal{{Shanten: 0, ThrowableVector: hand.TileCounts34{1: 1}}},
	)
	if got != service.InfinityShanten {
		t.Errorf("candidateShanten() = %d, want InfinityShanten", got)
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

func TestGetOtherDiscardReactionCandidates_BuildsPassAndCallDiscards(t *testing.T) {
	self := seat.MustSeat(0)
	target := seat.MustSeat(3)
	pon, err := action.NewPon(self, target, tile.MustTileFromCode("5p"), [2]tile.Tile{
		tile.MustTileFromCode("5p"),
		tile.MustTileFromCode("5p"),
	})
	if err != nil {
		t.Fatalf("NewPon() failed: %v", err)
	}
	pass := action.NewPass(self)

	got, err := getOtherDiscardReactionCandidates([]action.Action{pass, pon}, stubPlayerViewer{
		hand: hand.CodesToHand([]string{
			"1m", "2m", "3m", "4m", "5m", "6m", "7m",
			"1p", "2p", "5p", "5p", "5pr", "E",
		}),
		riichiState: player.NotRiichi,
	})
	if err != nil {
		t.Fatalf("getOtherDiscardReactionCandidates() failed: %v", err)
	}

	found := map[string]bool{"none": false, "0.4m": false}
	for _, candidate := range got {
		if _, ok := found[candidate.traceKey]; ok {
			found[candidate.traceKey] = true
		}
		if candidate.traceKey == "0.5pr" {
			t.Errorf("getOtherDiscardReactionCandidates() included kuikae discard %q", candidate.traceKey)
		}
		if strings.HasPrefix(candidate.traceKey, "0.") && candidate.action != pon {
			t.Errorf("call candidate action = %v, want pon", candidate.action)
		}
	}
	for traceKey, ok := range found {
		if !ok {
			t.Errorf("getOtherDiscardReactionCandidates() does not contain %q", traceKey)
		}
	}
}

func TestManueAgent_decideOtherDiscardReaction_EvaluatesCallCandidates(t *testing.T) {
	self := seat.MustSeat(0)
	target := seat.MustSeat(3)
	pon, err := action.NewPon(self, target, tile.MustTileFromCode("5p"), [2]tile.Tile{
		tile.MustTileFromCode("5p"),
		tile.MustTileFromCode("5p"),
	})
	if err != nil {
		t.Fatalf("NewPon() failed: %v", err)
	}
	pass := action.NewPass(self)
	selfPlayer := stubPlayerViewer{
		hand: hand.CodesToHand([]string{
			"1m", "2m", "3m", "4m", "5m", "6m", "7m",
			"1p", "2p", "5p", "5p", "5pr", "E",
		}),
		riichiState: player.NotRiichi,
	}
	state := stubCandidateEvaluationStateViewer{
		turn:         1,
		roundWind:    wind.East,
		seatWinds:    [common.NumPlayers]wind.Wind{wind.East, wind.South, wind.West, wind.North},
		dealer:       self,
		scores:       [common.NumPlayers]int{25000, 25000, 25000, 25000},
		startingSeat: self,
		players: [common.NumPlayers]player.PlayerViewer{
			selfPlayer,
			stubPlayerViewer{},
			stubPlayerViewer{},
			stubPlayerViewer{},
		},
	}

	decision, err := newTestManueAgent(t, 0).decideOtherDiscardReaction(
		[]action.Action{pass, pon},
		state,
		self,
		selfPlayer,
	)
	if err != nil {
		t.Fatalf("decideOtherDiscardReaction() failed: %v", err)
	}

	if decision.Action != pass && decision.Action != pon {
		t.Errorf("Action = %T %[1]v, want pass or pon", decision.Action)
	}
	if decision.Log == "" {
		t.Fatal("Log is empty; reaction candidates were not evaluated")
	}
	if !strings.Contains(decision.Log, "none") {
		t.Errorf("Log = %q, want pass candidate", decision.Log)
	}
	if !strings.Contains(decision.Log, "0.") {
		t.Errorf("Log = %q, want pon discard candidates", decision.Log)
	}
	if !strings.Contains(decision.Trace, "decidedKey ") {
		t.Errorf("Trace = %q, want decidedKey suffix", decision.Trace)
	}
}

func TestManueAgent_getSelfTurnCandidates_BuildsRiichiCandidate(t *testing.T) {
	self := seat.MustSeat(0)
	riichi := action.NewRiichi(self)
	discard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}

	got, err := getSelfTurnCandidates([]action.Action{discard, riichi}, stubPlayerViewer{
		hand: hand.CodesToHand([]string{
			"1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m",
			"1p", "1p", "E", "E", "5m",
		}),
		riichiState: player.NotRiichi,
		drawnTile:   nil,
	})
	if err != nil {
		t.Fatalf("getSelfTurnCandidates() failed: %v", err)
	}
	var gotRiichi *actionCandidate
	for i := range got {
		if got[i].riichi {
			gotRiichi = &got[i]
			break
		}
	}
	if gotRiichi == nil {
		t.Fatalf("getSelfTurnCandidates() contains no riichi candidate")
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

func TestManueAgent_getSelfTurnCandidates_FiltersNonTenpaiRiichiCandidate(t *testing.T) {
	self := seat.MustSeat(0)
	riichi := action.NewRiichi(self)
	discard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}

	got, err := getSelfTurnCandidates([]action.Action{discard, riichi}, stubPlayerViewer{
		hand: hand.CodesToHand([]string{
			"1m", "2m", "3m", "4m", "5m", "6m", "7m",
			"1p", "2p", "3p", "E", "E", "S", "S",
		}),
		riichiState: player.NotRiichi,
		drawnTile:   nil,
	})
	if err != nil {
		t.Fatalf("getSelfTurnCandidates() failed: %v", err)
	}
	for _, candidate := range got {
		if candidate.riichi {
			t.Fatalf("getSelfTurnCandidates() contains riichi candidate unexpectedly: %+v", candidate)
		}
	}
}

func TestManueAgent_getSelfTurnCandidates_RiichiDeclaredScoresDiscardAsRiichi(t *testing.T) {
	self := seat.MustSeat(0)
	tenpaiDiscard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard(tenpai) failed: %v", err)
	}
	nonTenpaiDiscard, err := action.NewDiscard(self, tile.MustTileFromCode("1p"), false)
	if err != nil {
		t.Fatalf("NewDiscard(non-tenpai) failed: %v", err)
	}

	got, err := getSelfTurnCandidates([]action.Action{tenpaiDiscard, nonTenpaiDiscard}, stubPlayerViewer{
		hand: hand.CodesToHand([]string{
			"1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m",
			"1p", "1p", "E", "E", "5m",
		}),
		riichiState: player.RiichiDeclared,
		drawnTile:   nil,
	})
	if err != nil {
		t.Fatalf("getSelfTurnCandidates() failed: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(getSelfTurnCandidates()) = %d, want 1", len(got))
	}
	if got[0].traceKey != "-1.5m" {
		t.Errorf("traceKey = %q, want -1.5m", got[0].traceKey)
	}
	if got[0].riichi {
		t.Errorf("riichi = true, want false because declaration action is already done")
	}
	if !got[0].scoreAsRiichi {
		t.Errorf("scoreAsRiichi = false, want true")
	}
	if got[0].action != tenpaiDiscard {
		t.Errorf("action = %v, want tenpai discard", got[0].action)
	}
}

func TestManueAgent_getSelfTurnCandidates_IncludesRiichiAndDiscardCandidates(t *testing.T) {
	self := seat.MustSeat(0)
	riichi := action.NewRiichi(self)
	discard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}

	got, err := getSelfTurnCandidates([]action.Action{discard, riichi}, stubPlayerViewer{
		hand: hand.CodesToHand([]string{
			"1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m",
			"1p", "1p", "E", "E", "5m",
		}),
		riichiState: player.NotRiichi,
		drawnTile:   nil,
	})
	if err != nil {
		t.Fatalf("getSelfTurnCandidates() failed: %v", err)
	}
	wantTraceKeys := map[string]bool{"0.5m": false, "-1.5m": false}
	for _, candidate := range got {
		if _, ok := wantTraceKeys[candidate.traceKey]; ok {
			wantTraceKeys[candidate.traceKey] = true
		}
	}
	for traceKey, found := range wantTraceKeys {
		if !found {
			t.Errorf("getSelfTurnCandidates() does not contain %q", traceKey)
		}
	}
}

func TestManueAgent_scoreDiscardCandidate_MarksRedDiscard(t *testing.T) {
	redDiscardTile := tile.MustTileFromCode("5mr")

	got := scoreDiscardCandidate(redDiscardTile, 1)
	if !got.red {
		t.Errorf("red = false, want true")
	}
	if got.shanten != 1 {
		t.Errorf("shanten = %d, want 1", got.shanten)
	}
}

func TestEvaluateCandidateScore(t *testing.T) {
	scoreChanges := newScoreDeltaProbDist(map[scoreDelta]float64{
		{1000, -1000, 0, 0}: 0.25,
		{-500, 500, 0, 0}:   0.75,
	})
	base := candidateScore{
		winProb:  0.2,
		drawProb: 0.3,
		shanten:  1,
	}

	got := evaluateCandidateScore(base, scoreChanges, 0, 25000, 0, []rankOpponent{
		{
			id:       1,
			score:    25000,
			position: 1,
			winProbs: relativeWinProbTable{
				"2000":  1.0,
				"-1000": 0.0,
			},
		},
		{
			id:       2,
			score:    24000,
			position: 2,
			winProbs: relativeWinProbTable{
				"2000": 1.0,
				"500":  1.0,
			},
		},
		{
			id:       3,
			score:    26000,
			position: 3,
			winProbs: relativeWinProbTable{
				"0":     0.5,
				"-1500": 0.0,
			},
		},
	})

	if !almostEqual(got.expPts, -125) {
		t.Errorf("expPts = %v, want -125", got.expPts)
	}
	if !almostEqual(got.avgRank, 2.625) {
		t.Errorf("avgRank = %v, want 2.625", got.avgRank)
	}
	if got.winProb != base.winProb {
		t.Errorf("winProb = %v, want %v", got.winProb, base.winProb)
	}
	if got.drawProb != base.drawProb {
		t.Errorf("drawProb = %v, want %v", got.drawProb, base.drawProb)
	}
	if got.shanten != base.shanten {
		t.Errorf("shanten = %v, want %v", got.shanten, base.shanten)
	}
}

func TestEvaluateCandidateScoreFromState(t *testing.T) {
	scoreChanges := newScoreDeltaProbDist(map[scoreDelta]float64{
		{1000, -1000, 0, 0}: 0.25,
		{-500, 500, 0, 0}:   0.75,
	})
	base := candidateScore{
		winProb:  0.2,
		drawProb: 0.3,
		shanten:  1,
	}

	got := evaluateCandidateScoreFromState(
		base,
		scoreChanges,
		stubManueStats{
			relativeWinProbs: map[string]map[string]float64{
				"E1,0,1": {
					"2000":  1.0,
					"-1000": 0.0,
				},
				"E1,0,2": {
					"2000": 1.0,
					"500":  1.0,
				},
				"E1,0,3": {
					"0":     0.5,
					"-1500": 0.0,
				},
			},
		},
		stubRankStateViewer{
			nextRoundWind:  wind.East,
			nextRoundNum:   1,
			scores:         [common.NumPlayers]int{25000, 25000, 24000, 26000},
			startingDealer: seat.MustSeat(0),
		},
		seat.MustSeat(0),
	)

	if !almostEqual(got.expPts, -125) {
		t.Errorf("expPts = %v, want -125", got.expPts)
	}
	if !almostEqual(got.avgRank, 2.625) {
		t.Errorf("avgRank = %v, want 2.625", got.avgRank)
	}
	if got.winProb != base.winProb {
		t.Errorf("winProb = %v, want %v", got.winProb, base.winProb)
	}
	if got.drawProb != base.drawProb {
		t.Errorf("drawProb = %v, want %v", got.drawProb, base.drawProb)
	}
	if got.shanten != base.shanten {
		t.Errorf("shanten = %v, want %v", got.shanten, base.shanten)
	}
}

func TestCandidateTotalScoreDeltaDist(t *testing.T) {
	score := candidateScore{
		winProb:       0.2,
		drawProb:      0.3,
		othersWinProb: 0.4,
	}
	immediateDist := scoreDeltaProbDist{
		{}:            0.75,
		{-1000, 1000}: 0.25,
	}
	selfWinDist := scoreDeltaProbDist{{1000, 0, 0, 0}: 1}
	exhaustiveDrawDist := scoreDeltaProbDist{{0, 1000, 0, 0}: 1}
	otherWinDists := []scoreDeltaProbDist{
		{{0, 0, 1000, 0}: 1},
		{{0, 0, 0, 1000}: 1},
	}

	got := candidateTotalScoreDeltaDist(
		score,
		immediateDist,
		selfWinDist,
		exhaustiveDrawDist,
		otherWinDists,
	)

	want := scoreDeltaProbDist{
		{-1000, 1000}:   0.25,
		{1000, 0, 0, 0}: 0.15,
		{0, 1000, 0, 0}: 0.225,
		{0, 0, 1000, 0}: 0.15,
		{0, 0, 0, 1000}: 0.15,
	}
	assertScoreDeltaProbDist(t, got, want)
}

func TestEvaluateCandidateFromPreparedDists(t *testing.T) {
	score := candidateScore{
		shanten: 1,
	}
	dealInEstimates := []dealInEstimate{
		{winnerID: 1, prob: 0.2},
		{winnerID: 2, prob: 0.25},
	}
	immediateDist := scoreDeltaProbDist{
		{}: 1,
	}
	selfWinDist := scoreDeltaProbDist{
		{1000, 0, 0, 0}: 1,
	}
	exhaustiveDrawDist := scoreDeltaProbDist{
		{0, 1000, 0, 0}: 1,
	}
	otherWinDists := []scoreDeltaProbDist{
		{{0, 0, 1000, 0}: 1},
		{{0, 0, 0, 1000}: 1},
	}

	got, err := evaluateCandidateFromPreparedDists(
		score,
		dealInEstimates,
		winEstimate{
			prob:   0.2,
			avgPts: 3900,
		},
		0.25,
		1200,
		immediateDist,
		selfWinDist,
		exhaustiveDrawDist,
		otherWinDists,
		stubManueStats{
			relativeWinProbs: map[string]map[string]float64{
				"E1,0,1": {"-1000": 0.5, "0": 0.5, "1000": 0.5},
				"E1,0,2": {"-1000": 0.5, "0": 0.5, "1000": 0.5},
				"E1,0,3": {"-1000": 0.5, "0": 0.5, "1000": 0.5},
			},
		},
		stubRankStateViewer{
			nextRoundWind:  wind.East,
			nextRoundNum:   1,
			scores:         [common.NumPlayers]int{25000, 25000, 25000, 25000},
			startingDealer: seat.MustSeat(0),
		},
		seat.MustSeat(0),
	)
	if err != nil {
		t.Fatalf("evaluateCandidateFromPreparedDists() failed: %v", err)
	}

	if !almostEqual(got.dealInProb, 0.4) {
		t.Errorf("dealInProb = %v, want 0.4", got.dealInProb)
	}
	if !almostEqual(got.drawProb, 0.2) {
		t.Errorf("drawProb = %v, want 0.2", got.drawProb)
	}
	if !almostEqual(got.othersWinProb, 0.6) {
		t.Errorf("othersWinProb = %v, want 0.6", got.othersWinProb)
	}
	if got.avgWinPts != 3900 {
		t.Errorf("avgWinPts = %v, want 3900", got.avgWinPts)
	}
	if got.avgDrawPts != 1200 {
		t.Errorf("avgDrawPts = %v, want 1200", got.avgDrawPts)
	}
	if !almostEqual(got.expPts, 200) {
		t.Errorf("expPts = %v, want 200", got.expPts)
	}
	if !almostEqual(got.avgRank, 2.5) {
		t.Errorf("avgRank = %v, want 2.5", got.avgRank)
	}
}

func TestEvaluateCandidateFromPreparedDists_ReturnsErrorWithInvalidEstimate(t *testing.T) {
	_, err := evaluateCandidateFromPreparedDists(
		candidateScore{},
		[]dealInEstimate{{winnerID: 1, prob: 1.1}},
		winEstimate{},
		0.25,
		0,
		nil,
		nil,
		nil,
		nil,
		stubManueStats{},
		stubRankStateViewer{},
		seat.MustSeat(0),
	)
	if err == nil {
		t.Fatal("evaluateCandidateFromPreparedDists() succeeded unexpectedly")
	}
}

func TestCandidateTraceKeys(t *testing.T) {
	got, err := candidateTraceKeys([]actionCandidate{
		{traceKey: "-1.5m"},
		{traceKey: "0.5m"},
	})
	if err != nil {
		t.Fatalf("candidateTraceKeys() failed: %v", err)
	}

	want := []string{"-1.5m", "0.5m"}
	if len(got) != len(want) {
		t.Fatalf("len(candidateTraceKeys()) = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("candidateTraceKeys()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestCandidateTraceKeys_ReturnsErrorWithInvalidKey(t *testing.T) {
	tests := []struct {
		name       string
		candidates []actionCandidate
	}{
		{
			name:       "empty",
			candidates: []actionCandidate{{traceKey: ""}},
		},
		{
			name: "duplicate",
			candidates: []actionCandidate{
				{traceKey: "-1.5m"},
				{traceKey: "-1.5m"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := candidateTraceKeys(tt.candidates)
			if err == nil {
				t.Fatal("candidateTraceKeys() succeeded unexpectedly")
			}
		})
	}
}

func TestWinEstimatesForCandidates(t *testing.T) {
	got, err := winEstimatesForCandidates(
		[]actionCandidate{
			{traceKey: "-1.5m"},
			{traceKey: "0.5m"},
		},
		[]map[string]float64{
			{"-1.5m": 2000},
			{"0.5m": 5000},
			{},
		},
	)
	if err != nil {
		t.Fatalf("winEstimatesForCandidates() failed: %v", err)
	}

	discardEstimate := got["-1.5m"]
	if !almostEqual(discardEstimate.prob, 1.0/3.0) {
		t.Errorf("discard prob = %v, want %v", discardEstimate.prob, 1.0/3.0)
	}
	if discardEstimate.avgPts != 2000 {
		t.Errorf("discard avgPts = %v, want 2000", discardEstimate.avgPts)
	}

	riichiEstimate := got["0.5m"]
	if !almostEqual(riichiEstimate.prob, 1.0/3.0) {
		t.Errorf("riichi prob = %v, want %v", riichiEstimate.prob, 1.0/3.0)
	}
	if riichiEstimate.avgPts != 5000 {
		t.Errorf("riichi avgPts = %v, want 5000", riichiEstimate.avgPts)
	}
}

func TestWinEstimatesForCandidates_ReturnsErrorWithInvalidCandidates(t *testing.T) {
	_, err := winEstimatesForCandidates(
		[]actionCandidate{{traceKey: ""}},
		[]map[string]float64{},
	)
	if err == nil {
		t.Fatal("winEstimatesForCandidates() succeeded unexpectedly")
	}
}

func TestWinEstimatesForCandidates_ReturnsErrorWithInvalidTrial(t *testing.T) {
	_, err := winEstimatesForCandidates(
		[]actionCandidate{{traceKey: "-1.5m"}},
		[]map[string]float64{{"unknown": 2000}},
	)
	if err == nil {
		t.Fatal("winEstimatesForCandidates() succeeded unexpectedly")
	}
}

func TestFilteredWinEstimateGoals(t *testing.T) {
	tests := []struct {
		name      string
		candidate actionCandidate
		want      []int
	}{
		{
			name: "keeps all goals normally",
			candidate: actionCandidate{
				discardTile: tile.MustTileFromCode("1m"),
				score:       candidateScore{shanten: 2},
				shantenGoals: []service.Goal{
					{Shanten: 1, ThrowableVector: hand.TileCounts34{0: 1}},
					{Shanten: 2, ThrowableVector: hand.TileCounts34{0: 1}},
					{Shanten: 3, ThrowableVector: hand.TileCounts34{0: 1}},
				},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "riichi keeps only ready goals",
			candidate: actionCandidate{
				riichi:        true,
				scoreAsRiichi: true,
				discardTile:   tile.MustTileFromCode("1m"),
				score:         candidateScore{shanten: 0},
				shantenGoals: []service.Goal{
					{Shanten: 0, ThrowableVector: hand.TileCounts34{0: 1}},
					{Shanten: 1, ThrowableVector: hand.TileCounts34{0: 1}},
				},
			},
			want: []int{0},
		},
		{
			name: "heavy shanten drops extra tile goals",
			candidate: actionCandidate{
				discardTile: tile.MustTileFromCode("1m"),
				baseShanten: 4,
				score:       candidateScore{shanten: 1},
				shantenGoals: []service.Goal{
					{Shanten: 3, ThrowableVector: hand.TileCounts34{0: 1}},
					{Shanten: 4, ThrowableVector: hand.TileCounts34{0: 1}},
					{Shanten: 5, ThrowableVector: hand.TileCounts34{0: 1}},
				},
			},
			want: []int{3, 4},
		},
		{
			name: "riichi and heavy shanten combine",
			candidate: actionCandidate{
				riichi:        true,
				scoreAsRiichi: true,
				discardTile:   tile.MustTileFromCode("1m"),
				baseShanten:   4,
				score:         candidateScore{shanten: 1},
				shantenGoals: []service.Goal{
					{Shanten: 0, ThrowableVector: hand.TileCounts34{0: 1}},
					{Shanten: 4, ThrowableVector: hand.TileCounts34{0: 1}},
					{Shanten: 5, ThrowableVector: hand.TileCounts34{0: 1}},
				},
			},
			want: []int{0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filteredWinEstimateGoals(tt.candidate)
			if len(got) != len(tt.want) {
				t.Fatalf("len(filteredWinEstimateGoals()) = %d, want %d", len(got), len(tt.want))
			}
			for i, want := range tt.want {
				if got[i].Shanten != want {
					t.Errorf("filteredWinEstimateGoals()[%d].Shanten = %d, want %d", i, got[i].Shanten, want)
				}
			}
		})
	}
}

func TestScoredWinEstimateGoals(t *testing.T) {
	turnHand := hand.MustVisibleHand([]tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("2m"),
		tile.MustTileFromCode("3m"),
		tile.MustTileFromCode("4m"),
		tile.MustTileFromCode("3p"),
		tile.MustTileFromCode("4p"),
		tile.MustTileFromCode("5p"),
		tile.MustTileFromCode("4s"),
		tile.MustTileFromCode("5s"),
		tile.MustTileFromCode("6s"),
		tile.MustTileFromCode("6s"),
		tile.MustTileFromCode("7s"),
		tile.MustTileFromCode("8s"),
		tile.MustTileFromCode("2s"),
	})
	afterDiscardHand, err := turnHand.Discard(tile.MustTileFromCode("1m"))
	if err != nil {
		t.Fatalf("Discard(1m) failed: %v", err)
	}
	candidate := actionCandidate{
		discardTile:      tile.MustTileFromCode("1m"),
		turnHand:         turnHand,
		afterDiscardHand: afterDiscardHand,
		shantenGoals: []service.Goal{
			{
				Blocks: []block.Block{
					block.MustSequence(tile.MustTileFromCode("2m")),
					block.MustSequence(tile.MustTileFromCode("3p")),
					block.MustSequence(tile.MustTileFromCode("4s")),
					block.MustSequence(tile.MustTileFromCode("6s")),
					block.MustPair(tile.MustTileFromCode("2s")),
				},
				RequiredVector:  hand.TileCounts34{19: 1},
				ThrowableVector: hand.TileCounts34{0: 1},
			},
			{
				Blocks: []block.Block{
					block.MustSequence(tile.MustTileFromCode("1m")),
					block.MustSequence(tile.MustTileFromCode("2p")),
					block.MustSequence(tile.MustTileFromCode("3s")),
					block.MustSequence(tile.MustTileFromCode("6s")),
					block.MustPair(tile.MustTileFromCode("E")),
				},
				RequiredVector:  hand.TileCounts34{1: 1},
				ThrowableVector: hand.TileCounts34{0: 1},
			},
		},
	}

	got, err := scoredWinEstimateGoals(candidate, winEstimateGoalContext{
		roundWind: wind.East,
		seatWind:  wind.South,
		dealer:    false,
	})
	if err != nil {
		t.Fatalf("scoredWinEstimateGoals() failed: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(scoredWinEstimateGoals()) = %d, want 1", len(got))
	}
	wantPoints := service.RonPoints(30, 2, false)
	if got[0].points != float64(wantPoints) {
		t.Errorf("scoredWinEstimateGoals()[0].points = %v, want %v", got[0].points, wantPoints)
	}
	if got[0].RequiredVector[19] != 1 {
		t.Errorf("scoredWinEstimateGoals()[0].RequiredVector[19] = %d, want 1", got[0].RequiredVector[19])
	}
}

func TestScoredWinEstimateGoalsRequiresTurnHand(t *testing.T) {
	_, err := scoredWinEstimateGoals(actionCandidate{}, winEstimateGoalContext{})
	if err == nil {
		t.Fatal("scoredWinEstimateGoals() succeeded unexpectedly")
	}
}

func TestScoredWinEstimateGoalsByKey(t *testing.T) {
	turnHand := hand.MustVisibleHand([]tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("2m"),
		tile.MustTileFromCode("3m"),
		tile.MustTileFromCode("4m"),
		tile.MustTileFromCode("3p"),
		tile.MustTileFromCode("4p"),
		tile.MustTileFromCode("5p"),
		tile.MustTileFromCode("4s"),
		tile.MustTileFromCode("5s"),
		tile.MustTileFromCode("6s"),
		tile.MustTileFromCode("6s"),
		tile.MustTileFromCode("7s"),
		tile.MustTileFromCode("8s"),
		tile.MustTileFromCode("2s"),
	})
	goals := []service.Goal{
		{
			Blocks: []block.Block{
				block.MustSequence(tile.MustTileFromCode("2m")),
				block.MustSequence(tile.MustTileFromCode("3p")),
				block.MustSequence(tile.MustTileFromCode("4s")),
				block.MustSequence(tile.MustTileFromCode("6s")),
				block.MustPair(tile.MustTileFromCode("2s")),
			},
			RequiredVector:  hand.TileCounts34{19: 1},
			ThrowableVector: hand.TileCounts34{0: 1},
		},
	}
	afterDiscard1m, err := turnHand.Discard(tile.MustTileFromCode("1m"))
	if err != nil {
		t.Fatalf("Discard(1m) failed: %v", err)
	}
	afterDiscard2m, err := turnHand.Discard(tile.MustTileFromCode("2m"))
	if err != nil {
		t.Fatalf("Discard(2m) failed: %v", err)
	}
	candidates := []actionCandidate{
		{
			traceKey:         "-1.1m",
			discardTile:      tile.MustTileFromCode("1m"),
			turnHand:         turnHand,
			afterDiscardHand: afterDiscard1m,
			shantenGoals:     goals,
		},
		{
			traceKey:         "-1.2m",
			discardTile:      tile.MustTileFromCode("2m"),
			turnHand:         turnHand,
			afterDiscardHand: afterDiscard2m,
		},
	}

	got, err := scoredWinEstimateGoalsByKey(candidates, winEstimateGoalContext{
		roundWind: wind.East,
		seatWind:  wind.South,
	})
	if err != nil {
		t.Fatalf("scoredWinEstimateGoalsByKey() failed: %v", err)
	}
	if len(got["-1.1m"]) != 1 {
		t.Fatalf("len(scored goals for -1.1m) = %d, want 1", len(got["-1.1m"]))
	}
	if got["-1.1m"][0].points != float64(service.RonPoints(30, 2, false)) {
		t.Errorf("points for -1.1m = %v, want %v", got["-1.1m"][0].points, service.RonPoints(30, 2, false))
	}
	if goals, ok := got["-1.2m"]; !ok {
		t.Fatal("scored goals for -1.2m missing")
	} else if len(goals) != 0 {
		t.Errorf("len(scored goals for -1.2m) = %d, want 0", len(goals))
	}
}

func TestScoredWinEstimateGoalsByKeyUsesCandidateMelds(t *testing.T) {
	turnHand := hand.MustVisibleHand([]tile.Tile{
		tile.MustTileFromCode("2m"),
		tile.MustTileFromCode("2m"),
		tile.MustTileFromCode("2m"),
		tile.MustTileFromCode("3p"),
		tile.MustTileFromCode("3p"),
		tile.MustTileFromCode("3p"),
		tile.MustTileFromCode("4s"),
		tile.MustTileFromCode("4s"),
		tile.MustTileFromCode("4s"),
		tile.MustTileFromCode("9s"),
		tile.MustTileFromCode("9s"),
	})
	goals := []service.Goal{
		{
			Blocks: []block.Block{
				block.MustTriplet(tile.MustTileFromCode("2m")),
				block.MustTriplet(tile.MustTileFromCode("3p")),
				block.MustTriplet(tile.MustTileFromCode("4s")),
				block.MustPair(tile.MustTileFromCode("9s")),
			},
			RequiredVector: hand.TileCounts34{0: 1},
		},
	}
	dragonPon := meld.MustPon(
		tile.MustTileFromCode("F"),
		[2]tile.Tile{tile.MustTileFromCode("F"), tile.MustTileFromCode("F")},
		seat.MustSeat(1),
	)
	candidates := []actionCandidate{
		{
			traceKey:         "none",
			discardTile:      tile.MustTileFromCode("?"),
			turnHand:         turnHand,
			afterDiscardHand: turnHand,
			shantenGoals:     goals,
		},
		{
			traceKey:         "0.1m",
			discardTile:      tile.MustTileFromCode("?"),
			turnHand:         turnHand,
			afterDiscardHand: turnHand,
			melds:            []meld.Meld{dragonPon},
			shantenGoals:     goals,
		},
	}

	got, err := scoredWinEstimateGoalsByKey(candidates, winEstimateGoalContext{
		roundWind: wind.East,
		seatWind:  wind.South,
	})
	if err != nil {
		t.Fatalf("scoredWinEstimateGoalsByKey() failed: %v", err)
	}

	if len(got["none"]) != 1 {
		t.Fatalf("len(scored goals for none) = %d, want 1", len(got["none"]))
	}
	if len(got["0.1m"]) != 1 {
		t.Fatalf("len(scored goals for 0.1m) = %d, want 1", len(got["0.1m"]))
	}
	want := float64(service.RonPoints(30, 3, false))
	if got["0.1m"][0].points != want {
		t.Errorf("points for 0.1m = %v, want %v", got["0.1m"][0].points, want)
	}
	if got["0.1m"][0].points <= got["none"][0].points {
		t.Errorf("points with candidate melds = %v, want greater than %v", got["0.1m"][0].points, got["none"][0].points)
	}
}

func TestScoredWinEstimateGoalsByKeyWrapsCandidateError(t *testing.T) {
	_, err := scoredWinEstimateGoalsByKey(
		[]actionCandidate{{traceKey: "-1.1m"}},
		winEstimateGoalContext{},
	)
	if err == nil {
		t.Fatal("scoredWinEstimateGoalsByKey() succeeded unexpectedly")
	}
	if !strings.Contains(err.Error(), "-1.1m") {
		t.Errorf("scoredWinEstimateGoalsByKey() error = %v, want trace key", err)
	}
}

func TestTrialTileCounts(t *testing.T) {
	got := trialTileCounts([]tile.Tile{
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("5mr"),
		tile.MustTileFromCode("E"),
	})

	if got[4] != 2 {
		t.Errorf("5m count = %d, want 2", got[4])
	}
	if got[27] != 1 {
		t.Errorf("E count = %d, want 1", got[27])
	}
}

func TestWallTilesFromCounts(t *testing.T) {
	got, err := wallTilesFromCounts(hand.TileCounts34{
		0:  2,
		4:  1,
		27: 1,
	})
	if err != nil {
		t.Fatalf("wallTilesFromCounts() failed: %v", err)
	}
	want := []tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("E"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("wallTilesFromCounts() = %v, want %v", got, want)
	}
}

func TestWallTilesFromCountsRejectsNegativeCount(t *testing.T) {
	_, err := wallTilesFromCounts(hand.TileCounts34{0: -1})
	if err == nil {
		t.Fatal("wallTilesFromCounts() succeeded unexpectedly")
	}
}

func TestUnseenWallFromVisibleTiles(t *testing.T) {
	got, err := unseenWallFromVisibleTiles([]tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("5mr"),
		tile.MustTileFromCode("E"),
	})
	if err != nil {
		t.Fatalf("unseenWallFromVisibleTiles() failed: %v", err)
	}
	counts := trialTileCounts(got)
	if counts[0] != 3 {
		t.Errorf("1m unseen count = %d, want 3", counts[0])
	}
	if counts[4] != 3 {
		t.Errorf("5m unseen count = %d, want 3", counts[4])
	}
	if counts[27] != 3 {
		t.Errorf("E unseen count = %d, want 3", counts[27])
	}
	if numTiles := (&counts).NumTiles(); numTiles != 133 {
		t.Errorf("unseen wall tile count = %d, want 133", numTiles)
	}
}

func TestUnseenWallFromVisibleTilesRejectsInvalidVisibleTiles(t *testing.T) {
	if _, err := unseenWallFromVisibleTiles([]tile.Tile{tile.MustTileFromCode("?")}); err == nil {
		t.Fatal("unseenWallFromVisibleTiles() accepted unknown tile")
	}

	_, err := unseenWallFromVisibleTiles([]tile.Tile{
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("5mr"),
	})
	if err == nil {
		t.Fatal("unseenWallFromVisibleTiles() accepted tile visible more than 4 times")
	}
}

func TestTrialTilesFromWall(t *testing.T) {
	wall := []tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("2m"),
		tile.MustTileFromCode("3m"),
	}

	got, err := trialTilesFromWall(wall, 2)
	if err != nil {
		t.Fatalf("trialTilesFromWall() failed: %v", err)
	}
	want := []tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("2m"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("trialTilesFromWall() = %v, want %v", got, want)
	}

	got[0] = tile.MustTileFromCode("9m")
	if wall[0] != tile.MustTileFromCode("1m") {
		t.Errorf("wall[0] = %s, want unchanged 1m", wall[0])
	}
}

func TestTrialTilesFromWallRejectsInvalidNumDraws(t *testing.T) {
	wall := []tile.Tile{tile.MustTileFromCode("1m")}
	if _, err := trialTilesFromWall(wall, -1); err == nil {
		t.Fatal("trialTilesFromWall(-1) succeeded unexpectedly")
	}
	if _, err := trialTilesFromWall(wall, 2); err == nil {
		t.Fatal("trialTilesFromWall(2) succeeded unexpectedly")
	}
}

func TestCanAchieveGoalWithTrialTiles(t *testing.T) {
	goal := service.Goal{
		RequiredVector: hand.TileCounts34{
			0:  1,
			4:  2,
			27: 1,
		},
	}

	if !canAchieveGoalWithTrialTiles(goal, trialTileCounts([]tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("5mr"),
		tile.MustTileFromCode("E"),
	})) {
		t.Fatal("canAchieveGoalWithTrialTiles() = false, want true")
	}
	if canAchieveGoalWithTrialTiles(goal, trialTileCounts([]tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("E"),
	})) {
		t.Fatal("canAchieveGoalWithTrialTiles() = true, want false")
	}
}

func TestTrialWinPts(t *testing.T) {
	goals := []winEstimateGoal{
		{
			Goal: service.Goal{
				RequiredVector: hand.TileCounts34{0: 1},
			},
			points: 1000,
		},
		{
			Goal: service.Goal{
				RequiredVector: hand.TileCounts34{0: 1, 4: 1},
			},
			points: 2000,
		},
		{
			Goal: service.Goal{
				RequiredVector: hand.TileCounts34{27: 1},
			},
			points: 8000,
		},
	}

	got, ok, err := trialWinPts(goals, trialTileCounts([]tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("5m"),
	}))
	if err != nil {
		t.Fatalf("trialWinPts() failed: %v", err)
	}
	if !ok {
		t.Fatal("trialWinPts() ok = false, want true")
	}
	if got != 2000 {
		t.Errorf("trialWinPts() = %v, want 2000", got)
	}

	got, ok, err = trialWinPts(goals, trialTileCounts([]tile.Tile{
		tile.MustTileFromCode("2m"),
	}))
	if err != nil {
		t.Fatalf("trialWinPts() failed: %v", err)
	}
	if ok {
		t.Fatal("trialWinPts() ok = true, want false")
	}
	if got != 0 {
		t.Errorf("trialWinPts() = %v, want 0", got)
	}
}

func TestTrialWinPtsRejectsNonPositivePoints(t *testing.T) {
	_, _, err := trialWinPts([]winEstimateGoal{
		{
			Goal: service.Goal{
				RequiredVector: hand.TileCounts34{0: 1},
			},
			points: 0,
		},
	}, trialTileCounts([]tile.Tile{tile.MustTileFromCode("1m")}))
	if err == nil {
		t.Fatal("trialWinPts() succeeded unexpectedly")
	}
}

func TestCandidateTrialWinPts(t *testing.T) {
	candidates := []actionCandidate{
		{traceKey: "-1.1m"},
		{traceKey: "-1.2m"},
		{traceKey: "0.3m"},
	}
	goalsByKey := map[string][]winEstimateGoal{
		"-1.1m": {
			{
				Goal: service.Goal{
					RequiredVector: hand.TileCounts34{0: 1},
				},
				points: 1000,
			},
			{
				Goal: service.Goal{
					RequiredVector: hand.TileCounts34{0: 1, 4: 1},
				},
				points: 3900,
			},
		},
		"-1.2m": {
			{
				Goal: service.Goal{
					RequiredVector: hand.TileCounts34{27: 1},
				},
				points: 8000,
			},
		},
		"0.3m": {},
	}

	got, err := candidateTrialWinPts(candidates, goalsByKey, trialTileCounts([]tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("5m"),
	}))
	if err != nil {
		t.Fatalf("candidateTrialWinPts() failed: %v", err)
	}
	want := map[string]float64{
		"-1.1m": 3900,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("candidateTrialWinPts() = %#v, want %#v", got, want)
	}
}

func TestCandidateTrialWinPtsRequiresGoalsForEveryCandidate(t *testing.T) {
	_, err := candidateTrialWinPts(
		[]actionCandidate{{traceKey: "-1.1m"}},
		map[string][]winEstimateGoal{},
		hand.TileCounts34{},
	)
	if err == nil {
		t.Fatal("candidateTrialWinPts() succeeded unexpectedly")
	}
}

func TestWinEstimatesFromTrialTiles(t *testing.T) {
	candidates := []actionCandidate{
		{traceKey: "-1.1m"},
		{traceKey: "-1.2m"},
	}
	goalsByKey := map[string][]winEstimateGoal{
		"-1.1m": {
			{
				Goal: service.Goal{
					RequiredVector: hand.TileCounts34{0: 1},
				},
				points: 1000,
			},
			{
				Goal: service.Goal{
					RequiredVector: hand.TileCounts34{0: 1, 4: 1},
				},
				points: 3900,
			},
		},
		"-1.2m": {
			{
				Goal: service.Goal{
					RequiredVector: hand.TileCounts34{27: 1},
				},
				points: 8000,
			},
		},
	}

	got, err := winEstimatesFromTrialTiles(candidates, goalsByKey, [][]tile.Tile{
		{
			tile.MustTileFromCode("1m"),
			tile.MustTileFromCode("5m"),
		},
		{
			tile.MustTileFromCode("E"),
		},
		{
			tile.MustTileFromCode("2m"),
		},
	})
	if err != nil {
		t.Fatalf("winEstimatesFromTrialTiles() failed: %v", err)
	}

	first := got["-1.1m"]
	if first.prob != float64(1)/3 {
		t.Errorf("first.prob = %v, want %v", first.prob, float64(1)/3)
	}
	if first.avgPts != 3900 {
		t.Errorf("first.avgPts = %v, want 3900", first.avgPts)
	}
	if first.expPts != 1300 {
		t.Errorf("first.expPts = %v, want 1300", first.expPts)
	}
	if got := first.pointsDist.expected(); got != 3900 {
		t.Errorf("first.pointsDist.expected() = %v, want 3900", got)
	}

	second := got["-1.2m"]
	if second.prob != float64(1)/3 {
		t.Errorf("second.prob = %v, want %v", second.prob, float64(1)/3)
	}
	if second.avgPts != 8000 {
		t.Errorf("second.avgPts = %v, want 8000", second.avgPts)
	}
	if second.expPts != float64(8000)/3 {
		t.Errorf("second.expPts = %v, want %v", second.expPts, float64(8000)/3)
	}
}

func TestWinEstimatesFromShuffledWall(t *testing.T) {
	candidates := []actionCandidate{
		{traceKey: "-1.1m"},
	}
	goalsByKey := map[string][]winEstimateGoal{
		"-1.1m": {
			{
				Goal: service.Goal{
					RequiredVector: hand.TileCounts34{0: 1, 4: 1},
				},
				points: 3900,
			},
		},
	}
	wall := []tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("5m"),
	}

	got, err := winEstimatesFromShuffledWall(
		candidates,
		goalsByKey,
		wall,
		2,
		3,
		rand.New(rand.NewPCG(1, 0)),
	)
	if err != nil {
		t.Fatalf("winEstimatesFromShuffledWall() failed: %v", err)
	}

	estimate := got["-1.1m"]
	if estimate.prob != 1 {
		t.Errorf("estimate.prob = %v, want 1", estimate.prob)
	}
	if estimate.avgPts != 3900 {
		t.Errorf("estimate.avgPts = %v, want 3900", estimate.avgPts)
	}
	if estimate.expPts != 3900 {
		t.Errorf("estimate.expPts = %v, want 3900", estimate.expPts)
	}
	wantWall := []tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("5m"),
	}
	if !reflect.DeepEqual(wall, wantWall) {
		t.Errorf("wall = %v, want unchanged %v", wall, wantWall)
	}
}

func TestWinEstimatesFromShuffledWallRejectsInvalidInputs(t *testing.T) {
	candidates := []actionCandidate{{traceKey: "-1.1m"}}
	goalsByKey := map[string][]winEstimateGoal{"-1.1m": {}}
	wall := []tile.Tile{tile.MustTileFromCode("1m")}

	if _, err := winEstimatesFromShuffledWall(candidates, goalsByKey, wall, 1, 0, rand.New(rand.NewPCG(1, 0))); err == nil {
		t.Fatal("winEstimatesFromShuffledWall() accepted zero numTries")
	}
	if _, err := winEstimatesFromShuffledWall(candidates, goalsByKey, wall, 1, 1, nil); err == nil {
		t.Fatal("winEstimatesFromShuffledWall() accepted nil rng")
	}
	if _, err := winEstimatesFromShuffledWall(candidates, goalsByKey, wall, 2, 1, rand.New(rand.NewPCG(1, 0))); err == nil {
		t.Fatal("winEstimatesFromShuffledWall() accepted too many draws")
	}
}

func TestWinEstimatesFromState(t *testing.T) {
	candidates := []actionCandidate{
		{traceKey: "-1.1m"},
	}
	goalsByKey := map[string][]winEstimateGoal{
		"-1.1m": {
			{
				Goal: service.Goal{
					RequiredVector: hand.TileCounts34{0: 1},
				},
				points: 1000,
			},
		},
	}
	visibleTiles := make([]tile.Tile, 0, 135)
	for id := range tile.NumTileType34 {
		count := 4
		if id == 0 {
			count = 3
		}
		for range count {
			visibleTiles = append(visibleTiles, tile.MustTileFromID(id))
		}
	}
	state := stubWinEstimateStateViewer{
		turn:         0,
		visibleTiles: visibleTiles,
	}
	turnDistribution := make([]float64, numTurnDistributionEntries)
	turnDistribution[0] = 1
	stats := stubManueStats{
		turnDistribution: turnDistribution,
	}

	got, err := winEstimatesFromState(
		stats,
		state,
		seat.MustSeat(0),
		candidates,
		goalsByKey,
		3,
		rand.New(rand.NewPCG(1, 0)),
	)
	if err != nil {
		t.Fatalf("winEstimatesFromState() failed: %v", err)
	}

	estimate := got["-1.1m"]
	if estimate.prob != 1 {
		t.Errorf("estimate.prob = %v, want 1", estimate.prob)
	}
	if estimate.avgPts != 1000 {
		t.Errorf("estimate.avgPts = %v, want 1000", estimate.avgPts)
	}
	if estimate.expPts != 1000 {
		t.Errorf("estimate.expPts = %v, want 1000", estimate.expPts)
	}
}

func TestWinEstimatesFromStateReturnsErrorWithInvalidVisibleTiles(t *testing.T) {
	_, err := winEstimatesFromState(
		stubManueStats{turnDistribution: fullTurnDistribution(1)},
		stubWinEstimateStateViewer{visibleTiles: []tile.Tile{tile.MustTileFromCode("?")}},
		seat.MustSeat(0),
		[]actionCandidate{{traceKey: "-1.1m"}},
		map[string][]winEstimateGoal{"-1.1m": {}},
		1,
		rand.New(rand.NewPCG(1, 0)),
	)
	if err == nil {
		t.Fatal("winEstimatesFromState() succeeded unexpectedly")
	}
}

func TestManueAgent_formatDiscardTraceKey(t *testing.T) {
	discardTile := tile.MustTileFromCode("5m")
	if got := formatDiscardTraceKey(false, discardTile); got != "-1.5m" {
		t.Errorf("formatDiscardTraceKey(false) = %q, want %q", got, "-1.5m")
	}
	if got := formatDiscardTraceKey(true, discardTile); got != "0.5m" {
		t.Errorf("formatDiscardTraceKey(true) = %q, want %q", got, "0.5m")
	}
}

func TestManueAgent_formatCandidateTrace(t *testing.T) {
	self := seat.MustSeat(0)
	discard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}

	got := formatCandidateTrace([]actionCandidate{
		{
			traceKey:    "-1.5m",
			action:      discard,
			riichi:      false,
			discardTile: discard.Tile(),
			score: candidateScore{
				avgRank:       2.25,
				expPts:        1200,
				dealInProb:    0.125,
				winProb:       0.25,
				drawProb:      0.375,
				othersWinProb: 0.5,
				avgWinPts:     3900,
				avgDrawPts:    1000,
				shanten:       1,
			},
		},
	})
	want := "| action | avgRank | expPt | hojuProb | myHoraProb | ryukyokuProb | otherHoraProb | avgHoraPt | ryukyokuAvgPt | shanten | \n" +
		"|  -1.5m |  2.2500 |  1200 |    0.125 |      0.250 |        0.375 |         0.500 |      3900 |          1000 |       1 | \n"
	if got != want {
		t.Errorf("formatCandidateTrace() =\n%q\nwant\n%q", got, want)
	}
}

func TestManueAgent_formatCandidateTrace_FormatsInfinityShanten(t *testing.T) {
	self := seat.MustSeat(0)
	discard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}

	got := formatCandidateTrace([]actionCandidate{
		{
			traceKey:    "-1.5m",
			action:      discard,
			riichi:      false,
			discardTile: discard.Tile(),
			score: candidateScore{
				shanten: service.InfinityShanten,
			},
		},
	})
	if !strings.Contains(got, "Inf") {
		t.Errorf("formatCandidateTrace() = %q, want it to contain Inf", got)
	}
	if strings.Contains(got, fmt.Sprintf("%d", service.InfinityShanten)) {
		t.Errorf("formatCandidateTrace() = %q, should not contain raw InfinityShanten integer", got)
	}
}

func TestManueAgent_formatDecisionTrace_AppendsDecidedKey(t *testing.T) {
	self := seat.MustSeat(0)
	discard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}
	selected := &actionCandidate{
		traceKey:    "-1.5m",
		action:      discard,
		riichi:      false,
		discardTile: discard.Tile(),
		score: candidateScore{
			shanten: 1,
		},
	}

	got := formatDecisionTrace(formatCandidateLog([]actionCandidate{*selected}), selected)
	if !strings.HasSuffix(got, "\n\n\ndecidedKey -1.5m\n") {
		t.Errorf("formatDecisionTrace() = %q, want two blank lines before decidedKey suffix", got)
	}
}

func TestManueAgent_compareCandidateScore(t *testing.T) {
	tests := []struct {
		name        string
		lhs         candidateScore
		rhs         candidateScore
		preferBlack bool
		want        int
	}{
		{
			name:        "better average rank wins",
			lhs:         candidateScore{avgRank: 1.9, expPts: 0},
			rhs:         candidateScore{avgRank: 2.0, expPts: 1000},
			preferBlack: true,
			want:        -1,
		},
		{
			name:        "higher expected points wins on rank tie",
			lhs:         candidateScore{avgRank: 2.0, expPts: 1000},
			rhs:         candidateScore{avgRank: 2.0, expPts: 900},
			preferBlack: true,
			want:        -1,
		},
		{
			name:        "black tile wins on expected value tie",
			lhs:         candidateScore{avgRank: 2.0, expPts: 1000, red: false},
			rhs:         candidateScore{avgRank: 2.0, expPts: 1000, red: true},
			preferBlack: true,
			want:        -1,
		},
		{
			name:        "red tie ignored when preferBlack is false",
			lhs:         candidateScore{avgRank: 2.0, expPts: 1000, red: false},
			rhs:         candidateScore{avgRank: 2.0, expPts: 1000, red: true},
			preferBlack: false,
			want:        0,
		},
		{
			name:        "complete tie returns zero",
			lhs:         candidateScore{avgRank: 2.0, expPts: 1000},
			rhs:         candidateScore{avgRank: 2.0, expPts: 1000},
			preferBlack: false,
			want:        0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := compareCandidateScore(&tt.lhs, &tt.rhs, tt.preferBlack)
			if got != tt.want {
				t.Errorf("compareCandidateScore() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestManueAgent_chooseBestCandidate_DoesNotPreferRiichiOnScoreTie(t *testing.T) {
	self := seat.MustSeat(0)
	riichi := action.NewRiichi(self)
	discard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}

	candidates := []actionCandidate{
		{action: discard, riichi: false, score: candidateScore{avgRank: 2.0, expPts: 1000}},
		{action: riichi, riichi: true, score: candidateScore{avgRank: 2.0, expPts: 1000}},
	}

	got := chooseBestCandidate(candidates, false)
	if got.action != discard {
		t.Errorf("chooseBestCandidate() = %T %[1]v, want first tied candidate", got.action)
	}
}
