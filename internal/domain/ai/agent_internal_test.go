package ai

import (
	"strings"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
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
				return agent.decideSelfTurn(actions, stubStateWithSelf(self), seat.MustSeat(0))
			},
		},
		{
			name:        "self turn evaluates riichi and discard candidates",
			actions:     []action.Action{handDiscard, riichi},
			hand:        hand.CodesToHand([]string{"1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m", "1p", "1p", "E", "E", "1m"}),
			riichiState: player.NotRiichi,
			drawnTile:   nil,
			want:        handDiscard,
			decide: func(agent *ManueAgent, actions []action.Action, self player.PlayerViewer) (Decision, error) {
				return agent.decideSelfTurn(actions, stubStateWithSelf(self), seat.MustSeat(0))
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
	if agent.evaluator.rng == nil {
		t.Errorf("evaluator.rng = nil, want initialized rng")
	}
	if agent.evaluator.stats == nil {
		t.Errorf("evaluator.stats = nil, want stats")
	}
	if agent.evaluator.danger == nil {
		t.Errorf("evaluator.danger = nil, want danger estimator")
	}
	firstRand := agent.evaluator.rng.Float64()
	agent.Reset()
	if secondRand := agent.evaluator.rng.Float64(); secondRand != firstRand {
		t.Errorf("Reset() first random value = %v, want %v", secondRand, firstRand)
	}
}

func TestManueAgent_decideSelfTurn_ReturnsOriginalStyleActionLog(t *testing.T) {
	self := seat.MustSeat(0)
	discard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}

	decision, err := newTestManueAgent(t, 0).decideSelfTurn([]action.Action{discard}, stubStateWithSelf(stubPlayerViewer{
		hand:        hand.CodesToHand([]string{"1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m", "1p", "1p", "E", "E", "5m"}),
		riichiState: player.NotRiichi,
		drawnTile:   nil,
	}), seat.MustSeat(0))
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

	_, err = newTestManueAgent(t, 0).decideSelfTurn([]action.Action{handDiscard}, stubStateWithSelf(stubPlayerViewer{
		riichiState: player.RiichiAccepted,
		drawnTile:   &drawnTile,
	}), seat.MustSeat(0))
	if err == nil {
		t.Fatal("selectAction() succeeded unexpectedly")
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

func TestBuildCandidateDecision(t *testing.T) {
	self := seat.MustSeat(0)
	redDiscard, err := action.NewDiscard(self, tile.MustTileFromCode("5mr"), false)
	if err != nil {
		t.Fatalf("NewDiscard(red) failed: %v", err)
	}
	blackDiscard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard(black) failed: %v", err)
	}

	decision := buildCandidateDecision([]actionCandidate{
		{
			traceKey:    "-1.5mr",
			action:      redDiscard,
			discardTile: redDiscard.Tile(),
			score: candidateScore{
				averageRank:    2.0,
				expectedPoints: 1000,
				red:            true,
			},
		},
		{
			traceKey:    "-1.5m",
			action:      blackDiscard,
			discardTile: blackDiscard.Tile(),
			score: candidateScore{
				averageRank:    2.0,
				expectedPoints: 1000,
				red:            false,
			},
		},
	}, true)

	if decision.Action != blackDiscard {
		t.Errorf("Action = %T %[1]v, want black discard", decision.Action)
	}
	if !strings.Contains(decision.Log, "| action |") {
		t.Errorf("Log = %q, want candidate table", decision.Log)
	}
	if !strings.Contains(decision.Trace, "decidedKey -1.5m\n") {
		t.Errorf("Trace = %q, want selected decidedKey", decision.Trace)
	}
}
