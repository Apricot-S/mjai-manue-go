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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <game_stats.json>\n", os.Args[0])
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
}
