package main

import "github.com/Apricot-S/mjai-manue-go/internal/base"

type MetaData struct {
	FeatureNames []string
}

type StoredKyoku struct {
	Scenes []StoredScene
}

type StoredScene struct {
	Candidates []CandidateData
}

type CandidateData struct {
	FeatureVector *BitVector
	Hit           bool
}

type CandidateInfo struct {
	Pai           base.Pai
	Hit           bool
	FeatureVector *BitVector
}
