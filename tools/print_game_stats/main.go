package main

import (
	"encoding/json/v2"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/Apricot-S/mjai-manue-go/configs"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round"
)

const (
	maxTurn    = 18
	maxNumMeld = 4
)

func loadStatsFromFile(path string) (*configs.GameStats, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var stats configs.GameStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

func printNumTurnsDistribution(w io.Writer, stats configs.GameStats) {
	if stats.NumTurnsDistribution == nil {
		return
	}

	fmt.Fprintln(w, "numTurnsDistribution:")
	for i, n := range stats.NumTurnsDistribution {
		fmt.Fprintf(w, "  %2d: %.3f\n", i, n)
	}
	fmt.Fprintln(w)
}

func printYamitenStats(w io.Writer, stats configs.GameStats) {
	if stats.YamitenStats == nil {
		return
	}

	fmt.Fprintln(w, "yamitenStats:")
	for remainTurns := range maxTurn {
		fmt.Fprintf(w, "  %2d: ", remainTurns)
		for numMelds := range maxNumMeld + 1 {
			key := fmt.Sprintf("%d,%d", remainTurns, numMelds)
			stat := stats.YamitenStats[key]
			ratio := ratioOrNaN(stat.Tenpai, stat.Total)
			fmt.Fprintf(w, "%s(%5d/%5d)  ", formatRatio(ratio), stat.Tenpai, stat.Total)
		}
		fmt.Fprintln(w)
	}
	fmt.Fprintln(w)
}

func printRyukyokuTenpaiStats(w io.Writer, stats configs.GameStats) {
	if stats.RyukyokuTenpaiStat.TenpaiTurnDistribution == nil {
		return
	}

	fmt.Fprintln(w, "ryukyokuTenpaiStat:")
	for turn := 0.0; turn <= round.FinalTurn; turn += 0.25 {
		key := strconv.FormatFloat(turn, 'f', -1, 64)
		freq := stats.RyukyokuTenpaiStat.TenpaiTurnDistribution[key]
		ratio := ratioOrNaN(freq, stats.RyukyokuTenpaiStat.Total)
		fmt.Fprintf(w, "  %5.2f: %s (%d)\n", turn, formatRatio(ratio), freq)
	}
	ratio := ratioOrNaN(stats.RyukyokuTenpaiStat.Noten, stats.RyukyokuTenpaiStat.Total)
	fmt.Fprintf(w, "  noten: %s (%d)\n", formatRatio(ratio), stats.RyukyokuTenpaiStat.Noten)
	fmt.Fprintln(w)
}

func printStats(w io.Writer, stats configs.GameStats) {
	printNumTurnsDistribution(w, stats)
	printYamitenStats(w, stats)
	printRyukyokuTenpaiStats(w, stats)
}

func ratioOrNaN(numerator int, denominator int) float64 {
	if denominator == 0 {
		return math.NaN()
	}
	return float64(numerator) / float64(denominator)
}

func formatRatio(ratio float64) string {
	if math.IsNaN(ratio) {
		return "  NaN"
	}
	return fmt.Sprintf("%.3f", ratio)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <PATH TO game_stats.json>\n", os.Args[0])
		os.Exit(2)
	}

	stats, err := loadStatsFromFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	printStats(os.Stdout, *stats)
}
