package main

import (
	"fmt"
	"strings"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/game"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
)

type DumpListener struct {
	filter map[string]string
}

func NewDumpListener(filterSpec string) *DumpListener {
	filter := make(map[string]string)
	fields := strings.SplitSeq(filterSpec, "&")
	for field := range fields {
		k, v, found := strings.Cut(field, ":")
		if found {
			filter[k] = v
		}
	}
	return &DumpListener{filter: filter}
}

func (dl *DumpListener) OnDahai(
	path string,
	state game.StateViewer,
	action inbound.Event,
	reacher *base.Player,
	candidates []CandidateInfo,
) {
	var cands []CandidateInfo
	for _, c := range candidates {
		if dl.meetFilter(&c) {
			cands = append(cands, c)
		}
	}

	if len(cands) == 0 {
		return
	}

	fmt.Println(path)
	// state.DumpAction(action)
	fmt.Printf("reacher: %d\n", reacher.ID())
	for _, cand := range cands {
		h := 0
		if cand.Hit {
			h = 1
		}
		fmt.Printf("candidate %s: hit=%d, %s\n", cand.Pai.ToString(), h, FeatureVectorToStr(cand.FeatureVector))
	}
	fmt.Println(strings.Repeat("=", 80))
}

func (dl *DumpListener) meetFilter(cand *CandidateInfo) bool {
	for k, v := range dl.filter {
		expected := v == "1"
		if k == "hit" {
			if cand.Hit != expected {
				return false
			}
		} else {
			actual := GetFeatureValue(cand.FeatureVector, k)
			if actual != expected {
				return false
			}
		}
	}
	return true
}
