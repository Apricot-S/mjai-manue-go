package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Apricot-S/mjai-manue-go/configs"
	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/game"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/protocol/mjai"
	"github.com/Apricot-S/mjai-manue-go/tools/shared"
	"github.com/go-json-experiment/json"
)

// TODO Support kokushimuso and chitoitsu
func isTenpai(actor *base.Player) (bool, error) {
	tehaiSet, err := base.NewPaiSet(actor.Tehais())
	if err != nil {
		return false, fmt.Errorf("failed to get the hand: %w", err)
	}
	return game.IsTenpaiGeneral(tehaiSet)
}

type Counter interface {
	OnAction(action inbound.Event, g game.StateViewer) error
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

func (bc *BasicCounter) OnAction(action inbound.Event, g game.StateViewer) error {
	switch a := action.(type) {
	case *inbound.Hora:
		bc.NumHoras++
		if a.Actor == a.Target {
			bc.NumTsumoHoras++
		}
		if a.HoraPoints != nil {
			bc.TotalHoraPoints += *a.HoraPoints
		}
	case *inbound.Ryukyoku:
		bc.NumRyukyokus++
	case *inbound.EndKyoku:
		bc.NumKyokus++
		idx := (game.NumInitPipais - g.NumPipais()) / 4
		bc.NumTurnsFreqs[idx]++
	}

	return nil
}

type HoraPointsCounter struct {
	KoFreqs  map[string]int
	OyaFreqs map[string]int
}

func NewHoraPointsCounter() *HoraPointsCounter {
	return &HoraPointsCounter{
		KoFreqs:  map[string]int{"total": 0},
		OyaFreqs: map[string]int{"total": 0},
	}
}

func (hpc *HoraPointsCounter) OnAction(action inbound.Event, g game.StateViewer) error {
	hora, ok := action.(*inbound.Hora)
	if !ok {
		return nil
	}

	var freqs map[string]int
	if hora.Actor == g.Oya().ID() {
		freqs = hpc.OyaFreqs
	} else {
		freqs = hpc.KoFreqs
	}

	freqs["total"]++
	if hora.HoraPoints != nil {
		key := strconv.Itoa(*hora.HoraPoints)
		freqs[key]++
	}

	return nil
}

type YamitenCounter struct {
	Stats map[string]*configs.YamitenStat
}

func NewYamitenCounter() *YamitenCounter {
	return &YamitenCounter{
		Stats: make(map[string]*configs.YamitenStat),
	}
}

func (yc *YamitenCounter) OnAction(action inbound.Event, g game.StateViewer) error {
	dahai, ok := action.(*inbound.Dahai)
	if !ok {
		return nil
	}

	actor := g.Players()[dahai.Actor]
	if actor.ReachState() != base.NotReach {
		return nil
	}

	numTurns := g.NumPipais() / 4
	numFuros := len(actor.Furos())
	key := fmt.Sprintf("%d,%d", numTurns, numFuros)

	if _, ok := yc.Stats[key]; !ok {
		yc.Stats[key] = new(configs.YamitenStat)
	}
	yc.Stats[key].Total++

	isTenpai, err := isTenpai(&actor)
	if err != nil {
		return err
	}
	if isTenpai {
		yc.Stats[key].Tenpai++
	}

	return nil
}

type RyukyokuTenpaiCounter struct {
	Stats       *configs.RyukyokuTenpaiStat
	tenpaiTurns [4]*float64
}

func NewRyukyokuTenpaiCounter() *RyukyokuTenpaiCounter {
	tenpaiTurnDistribution := make(map[string]int)
	for i := 0.0; i <= game.FinalTurn; i += 1.0 / 4.0 {
		key := strconv.FormatFloat(i, 'f', -1, 64)
		tenpaiTurnDistribution[key] = 0
	}

	return &RyukyokuTenpaiCounter{
		Stats: &configs.RyukyokuTenpaiStat{
			Total:                  0,
			Tenpai:                 0,
			Noten:                  0,
			TenpaiTurnDistribution: tenpaiTurnDistribution,
		},
		tenpaiTurns: [4]*float64{nil, nil, nil, nil},
	}
}

func (rtc *RyukyokuTenpaiCounter) OnAction(action inbound.Event, g game.StateViewer) error {
	switch a := action.(type) {
	case *inbound.StartKyoku:
		rtc.tenpaiTurns = [4]*float64{nil, nil, nil, nil}
	case *inbound.Dahai:
		actor := g.Players()[a.Actor]
		isTenpai, err := isTenpai(&actor)
		if err != nil {
			return err
		}
		if rtc.tenpaiTurns[a.Actor] == nil && isTenpai {
			turn := g.Turn()
			rtc.tenpaiTurns[a.Actor] = &turn
		}
	case *inbound.Ryukyoku:
		for _, player := range g.Players() {
			isTenpai, err := isTenpai(&player)
			if err != nil {
				return err
			}

			rtc.Stats.Total++
			// Temporary handling:
			// In the case of Nine Different Terminals and Honors,
			// the round may end in a draw before any discard, so it's possible for
			// isTenpai && rtc.tenpaiTurns[player.ID()] == nil.
			// However, since there's no point exchange, all players are treated as Noten.
			if isTenpai && rtc.tenpaiTurns[player.ID()] != nil {
				rtc.Stats.Tenpai++
				key := strconv.FormatFloat(*rtc.tenpaiTurns[player.ID()], 'f', -1, 64)
				rtc.Stats.TenpaiTurnDistribution[key]++
			} else {
				rtc.Stats.Noten++
			}
		}
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
	archive := shared.NewArchive(paths, &mjai.MjaiAdapter{})

	basic := NewBasicCounter()
	horaPoints := NewHoraPointsCounter()
	yamiten := NewYamitenCounter()
	ryukyokuTenpai := NewRyukyokuTenpaiCounter()
	counters := []Counter{basic, horaPoints, yamiten, ryukyokuTenpai}

	onAction := func(action inbound.Event) error {
		if _, ok := action.(*inbound.Error); ok {
			return fmt.Errorf("error in the log")
		}
		for _, counter := range counters {
			if err := counter.OnAction(action, archive.StateViewer()); err != nil {
				return err
			}
		}
		return nil
	}

	if err := archive.Play(onAction); err != nil {
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
		KoHoraPointsFreqs:    horaPoints.KoFreqs,
		OyaHoraPointsFreqs:   horaPoints.OyaFreqs,
		YamitenStats:         yamiten.Stats,
		RyukyokuTenpaiStat:   ryukyokuTenpai.Stats,
	}
	if err := json.MarshalWrite(os.Stdout, output, json.Deterministic(true)); err != nil {
		log.Fatalf("failed to output result: %v", err)
	}
	fmt.Print("\n")
}
