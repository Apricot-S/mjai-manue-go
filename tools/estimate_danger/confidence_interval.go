package main

import (
	"math/rand/v2"
	"slices"
)

const numTries = 1000
const numTriesFloat = float64(numTries)

// CalculateConfidenceInterval calculates confidence intervals using bootstrap resampling.
func CalculateConfidenceInterval(samples []float64, min, max, confLevel float64) (lower, upper float64) {
	averages := make([]float64, numTries)
	sampleSize := len(samples)
	resamplePoolSize := sampleSize + 2
	for i := range numTries {
		sum := 0.0
		for range resamplePoolSize {
			idx := rand.IntN(resamplePoolSize)
			switch idx {
			case sampleSize:
				sum += min
			case sampleSize + 1:
				sum += max
			default:
				sum += samples[idx]
			}
		}
		averages[i] = sum / float64(resamplePoolSize)
	}
	slices.Sort(averages)

	margin := (1.0 - confLevel) / 2.0
	lower = averages[int(numTriesFloat*margin)]
	upper = averages[int(numTriesFloat*(1.0-margin))]
	return
}
