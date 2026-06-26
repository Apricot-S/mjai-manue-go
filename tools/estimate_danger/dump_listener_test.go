package main

import "testing"

func TestDumpListenerEmptyFilterMatchesAllCandidates(t *testing.T) {
	listener := NewDumpListener("")

	if !listener.meetFilter(&CandidateInfo{}) {
		t.Error("meetFilter() = false, want true")
	}
}
