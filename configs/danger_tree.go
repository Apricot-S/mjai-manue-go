package configs

import (
	_ "embed"
	"encoding/json/v2"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/ai"
)

// DecisionNode represents a node of a decision tree for danger estimation.
type DecisionNode struct {
	AverageProb  float64    `json:"average_prob"`
	ConfInterval [2]float64 `json:"conf_interval"`
	NumSamples   int        `json:"num_samples"`
	// Name of the feature (nil if the node is a leaf node).
	FeatureName *string `json:"feature_name"`
	// Child node if the feature is false (nil if the node is a leaf node).
	Negative *DecisionNode `json:"negative"`
	// Child node if the feature is true (nil if the node is a leaf node).
	Positive *DecisionNode `json:"positive"`
}

//go:embed danger_tree.all.json
var rawDangerTree []byte

func LoadDangerTree() (*DecisionNode, error) {
	var root DecisionNode
	if err := json.Unmarshal(rawDangerTree, &root); err != nil {
		return nil, err
	}
	return &root, nil
}

func (n *DecisionNode) LeafProb() (float64, bool) {
	if n == nil || n.FeatureName != nil {
		return 0, false
	}
	return n.AverageProb, true
}

func (n *DecisionNode) Feature() (string, bool) {
	if n == nil || n.FeatureName == nil {
		return "", false
	}
	return *n.FeatureName, true
}

func (n *DecisionNode) NegativeNode() ai.DangerTreeNode {
	if n == nil {
		return nil
	}
	return n.Negative
}

func (n *DecisionNode) PositiveNode() ai.DangerTreeNode {
	if n == nil {
		return nil
	}
	return n.Positive
}
