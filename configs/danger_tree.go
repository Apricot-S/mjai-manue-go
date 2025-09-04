package configs

import (
	_ "embed"
	"encoding/json/v2"
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

func GetDangerTree() (*DecisionNode, error) {
	var root DecisionNode
	if err := json.Unmarshal(rawDangerTree, &root); err != nil {
		return nil, err
	}
	return &root, nil
}
