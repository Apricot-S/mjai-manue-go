package main

import (
	"cmp"
	"encoding/json/v2"
	"fmt"
	"log"
	"maps"
	"os"
	"slices"
	"strconv"

	"github.com/Apricot-S/mjai-manue-go/configs"
)

type scoreStats = map[string]map[string]float64

type input struct {
	ScoreStats scoreStats `json:"scoreStats"`
}

type ratiosMapEntry = map[int]float64
type ratiosMap = map[string]ratiosMapEntry
type winProbsMapEntry = map[string]float64
type winProbsMap = map[string]winProbsMapEntry

func loadStatsFromFile(path string) (*input, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var in input
	if err := json.Unmarshal(data, &in); err != nil {
		return nil, err
	}
	return &in, nil
}

func computeRatios(stats scoreStats) (ratiosMap, error) {
	result := make(ratiosMap)
	for key, freqs := range stats {
		total := 0.0
		for _, freq := range freqs {
			total += freq
		}
		if total == 0 {
			return nil, fmt.Errorf("scoreStats[%q] has zero total", key)
		}

		entry := make(ratiosMapEntry)
		for scoreDiffStr, freq := range freqs {
			scoreDiff, err := strconv.Atoi(scoreDiffStr)
			if err != nil {
				return nil, fmt.Errorf("invalid score key %q: %w", scoreDiffStr, err)
			}
			entry[scoreDiff] = freq / total
		}
		result[key] = entry
	}
	return result, nil
}

func computeWinProbs(ratios ratiosMap) winProbsMap {
	result := make(winProbsMap)

	rounds := []string{"E1", "E2", "E3", "E4", "S1", "S2", "S3", "S4"}
	for _, roundName := range rounds {
		for selfPosition := range 4 {
			for otherPosition := range 4 {
				if selfPosition == otherPosition {
					continue
				}

				selfKey := fmt.Sprintf("%s,%d", roundName, selfPosition)
				otherKey := fmt.Sprintf("%s,%d", roundName, otherPosition)
				relativeRatios := make(map[int]float64)
				for selfDiff, selfRatio := range ratios[selfKey] {
					for otherDiff, otherRatio := range ratios[otherKey] {
						relativeRatios[selfDiff-otherDiff] += selfRatio * otherRatio
					}
				}

				delta := 0
				if selfPosition <= otherPosition {
					delta = 100
				}
				key := fmt.Sprintf("%s,%d,%d", roundName, selfPosition, otherPosition)
				result[key] = buildEntry(relativeRatios, delta)
			}
		}
	}

	return result
}

func buildEntry(relativeRatios map[int]float64, delta int) winProbsMapEntry {
	result := make(winProbsMapEntry)
	relativeScores := slices.SortedFunc(maps.Keys(relativeRatios), func(a, b int) int {
		return cmp.Compare(b, a)
	})

	accumProb := 0.0
	for _, relative := range relativeScores {
		accumProb += relativeRatios[relative]
		result[strconv.Itoa(delta-relative)] = accumProb
	}
	return result
}

func run(inputPath string) (*configs.LightGameStats, error) {
	input, err := loadStatsFromFile(inputPath)
	if err != nil {
		return nil, err
	}
	ratios, err := computeRatios(input.ScoreStats)
	if err != nil {
		return nil, err
	}
	return &configs.LightGameStats{WinProbsMap: computeWinProbs(ratios)}, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <PATH TO score_stats.json>\n", os.Args[0])
		os.Exit(2)
	}

	output, err := run(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	if err := json.MarshalWrite(os.Stdout, output, json.Deterministic(true)); err != nil {
		log.Fatalf("failed to output result: %v", err)
	}
	fmt.Println()
}
