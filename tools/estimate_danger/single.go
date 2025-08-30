package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"maps"
	"os"
	"slices"

	"github.com/go-json-experiment/json"
)

func getCriteria(featureNames []string) []Criterion {
	criteria := make([]Criterion, len(featureNames)*2)
	for _, s := range featureNames {
		criteria = append(criteria, Criterion{s: false}, Criterion{s: true})
	}
	return criteria
}

func createCriterionMasks(featureNames []string, criteria []Criterion) (CriterionMasks, error) {
	positiveAry := make([]bool, len(featureNames))
	negativeAry := make([]bool, len(featureNames))
	for i := range len(negativeAry) {
		negativeAry[i] = true
	}

	criterionMasks := make(CriterionMasks, len(criteria))
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
			return nil, fmt.Errorf("failed to encode criterion: %w", err)
		}
		key := string(keyBytes)
		criterionMasks[key] = [2]*BitVector{BoolArrayToBitVector(pa), BoolArrayToBitVector(na)}
	}

	return criterionMasks, nil
}

func loadStoredKyokus(r io.Reader, featureNames []string) ([]StoredKyoku, error) {
	decoder := gob.NewDecoder(r)

	var metaData MetaData
	if err := decoder.Decode(&metaData); err != nil {
		return nil, fmt.Errorf("failed to load features %w", err)
	}
	if slices.Compare(metaData.FeatureNames, featureNames) != 0 {
		return nil, fmt.Errorf("feature set has been changed")
	}

	// bar := progressbar.DefaultBytes(-1)

	var storedKyokus []StoredKyoku
	for {
		var sks []StoredKyoku
		err := decoder.Decode(&sks)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		storedKyokus = slices.Concat(storedKyokus, sks)
	}

	return storedKyokus, nil
}

func createMetricsForKyoku(storedKyoku StoredKyoku, criterionMasks CriterionMasks) map[string][]float64 {
	sceneProbSums := make(map[string]float64)
	sceneCounts := make(map[string]int)
	for _, storedScene := range storedKyoku.Scenes {
		paiFreqs := make(map[string]map[bool]int)
		for _, candidate := range storedScene.Candidates {
			for criterion, masks := range criterionMasks {
				if Matches(candidate.FeatureVector, masks[0], masks[1]) {
					if _, ok := paiFreqs[criterion]; !ok {
						paiFreqs[criterion] = make(map[bool]int)
					}
					paiFreqs[criterion][candidate.Hit] += 1
				}
			}
		}

		for criterion, freqs := range paiFreqs {
			sceneProb := float64(freqs[true]) / float64(freqs[false]+freqs[true])
			sceneProbSums[criterion] += sceneProb
			sceneCounts[criterion] += 1
		}
	}

	kyokuProbsMap := make(map[string][]float64)
	for criterion, count := range sceneCounts {
		kyokuProb := sceneProbSums[criterion] / float64(count)
		kyokuProbsMap[criterion] = append(kyokuProbsMap[criterion], kyokuProb)
	}
	return kyokuProbsMap
}

func createKyokuProbsMap(r io.Reader, featureNames []string, criteria []Criterion) (map[string][]float64, error) {
	criterionMasks, err := createCriterionMasks(featureNames, criteria)
	if err != nil {
		return nil, err
	}

	storedKyokus, err := loadStoredKyokus(r, featureNames)
	if err != nil {
		return nil, err
	}

	kyokuProbsMap := make(map[string][]float64)
	for _, sk := range storedKyokus {
		kpm := createMetricsForKyoku(sk, criterionMasks)
		maps.Copy(kyokuProbsMap, kpm)
	}

	return kyokuProbsMap, nil
}

func aggregateProbabilities(w io.Writer, kyokuProbsMap map[string][]float64, criteria []Criterion) (any, error) {
	return nil, nil
}

func calculateProbabilities(r io.Reader, w io.Writer, featureNames []string, criteria []Criterion) (any, error) {
	kyokuProbsMap, err := createKyokuProbsMap(r, featureNames, criteria)
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
