package main

import (
	"bytes"
	"encoding/json/v2"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/configs"
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

func TestDumpAndLoadDecisionTreeRoundTrip(t *testing.T) {
	featureName := "safe"
	root := &configs.DecisionNode{
		AverageProb:  0.25,
		ConfInterval: [2]float64{0.1, 0.4},
		NumSamples:   8,
		FeatureName:  &featureName,
		Negative: &configs.DecisionNode{
			AverageProb:  0.5,
			ConfInterval: [2]float64{0.2, 0.8},
			NumSamples:   4,
		},
		Positive: &configs.DecisionNode{
			AverageProb:  0.0,
			ConfInterval: [2]float64{0.0, 0.1},
			NumSamples:   4,
		},
	}

	dir := t.TempDir()
	treePath := filepath.Join(dir, "tree.gob")
	if err := DumpDecisionTree(root, treePath); err != nil {
		t.Fatalf("DumpDecisionTree() error = %v", err)
	}

	loaded, err := LoadDecisionTree(treePath)
	if err != nil {
		t.Fatalf("LoadDecisionTree() error = %v", err)
	}

	if !reflect.DeepEqual(loaded, root) {
		t.Errorf("LoadDecisionTree() = %#v, want %#v", loaded, root)
	}
}

func TestDumpDecisionTreeJSONWritesConfigsCompatibleJSON(t *testing.T) {
	featureName := "safe"
	root := &configs.DecisionNode{
		AverageProb:  0.25,
		ConfInterval: [2]float64{0.1, 0.4},
		NumSamples:   8,
		FeatureName:  &featureName,
		Negative: &configs.DecisionNode{
			AverageProb:  0.5,
			ConfInterval: [2]float64{0.2, 0.8},
			NumSamples:   4,
		},
		Positive: &configs.DecisionNode{
			AverageProb:  0.0,
			ConfInterval: [2]float64{0.0, 0.1},
			NumSamples:   4,
		},
	}

	jsonPath := filepath.Join(t.TempDir(), "danger_tree.all.json")
	if err := DumpDecisionTreeJSON(root, jsonPath); err != nil {
		t.Fatalf("DumpDecisionTreeJSON() error = %v", err)
	}

	b, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	var got configs.DecisionNode
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if !reflect.DeepEqual(&got, root) {
		t.Errorf("DumpDecisionTreeJSON() = %#v, want %#v", &got, root)
	}
}

func TestRunDumpTreeRendersSavedTree(t *testing.T) {
	featureName := "safe"
	root := &configs.DecisionNode{
		AverageProb:  0.25,
		ConfInterval: [2]float64{0.1, 0.4},
		NumSamples:   8,
		FeatureName:  &featureName,
		Negative: &configs.DecisionNode{
			AverageProb:  0.5,
			ConfInterval: [2]float64{0.2, 0.8},
			NumSamples:   4,
		},
		Positive: &configs.DecisionNode{
			AverageProb:  0.0,
			ConfInterval: [2]float64{0.0, 0.1},
			NumSamples:   4,
		},
	}

	treePath := filepath.Join(t.TempDir(), "tree.gob")
	if err := DumpDecisionTree(root, treePath); err != nil {
		t.Fatalf("DumpDecisionTree() error = %v", err)
	}

	var out bytes.Buffer
	if err := runDumpTree(treePath, &out); err != nil {
		t.Fatalf("runDumpTree() error = %v", err)
	}

	want := "all : 25.00 [10.00, 40.00] (8 samples)\n" +
		"  safe = true : 0.00 [0.00, 10.00] (4 samples)\n" +
		"  safe = false : 50.00 [20.00, 80.00] (4 samples)\n"
	if got := out.String(); got != want {
		t.Errorf("runDumpTree() output = %q, want %q", got, want)
	}
}
