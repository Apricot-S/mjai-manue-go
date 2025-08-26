package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/game"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/protocol/mjai"
	"github.com/schollz/progressbar/v3"
)

type CandidateInfo struct {
	Pai           base.Pai
	Hit           bool
	FeatureVector *BitVector
}

type Listener interface {
	OnDahai(
		logger io.Writer,
		state game.StateViewer,
		action inbound.Event,
		reacher *base.Player,
		candidates []CandidateInfo,
		path string,
		rawAction []byte,
	)
}

const batchSize = 100

var excludedPlayers = []string{"ASAPIN", "（≧▽≦）"}

func extractFeaturesSingle(inputPath string, listener Listener, verbose bool, logger io.Writer) ([]StoredKyoku, error) {
	r, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %w", err)
	}
	defer r.Close()

	lines, err := readAllLines(r)
	if err != nil {
		return nil, err
	}

	var actions []inbound.Event
	for _, line := range lines {
		as, err := mjai.Adapter.DecodeMessages(line)
		if err != nil {
			return nil, fmt.Errorf("failed to decode: %w", err)
		}
		actions = slices.Concat(actions, as)
	}

	storedKyokus, err := processActions(actions, listener, verbose, logger, inputPath, lines)
	if err != nil {
		return nil, err
	}
	return storedKyokus, nil
}

func readAllLines(r io.Reader) ([][]byte, error) {
	var lines [][]byte
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) > 0 {
			lines = append(lines, bytes.Clone(line))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return lines, nil
}

func processActions(
	actions []inbound.Event,
	listener Listener,
	verbose bool,
	logger io.Writer,
	path string,
	rawActions [][]byte,
) ([]StoredKyoku, error) {
	state := &game.StateImpl{}
	var kyokus []StoredKyoku
	var current *StoredKyoku
	var reacher *base.Player
	var waited *base.PaiSet
	skip := false

	for i, action := range actions {
		if err := state.Update(action); err != nil {
			return nil, fmt.Errorf("failed to update state: %w", err)
		}

		if verbose {
			logger.Write(rawActions[i])
			fmt.Fprintln(logger)
			fmt.Fprint(logger, state.RenderBoard())
		}

		switch a := action.(type) {
		case *inbound.StartKyoku:
			current, reacher, waited, skip = onStartKyoku()
		case *inbound.EndKyoku:
			var err error
			if kyokus, current, err = onEndKyoku(kyokus, current, skip); err != nil {
				return nil, err
			}
		case *inbound.ReachAccepted:
			var err error
			if reacher, waited, skip, err = onReachAccepted(a, state, skip, reacher); err != nil {
				return nil, err
			}
		case *inbound.Dahai:
			scene, candidates, err := onDahai(a, state, skip, reacher, waited, verbose, logger)
			if err != nil {
				return nil, err
			}
			if scene != nil {
				current.Scenes = append(current.Scenes, *scene)
			}
			if listener != nil && candidates != nil {
				listener.OnDahai(logger, state, action, reacher, candidates, path, rawActions[i])
			}
		}
	}
	if current != nil {
		return nil, fmt.Errorf(`game log ended without "end_kyoku"`)
	}
	return kyokus, nil
}

func onStartKyoku() (current *StoredKyoku, reacher *base.Player, waited *base.PaiSet, skip bool) {
	return &StoredKyoku{Scenes: nil}, nil, nil, false
}

func onEndKyoku(kyokus []StoredKyoku, current *StoredKyoku, skip bool) ([]StoredKyoku, *StoredKyoku, error) {
	if current == nil {
		return nil, nil, fmt.Errorf(`"end_kyoku" exists before "start_kyoku"`)
	}
	if skip {
		return kyokus, nil, nil
	}
	kyokus = append(kyokus, *current)
	return kyokus, nil, nil
}

