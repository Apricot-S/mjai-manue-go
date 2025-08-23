package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Apricot-S/mjai-manue-go/internal/game"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/protocol/mjai"
	"github.com/Apricot-S/mjai-manue-go/tools/shared"
	"github.com/go-json-experiment/json"
)

var InitScores = [game.NumPlayers]int{game.InitScore, game.InitScore, game.InitScore, game.InitScore}

type Stats = map[string]map[int]int

type KyokuStat struct {
	KyokuName string
	Scores    [game.NumPlayers]int
}

type Output struct {
	ScoreStats Stats `json:"scoreStats"`
}

type ScoreCounter struct {
	scores     [game.NumPlayers]int
	stats      Stats
	kyokuStats []KyokuStat
	chichaId   int
}

func NewScoreCounter() *ScoreCounter {
	return &ScoreCounter{
		scores:     InitScores,
		stats:      make(Stats),
		kyokuStats: []KyokuStat{},
		chichaId:   0,
	}
}

func (sc *ScoreCounter) OnAction(action inbound.Event) error {
	// Get scores
	switch a := action.(type) {
	case *inbound.StartKyoku:
		if a.Scores != nil {
			sc.scores = *a.Scores
		}
	case *inbound.Hora:
		if a.Scores != nil {
			sc.scores = *a.Scores
		}
	case *inbound.Ryukyoku:
		if a.Scores != nil {
			sc.scores = *a.Scores
		}
	}

	switch a := action.(type) {
	case *inbound.StartGame:
		sc.scores = InitScores
		sc.kyokuStats = []KyokuStat{}
	case *inbound.StartKyoku:
		kyokuName := fmt.Sprintf("%s%d", a.Bakaze.ToString(), a.Kyoku)
		snapshot := KyokuStat{
			KyokuName: kyokuName,
			Scores:    sc.scores,
		}
		sc.kyokuStats = append(sc.kyokuStats, snapshot)
	case *inbound.EndGame:
		for playerId := range game.NumPlayers {
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

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <log_glob_patterns...>\n", os.Args[0])
		os.Exit(2)
	}

	args := os.Args[1:]
	paths, err := shared.GlobAll(args)
	if err != nil {
		log.Fatalf("error in glob: %v", err)
	}
	archive := shared.NewArchive(paths, mjai.Adapter)
	counter := NewScoreCounter()

	onAction := func(action inbound.Event) error {
		if _, ok := action.(*inbound.Error); ok {
			return fmt.Errorf("error in the log")
		}
		return counter.OnAction(action)
	}

	if err := archive.PlayLight(onAction); err != nil {
		log.Fatalf("error in processing log: %v", err)
	}

	output := Output{ScoreStats: counter.stats}
	if err := json.MarshalWrite(os.Stdout, output, json.Deterministic(true)); err != nil {
		log.Fatalf("failed to output result: %v", err)
	}
	fmt.Print("\n")
}
