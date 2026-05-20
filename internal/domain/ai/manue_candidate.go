package ai

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type candidateScore struct {
	// avgRank is the average rank.
	avgRank float64
	// expPts is the expected points.
	expPts float64
	// dealInProb is the deal-in probability.
	dealInProb float64
	// winProb is the win probability.
	winProb float64
	// drawProb is the draw probability.
	drawProb float64
	// othersWinProb is the other players' win probability.
	othersWinProb float64
	// avgWinPts is the average win points.
	avgWinPts float64
	// avgDrawPts is the average draw points.
	avgDrawPts float64
	// shanten is the shanten number.
	shanten int
	// red indicates whether the candidate discards a red tile.
	red bool
}

type actionCandidate struct {
	traceKey    string
	action      action.Action
	riichi      bool
	discardTile tile.Tile
	score       candidateScore
}

func chooseBestCandidate(candidates []actionCandidate, preferBlack bool) actionCandidate {
	best := candidates[0]
	for _, candidate := range candidates[1:] {
		if compareCandidate(candidate, best, preferBlack) < 0 {
			best = candidate
		}
	}
	return best
}

func compareCandidate(lhs, rhs actionCandidate, preferBlack bool) int {
	if r := compareCandidateScore(&lhs.score, &rhs.score, preferBlack); r != 0 {
		return r
	}
	return compareCandidateFallback(lhs, rhs)
}

// compareCandidateFallback preserves Phase 1 behavior until score calculation is migrated.
func compareCandidateFallback(lhs, rhs actionCandidate) int {
	if lhs.riichi && !rhs.riichi {
		return -1
	}
	if !lhs.riichi && rhs.riichi {
		return 1
	}
	return 0
}

func compareCandidateScore(lhs, rhs *candidateScore, preferBlack bool) int {
	if lhs.avgRank < rhs.avgRank {
		return -1
	}
	if lhs.avgRank > rhs.avgRank {
		return 1
	}
	if lhs.expPts > rhs.expPts {
		return -1
	}
	if lhs.expPts < rhs.expPts {
		return 1
	}
	if preferBlack {
		if !lhs.red && rhs.red {
			return -1
		}
		if lhs.red && !rhs.red {
			return 1
		}
	}
	return 0
}

func getSelfTurnCandidates(actions []action.Action, self player.PlayerViewer) ([]actionCandidate, error) {
	h, err := selfTurnHand(self)
	if err != nil {
		return nil, fmt.Errorf("cannot build self-turn candidates: %w", err)
	}

	// Self-turn candidates currently cover discard and riichi+discard only.
	// Concealed kan, promoted kan, and kyushukyuhai are intentionally not
	// selected.
	riichi := firstActionOfType[*action.Riichi](actions)
	var candidates []actionCandidate
	for _, a := range actions {
		discard, ok := a.(*action.Discard)
		if !ok {
			continue
		}

		afterDiscard, err := h.Discard(discard.Tile())
		if err != nil {
			return nil, fmt.Errorf("cannot build self-turn candidate for %s: %w", discard.Tile(), err)
		}
		shanten, _ := service.AnalyzeShanten(afterDiscard, service.AllowedExtraTiles(1))
		if riichi != nil && shanten <= 0 {
			candidates = append(candidates, buildSelfTurnCandidate(riichi, discard.Tile(), shanten, true))
		}
		candidates = append(candidates, buildSelfTurnCandidate(discard, discard.Tile(), shanten, false))
	}
	return candidates, nil
}

func selfTurnHand(self player.PlayerViewer) (*hand.VisibleHand, error) {
	h, ok := self.Hand()
	if !ok {
		return nil, fmt.Errorf("self hand is not visible")
	}
	drawnTile := self.DrawnTile()
	if drawnTile == nil {
		return h, nil
	}
	withDrawnTile, err := h.Draw(*drawnTile)
	if err != nil {
		return nil, fmt.Errorf("cannot add drawn tile %s to self hand: %w", *drawnTile, err)
	}
	return withDrawnTile, nil
}

func buildSelfTurnCandidate(
	immediateAction action.Action,
	discardTile tile.Tile,
	shanten int,
	riichi bool,
) actionCandidate {
	return actionCandidate{
		traceKey:    formatDiscardTraceKey(riichi, discardTile),
		action:      immediateAction,
		riichi:      riichi,
		discardTile: discardTile,
		score:       scoreDiscardCandidate(discardTile, shanten),
	}
}

func scoreDiscardCandidate(discardTile tile.Tile, shanten int) candidateScore {
	return candidateScore{
		// Phase 2 scaffold: the full expected-value fields are filled by later
		// configs/estimator migration. Keep all choices tied except red fallback.
		avgRank: 0,
		expPts:  0,
		shanten: shanten,
		red:     discardTile.IsRed(),
	}
}
