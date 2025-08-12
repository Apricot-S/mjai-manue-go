package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/schollz/progressbar/v3"
)

func extractFeaturesSingle(input io.Reader, listener any) ([]StoredKyoku, error) {
	panic("not implemented")
}

func extractFeaturesBatch(
	inputPaths []string,
	output io.Writer,
	featureExtractor func(io.Reader) ([]StoredKyoku, error),
) error {
	numInputs := len(inputPaths)
	fmt.Fprintf(os.Stderr, "%d files.\n", numInputs)
	bar := progressbar.Default(int64(numInputs))

	encoder := gob.NewEncoder(output)

	metaData := MetaData{FeatureNames: FeatureNames()}
	if err := encoder.Encode(metaData); err != nil {
		return fmt.Errorf("failed to encode metadata: %w", err)
	}

	storedKyokus := make([]StoredKyoku, 0)

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
			storedKyokus = make([]StoredKyoku, 0)
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

func ExtractFeaturesFromFiles(inputPaths []string, outputPath string, listener any) error {
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