func onReachAccepted(
	action *inbound.ReachAccepted,
	state game.State,
	skip bool,
	reacher *base.Player,
) (*base.Player, *base.PaiSet, bool, error) {
	if slices.Contains(excludedPlayers, state.Players()[action.Actor].Name()) {
		skip = true
	}
	if reacher != nil {
		// Skip if the second player has declared Riichi
		skip = true
	}
	if skip {
		return reacher, nil, true, nil
	}

	reacher = &state.Players()[action.Actor]
	tehaiSet, err := base.NewPaiSet(reacher.Tehais())
	if err != nil {
		return nil, nil, false, err
	}
	waited, err := game.GetWaitedPaisAll(tehaiSet)
	if err != nil {
		return nil, nil, false, err
	}
	return reacher, waited, false, nil
}

func onDahai(
	action *inbound.Dahai,
	state game.State,
	skip bool,
	reacher *base.Player,
	waited *base.PaiSet,
	verbose bool,
	logger io.Writer,
) (*StoredScene, []CandidateInfo, error) {
	me := &state.Players()[action.Actor]
	if skip || reacher == nil || me.ReachState() == base.ReachAccepted {
		// Skip if:
		// - No player has declared Riichi
		// - Dahai by the Riichi declarer itself
		return nil, nil, nil
	}

	scene, err := NewScene(state, me, &action.Pai, reacher)
	if err != nil {
		return nil, nil, err
	}

	if verbose {
		fmt.Fprintf(logger, "reacher: %d\n", reacher.ID())
	}

	storedScene := &StoredScene{Candidates: nil}
	var candidates []CandidateInfo
	for _, pai := range scene.Candidates() {
		hit, err := waited.Has(&pai)
		if err != nil {
			return nil, nil, err
		}
		feature, err := scene.FeatureVector(&pai)
		if err != nil {
			return nil, nil, err
		}
		storedScene.Candidates = append(storedScene.Candidates, Candidate{
			FeatureVector: feature,
			Hit:           hit,
		})
		candidates = append(candidates, CandidateInfo{
			Pai:           pai,
			Hit:           hit,
			FeatureVector: feature,
		})

		if verbose {
			h := 0
			if hit {
				h = 1
			}
			fmt.Fprintf(logger, "candidate %s: hit=%d, %s\n", pai.ToString(), h, FeatureVectorToStr(feature))
		}
	}
	return storedScene, candidates, nil
}

func extractFeaturesBatch(
	inputPaths []string,
	writer io.Writer,
	featureExtractor func(string) ([]StoredKyoku, error),
) error {
	numInputs := len(inputPaths)
	fmt.Fprintf(os.Stderr, "%d files.\n", numInputs)
	bar := progressbar.Default(int64(numInputs))

	encoder := gob.NewEncoder(writer)

	metaData := MetaData{FeatureNames: FeatureNames()}
	if err := encoder.Encode(metaData); err != nil {
		return fmt.Errorf("failed to encode metadata: %w", err)
	}

	var storedKyokus []StoredKyoku

	for i, path := range inputPaths {
		sks, err := featureExtractor(path)
		if err != nil {
			return fmt.Errorf("error at %s: %w", path, err)
		}
		storedKyokus = slices.Concat(storedKyokus, sks)

		if err := bar.Add(1); err != nil {
			return err
		}

		if i%batchSize == batchSize-1 {
			// Dump every batchSize games
			if err := encoder.Encode(storedKyokus); err != nil {
				return fmt.Errorf("failed to encode storedKyokus: %w", err)
			}
			storedKyokus = nil
		}
	}

	if len(storedKyokus) > 0 {
		// Dump the rest
		if err := encoder.Encode(storedKyokus); err != nil {
			return fmt.Errorf("failed to encode storedKyokus: %w", err)
		}
	}

	return nil
}

func ExtractFeaturesFromFiles(
	inputPaths []string,
	outputPath string,
	listener Listener,
	verbose bool,
	logger io.Writer,
) error {
	featureExtractor := func(inputPath string) ([]StoredKyoku, error) {
		return extractFeaturesSingle(inputPath, listener, verbose, logger)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer f.Close()

	return extractFeaturesBatch(inputPaths, f, featureExtractor)
}
