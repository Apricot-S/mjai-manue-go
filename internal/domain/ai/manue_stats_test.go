package ai

import "testing"

func TestValidateManueStats(t *testing.T) {
	stats := validStubManueStats()

	if err := validateManueStats(stats); err != nil {
		t.Errorf("validateManueStats() failed: %v", err)
	}
}

func TestValidateWinScoreStats_ReturnsErrorWithInvalidNumWins(t *testing.T) {
	stats := validStubManueStats()
	stats.numWins = 0

	if err := validateWinScoreStats(stats); err == nil {
		t.Fatal("validateWinScoreStats() succeeded unexpectedly")
	}
}

func TestValidateWinScoreStats_ReturnsErrorWithTooManySelfDrawWins(t *testing.T) {
	stats := validStubManueStats()
	stats.numSelfDrawWins = stats.numWins + 1

	if err := validateWinScoreStats(stats); err == nil {
		t.Fatal("validateWinScoreStats() succeeded unexpectedly")
	}
}

func TestValidateWinScoreStats_ReturnsErrorWithInvalidPointFreqs(t *testing.T) {
	tests := []struct {
		name       string
		pointFreqs map[string]int
	}{
		{
			name:       "missing total",
			pointFreqs: map[string]int{"1000": 1},
		},
		{
			name:       "invalid total",
			pointFreqs: map[string]int{"1000": 1, "total": 0},
		},
		{
			name:       "invalid point key",
			pointFreqs: map[string]int{"bad": 1, "total": 1},
		},
		{
			name:       "fractional point key",
			pointFreqs: map[string]int{"1000.5": 1, "total": 1},
		},
		{
			name:       "zero point key",
			pointFreqs: map[string]int{"0": 1, "total": 1},
		},
		{
			name:       "negative point key",
			pointFreqs: map[string]int{"-1000": 1, "total": 1},
		},
		{
			name:       "not multiple of 100",
			pointFreqs: map[string]int{"1050": 1, "total": 1},
		},
		{
			name:       "negative frequency",
			pointFreqs: map[string]int{"1000": -1, "total": 1},
		},
		{
			name:       "sum mismatch",
			pointFreqs: map[string]int{"1000": 1, "total": 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := validStubManueStats()
			stats.nonDealerWinPointFreqs = tt.pointFreqs

			if err := validateWinScoreStats(stats); err == nil {
				t.Fatal("validateWinScoreStats() succeeded unexpectedly")
			}
		})
	}
}

func TestValidateRoundEndStats_ReturnsErrorWithInvalidTurnDistributionLength(t *testing.T) {
	stats := validStubManueStats()
	stats.turnDistribution = stats.turnDistribution[:numTurnDistributionEntries-1]

	if err := validateRoundEndStats(stats); err == nil {
		t.Fatal("validateRoundEndStats() succeeded unexpectedly")
	}
}

func TestValidateRoundEndStats_ReturnsErrorWithNegativeTurnProb(t *testing.T) {
	stats := validStubManueStats()
	stats.turnDistribution[1] = -0.1

	if err := validateRoundEndStats(stats); err == nil {
		t.Fatal("validateRoundEndStats() succeeded unexpectedly")
	}
}

func TestValidateRoundEndStats_ReturnsErrorWithInvalidExhaustiveDrawRatio(t *testing.T) {
	stats := validStubManueStats()
	stats.exhaustiveDrawRatio = -0.1

	if err := validateRoundEndStats(stats); err == nil {
		t.Fatal("validateRoundEndStats() succeeded unexpectedly")
	}
}

func TestValidateRoundEndStats_ReturnsErrorWhenExhaustiveDrawRatioExceedsTotal(t *testing.T) {
	stats := validStubManueStats()
	stats.exhaustiveDrawRatio = 2

	if err := validateRoundEndStats(stats); err == nil {
		t.Fatal("validateRoundEndStats() succeeded unexpectedly")
	}
}

func TestValidateDrawTenpaiStats_ReturnsErrorWithNegativeNotenCount(t *testing.T) {
	stats := validStubManueStats()
	stats.exhaustiveDrawNotenCount = -1

	if err := validateDrawTenpaiStats(stats); err == nil {
		t.Fatal("validateDrawTenpaiStats() succeeded unexpectedly")
	}
}

func TestValidateDrawTenpaiStats_ReturnsErrorWithMissingTurnFreq(t *testing.T) {
	stats := validStubManueStats()
	delete(stats.exhaustiveDrawTenpaiTurnFreqs, "17.5")

	if err := validateDrawTenpaiStats(stats); err == nil {
		t.Fatal("validateDrawTenpaiStats() succeeded unexpectedly")
	}
}

