package ai

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestManueAgent_DecideActionSkeleton(t *testing.T) {
	self := seat.MustSeat(0)
	target := seat.MustSeat(1)
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
	chii, err := action.NewChii(self, target, tile.MustTileFromCode("3m"), [2]tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("2m"),
	})
	if err != nil {
		t.Fatalf("NewChii() failed: %v", err)
	}
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
				return agent.decideSelfTurn(actions, self)
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
				return agent.decideSelfTurn(actions, self)
			},
		},
		{
			name:        "discard before call",
			actions:     []action.Action{chii, handDiscard, pass},
			hand:        hand.CodesToHand([]string{"1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m", "1p", "1p", "E", "S"}),
			riichiState: player.NotRiichi,
			drawnTile:   &drawnTile,
			want:        handDiscard,
			decide: func(agent *ManueAgent, actions []action.Action, self player.PlayerViewer) (Decision, error) {
				return agent.decideSelfTurn(actions, self)
			},
		},
		{
			name:        "call before pass",
			actions:     []action.Action{pass, chii},
			riichiState: player.NotRiichi,
			drawnTile:   nil,
			want:        chii,
			decide: func(agent *ManueAgent, actions []action.Action, self player.PlayerViewer) (Decision, error) {
				return agent.decideOtherDiscardReaction(actions)
			},
		},
		{
			name:        "pass only",
			actions:     []action.Action{pass},
			riichiState: player.NotRiichi,
			drawnTile:   nil,
			want:        pass,
			decide: func(agent *ManueAgent, actions []action.Action, self player.PlayerViewer) (Decision, error) {
				return agent.decideOtherDiscardReaction(actions)
			},
		},
	}

	agent := NewManueAgent(0)
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

func TestManueAgent_decideSelfTurn_ReturnsOriginalStyleActionLog(t *testing.T) {
	self := seat.MustSeat(0)
	discard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}

	decision, err := NewManueAgent(0).decideSelfTurn([]action.Action{discard}, stubPlayerViewer{
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

	_, err = NewManueAgent(0).decideSelfTurn([]action.Action{handDiscard}, stubPlayerViewer{
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
	if !strings.Contains(got, "Infinity") {
		t.Errorf("formatCandidateTrace() = %q, want it to contain Infinity", got)
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

func TestManueAgent_compareCandidateFallback_PreservesPhase1RiichiPreference(t *testing.T) {
	self := seat.MustSeat(0)
	riichi := action.NewRiichi(self)
	discard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}

	lhs := actionCandidate{action: riichi, riichi: true, discardTile: discard.Tile()}
	rhs := actionCandidate{action: discard, riichi: false, discardTile: discard.Tile()}
	if got := compareCandidateFallback(lhs, rhs); got != -1 {
		t.Errorf("compareCandidateFallback(riichi, discard) = %d, want -1", got)
	}
}

type stubPlayerViewer struct {
	hand        *hand.VisibleHand
	riichiState player.RiichiState
	drawnTile   *tile.Tile
}

func (p stubPlayerViewer) Hand() (*hand.VisibleHand, bool) {
	if p.hand == nil {
		return nil, false
	}
	return p.hand, true
}
func (p stubPlayerViewer) HandTiles() []tile.Tile          { return nil }
func (p stubPlayerViewer) DrawnTile() *tile.Tile           { return p.drawnTile }
func (p stubPlayerViewer) Melds() []meld.Meld              { return nil }
func (p stubPlayerViewer) River() []tile.Tile              { return nil }
func (p stubPlayerViewer) DiscardedTiles() []tile.Tile     { return nil }
func (p stubPlayerViewer) ExtraSafeTiles() []tile.Tile     { return nil }
func (p stubPlayerViewer) IsFuriten() bool                 { return false }
func (p stubPlayerViewer) CanRonBy(*tile.Tile) bool        { return true }
func (p stubPlayerViewer) RiichiState() player.RiichiState { return p.riichiState }
func (p stubPlayerViewer) RiichiRiverIndex() int           { return -1 }
func (p stubPlayerViewer) RiichiDiscardedTilesIndex() int  { return -1 }
func (p stubPlayerViewer) CanDiscard() bool                { return p.drawnTile != nil }
func (p stubPlayerViewer) CanChiiPonKan() bool             { return p.drawnTile == nil }
func (p stubPlayerViewer) IsConcealed() bool               { return true }
func (p stubPlayerViewer) SwapCallTiles() []tile.Tile      { return nil }
