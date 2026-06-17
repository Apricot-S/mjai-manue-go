package ai

import (
	"fmt"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type actionCandidate struct {
	traceKey         string
	evaluationGroup  int
	action           action.Action
	riichi           bool
	scoreAsRiichi    bool
	discardTile      tile.Tile
	melds            []meld.Meld
	afterDiscardHand *hand.VisibleHand
	baseShanten      int
	shanten          int
	shantenGoals     []service.Goal
	red              bool
}

type evaluatedActionCandidate struct {
	candidate actionCandidate
	score     candidateScore
}

func chooseBestCandidate(candidates []evaluatedActionCandidate, preferBlack bool) evaluatedActionCandidate {
	best := candidates[0]
	for _, candidate := range candidates[1:] {
		if compareCandidates(candidate, best, preferBlack) < 0 {
			best = candidate
		}
	}
	return best
}

func sortedCandidates(candidates []evaluatedActionCandidate, preferBlack bool) []evaluatedActionCandidate {
	sortedCandidates := slices.Clone(candidates)
	slices.SortFunc(sortedCandidates, func(lhs, rhs evaluatedActionCandidate) int {
		return compareCandidates(lhs, rhs, preferBlack)
	})
	return sortedCandidates
}

func compareCandidates(lhs, rhs evaluatedActionCandidate, preferBlack bool) int {
	if result := compareCandidateScore(&lhs.score, &rhs.score); result != 0 {
		return result
	}
	if preferBlack {
		if !lhs.candidate.red && rhs.candidate.red {
			return -1
		}
		if lhs.candidate.red && !rhs.candidate.red {
			return 1
		}
	}
	return 0
}

func candidateTraceKeys(candidates []actionCandidate) ([]string, error) {
	keys := make([]string, 0, len(candidates))
	seen := make(map[string]struct{}, len(candidates))
	for _, candidate := range candidates {
		if candidate.traceKey == "" {
			return nil, fmt.Errorf("cannot build candidate keys: trace key must not be empty")
		}
		if _, ok := seen[candidate.traceKey]; ok {
			return nil, fmt.Errorf("cannot build candidate keys: duplicate trace key %q", candidate.traceKey)
		}
		seen[candidate.traceKey] = struct{}{}
		keys = append(keys, candidate.traceKey)
	}
	return keys, nil
}
