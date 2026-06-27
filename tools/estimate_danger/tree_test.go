package main

import (
	"bytes"
	"testing"
)

func TestGenerateDecisionTreeImplSplitsByLargestGap(t *testing.T) {
	featureNames := []string{"safe"}
	storedKyokus := []StoredKyoku{
		{
			Scenes: []StoredScene{
				{
					Candidates: []Candidate{
						{FeatureVector: BoolArrayToBitVector([]bool{true}), Hit: false},
						{FeatureVector: BoolArrayToBitVector([]bool{false}), Hit: true},
					},
				},
			},
		},
		{
			Scenes: []StoredScene{
				{
					Candidates: []Candidate{
						{FeatureVector: BoolArrayToBitVector([]bool{true}), Hit: false},
						{FeatureVector: BoolArrayToBitVector([]bool{false}), Hit: true},
					},
				},
			},
		},
	}

	var out bytes.Buffer
	root, err := generateDecisionTreeImpl(&out, storedKyokus, featureNames, make(Criterion), nil, nil, -1.0)
	if err != nil {
		t.Fatalf("generateDecisionTreeImpl() error = %v", err)
	}

	if root.FeatureName == nil || *root.FeatureName != "safe" {
		t.Fatalf("root.FeatureName = %v, want safe", root.FeatureName)
	}
	if root.Negative == nil {
		t.Fatal("root.Negative = nil")
	}
	if root.Positive == nil {
		t.Fatal("root.Positive = nil")
	}
	if got := root.Negative.AverageProb; got != 1.0 {
		t.Errorf("root.Negative.AverageProb = %v, want 1", got)
	}
	if got := root.Positive.AverageProb; got != 0.0 {
		t.Errorf("root.Positive.AverageProb = %v, want 0", got)
	}
}
