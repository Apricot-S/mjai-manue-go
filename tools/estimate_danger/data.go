package main

type MetaData struct {
	FeatureNames []string
}

type StoredKyoku struct {
	Scenes []StoredScene
}

type StoredScene struct {
	Candidates []Candidate
}

type Candidate struct {
	FeatureVector *BitVector
	Hit           bool
}

type Criterion map[string]bool

type CriterionMasks map[string][2]*BitVector
