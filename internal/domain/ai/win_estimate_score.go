package ai

import (
	"fmt"
	"maps"
)

type winEstimate struct {
	// prob is the win probability.
	prob float64
	// avgPts is the average win points when self wins.
	avgPts float64
	// expectedPoints is the expected win points over all estimation trials.
	expectedPoints float64
	// pointsDist is the win-points distribution when self wins.
	pointsDist scalarProbDist
}

type winEstimateAccumulator struct {
	numTries   int
	totalWins  int
	totalPts   float64
	pointFreqs map[float64]int
}

func (a *winEstimateAccumulator) addNoWinTrial() {
	a.numTries++
}

func (a *winEstimateAccumulator) addWinTrial(points float64) error {
	if points <= 0 {
		return fmt.Errorf("cannot add win trial: points must be positive")
	}
	a.numTries++
	a.totalWins++
	a.totalPts += points
	if a.pointFreqs == nil {
		a.pointFreqs = make(map[float64]int)
	}
	a.pointFreqs[points]++
	return nil
}

func (a *winEstimateAccumulator) merge(other winEstimateAccumulator) error {
	if other.numTries < 0 || other.totalWins < 0 || other.totalPts < 0 {
		return fmt.Errorf("cannot merge win estimate accumulator: values must be non-negative")
	}
	if other.totalWins > other.numTries {
		return fmt.Errorf("cannot merge win estimate accumulator: totalWins must not exceed numTries")
	}

	countedWins := 0
	for points, freq := range other.pointFreqs {
		if points <= 0 {
			return fmt.Errorf("cannot merge win estimate accumulator: points must be positive")
		}
		if freq < 0 {
			return fmt.Errorf("cannot merge win estimate accumulator: point frequency must be non-negative")
		}
		countedWins += freq
	}
	if countedWins != other.totalWins {
		return fmt.Errorf("cannot merge win estimate accumulator: point frequencies must sum to totalWins")
	}

	a.numTries += other.numTries
	a.totalWins += other.totalWins
	a.totalPts += other.totalPts
	if len(other.pointFreqs) > 0 && a.pointFreqs == nil {
		a.pointFreqs = make(map[float64]int, len(other.pointFreqs))
	}
	for points, freq := range other.pointFreqs {
		a.pointFreqs[points] += freq
	}
	return nil
}

func (a winEstimateAccumulator) estimate() (winEstimate, error) {
	return newWinEstimate(a.numTries, a.totalWins, a.totalPts, a.pointFreqs)
}

func (a winEstimateAccumulator) clone() winEstimateAccumulator {
	return winEstimateAccumulator{
		numTries:   a.numTries,
		totalWins:  a.totalWins,
		totalPts:   a.totalPts,
		pointFreqs: maps.Clone(a.pointFreqs),
	}
}

type winEstimateAccumulatorSet map[string]winEstimateAccumulator

func newWinEstimateAccumulatorSet(keys []string) winEstimateAccumulatorSet {
	accumulators := make(winEstimateAccumulatorSet, len(keys))
	for _, key := range keys {
		accumulators[key] = winEstimateAccumulator{}
	}
	return accumulators
}

func (s winEstimateAccumulatorSet) addNoWinTrial(key string) {
	accumulator := s[key]
	accumulator.addNoWinTrial()
	s[key] = accumulator
}

func (s winEstimateAccumulatorSet) addWinTrial(key string, points float64) error {
	accumulator := s[key]
	if err := accumulator.addWinTrial(points); err != nil {
		return err
	}
	s[key] = accumulator
	return nil
}

func (s winEstimateAccumulatorSet) addTrial(winPtsByKey map[string]float64) error {
	for key, points := range winPtsByKey {
		if _, ok := s[key]; !ok {
			return fmt.Errorf("cannot add win estimate trial: unknown candidate key %q", key)
		}
		if points <= 0 {
			return fmt.Errorf("cannot add win estimate trial: points must be positive")
		}
	}

	for key := range s {
		points, ok := winPtsByKey[key]
		if ok {
			if err := s.addWinTrial(key, points); err != nil {
				return err
			}
			continue
		}
		s.addNoWinTrial(key)
	}
	return nil
}

func (s winEstimateAccumulatorSet) merge(other winEstimateAccumulatorSet) error {
	if len(s) != len(other) {
		return fmt.Errorf("cannot merge win estimate accumulator sets: candidate key count mismatch")
	}

	merged := make(winEstimateAccumulatorSet, len(other))
	for key, otherAccumulator := range other {
		current, ok := s[key]
		if !ok {
			return fmt.Errorf("cannot merge win estimate accumulator sets: unknown candidate key %q", key)
		}
		accumulator := current.clone()
		if err := accumulator.merge(otherAccumulator); err != nil {
			return fmt.Errorf("cannot merge win estimate accumulator for %q: %w", key, err)
		}
		merged[key] = accumulator
	}
	maps.Copy(s, merged)
	return nil
}

func (s winEstimateAccumulatorSet) estimates() (map[string]winEstimate, error) {
	estimates := make(map[string]winEstimate, len(s))
	for key, accumulator := range s {
		estimate, err := accumulator.estimate()
		if err != nil {
			return nil, fmt.Errorf("cannot build win estimate for %q: %w", key, err)
		}
		estimates[key] = estimate
	}
	return estimates, nil
}

func winEstimatesFromTrials(keys []string, trials []map[string]float64) (map[string]winEstimate, error) {
	accumulators := newWinEstimateAccumulatorSet(keys)
	for _, trial := range trials {
		if err := accumulators.addTrial(trial); err != nil {
			return nil, err
		}
	}
	return accumulators.estimates()
}

func newWinEstimate(
	numTries int,
	totalWins int,
	totalPts float64,
	pointFreqs map[float64]int,
) (winEstimate, error) {
	if numTries <= 0 {
		return winEstimate{}, fmt.Errorf("cannot build win estimate: numTries must be positive")
	}
	if totalWins < 0 {
		return winEstimate{}, fmt.Errorf("cannot build win estimate: totalWins must be non-negative")
	}
	if totalWins > numTries {
		return winEstimate{}, fmt.Errorf("cannot build win estimate: totalWins must not exceed numTries")
	}
	if totalPts < 0 {
		return winEstimate{}, fmt.Errorf("cannot build win estimate: totalPts must be non-negative")
	}
	if totalWins == 0 {
		if totalPts != 0 || len(pointFreqs) != 0 {
			return winEstimate{}, fmt.Errorf("cannot build win estimate: zero wins must have no points")
		}
		return winEstimate{
			prob:           0,
			avgPts:         0,
			expectedPoints: 0,
			pointsDist:     scalarProbDist{},
		}, nil
	}

	totalWinsFloat := float64(totalWins)
	dist := make(map[float64]float64, len(pointFreqs))
	countedWins := 0
	for points, freq := range pointFreqs {
		if points <= 0 {
			return winEstimate{}, fmt.Errorf("cannot build win estimate: points must be positive")
		}
		if freq < 0 {
			return winEstimate{}, fmt.Errorf("cannot build win estimate: point frequency must be non-negative")
		}
		countedWins += freq
		dist[points] = float64(freq) / totalWinsFloat
	}
	if countedWins != totalWins {
		return winEstimate{}, fmt.Errorf("cannot build win estimate: point frequencies must sum to totalWins")
	}

	return winEstimate{
		prob:           totalWinsFloat / float64(numTries),
		avgPts:         totalPts / totalWinsFloat,
		expectedPoints: totalPts / float64(numTries),
		pointsDist:     newScalarProbDist(dist),
	}, nil
}
