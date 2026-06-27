package main

import "io"

func buildSingleCriteria(featureNames []string) []Criterion {
	criteria := make([]Criterion, 0, len(featureNames)*2)
	for _, name := range featureNames {
		criteria = append(criteria, Criterion{name: false}, Criterion{name: true})
	}
	return criteria
}

func CalculateSingleProbabilities(featuresPath string, w io.Writer) error {
	featureNames := FeatureNames()
	storedKyokus, err := LoadStoredKyokus(featuresPath, featureNames)
	if err != nil {
		return err
	}

	criteria := buildSingleCriteria(featureNames)
	_, err = CalculateProbabilities(w, storedKyokus, featureNames, criteria)
	return err
}
