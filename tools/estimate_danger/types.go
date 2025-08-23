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
