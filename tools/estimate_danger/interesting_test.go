package main

import (
	"bytes"
	"encoding/gob"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestGetNumberCriteria(t *testing.T) {
	tests := []struct {
		name string
		base Criterion
		want []Criterion
	}{
		{
			name: "no conflict",
			base: Criterion{"tsupai": false},
			want: []Criterion{
				{"tsupai": false, "2<=n<=8": false},
				{"tsupai": false, "2<=n<=8": true, "3<=n<=7": false},
				{"tsupai": false, "3<=n<=7": true, "4<=n<=6": false},
				{"tsupai": false, "4<=n<=6": true, "5<=n<=5": false},
				{"tsupai": false, "5<=n<=5": true},
			},
		},
		{
			name: "conflicting existing upper range",
			base: Criterion{"tsupai": false, "4<=n<=6": true},
			want: []Criterion{
				{"tsupai": false, "4<=n<=6": true, "2<=n<=8": false},
				{"tsupai": false, "4<=n<=6": true, "2<=n<=8": true, "3<=n<=7": false},
				nil,
				{"tsupai": false, "4<=n<=6": true, "5<=n<=5": false},
				{"tsupai": false, "4<=n<=6": true, "5<=n<=5": true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getNumberCriteria(tt.base)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getNumberCriteria() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildInterestingCriteriaRemovesNilNumberCriteria(t *testing.T) {
	got := buildInterestingCriteria()
	for i, criterion := range got {
		if criterion == nil {
			t.Fatalf("buildInterestingCriteria()[%d] = nil", i)
		}
	}
	if len(got) == 0 {
		t.Fatal("buildInterestingCriteria() returned no criteria")
	}
}

func TestCalculateInterestingProbabilities(t *testing.T) {
	featuresPath := writeMinimalFeaturesFileForTest(t)

	var out bytes.Buffer
	probs, err := CalculateInterestingProbabilities(featuresPath, &out)
	if err != nil {
		t.Fatalf("CalculateInterestingProbabilities() error = %v", err)
	}

	key, err := encodeCriterion(Criterion{"tsupai": false, "suji": false})
	if err != nil {
		t.Fatalf("encodeCriterion() error = %v", err)
	}
	if probs[key] == nil {
		t.Fatalf("probs[%q] = nil", key)
	}
	if got := out.String(); !strings.Contains(got, "map[suji:false tsupai:false]\n  0.00 ") {
		t.Errorf("CalculateInterestingProbabilities() output does not contain base supai probability: %q", got)
	}
}

func TestRunBenchmark(t *testing.T) {
	featuresPath := writeMinimalFeaturesFileForTest(t)
	if err := RunBenchmark(featuresPath); err != nil {
		t.Fatalf("RunBenchmark() error = %v", err)
	}
}

func writeMinimalFeaturesFileForTest(t *testing.T) string {
	t.Helper()

	featureNames := FeatureNames()
	featuresPath := filepath.Join(t.TempDir(), "features.gob")
	f, err := os.Create(featuresPath)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	encoder := gob.NewEncoder(f)
	if err := encoder.Encode(MetaData{FeatureNames: featureNames}); err != nil {
		_ = f.Close()
		t.Fatalf("Encode(metadata) error = %v", err)
	}
	if err := encoder.Encode([]StoredKyoku{
		{
			Scenes: []StoredScene{
				{
					Candidates: []Candidate{
						{
							FeatureVector: BoolArrayToBitVector(make([]bool, len(featureNames))),
							Hit:           false,
						},
					},
				},
			},
		},
	}); err != nil {
		_ = f.Close()
		t.Fatalf("Encode(kyokus) error = %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	return featuresPath
}
