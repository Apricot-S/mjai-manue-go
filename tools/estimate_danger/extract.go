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
	onDahai(state game.StateViewer, action inbound.Event, reacher *base.Player, candidates []CandidateInfo)
}

var adapter = mjai.MjaiAdapter{}

var excludedPlayers = []string{"ASAPIN", "（≧▽≦）"}

func extractFeaturesSingle(reader io.Reader, listener Listener) ([]StoredKyoku, error) {
	state := game.StateImpl{}
	var storedKyokus []StoredKyoku
	var reacher *base.Player = nil
	var waited *base.PaiSet = nil
	skip := false

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}

		actions, err := adapter.DecodeMessages(line)
		if err != nil {
			return nil, fmt.Errorf("failed to decode: %w", err)
		}

		storedKyoku := &StoredKyoku{Scenes: nil}
		for _, action := range actions {
			if err := state.Update(action); err != nil {
				return nil, fmt.Errorf("failed to update state: %w", err)
			}

			switch a := action.(type) {
			case *inbound.StartKyoku:
				storedKyoku.Scenes = nil
				reacher = nil
				skip = false
			case *inbound.EndKyoku:
				if skip {
					continue
				}
				if storedKyoku == nil {
					return nil, fmt.Errorf("should not happen")
				}
				storedKyokus = append(storedKyokus, *storedKyoku)
				storedKyoku = nil
			case *inbound.ReachAccepted:
				if slices.Contains(excludedPlayers, state.Players()[a.Actor].Name()) {
					skip = true
				}

				if reacher != nil {
					skip = true
				}
				if skip {
					continue
				}

				reacher = &state.Players()[a.Actor]
				tehaiSet, err := base.NewPaiSet(reacher.Tehais())
				if err != nil {
					return nil, err
				}
				waited, err = game.GetWaitedPaisAll(tehaiSet)
				if err != nil {
					return nil, err
				}
			case *inbound.Dahai:
				me := &state.Players()[a.Actor]
				if skip || reacher == nil || me.ReachState() == base.ReachAccepted {
					continue
				}

				scene, err := NewSceneWithState(&state, me, reacher)
				if err != nil {
					return nil, err
				}

				storedScene := StoredScene{Candidates: nil}
				var candidates []CandidateInfo
				for _, pai := range scene.Candidates() {
					hit, err := waited.Has(&pai)
					if err != nil {
						return nil, err
					}
					featureVector, err := scene.FeatureVector(&pai)
					if err != nil {
						return nil, err
					}

					candidate := Candidate{
						FeatureVector: featureVector,
						Hit:           hit,
					}
					storedScene.Candidates = append(storedScene.Candidates, candidate)

					candidateInfo := CandidateInfo{
						Pai:           pai,
						Hit:           hit,
						FeatureVector: featureVector,
					}
					candidates = append(candidates, candidateInfo)
				}

				storedKyoku.Scenes = append(storedKyoku.Scenes, storedScene)

				if listener != nil {
					listener.onDahai(&state, action, reacher, candidates)
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}
	return storedKyokus, nil
}

func extractFeaturesBatch(
	inputPaths []string,
	writer io.Writer,
	featureExtractor func(io.Reader) ([]StoredKyoku, error),
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
		r, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open input file: %w", err)
		}

		storedKyoku, err := featureExtractor(r)
		r.Close()
		if err != nil {
			return err
		}

		storedKyokus = slices.Concat(storedKyokus, storedKyoku)

		if err := bar.Add(1); err != nil {
			return err
		}

		if i%100 == 99 {
			// Dump every 100 games
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

func ExtractFeaturesFromFiles(inputPaths []string, outputPath string, listener Listener) error {
	featureExtractor := func(input io.Reader) ([]StoredKyoku, error) {
		return extractFeaturesSingle(input, listener)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer f.Close()

	return extractFeaturesBatch(inputPaths, f, featureExtractor)
}
