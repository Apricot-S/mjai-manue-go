package ai

import (
	"fmt"
	"strconv"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round"
)

// validateManueStats checks structural invariants of stats before they are used
// by ManueAgent. The validation assumes stats is immutable; implementations
// must not change returned values after validation.
func validateManueStats(stats ManueStats) error {
	if err := validateWinScoreStats(stats); err != nil {
		return err
	}
	if err := validateDrawTenpaiStats(stats); err != nil {
		return err
	}
	return nil
}

func validateWinScoreStats(stats WinScoreStats) error {
	numWins := stats.NumWins()
	if numWins <= 0 {
		return fmt.Errorf("invalid win score stats: numWins must be positive")
	}
	numSelfDrawWins := stats.NumSelfDrawWins()
	if numSelfDrawWins < 0 || numSelfDrawWins > numWins {
		return fmt.Errorf("invalid win score stats: numSelfDrawWins must be between 0 and numWins")
	}
	if err := validateWinPointFreqs("non-dealer", stats.NonDealerWinPointFreqs()); err != nil {
		return err
	}
	if err := validateWinPointFreqs("dealer", stats.DealerWinPointFreqs()); err != nil {
		return err
	}
	return nil
}

func validateWinPointFreqs(label string, pointFreqs map[string]int) error {
	totalFreqs, ok := pointFreqs["total"]
	if !ok {
		return fmt.Errorf("invalid %s win point frequencies: total is missing", label)
	}
	if totalFreqs <= 0 {
		return fmt.Errorf("invalid %s win point frequencies: total must be positive", label)
	}

	sumFreqs := 0
	for points, freq := range pointFreqs {
		if points == "total" {
			continue
		}
		parsedPoints, err := strconv.Atoi(points)
		if err != nil {
			return fmt.Errorf("invalid %s win point frequency key %q: %w", label, points, err)
		}
		if parsedPoints <= 0 {
			return fmt.Errorf("invalid %s win point frequency key %q: points must be positive", label, points)
		}
		if parsedPoints%100 != 0 {
			return fmt.Errorf("invalid %s win point frequency key %q: points must be a multiple of 100", label, points)
		}
		if freq < 0 {
			return fmt.Errorf("invalid %s win point frequency %q: frequency must be non-negative", label, points)
		}
		sumFreqs += freq
	}
	if sumFreqs != totalFreqs {
		return fmt.Errorf("invalid %s win point frequencies: total must equal frequency sum", label)
	}
	return nil
}

func validateDrawTenpaiStats(stats DrawTenpaiStats) error {
	notenFreq := stats.ExhaustiveDrawNotenCount()
	if notenFreq < 0 {
		return fmt.Errorf("invalid draw tenpai stats: noten count must be non-negative")
	}

	sumTenpaiFreqs := 0
	for turn := 0.0; turn <= round.FinalTurn; turn += 0.25 {
		key := strconv.FormatFloat(turn, 'f', -1, 64)
		freq, ok := stats.ExhaustiveDrawTenpaiTurnFreq(key)
		if !ok {
			return fmt.Errorf("invalid draw tenpai stats: missing tenpai turn frequency for turn %s", key)
		}
		if freq < 0 {
			return fmt.Errorf("invalid draw tenpai stats: tenpai turn frequency for turn %s must be non-negative", key)
		}
		sumTenpaiFreqs += freq
	}

	totalFreqs := sumTenpaiFreqs + notenFreq
	if totalFreqs <= 0 {
		return fmt.Errorf("invalid draw tenpai stats: frequency total must be positive")
	}
	return nil
}
