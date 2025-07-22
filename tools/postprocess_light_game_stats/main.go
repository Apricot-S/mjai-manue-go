package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/Apricot-S/mjai-manue-go/configs"
	"github.com/go-json-experiment/json"
)

type Input struct {
	ScoreStats map[string]map[string]float64 `json:"scoreStats"`
}

var kyokus = []string{"E1", "E2", "E3", "E4", "S1", "S2", "S3", "S4"}

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

func computeRatios(scoreStats map[string]map[string]float64) map[string]map[int]float64 {
	ratiosMap := make(map[string]map[int]float64)

	for key, freqs := range scoreStats {
		intFreqs := make(map[int]float64)
		total := 0.0

		for scoreStr, freq := range freqs {
			scoreDiff, err := strconv.Atoi(scoreStr)
			if err != nil {
				log.Printf("invalid score key %q: %v", scoreStr, err)
				continue
			}
			intFreqs[scoreDiff] = freq
			total += freq
		}

		ratios := make(map[int]float64)
		for scoreDiff, freq := range intFreqs {
			ratios[scoreDiff] = freq / total
		}

		ratiosMap[key] = ratios
	}

	return ratiosMap
}

func computeWinProbabilities(kyokus []string, ratiosMap map[string]map[int]float64) map[string]map[string]float64 {
	winProbsMap := make(map[string]map[string]float64)

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
				winProbs := buildWinProbabilities(relativeScoreRatios, delta)
				key := fmt.Sprintf("%s,%d,%d", kyoku, i, j)
				winProbsMap[key] = winProbs
			}
		}
	}

	return winProbsMap
}

func buildWinProbabilities(relativeRatios map[int]float64, delta int) map[string]float64 {
	winProbs := make(map[string]float64)

	relativeScores := make([]int, 0, len(relativeRatios))
	for score := range relativeRatios {
		relativeScores = append(relativeScores, score)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(relativeScores)))

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
		fmt.Fprintf(os.Stderr, "Usage: %s <light_game_stats.json>\n", os.Args[0])
		os.Exit(2)
	}

	stats, err := loadStatsFromFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	ratiosMap := computeRatios(stats.ScoreStats)
	winProbsMap := computeWinProbabilities(kyokus, ratiosMap)

	output := configs.LightGameStats{WinProbsMap: winProbsMap}
	if err := json.MarshalWrite(os.Stdout, output); err != nil {
		log.Fatalf("failed to output result: %v", err)
	}
	fmt.Print("\n")
}