func TestValidateDrawTenpaiStats_ReturnsErrorWithNegativeTurnFreq(t *testing.T) {
	stats := validStubManueStats()
	stats.exhaustiveDrawTenpaiTurnFreqs["17.5"] = -1

	if err := validateDrawTenpaiStats(stats); err == nil {
		t.Fatal("validateDrawTenpaiStats() succeeded unexpectedly")
	}
}

func TestValidateDrawTenpaiStats_ReturnsErrorWithoutFreqs(t *testing.T) {
	stats := validStubManueStats()
	for key := range stats.exhaustiveDrawTenpaiTurnFreqs {
		stats.exhaustiveDrawTenpaiTurnFreqs[key] = 0
	}
	stats.exhaustiveDrawNotenCount = 0

	if err := validateDrawTenpaiStats(stats); err == nil {
		t.Fatal("validateDrawTenpaiStats() succeeded unexpectedly")
	}
}

func TestValidateTenpaiEstimatorStats_ReturnsErrorWithInvalidYamitenCounts(t *testing.T) {
	tests := []struct {
		name  string
		count yamitenCount
	}{
		{
			name:  "invalid total",
			count: yamitenCount{total: 0, tenpai: 0},
		},
		{
			name:  "negative tenpai",
			count: yamitenCount{total: 10, tenpai: -1},
		},
		{
			name:  "too many tenpai",
			count: yamitenCount{total: 10, tenpai: 11},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := validStubManueStats()
			stats.yamitenCounts["1,0"] = tt.count

			if err := validateTenpaiEstimatorStats(stats); err == nil {
				t.Fatal("validateTenpaiEstimatorStats() succeeded unexpectedly")
			}
		})
	}
}

func TestValidateDealInStats_ReturnsErrorWithInvalidAvgWinPts(t *testing.T) {
	stats := validStubManueStats()
	stats.avgWinPointsValue = 0

	if err := validateDealInStats(stats); err == nil {
		t.Fatal("validateDealInStats() succeeded unexpectedly")
	}
}

func validStubManueStats() stubManueStats {
	return stubManueStats{
		numWins:         10,
		numSelfDrawWins: 4,
		nonDealerWinPointFreqs: map[string]int{
			"1000":  1,
			"2000":  2,
			"total": 3,
		},
		dealerWinPointFreqs: map[string]int{
			"1500":  1,
			"3000":  2,
			"total": 3,
		},
		turnDistribution:              fullTurnDistribution(0.01),
		exhaustiveDrawRatio:           0.1,
		avgWinPointsValue:             5500,
		exhaustiveDrawNotenCount:      100,
		exhaustiveDrawTenpaiTurnFreqs: fullTurnFreqs(1),
		yamitenCounts: map[string]yamitenCount{
			"1,0": {total: 10, tenpai: 3},
		},
	}
}

func fullTurnDistribution(prob float64) []float64 {
	distribution := make([]float64, numTurnDistributionEntries)
	for i := range distribution {
		distribution[i] = prob
	}
	return distribution
}

func fullTurnFreqs(freq int) map[string]int {
	return map[string]int{
		"0":     freq,
		"0.25":  freq,
		"0.5":   freq,
		"0.75":  freq,
		"1":     freq,
		"1.25":  freq,
		"1.5":   freq,
		"1.75":  freq,
		"2":     freq,
		"2.25":  freq,
		"2.5":   freq,
		"2.75":  freq,
		"3":     freq,
		"3.25":  freq,
		"3.5":   freq,
		"3.75":  freq,
		"4":     freq,
		"4.25":  freq,
		"4.5":   freq,
		"4.75":  freq,
		"5":     freq,
		"5.25":  freq,
		"5.5":   freq,
		"5.75":  freq,
		"6":     freq,
		"6.25":  freq,
		"6.5":   freq,
		"6.75":  freq,
		"7":     freq,
		"7.25":  freq,
		"7.5":   freq,
		"7.75":  freq,
		"8":     freq,
		"8.25":  freq,
		"8.5":   freq,
		"8.75":  freq,
		"9":     freq,
		"9.25":  freq,
		"9.5":   freq,
		"9.75":  freq,
		"10":    freq,
		"10.25": freq,
		"10.5":  freq,
		"10.75": freq,
		"11":    freq,
		"11.25": freq,
		"11.5":  freq,
		"11.75": freq,
		"12":    freq,
		"12.25": freq,
		"12.5":  freq,
		"12.75": freq,
		"13":    freq,
		"13.25": freq,
		"13.5":  freq,
		"13.75": freq,
		"14":    freq,
		"14.25": freq,
		"14.5":  freq,
		"14.75": freq,
		"15":    freq,
		"15.25": freq,
		"15.5":  freq,
		"15.75": freq,
		"16":    freq,
		"16.25": freq,
		"16.5":  freq,
		"16.75": freq,
		"17":    freq,
		"17.25": freq,
		"17.5":  freq,
	}
}
