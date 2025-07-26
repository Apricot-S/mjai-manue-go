package main

import (
	"cmp"
	"fmt"
	"log"
	"maps"
	"os"
	"slices"
	"strconv"

	"github.com/Apricot-S/mjai-manue-go/configs"
	"github.com/go-json-experiment/json"
)

type Stats = map[string]map[string]float64

type Input struct {
	ScoreStats Stats `json:"scoreStats"`
}

type RatiosMapEntry = map[int]float64
type RatiosMap = map[string]RatiosMapEntry

type WinProbsMapEntry = map[string]float64
type WinProbsMap = map[string]WinProbsMapEntry

func loadStatsFromFile(path string) (*Input, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var in Input
	if err := json.Unmarshal(data, &in); err != nil {
		return nil, err
	}
	return &in, nil
}

func computeRatios(scoreStats Stats) (RatiosMap, error) {
	ratiosMap := make(RatiosMap)

	for key, freqs := range scoreStats {
		total := 0.0
		for _, freq := range freqs {
			total += freq
		}

		ratiosMap[key] = make(RatiosMapEntry)
		for scoreDiffStr, freq := range freqs {
			scoreDiff, err := strconv.Atoi(scoreDiffStr)
			if err != nil {
				return nil, fmt.Errorf("invalid score key %q: %w", scoreDiffStr, err)
			}
			ratiosMap[key][scoreDiff] = freq / total
		}
	}

	return ratiosMap, nil
}

func computeWinProbs(ratiosMap RatiosMap) WinProbsMap {
	winProbsMap := make(WinProbsMap)

	var kyokus = []string{"E1", "E2", "E3", "E4", "S1", "S2", "S3", "S4"}
	for _, kyoku := range kyokus {
		for i := range 4 {
			for j := range 4 {
				if i == j {
					continue
				}
				keyI := fmt.Sprintf("%s,%d", kyoku, i)
				keyJ := fmt.Sprintf("%s,%d", kyoku, j)

				relativeScoreRatios := make(map[int]float64)
				for scoreDiff1, ratio1 := range ratiosMap[keyI] {
					for scoreDiff2, ratio2 := range ratiosMap[keyJ] {
						relative := scoreDiff1 - scoreDiff2
						relativeScoreRatios[relative] += ratio1 * ratio2
					}
				}

				delta := 0
				if i <= j {
					delta = 100
				}
				winProbs := buildEntry(relativeScoreRatios, delta)
				key := fmt.Sprintf("%s,%d,%d", kyoku, i, j)
				winProbsMap[key] = winProbs
			}
		}
	}

	return winProbsMap
}

func buildEntry(relativeRatios map[int]float64, delta int) WinProbsMapEntry {
	winProbs := make(WinProbsMapEntry)

	relativeScores := slices.SortedFunc(maps.Keys(relativeRatios), func(a, b int) int {
		return cmp.Compare(b, a)
	})

	accumProb := 0.0
	for _, relative := range relativeScores {
		accumProb += relativeRatios[relative]
		key := strconv.Itoa(delta - relative)
		winProbs[key] = accumProb
	}

	return winProbs
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <PATH TO light_game_stats.json>\n", os.Args[0])
		os.Exit(2)
	}

	input, err := loadStatsFromFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	ratiosMap, err := computeRatios(input.ScoreStats)
	if err != nil {
		log.Fatal(err)
	}
	winProbsMap := computeWinProbs(ratiosMap)

	output := configs.LightGameStats{WinProbsMap: winProbsMap}
	if err := json.MarshalWrite(os.Stdout, output, json.Deterministic(true)); err != nil {
		log.Fatalf("failed to output result: %v", err)
	}
	fmt.Print("\n")
}
