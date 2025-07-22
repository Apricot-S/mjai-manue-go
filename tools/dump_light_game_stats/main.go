package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/message"
	"github.com/Apricot-S/mjai-manue-go/tools/shared"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Stats = map[string]map[int]int

type ActionWithScores struct {
	message.Message
	Scores []int `json:"scores"`
}

type KyokuStat struct {
	KyokuName string
	Scores    []int
}

type Output struct {
	ScoreStats Stats `json:"scoreStats"`
}

type ScoreCounter struct {
	scores     []int
	stats      Stats
	kyokuStats []KyokuStat
	chichaId   int
}

func NewScoreCounter() *ScoreCounter {
	return &ScoreCounter{
		stats:      make(Stats),
		kyokuStats: []KyokuStat{},
		scores:     nil,
		chichaId:   0,
	}
}

func (sc *ScoreCounter) OnAction(action jsontext.Value) error {
	var msg ActionWithScores
	if err := json.Unmarshal(action, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	if len(msg.Scores) != 0 {
		sc.scores = msg.Scores
	}

	switch msg.Type {
	case message.TypeStartGame:
		sc.scores = []int{25000, 25000, 25000, 25000}
		sc.kyokuStats = []KyokuStat{}
	case message.TypeStartKyoku:
		var e message.StartKyoku
		if err := json.Unmarshal(action, &e); err != nil {
			return fmt.Errorf("failed to unmarshal start_kyoku: %w", err)
		}

		kyokuName := fmt.Sprintf("%s%d", e.Bakaze, e.Kyoku)
		snapshot := KyokuStat{
			KyokuName: kyokuName,
			Scores:    slices.Clone(sc.scores),
		}
		sc.kyokuStats = append(sc.kyokuStats, snapshot)
	case message.TypeEndGame:
		for playerId := range 4 {
			position := getDistance(playerId, sc.chichaId)
			for _, stat := range sc.kyokuStats {
				scoreDiff := sc.scores[playerId] - stat.Scores[playerId]
				key := fmt.Sprintf("%s,%d", stat.KyokuName, position)
				if _, ok := sc.stats[key]; !ok {
					sc.stats[key] = make(map[int]int)
				}
				sc.stats[key][scoreDiff]++
			}
		}
	}

	return nil
}

func getDistance(playerId1, playerId2 int) int {
	return (4 + playerId1 - playerId2) % 4
}

func GlobAll(patterns []string) ([]string, error) {
	var result []string
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid glob pattern %q: %w", pattern, err)
		}
		result = slices.Concat(result, matches)
	}
	return result, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <log_glob_patterns...>\n", os.Args[0])
		os.Exit(2)
	}

	args := os.Args[1:]
	paths, err := GlobAll(args)
	if err != nil {
		log.Fatalf("error in glob: %v", err)
	}
	a := shared.NewArchive(paths)
	counter := NewScoreCounter()

	onAction := func(action jsontext.Value) error {
		var msg message.Message
		if err := json.Unmarshal(action, &msg); err != nil {
			return err
		}
		if msg.Type == message.TypeError {
			return fmt.Errorf("error in the log")
		}
		return counter.OnAction(action)
	}

	if err := a.PlayLight(onAction); err != nil {
		log.Fatalf("error in processing log: %v", err)
	}

	output := Output{ScoreStats: counter.stats}
	if err := json.MarshalWrite(os.Stdout, output, json.Deterministic(true)); err != nil {
		log.Fatalf("failed to output result: %v", err)
	}
	fmt.Print("\n")
}
