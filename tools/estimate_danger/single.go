package main

import (
	"fmt"
	"io"
	"os"
)

func getCriteria() []map[string]bool {
	fn := FeatureNames()
	criteria := make([]map[string]bool, len(fn))
	for _, s := range fn {
		criteria = append(criteria, map[string]bool{s: false})
		criteria = append(criteria, map[string]bool{s: true})
	}
	return criteria
}

func createKyokuProbsMap(r io.Reader, w io.Writer, criteria []map[string]bool) (any, error) {
	return nil, nil
}

func aggregateProbabilities(w io.Writer, kyokuProbsMap any, criteria []map[string]bool) (any, error) {
	return nil, nil
}

func calculateProbabilities(r io.Reader, w io.Writer, criteria []map[string]bool) (any, error) {
	kyokuProbsMap, err := createKyokuProbsMap(r, w, criteria)
	if err != nil {
		return nil, err
	}

	return aggregateProbabilities(w, kyokuProbsMap, criteria)
}

func CalculateSingleProbabilities(featuresPath string, w io.Writer) error {
	r, err := os.Open(featuresPath)
	if err != nil {
		return fmt.Errorf("failed to open features file: %w", err)
	}
	defer r.Close()

	if _, err := calculateProbabilities(r, w, getCriteria()); err != nil {
		return err
	}
	return nil
}
