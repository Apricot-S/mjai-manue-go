package main

import (
	"encoding/gob"
	"encoding/json/v2"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/configs"
	"github.com/schollz/progressbar/v3"
)

func createCriterionMasks(featureNames []string, criteria []Criterion) (CriterionMasks, error) {
	positiveAry := make([]bool, len(featureNames))
	negativeAry := make([]bool, len(featureNames))
	for i := range negativeAry {
		negativeAry[i] = true
	}

	criterionMasks := make(CriterionMasks, len(criteria))
	for _, criterion := range criteria {
		pos := slices.Clone(positiveAry)
		neg := slices.Clone(negativeAry)

		for name, value := range criterion {
			index := slices.Index(featureNames, name)
			if index == -1 {
				return nil, fmt.Errorf("no such feature: %s", name)
			}

			if value {
				pos[index] = true
			} else {
				neg[index] = false
			}
		}

		key, err := encodeCriterion(criterion)
		if err != nil {
			return nil, err
		}
		criterionMasks[key] = [2]*BitVector{BoolArrayToBitVector(pos), BoolArrayToBitVector(neg)}
	}

	return criterionMasks, nil
}

// encodeCriterion returns a deterministic key for a criterion.
//
// Ruby's original implementation uses criterion.object_id as an internal hash
// key. Go needs a value key because maps cannot be keyed by another map, so use
// deterministic JSON for stable lookup while keeping the external behavior.
func encodeCriterion(criterion Criterion) (string, error) {
	key, err := json.Marshal(criterion, json.Deterministic(true))
	if err != nil {
		return "", fmt.Errorf("failed to encode criterion: %w", err)
	}
	return string(key), nil
}

func loadStoredKyokusImpl(r io.Reader, featureNames []string) ([]StoredKyoku, error) {
	decoder := gob.NewDecoder(r)

	var metaData MetaData
	if err := decoder.Decode(&metaData); err != nil {
		return nil, fmt.Errorf("failed to load features: %w", err)
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

		storedKyokus = append(storedKyokus, sks...)
	}

	return storedKyokus, nil
}

func LoadStoredKyokus(featuresPath string, featureNames []string) ([]StoredKyoku, error) {
	f, err := os.Open(featuresPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open features file: %w", err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat features file: %w", err)
	}

	bar := progressbar.DefaultBytes(stat.Size(), "loading   ")
	return loadStoredKyokusImpl(io.TeeReader(f, bar), featureNames)
}

// createKyokuProbabilities ports the original Ruby aggregation order:
//
//  1. calculate a hit probability for each scene and criterion,
//  2. average scene probabilities within one kyoku,
//  3. let aggregateProbabilities average those kyoku probabilities.
//
// This intentionally does not pool all discard candidates directly; doing so
// would give scenes with more matching candidates more weight than the original.
func createKyokuProbabilities(storedKyoku StoredKyoku, criterionMasks CriterionMasks) map[string]float64 {
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
					paiFreqs[criterion][candidate.Hit]++
				}
			}
		}

		for criterion, freqs := range paiFreqs {
			sceneProb := float64(freqs[true]) / float64(freqs[false]+freqs[true])
			sceneProbSums[criterion] += sceneProb
			sceneCounts[criterion]++
		}
	}

	kyokuProbs := make(map[string]float64)
	for criterion, count := range sceneCounts {
		kyokuProbs[criterion] = sceneProbSums[criterion] / float64(count)
	}
	return kyokuProbs
}

func CreateKyokuProbsMap(
	storedKyokus []StoredKyoku,
	featureNames []string,
	criteria []Criterion,
) (map[string][]float64, error) {
	criterionMasks, err := createCriterionMasks(featureNames, criteria)
	if err != nil {
		return nil, err
	}

	kyokuProbsMap := make(map[string][]float64)
	for _, sk := range storedKyokus {
		kyokuProbs := createKyokuProbabilities(sk, criterionMasks)
		for key, value := range kyokuProbs {
			kyokuProbsMap[key] = append(kyokuProbsMap[key], value)
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
		key, err := encodeCriterion(criterion)
		if err != nil {
			return nil, err
		}

		kyokuProbs := kyokuProbsMap[key]
		if len(kyokuProbs) == 0 {
			continue
		}

		sum := 0.0
		for _, p := range kyokuProbs {
			sum += p
		}
		numSamples := len(kyokuProbs)
		lower, upper := CalculateConfidenceInterval(kyokuProbs, 0.0, 1.0, 0.95)

		results[key] = &configs.DecisionNode{
			AverageProb:  sum / float64(numSamples),
			ConfInterval: [2]float64{lower, upper},
			NumSamples:   numSamples,
		}
	}
	return results, nil
}

func printAggregateResults(w io.Writer, criteria []Criterion, result map[string]*configs.DecisionNode) error {
	for _, criterion := range criteria {
		key, err := encodeCriterion(criterion)
		if err != nil {
			return err
		}

		node, found := result[key]
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

func CalculateProbabilities(
	w io.Writer,
	storedKyokus []StoredKyoku,
	featureNames []string,
	criteria []Criterion,
) (map[string]*configs.DecisionNode, error) {
	kyokuProbsMap, err := CreateKyokuProbsMap(storedKyokus, featureNames, criteria)
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

func DumpProbabilities(probs map[string]*configs.DecisionNode, outputPath string) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer f.Close()

	encoder := gob.NewEncoder(f)
	if err := encoder.Encode(probs); err != nil {
		return fmt.Errorf("failed to encode probabilities file: %w", err)
	}
	return nil
}
