package main

import (
	"reflect"
	"testing"
)

func TestParseOptionsExtractExcludePlayers(t *testing.T) {
	opts, paths, err := parseOptions("extract", []string{
		"-o", "features.gob",
		"-exclude_player", "ASAPIN",
		"-exclude_player", "（≧▽≦）",
		"logs/game.mjson",
	})
	if err != nil {
		t.Fatalf("parseOptions() error = %v", err)
	}

	wantExcluded := []string{"ASAPIN", "（≧▽≦）"}
	if got := []string(opts.ExcludePlayers); !reflect.DeepEqual(got, wantExcluded) {
		t.Errorf("ExcludePlayers = %v, want %v", got, wantExcluded)
	}
	wantPaths := []string{"logs/game.mjson"}
	if !reflect.DeepEqual(paths, wantPaths) {
		t.Errorf("paths = %v, want %v", paths, wantPaths)
	}
}
