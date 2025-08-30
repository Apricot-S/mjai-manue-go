package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/configs"
	"github.com/go-json-experiment/json"
	"github.com/schollz/progressbar/v3"
)

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

		key, err := json.Marshal(criterion, json.Deterministic(true))
		if err != nil {
			return nil, fmt.Errorf("failed to encode criterion: %w", err)
		}
		criterionMasks[string(key)] = [2]*BitVector{BoolArrayToBitVector(pa), BoolArrayToBitVector(na)}
	}

	return criterionMasks, nil
}

func loadStoredKyokus(r io.Reader, fileSize int64, featureNames []string) ([]StoredKyoku, error) {
	bar := progressbar.DefaultBytes(fileSize, "loading")
	tr := io.TeeReader(r, bar)

	decoder := gob.NewDecoder(tr)

	var metaData MetaData
	if err := decoder.Decode(&metaData); err != nil {
		return nil, fmt.Errorf("failed to load features %w", err)
	}
	if slices.Compare(metaData.FeatureNames, featureNames) != 0 {
		return nil, fmt.Errorf("feature set has been changed")
	}

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
					if _, found := paiFreqs[criterion]; !found {
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

func createKyokuProbsMap(
	r io.Reader,
	fileSize int64,
	featureNames []string,
	criteria []Criterion,
) (map[string][]float64, error) {
	criterionMasks, err := createCriterionMasks(featureNames, criteria)
	if err != nil {
		return nil, err
	}

	storedKyokus, err := loadStoredKyokus(r, fileSize, featureNames)
	if err != nil {
		return nil, err
	}

	kyokuProbsMap := make(map[string][]float64)
	for _, sk := range storedKyokus {
		kpm := createMetricsForKyoku(sk, criterionMasks)
		for key, value := range kpm {
			kyokuProbsMap[key] = slices.Concat(kyokuProbsMap[key], value)
		}
	}

	return kyokuProbsMap, nil
}

func aggregateProbabilities(
	kyokuProbsMap map[string][]float64,
	criteria []Criterion,
) (map[string]*configs.DecisionNode, error) {
	results := make(map[string]*configs.DecisionNode)
	for _, criterion := range criteria {
		key, err := json.Marshal(criterion, json.Deterministic(true))
		if err != nil {
			return nil, fmt.Errorf("failed to encode criterion: %w", err)
		}

		kyokuProbs := kyokuProbsMap[string(key)]
		if len(kyokuProbs) == 0 {
			continue
		}

		n := 0.0
		for _, p := range kyokuProbs {
			n += p
		}
		numSamples := len(kyokuProbs)
		lower, upper := CalculateConfidenceInterval(kyokuProbs, 0.0, 1.0, 0.95)

		node := &configs.DecisionNode{
			AverageProb:  n / float64(numSamples),
			ConfInterval: [2]float64{lower, upper},
			NumSamples:   numSamples,
		}
		results[string(key)] = node
	}
	return results, nil
}

func printAggregateResults(w io.Writer, criteria []Criterion, result map[string]*configs.DecisionNode) error {
	for _, criterion := range criteria {
		key, err := json.Marshal(criterion, json.Deterministic(true))
		if err != nil {
			return fmt.Errorf("failed to encode criterion: %w", err)
		}

		node, found := result[string(key)]
		if !found {
			continue
		}

		fmt.Fprintf(
			w,
			"%v\n  %.2f [%.2f, %.2f] (%d samples)\n\n",
			criterion,
			node.AverageProb*100.0,
			node.ConfInterval[0]*100.0,
			node.ConfInterval[1]*100.0,
			node.NumSamples,
		)
	}
	return nil
}

func calculateProbabilities(
	r io.Reader,
	w io.Writer,
	fileSize int64,
	featureNames []string,
	criteria []Criterion,
) (map[string]*configs.DecisionNode, error) {
	kyokuProbsMap, err := createKyokuProbsMap(r, fileSize, featureNames, criteria)
	if err != nil {
		return nil, err
	}

	results, err := aggregateProbabilities(kyokuProbsMap, criteria)
	if err != nil {
		return nil, err
	}

	if err := printAggregateResults(w, criteria, results); err != nil {
		return nil, err
	}

	return results, nil
}
