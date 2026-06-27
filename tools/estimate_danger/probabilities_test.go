package main

import "testing"

func TestCreateKyokuProbabilitiesAveragesScenes(t *testing.T) {
	key, err := encodeCriterion(Criterion{})
	if err != nil {
		t.Fatalf("encodeCriterion() error = %v", err)
	}
	masks, err := createCriterionMasks([]string{"safe"}, []Criterion{{}})
	if err != nil {
		t.Fatalf("createCriterionMasks() error = %v", err)
	}

	kyokuProbs := createKyokuProbabilities(StoredKyoku{
		Scenes: []StoredScene{
			{
				Candidates: []Candidate{
					{FeatureVector: BoolArrayToBitVector([]bool{false}), Hit: true},
					{FeatureVector: BoolArrayToBitVector([]bool{false}), Hit: false},
				},
			},
			{
				Candidates: []Candidate{
					{FeatureVector: BoolArrayToBitVector([]bool{false}), Hit: false},
				},
			},
		},
	}, masks)

	if got := kyokuProbs[key]; got != 0.25 {
		t.Errorf("kyokuProbs[%q] = %v, want 0.25", key, got)
	}
}
