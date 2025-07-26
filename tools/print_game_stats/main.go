package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Apricot-S/mjai-manue-go/configs"
	"github.com/go-json-experiment/json"
)

func printNumTurnsDistribution(stats configs.GameStats) {
	if stats.NumTurnsDistribution != nil {
		fmt.Println("numTurnsDistribution:")
		for i, n := range stats.NumTurnsDistribution {
			fmt.Printf("  %2d: %.3f\n", i, n)
		}
	}
}

func printYamitenStats(stats configs.GameStats) {
	const maxTurn = 18
	if stats.YamitenStats != nil {
		fmt.Println("yamitenStats:")
		for i := range maxTurn {
			line := fmt.Sprintf("  %2d: ", i)
			for j := range 5 {
				key := fmt.Sprintf("%d,%d", i, j)
				stat, ok := stats.YamitenStats[key]
				if !ok || stat.Total == 0 {
					stat = &configs.YamitenStat{}
				}

				ratio := 0.0
				if stat.Total > 0 {
					ratio = float64(stat.Tenpai) / float64(stat.Total)
				}
				line += fmt.Sprintf("%.3f(%5d/%5d)  ", ratio, stat.Tenpai, stat.Total)
			}
			fmt.Println(line)
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <PATH TO game_stats.json>\n", os.Args[0])
		os.Exit(2)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	var stats configs.GameStats
	if err := json.Unmarshal(data, &stats); err != nil {
		log.Fatalf("failed to unmarshal GameStats: %v", err)
	}

	printNumTurnsDistribution(stats)
	fmt.Println()
	printYamitenStats(stats)
	fmt.Println()
}
