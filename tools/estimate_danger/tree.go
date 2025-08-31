package main

import (
	"fmt"
	"io"
	"maps"
	"os"

	"github.com/Apricot-S/mjai-manue-go/configs"
	"github.com/go-json-experiment/json"
)

func generateDecisionTreeImpl(
	r io.ReadSeeker,
	w io.Writer,
	fileSize int64,
	featureNames []string,
	baseCriterion Criterion,
	baseNode, root *configs.DecisionNode,
	minGap float64,
) (*configs.DecisionNode, error) {
	r.Seek(0, io.SeekStart)

	targets := make(map[string][2]*Criterion)

	var criteria []Criterion
	if baseNode == nil {
		criteria = append(criteria, baseCriterion)
	}

	for _, name := range featureNames {
		if _, ok := baseCriterion[name]; ok {
			continue
		}

		negativeCriterion := maps.Clone(baseCriterion)
		maps.Copy(negativeCriterion, Criterion{name: false})
		positiveCriterion := maps.Clone(baseCriterion)
		maps.Copy(positiveCriterion, Criterion{name: true})

		targets[name] = [2]*Criterion{&negativeCriterion, &positiveCriterion}
		criteria = append(criteria, negativeCriterion, positiveCriterion)
	}

	nodeMap, err := CalculateProbabilities(r, w, fileSize, featureNames, criteria)
	if err != nil {
		return nil, err
	}
	if baseNode == nil {
		key, err := json.Marshal(baseCriterion, json.Deterministic(true))
		if err != nil {
			return nil, fmt.Errorf("failed to encode criterion: %w", err)
		}
		baseNode = nodeMap[string(key)]
	}
	if root == nil {
		root = baseNode
	}

	gaps := make(map[string]float64)
	for name, c := range targets {
		negativeKey, err := json.Marshal(c[0], json.Deterministic(true))
		if err != nil {
			return nil, fmt.Errorf("failed to encode criterion: %w", err)
		}
		negative := nodeMap[string(negativeKey)]
		if negative == nil {
			continue
		}

		positiveKey, err := json.Marshal(c[1], json.Deterministic(true))
		if err != nil {
			return nil, fmt.Errorf("failed to encode criterion: %w", err)
		}
		positive := nodeMap[string(positiveKey)]
		if positive == nil {
			continue
		}

		var gap float64
		if positive.AverageProb >= negative.AverageProb {
			gap = positive.ConfInterval[0] - negative.ConfInterval[1]
		} else {
			gap = negative.ConfInterval[0] - positive.ConfInterval[1]
		}
		if gap > minGap {
			gaps[name] = gap
		}

		var maxName string
		var maxValue float64
		first := true
		for name, value := range gaps {
			if first || value > maxValue {
				maxName = name
				maxValue = value
				first = false
			}
		}

		if maxName != "" {

		}
	}

	return baseNode, nil
}

func GenerateDecisionTree(featuresPath string, w io.Writer, minGap float64) (*configs.DecisionNode, error) {
	r, err := os.Open(featuresPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open features file: %w", err)
	}
	defer r.Close()

	stat, err := r.Stat()
	if err != nil {
		return nil, err
	}

	return generateDecisionTreeImpl(r, w, stat.Size(), FeatureNames(), nil, nil, nil, minGap)
}
