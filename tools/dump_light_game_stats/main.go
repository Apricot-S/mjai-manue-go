package main

import (
	"encoding/json/v2"
	"fmt"
	"log"
	"os"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/tools/internal/archive"
	"github.com/schollz/progressbar/v3"
)

type scoreStats = map[string]map[string]int

type output struct {
	ScoreStats scoreStats `json:"scoreStats"`
}

type roundSnapshot struct {
	name   string
	scores [common.NumPlayers]int
}

type scoreCounter struct {
	stats     scoreStats
	snapshots []roundSnapshot
	scores    [common.NumPlayers]int
}

func newScoreCounter() *scoreCounter {
	c := &scoreCounter{stats: make(scoreStats)}
	c.reset()
	return c
}

func (c *scoreCounter) reset() {
	c.snapshots = nil
	for i := range c.scores {
		c.scores[i] = archive.InitialScore
	}
}

func (c *scoreCounter) onMessage(msg inbound.Message) error {
	if _, ok := msg.(*inbound.Error); ok {
		return fmt.Errorf("error in the log")
	}
	if _, ok := msg.(*inbound.StartGame); ok {
		c.reset()
		return nil
	}

	c.updateScores(msg)

	switch m := msg.(type) {
	case *inbound.StartKyoku:
		c.snapshots = append(c.snapshots, roundSnapshot{
			name:   fmt.Sprintf("%s%d", m.Bakaze, m.Kyoku),
			scores: c.scores,
		})
	case *inbound.EndGame:
		c.finishGame()
	}
	return nil
}

func (c *scoreCounter) updateScores(msg inbound.Message) {
	switch m := msg.(type) {
	case *inbound.StartKyoku:
		if len(m.Scores) == common.NumPlayers {
			c.scores = [common.NumPlayers]int(m.Scores)
		}
	case *inbound.Hora:
		if len(m.Scores) == common.NumPlayers {
			c.scores = [common.NumPlayers]int(m.Scores)
		}
	case *inbound.Ryukyoku:
		if len(m.Scores) == common.NumPlayers {
			c.scores = [common.NumPlayers]int(m.Scores)
		}
	case *inbound.EndGame:
		if len(m.Scores) == common.NumPlayers {
			c.scores = [common.NumPlayers]int(m.Scores)
		}
	}
}

func (c *scoreCounter) finishGame() {
	for playerID := range common.NumPlayers {
		for _, snapshot := range c.snapshots {
			scoreDiff := c.scores[playerID] - snapshot.scores[playerID]
			// Mjai logs treated by this tool use player 0 as chicha, so playerID is the relative seat position.
			key := fmt.Sprintf("%s,%d", snapshot.name, playerID)
			if _, ok := c.stats[key]; !ok {
				c.stats[key] = make(map[string]int)
			}
			c.stats[key][fmt.Sprintf("%d", scoreDiff)]++
		}
	}
}

func run(patterns []string) (*output, error) {
	paths, err := archive.GlobAll(patterns)
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("no input files matched")
	}

	bar := progressbar.Default(int64(len(paths)))
	counter := newScoreCounter()
	a := archive.NewArchive()
	err = a.PlayPaths(paths, archive.Handlers{
		OnMessage: counter.onMessage,
		OnFileDone: func(string) error {
			return bar.Add(1)
		},
	})
	if err != nil {
		return nil, err
	}
	return &output{ScoreStats: counter.stats}, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <LOG_GLOB_PATTERNS>...\n", os.Args[0])
		os.Exit(2)
	}

	output, err := run(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	if err := json.MarshalWrite(os.Stdout, output, json.Deterministic(true)); err != nil {
		log.Fatalf("failed to output result: %v", err)
	}
	fmt.Println()
}
