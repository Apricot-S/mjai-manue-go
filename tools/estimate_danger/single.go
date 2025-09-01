package main

import (
	"io"
)

func buildSingleCriteria(featureNames []string) []Criterion {
	criteria := make([]Criterion, 0, len(featureNames)*2)
	for _, s := range featureNames {
		criteria = append(criteria, Criterion{s: false}, Criterion{s: true})
	}
	return criteria
}

func CalculateSingleProbabilities(featuresPath string, w io.Writer) error {
	fn := FeatureNames()
	storedKyokus, err := LoadStoredKyokus(featuresPath, fn)
	if err != nil {
		return err
	}

	criteria := buildSingleCriteria(fn)
	if _, err := CalculateProbabilities(w, storedKyokus, fn, criteria); err != nil {
		return err
	}
	return nil
}
