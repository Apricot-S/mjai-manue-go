package main

import (
	"encoding/json/v2"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Apricot-S/mjai-manue-go/configs"
	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/tools/internal/archive"
	"github.com/schollz/progressbar/v3"
)

const maxTurn = 18

type counter interface {
	onEvent(ev event.Event, state round.StateViewer) error
}

type basicCounter struct {
	numRounds       int
	numTurnFreqs    [maxTurn]int
	numDrawRounds   int
	totalWinPoints  int
	numWins         int
	numSelfDrawWins int
}

func (c *basicCounter) onEvent(ev event.Event, state round.StateViewer) error {
	switch ev := ev.(type) {
	case *event.Win:
		c.numWins++
		if ev.Actor() == ev.Target() {
			c.numSelfDrawWins++
		}
		c.totalWinPoints += ev.WinningPoints()
	case *event.DrawRound:
		c.numDrawRounds++
	case *event.EndRound:
		c.numRounds++
		turnIndex := (round.NumInitWall - state.NumLeftTiles()) / common.NumPlayers
		if 0 <= turnIndex && turnIndex < len(c.numTurnFreqs) {
			c.numTurnFreqs[turnIndex]++
		}
	}
	return nil
}

type winPointsCounter struct {
	koFreqs  map[string]int
	oyaFreqs map[string]int
}

func newWinPointsCounter() *winPointsCounter {
	return &winPointsCounter{
		koFreqs:  map[string]int{"total": 0},
		oyaFreqs: map[string]int{"total": 0},
	}
}

func (c *winPointsCounter) onEvent(ev event.Event, state round.StateViewer) error {
	win, ok := ev.(*event.Win)
	if !ok {
		return nil
	}

	freqs := c.koFreqs
	if win.Actor() == state.Dealer() {
		freqs = c.oyaFreqs
	}
	freqs["total"]++
	freqs[strconv.Itoa(win.WinningPoints())]++
	return nil
}

type yamitenCounter struct {
	stats map[string]configs.YamitenStat
}

func newYamitenCounter() *yamitenCounter {
	return &yamitenCounter{stats: make(map[string]configs.YamitenStat)}
}

func (c *yamitenCounter) onEvent(ev event.Event, state round.StateViewer) error {
	discard, ok := ev.(*event.Discard)
	if !ok {
		return nil
	}

	actor := state.Player(discard.Actor())
	if actor.RiichiState() != player.NotRiichi {
		return nil
	}

	key := fmt.Sprintf("%d,%d", state.NumLeftTiles()/common.NumPlayers, len(actor.Melds()))
	stat := c.stats[key]
	stat.Total++
	if isTenpai(actor) {
		stat.Tenpai++
	}
	c.stats[key] = stat
	return nil
}

type drawTenpaiCounter struct {
	stats       configs.RyukyokuTenpaiStat
	tenpaiTurns [common.NumPlayers]*float64
}

func newDrawTenpaiCounter() *drawTenpaiCounter {
	turnDistribution := make(map[string]int)
	for turn := 0.0; turn <= round.FinalTurn; turn += 0.25 {
		turnDistribution[strconv.FormatFloat(turn, 'f', -1, 64)] = 0
	}
	return &drawTenpaiCounter{
		stats: configs.RyukyokuTenpaiStat{
			TenpaiTurnDistribution: turnDistribution,
		},
	}
}

func (c *drawTenpaiCounter) onEvent(ev event.Event, state round.StateViewer) error {
	switch ev := ev.(type) {
	case *event.StartRound:
		c.tenpaiTurns = [common.NumPlayers]*float64{}
	case *event.Discard:
		actorIndex := ev.Actor().Index()
		if c.tenpaiTurns[actorIndex] == nil && isTenpai(state.Player(ev.Actor())) {
			c.tenpaiTurns[actorIndex] = new(state.Turn())
		}
	case *event.DrawRound:
		tenpais := ev.Tenpais()
		for playerID := range common.NumPlayers {
			c.stats.Total++
			if tenpais != nil && tenpais[playerID] {
				c.stats.Tenpai++
				if turn := c.tenpaiTurns[playerID]; turn != nil {
					key := strconv.FormatFloat(*turn, 'f', -1, 64)
					c.stats.TenpaiTurnDistribution[key]++
				}
			} else {
				c.stats.Noten++
			}
		}
	}
	return nil
}

func isTenpai(p player.PlayerViewer) bool {
	hand, ok := p.Hand()
	// Match the original dump_game_stats implementation: it intentionally
	// does not count Chiitoitsu or Kokushi Musou tenpai here.
	return ok && service.IsTenpaiGeneral(hand)
}

func run(patterns []string) (*configs.GameStats, error) {
	paths, err := archive.GlobAll(patterns)
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("no input files matched")
	}

	basic := &basicCounter{}
	winPoints := newWinPointsCounter()
	yamiten := newYamitenCounter()
	drawTenpai := newDrawTenpaiCounter()
	counters := []counter{basic, winPoints, yamiten, drawTenpai}

	bar := progressbar.Default(int64(len(paths)))
	a := archive.NewArchive()
	err = a.PlayPaths(paths, archive.Handlers{
		OnMessage: func(msg inbound.Message) error {
			if _, ok := msg.(*inbound.Error); ok {
				return fmt.Errorf("error in the log")
			}
			return nil
		},
		OnEvent: func(ev event.Event, a *archive.Archive) error {
			state, ok := a.StateViewer()
			if !ok {
				return nil
			}
			for _, c := range counters {
				if err := c.onEvent(ev, state); err != nil {
					return err
				}
			}
			return nil
		},
		OnFileDone: func(string) error {
			return bar.Add(1)
		},
	})
	if err != nil {
		return nil, err
	}

	return buildOutput(basic, winPoints, yamiten, drawTenpai), nil
}

func buildOutput(
	basic *basicCounter,
	winPoints *winPointsCounter,
	yamiten *yamitenCounter,
	drawTenpai *drawTenpaiCounter,
) *configs.GameStats {
	turnDistribution := make([]float64, maxTurn)
	if basic.numRounds > 0 {
		for i, freq := range basic.numTurnFreqs {
			turnDistribution[i] = float64(freq) / float64(basic.numRounds)
		}
	}

	var drawRoundRatio float64
	if basic.numRounds > 0 {
		drawRoundRatio = float64(basic.numDrawRounds) / float64(basic.numRounds)
	}

	var averageWinPoints float64
	if basic.numWins > 0 {
		averageWinPoints = float64(basic.totalWinPoints) / float64(basic.numWins)
	}

	return &configs.GameStats{
		NumHoras:             basic.numWins,
		NumTsumoHoras:        basic.numSelfDrawWins,
		NumTurnsDistribution: turnDistribution,
		RyukyokuRatio:        drawRoundRatio,
		AverageHoraPoints:    averageWinPoints,
		KoHoraPointsFreqs:    winPoints.koFreqs,
		OyaHoraPointsFreqs:   winPoints.oyaFreqs,
		YamitenStats:         yamiten.stats,
		RyukyokuTenpaiStat:   drawTenpai.stats,
	}
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
