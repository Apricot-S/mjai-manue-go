package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Apricot-S/mjai-manue-go/configs"
	"github.com/Apricot-S/mjai-manue-go/internal/game"
	"github.com/Apricot-S/mjai-manue-go/internal/message"
	"github.com/Apricot-S/mjai-manue-go/tools/shared"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// TODO Support kokushimuso and chitoitsu
func isTenpai(actor *game.Player) (bool, error) {
	tehaiSet, err := game.NewPaiSet(actor.Tehais())
	if err != nil {
		return false, fmt.Errorf("failed to get the hand: %w", err)
	}
	shantenNumber, _, err := game.AnalyzeShantenWithOption(tehaiSet, 0, 0)
	if err != nil {
		return false, fmt.Errorf("failed to calculate shanten number: %w", err)
	}
	return shantenNumber <= 0, nil
}

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

func (hpc *HoraPointsCounter) OnAction(action jsontext.Value, g game.StateViewer) error {
	var msg message.Message
	if err := json.Unmarshal(action, &msg); err != nil {
		return err
	}

	if msg.Type != message.TypeHora {
		return nil
	}

	var hora message.Hora
	if err := json.Unmarshal(action, &hora); err != nil {
		return fmt.Errorf("failed to unmarshal hora: %w", err)
	}

	var freqs map[string]int
	if hora.Actor == g.Oya().ID() {
		freqs = hpc.OyaFreqs
	} else {
		freqs = hpc.KoFreqs
	}

	freqs["total"]++
	key := strconv.Itoa(hora.HoraPoints)
	freqs[key]++

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

func (yc *YamitenCounter) OnAction(action jsontext.Value, g game.StateViewer) error {
	var msg message.Message
	if err := json.Unmarshal(action, &msg); err != nil {
		return err
	}

	if msg.Type != message.TypeDahai {
		return nil
	}

	var dahai message.Dahai
	if err := json.Unmarshal(action, &dahai); err != nil {
		return fmt.Errorf("failed to unmarshal dahai: %w", err)
	}

	actor := g.Players()[dahai.Actor]
	if actor.ReachState() != game.NotReach {
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

func (rtc *RyukyokuTenpaiCounter) OnAction(action jsontext.Value, g game.StateViewer) error {
	var msg message.Message
	if err := json.Unmarshal(action, &msg); err != nil {
		return err
	}

	switch msg.Type {
	case message.TypeStartKyoku:
		rtc.tenpaiTurns = [4]*float64{nil, nil, nil, nil}
	case message.TypeDahai:
		var dahai message.Dahai
		if err := json.Unmarshal(action, &dahai); err != nil {
			return fmt.Errorf("failed to unmarshal dahai: %w", err)
		}

		actor := g.Players()[dahai.Actor]
		isTenpai, err := isTenpai(&actor)
		if err != nil {
			return err
		}
		if rtc.tenpaiTurns[dahai.Actor] == nil && isTenpai {
			// Note:
			// This branch isn't reached if the player entered Tenpai once,
			// then broke Tenpai before the end of the round,
			// but it's unclear if this is intentional.
			turn := g.Turn()
			rtc.tenpaiTurns[dahai.Actor] = &turn
		}
	case message.TypeRyukyoku:
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
	a := shared.NewArchive(paths)

	basic := NewBasicCounter()
	horaPoints := NewHoraPointsCounter()
	yamiten := NewYamitenCounter()
	ryukyokuTenpai := NewRyukyokuTenpaiCounter()
	counters := []Counter{basic, horaPoints, yamiten, ryukyokuTenpai}

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
