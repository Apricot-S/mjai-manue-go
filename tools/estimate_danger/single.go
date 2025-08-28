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

func calculateProbabilities(r io.Reader, w io.Writer, criteria []map[string]bool) error {
	return nil
}

func CalculateSingleProbabilities(featuresPath string, w io.Writer) error {
	r, err := os.Open(featuresPath)
	if err != nil {
		return fmt.Errorf("failed to open features file: %w", err)
	}
	defer r.Close()

	return calculateProbabilities(r, w, getCriteria())
}
