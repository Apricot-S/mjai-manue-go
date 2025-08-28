package main

import (
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/go-json-experiment/json"
)

func getCriteria(featureNames []string) []map[string]bool {
	criteria := make([]map[string]bool, len(featureNames)*2)
	for _, s := range featureNames {
		criteria = append(criteria, map[string]bool{s: false}, map[string]bool{s: true})
	}
	return criteria
}

func createCriterionMasks(featureNames []string, criteria []map[string]bool) (map[string][2]*BitVector, error) {
	positiveAry := make([]bool, len(featureNames))
	negativeAry := make([]bool, len(featureNames))
	for i := range len(negativeAry) {
		negativeAry[i] = true
	}

	criterionMasks := make(map[string][2]*BitVector, len(criteria))
	for _, criterion := range criteria {
		pa := slices.Clone(positiveAry)
		na := slices.Clone(negativeAry)

		for name, value := range criterion {
			index := slices.Index(featureNames, name)
			if index == -1 {
				return nil, fmt.Errorf("no such feature: %s", name)
			}

			if value {
				pa[index] = true
			} else {
				na[index] = false
			}
		}

		keyBytes, err := json.Marshal(criterion, json.Deterministic(true))
		if err != nil {
			return nil, fmt.Errorf("failed to decode criterion: %w", err)
		}
		key := string(keyBytes)
		criterionMasks[key] = [2]*BitVector{BoolArrayToBitVector(pa), BoolArrayToBitVector(na)}
	}

	return criterionMasks, nil
}

func processFeaturesFile(r io.Reader, w io.Writer, featureNames []string, criterionMasks map[string][2]*BitVector) error {
	return nil
}

func createKyokuProbsMap(r io.Reader, w io.Writer, featureNames []string, criteria []map[string]bool) (any, error) {
	criterionMasks, err := createCriterionMasks(featureNames, criteria)
	if err != nil {
		return nil, err
	}
	return nil, processFeaturesFile(r, w, featureNames, criterionMasks)
}

func aggregateProbabilities(w io.Writer, kyokuProbsMap any, criteria []map[string]bool) (any, error) {
	return nil, nil
}

func calculateProbabilities(r io.Reader, w io.Writer, featureNames []string, criteria []map[string]bool) (any, error) {
	kyokuProbsMap, err := createKyokuProbsMap(r, w, featureNames, criteria)
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

	fn := FeatureNames()
	criteria := getCriteria(fn)
	if _, err := calculateProbabilities(r, w, fn, criteria); err != nil {
		return err
	}
	return nil
}
