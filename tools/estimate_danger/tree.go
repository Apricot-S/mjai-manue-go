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
)

type splitTargets map[string][2]*Criterion

func buildSplitTargets(featureNames []string, baseCriterion Criterion, includeBase bool) (splitTargets, []Criterion) {
	targets := make(splitTargets)
	var criteria []Criterion
	if includeBase {
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

	return targets, criteria
}

func selectBestSplit(
	w io.Writer,
	nodeMap map[string]*configs.DecisionNode,
	targets splitTargets,
	minGap float64,
) (string, error) {
	gaps := make(map[string]float64)
	for name, c := range targets {
		negativeCriterion := c[0]
		positiveCriterion := c[1]

		negativeKey, err := encodeCriterion(*negativeCriterion)
		if err != nil {
			return "", err
		}
		negative := nodeMap[negativeKey]
		if negative == nil {
			continue
		}

		positiveKey, err := encodeCriterion(*positiveCriterion)
		if err != nil {
			return "", err
		}
		positive := nodeMap[positiveKey]
		if positive == nil {
			continue
		}

		var gap float64
		if positive.AverageProb >= negative.AverageProb {
			gap = positive.ConfInterval[0] - negative.ConfInterval[1]
		} else {
			gap = negative.ConfInterval[0] - positive.ConfInterval[1]
		}
		fmt.Fprintf(w, "%#v, %#v\n", name, gap)

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
	fmt.Fprintf(w, ":max_name, %#v\n", maxName)
	return maxName, nil
}

func generateDecisionTreeImpl(
	w io.Writer,
	storedKyokus []StoredKyoku,
	featureNames []string,
	baseCriterion Criterion,
	baseNode, root *configs.DecisionNode,
	minGap float64,
) (*configs.DecisionNode, error) {
	fmt.Fprintf(w, ":generate_decision_tree, %#v\n", baseCriterion)

	// Build the base criterion and every one-feature split to evaluate at this node.
	targets, criteria := buildSplitTargets(featureNames, baseCriterion, baseNode == nil)
	nodeMap, err := CalculateProbabilities(w, storedKyokus, featureNames, criteria)
	if err != nil {
		return nil, err
	}

	// Initialize the current/root nodes on the first recursion.
	if baseNode == nil {
		key, err := encodeCriterion(baseCriterion)
		if err != nil {
			return nil, err
		}
		baseNode = nodeMap[key]
		if baseNode == nil {
			return nil, fmt.Errorf("base criterion has no samples")
		}
	}
	if root == nil {
		root = baseNode
	}

	// Pick the feature whose true/false branches have the largest confidence gap.
	maxName, err := selectBestSplit(w, nodeMap, targets, minGap)
	if err != nil {
		return nil, err
	}
	if maxName != "" {
		c := targets[maxName]
		negativeCriterion := c[0]
		positiveCriterion := c[1]
		baseNode.FeatureName = &maxName

		// Attach the selected child nodes and continue recursively down both branches.
		negativeKey, err := encodeCriterion(*negativeCriterion)
		if err != nil {
			return nil, err
		}
		baseNode.Negative = nodeMap[negativeKey]

		positiveKey, err := encodeCriterion(*positiveCriterion)
		if err != nil {
			return nil, err
		}
		baseNode.Positive = nodeMap[positiveKey]

		RenderDecisionTree(w, root, "all", 0)

		if _, err := generateDecisionTreeImpl(w, storedKyokus, featureNames, *negativeCriterion, baseNode.Negative, root, minGap); err != nil {
			return nil, err
		}
		if _, err := generateDecisionTreeImpl(w, storedKyokus, featureNames, *positiveCriterion, baseNode.Positive, root, minGap); err != nil {
			return nil, err
		}
	}

	return baseNode, nil
}

func GenerateDecisionTree(featuresPath string, w io.Writer, minGap float64) (*configs.DecisionNode, error) {
	featureNames := FeatureNames()
	storedKyokus, err := LoadStoredKyokus(featuresPath, featureNames)
	if err != nil {
		return nil, err
	}

	return generateDecisionTreeImpl(w, storedKyokus, featureNames, make(Criterion), nil, nil, minGap)
}

func RenderDecisionTree(w io.Writer, node *configs.DecisionNode, label string, indent int) {
	if node == nil {
		return
	}

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
	if err := encoder.Encode(node); err != nil {
		return fmt.Errorf("failed to encode tree file: %w", err)
	}
	return nil
}
