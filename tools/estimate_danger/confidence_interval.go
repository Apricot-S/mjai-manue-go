package main

import (
	"math/rand/v2"
	"slices"
)

const numBootstrapTries = 1000

// CalculateConfidenceInterval ports mjai's confidence interval calculation as-is.
//
// The original calls this bootstrap resampling, but it resamples from the input
// samples plus artificial min/max boundary values. That is unusual for a
// standard bootstrap confidence interval; keep the behavior for compatibility.
// See: https://github.com/gimite/mjai/blob/master/lib/mjai/confidence_interval.rb
func CalculateConfidenceInterval(samples []float64, minValue, maxValue, confLevel float64) (lower, upper float64) {
	averages := make([]float64, numBootstrapTries)
	sampleSize := len(samples)
	resamplePoolSize := sampleSize + 2
	rng := rand.New(rand.NewPCG(0, 0))

	for i := range numBootstrapTries {
		sum := 0.0
		for range resamplePoolSize {
			idx := rng.IntN(resamplePoolSize)
			switch idx {
			case sampleSize:
				sum += minValue
			case sampleSize + 1:
				sum += maxValue
			default:
				sum += samples[idx]
			}
		}
		averages[i] = sum / float64(resamplePoolSize)
	}
	slices.Sort(averages)

	margin := (1.0 - confLevel) / 2.0
	lower = averages[int(float64(numBootstrapTries)*margin)]
	upper = averages[int(float64(numBootstrapTries)*(1.0-margin))]
	return lower, upper
}
