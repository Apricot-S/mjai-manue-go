package main

import (
	"cmp"
	"encoding/gob"
	"fmt"
	"io"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/Apricot-S/mjai-manue-go/configs"
	"github.com/go-json-experiment/json"
)

func generateDecisionTreeImpl(
	w io.Writer,
	storedKyokus []StoredKyoku,
	featureNames []string,
	baseCriterion Criterion,
	baseNode, root *configs.DecisionNode,
	minGap float64,
) (*configs.DecisionNode, error) {
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
		negativeCriterion[name] = false
		positiveCriterion := maps.Clone(baseCriterion)
		positiveCriterion[name] = true

		targets[name] = [2]*Criterion{&negativeCriterion, &positiveCriterion}
		criteria = append(criteria, negativeCriterion, positiveCriterion)
	}

	nodeMap, err := CalculateProbabilities(w, storedKyokus, featureNames, criteria)
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
		c := targets[maxName]
		baseNode.FeatureName = &maxName

		negativeKey, err := json.Marshal(c[0], json.Deterministic(true))
		if err != nil {
			return nil, fmt.Errorf("failed to encode criterion: %w", err)
		}
		baseNode.Negative = nodeMap[string(negativeKey)]

		positiveKey, err := json.Marshal(c[1], json.Deterministic(true))
		if err != nil {
			return nil, fmt.Errorf("failed to encode criterion: %w", err)
		}
		baseNode.Positive = nodeMap[string(positiveKey)]

		RenderDecisionTree(w, root, "all", 0)

		_, err = generateDecisionTreeImpl(w, storedKyokus, featureNames, *c[0], baseNode.Negative, root, minGap)
		if err != nil {
			return nil, err
		}
		_, err = generateDecisionTreeImpl(w, storedKyokus, featureNames, *c[1], baseNode.Positive, root, minGap)
		if err != nil {
			return nil, err
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

	fn := FeatureNames()
	storedKyokus, err := LoadStoredKyokus(r, stat.Size(), fn)
	if err != nil {
		return nil, err
	}

	return generateDecisionTreeImpl(w, storedKyokus, FeatureNames(), make(Criterion), nil, nil, minGap)
}

func RenderDecisionTree(w io.Writer, node *configs.DecisionNode, label string, indent int) {
	fmt.Fprintf(
		w,
		"%s%s : %.2f [%.2f, %.2f] (%d samples)\n",
		strings.Repeat("  ", indent),
		label,
		node.AverageProb*100.0,
		node.ConfInterval[0]*100.0,
		node.ConfInterval[1]*100.0,
		node.NumSamples,
	)

	if node.FeatureName != nil {
		type childNode struct {
			Value bool
			Node  *configs.DecisionNode
		}
		children := []childNode{{false, node.Negative}, {true, node.Positive}}
		slices.SortFunc(children, func(a, b childNode) int {
			return cmp.Compare(a.Node.AverageProb, b.Node.AverageProb)
		})

		for _, child := range children {
			RenderDecisionTree(
				w,
				child.Node,
				fmt.Sprintf("%s = %v", *node.FeatureName, child.Value),
				indent+1,
			)
		}
	}
}

func DumpDecisionTree(node *configs.DecisionNode, outputPath string) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer f.Close()

	encoder := gob.NewEncoder(f)
	return encoder.Encode(node)
}

func LoadDecisionTree(treePath string) (*configs.DecisionNode, error) {
	f, err := os.Open(treePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open tree file: %w", err)
	}
	defer f.Close()

	decoder := gob.NewDecoder(f)
	var node configs.DecisionNode
	if err := decoder.Decode(&node); err != nil {
		return nil, fmt.Errorf("failed to decode tree file: %w", err)
	}
	return &node, nil
}
