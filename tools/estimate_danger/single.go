package main

import (
	"fmt"
	"io"
	"os"
)

func buildSingleCriteria(featureNames []string) []Criterion {
	criteria := make([]Criterion, 0, len(featureNames)*2)
	for _, s := range featureNames {
		criteria = append(criteria, Criterion{s: false}, Criterion{s: true})
	}
	return criteria
}

func CalculateSingleProbabilities(featuresPath string, w io.Writer) error {
	r, err := os.Open(featuresPath)
	if err != nil {
		return fmt.Errorf("failed to open features file: %w", err)
	}
	defer r.Close()

	stat, err := r.Stat()
	if err != nil {
		return err
	}

	fn := FeatureNames()
	criteria := buildSingleCriteria(fn)
	if _, err := CalculateProbabilities(r, w, stat.Size(), fn, criteria); err != nil {
		return err
	}
	return nil
}
