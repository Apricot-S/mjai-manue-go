package main

import "testing"

func TestFindPlayer(t *testing.T) {
	got, err := findPlayer([]string{"a", defaultPlayerName, "b", "c"}, defaultPlayerName)
	if err != nil {
		t.Fatalf("findPlayer() failed: %v", err)
	}
	if got != 1 {
		t.Errorf("player index = %d, want 1", got)
	}

	if _, err := findPlayer([]string{"a", "b", "c", "d"}, defaultPlayerName); err == nil {
		t.Error("findPlayer() with no match succeeded, want error")
	}
	if _, err := findPlayer([]string{defaultPlayerName, "b", defaultPlayerName, "d"}, defaultPlayerName); err == nil {
		t.Error("findPlayer() with duplicate matches succeeded, want error")
	}
}
