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

type Stats struct {
	ScoreStats map[string]map[string]float64 `json:"scoreStats"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <light_game_stats.json>\n", os.Args[0])
		os.Exit(2)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	var stats Stats
	if err := json.Unmarshal(data, &stats); err != nil {
		log.Fatal(err)
	}

	ratiosMap := make(map[string]map[int]float64)
	for key, freqs := range stats.ScoreStats {
		total := 0.0
		intFreqs := make(map[int]float64)
		for scoreStr, freq := range freqs {
			scoreDiff, _ := strconv.Atoi(scoreStr)
			intFreqs[scoreDiff] = freq
			total += freq
		}
		ratios := make(map[int]float64)
		for scoreDiff, freq := range intFreqs {
			ratios[scoreDiff] = freq / total
		}
		ratiosMap[key] = ratios
	}

	winProbsMap := make(map[string]map[string]float64)
	kyokus := []string{"E1", "E2", "E3", "E4", "S1", "S2", "S3", "S4"}

	for _, kyokuName := range kyokus {
		for i := range 4 {
			for j := range 4 {
				if i == j {
					continue
				}
				keyI := fmt.Sprintf("%s,%d", kyokuName, i)
				keyJ := fmt.Sprintf("%s,%d", kyokuName, j)

				relativeScoreRatios := make(map[int]float64)
				for scoreDiff1, ratio1 := range ratiosMap[keyI] {
					for scoreDiff2, ratio2 := range ratiosMap[keyJ] {
						relative := scoreDiff1 - scoreDiff2
						relativeScoreRatios[relative] += ratio1 * ratio2
					}
				}

				var relativeScores []int
				for score := range relativeScoreRatios {
					relativeScores = append(relativeScores, score)
				}
				sort.Sort(sort.Reverse(sort.IntSlice(relativeScores)))

				delta := 0
				if i <= j {
					delta = 100
				}
				winProbs := make(map[string]float64)
				accumProb := 0.0
				for _, relative := range relativeScores {
					accumProb += relativeScoreRatios[relative]
					winProbs[strconv.Itoa(delta-relative)] = accumProb
				}
				winProbsMap[fmt.Sprintf("%s,%d,%d", kyokuName, i, j)] = winProbs
			}
		}
	}

	output := configs.LightGameStats{
		WinProbsMap: winProbsMap,
	}

	if err := json.MarshalWrite(os.Stdout, output); err != nil {
		log.Fatalf("failed to output result: %v", err)
	}
	fmt.Print("\n")
}
