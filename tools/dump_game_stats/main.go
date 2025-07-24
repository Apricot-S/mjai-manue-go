package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Apricot-S/mjai-manue-go/configs"
	"github.com/Apricot-S/mjai-manue-go/internal/game"
	"github.com/Apricot-S/mjai-manue-go/internal/message"
	"github.com/Apricot-S/mjai-manue-go/tools/shared"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Counter interface {
	OnAction(action jsontext.Value, g game.StateViewer) error
}

const maxTurn = 18

type BasicCounter struct {
	NumKyokus       int
	NumTurnsFreqs   [maxTurn]int
	NumRyukyokus    int
	TotalHoraPoints int
	NumHoras        int
	NumTsumoHoras   int
}

func NewBasicCounter() *BasicCounter {
	return &BasicCounter{}
}

func (bc *BasicCounter) OnAction(action jsontext.Value, g game.StateViewer) error {
	var msg message.Message
	if err := json.Unmarshal(action, &msg); err != nil {
		return err
	}

	switch msg.Type {
	case message.TypeHora:
		bc.NumHoras++

		var h message.Hora
		if err := json.Unmarshal(action, &h); err != nil {
			return fmt.Errorf("failed to unmarshal hora: %w", err)
		}

		if h.Actor == h.Target {
			bc.NumTsumoHoras++
		}
		bc.TotalHoraPoints += h.HoraPoints
	case message.TypeRyukyoku:
		bc.NumRyukyokus++
	case message.TypeEndKyoku:
		bc.NumKyokus++
		idx := (game.NumInitPipais - g.NumPipais()) / 4
		bc.NumTurnsFreqs[idx]++
	}

	return nil
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
	a := shared.NewArchive(paths)

	basic := NewBasicCounter()
	counters := []Counter{
		basic,
	}

	onAction := func(action jsontext.Value) error {
		var msg message.Message
		if err := json.Unmarshal(action, &msg); err != nil {
			return err
		}
		if msg.Type == message.TypeError {
			return fmt.Errorf("error in the log")
		}
		for _, counter := range counters {
			if err := counter.OnAction(action, a.StateViewer()); err != nil {
				return err
			}
		}
		return nil
	}

	if err := a.Play(onAction); err != nil {
		log.Fatalf("error in processing log: %v", err)
	}

	numTurnsDistribution := make([]float64, maxTurn)
	for i, f := range basic.NumTurnsFreqs {
		numTurnsDistribution[i] = float64(f) / float64(basic.NumKyokus)
	}

	output := configs.GameStats{
		NumHoras:             basic.NumHoras,
		NumTsumoHoras:        basic.NumTsumoHoras,
		NumTurnsDistribution: numTurnsDistribution,
		RyukyokuRatio:        float64(basic.NumRyukyokus) / float64(basic.NumKyokus),
		AverageHoraPoints:    float64(basic.TotalHoraPoints) / float64(basic.NumHoras),
	}
	if err := json.MarshalWrite(os.Stdout, output, json.Deterministic(true)); err != nil {
		log.Fatalf("failed to output result: %v", err)
	}
	fmt.Print("\n")
}
