package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/tools/internal/archive"
	"github.com/schollz/progressbar/v3"
)

type CandidateInfo struct {
	Tile          tile.Tile
	Hit           bool
	FeatureVector *BitVector
}

type Listener interface {
	OnDahai(logger io.Writer, reacher seat.Seat, candidates []CandidateInfo, path string, rawAction []byte)
}

const batchSize = 100

var excludedPlayers = []string{"ASAPIN", "（≧▽≦）"}

type extractor struct {
	listener Listener
	verbose  bool
	logger   io.Writer

	encoder      *gob.Encoder
	storedKyokus []StoredKyoku
	current      *StoredKyoku
	reacher      *seat.Seat
	waits        service.WaitSet
	skip         bool
	currentPath  string
	rawAction    []byte
	names        []string
}

func ExtractFeaturesFromFiles(paths []string, outputPath string, listener Listener, verbose bool, logger io.Writer) error {
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer out.Close()

	e := &extractor{
		listener: listener,
		verbose:  verbose,
		logger:   logger,
		encoder:  gob.NewEncoder(out),
	}

	if err := e.encoder.Encode(MetaData{FeatureNames: FeatureNames()}); err != nil {
		return fmt.Errorf("failed to write feature metadata: %w", err)
	}

	bar := progressbar.Default(int64(len(paths)))
	fileIndex := 0
	e.currentPath = paths[0]
	a := archive.NewArchive()
	if err := a.PlayPaths(paths, archive.Handlers{
		OnRaw: func(line []byte) error {
			e.rawAction = line
			return nil
		},
		OnMessage: func(msg inbound.Message) error {
			if start, ok := msg.(*inbound.StartGame); ok {
				e.names = slices.Clone(start.Names)
			}
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
			return e.onEvent(ev, state)
		},
		OnFileDone: func(string) error {
			fileIndex++
			if fileIndex%batchSize == 0 {
				if err := e.flush(); err != nil {
					return err
				}
			}
			if fileIndex < len(paths) {
				e.currentPath = paths[fileIndex]
			}
			return bar.Add(1)
		},
	}); err != nil {
		return err
	}

	if e.current != nil {
		return fmt.Errorf(`game log ended without "end_kyoku"`)
	}

	return e.flush()
}

func (e *extractor) flush() error {
	if len(e.storedKyokus) == 0 {
		return nil
	}
	if err := e.encoder.Encode(e.storedKyokus); err != nil {
		return fmt.Errorf("failed to write extracted features: %w", err)
	}
	e.storedKyokus = nil
	return nil
}

func (e *extractor) onEvent(ev event.Event, state round.StateViewer) error {
	if e.verbose && len(e.rawAction) > 0 {
		if _, err := e.logger.Write(e.rawAction); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(e.logger); err != nil {
			return err
		}
	}

	switch ev := ev.(type) {
	case *event.StartRound:
		e.current = &StoredKyoku{}
		e.reacher = nil
		e.waits = 0
		e.skip = false
	case *event.EndRound:
		if e.current == nil {
			return fmt.Errorf(`"end_kyoku" exists before "start_kyoku"`)
		}
		if !e.skip {
			e.storedKyokus = append(e.storedKyokus, *e.current)
		}
		e.current = nil
	case *event.RiichiAccepted:
		return e.onRiichiAccepted(ev, state)
	case *event.Discard:
		return e.onDiscard(ev, state)
	}
	return nil
}

func (e *extractor) onRiichiAccepted(ev *event.RiichiAccepted, state round.StateViewer) error {
	actor := ev.Actor()
	if actor.Index() < len(e.names) && slices.Contains(excludedPlayers, e.names[actor.Index()]) {
		// Logs from known non-standard players were excluded from the danger training data.
		e.skip = true
	}
	if e.reacher != nil {
		// The danger tree is trained only from scenes with exactly one riichi player.
		// Once a second riichi is accepted, skip the whole round.
		e.skip = true
	}
	if e.skip {
		return nil
	}

	p := state.Player(actor)
	hand, ok := p.Hand()
	if !ok {
		return fmt.Errorf("riichi actor hand is not visible")
	}
	e.reacher = &actor
	e.waits = service.WaitsFor(hand)
	return nil
}

func (e *extractor) onDiscard(ev *event.Discard, state round.StateViewer) error {
	if e.skip || e.reacher == nil || ev.Actor() == *e.reacher {
		return nil
	}
	if state.Player(ev.Actor()).RiichiState() == player.RiichiAccepted {
		return nil
	}

	scene := NewScene(state, ev.Actor(), *e.reacher, ev.Tile())
	if e.verbose {
		fmt.Fprintf(e.logger, "reacher: %d\n", e.reacher.Index())
	}

	storedScene := StoredScene{}
	candidates := make([]CandidateInfo, 0, len(scene.Candidates()))
	for _, candidate := range scene.Candidates() {
		hit := e.waits.Has(candidate)
		featureVector, err := scene.FeatureVector(candidate)
		if err != nil {
			return err
		}
		storedScene.Candidates = append(storedScene.Candidates, Candidate{
			FeatureVector: featureVector,
			Hit:           hit,
		})
		candidates = append(candidates, CandidateInfo{
			Tile:          candidate,
			Hit:           hit,
			FeatureVector: featureVector,
		})
		if e.verbose {
			h := 0
			if hit {
				h = 1
			}
			fmt.Fprintf(e.logger, "candidate %s: hit=%d, %s\n", candidate, h, FeatureVectorToStr(featureVector))
		}
	}
	e.current.Scenes = append(e.current.Scenes, storedScene)
	if e.listener != nil {
		e.listener.OnDahai(e.logger, *e.reacher, candidates, e.currentPath, e.rawAction)
	}
	return nil
}
