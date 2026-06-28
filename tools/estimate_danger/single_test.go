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

func TestBuildSingleCriteria(t *testing.T) {
	got := buildSingleCriteria([]string{"safe", "dora"})
	want := []Criterion{
		{"safe": false},
		{"safe": true},
		{"dora": false},
		{"dora": true},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("buildSingleCriteria() = %v, want %v", got, want)
	}
}

func TestCalculateSingleProbabilities(t *testing.T) {
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

	var out bytes.Buffer
	if err := CalculateSingleProbabilities(featuresPath, &out); err != nil {
		t.Fatalf("CalculateSingleProbabilities() error = %v", err)
	}

	if got := out.String(); !strings.Contains(got, "map[tsupai:false]\n  0.00 ") {
		t.Errorf("CalculateSingleProbabilities() output does not contain tsupai false probability: %q", got)
	}
}
