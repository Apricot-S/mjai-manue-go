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

func TestParseOptionsExtractEmptyFilterIsSet(t *testing.T) {
	opts, paths, err := parseOptions("extract", []string{
		"-o", "features.gob",
		"-filter", "",
		"logs/game.mjson",
	})
	if err != nil {
		t.Fatalf("parseOptions() error = %v", err)
	}

	if !opts.FilterSet {
		t.Error("FilterSet = false, want true")
	}
	if opts.Filter != "" {
		t.Errorf("Filter = %q, want empty string", opts.Filter)
	}
	wantPaths := []string{"logs/game.mjson"}
	if !reflect.DeepEqual(paths, wantPaths) {
		t.Errorf("paths = %v, want %v", paths, wantPaths)
	}
}

func TestParseOptionsExtractFilterUnset(t *testing.T) {
	opts, _, err := parseOptions("extract", []string{
		"-o", "features.gob",
		"logs/game.mjson",
	})
	if err != nil {
		t.Fatalf("parseOptions() error = %v", err)
	}

	if opts.FilterSet {
		t.Error("FilterSet = true, want false")
	}
}

func TestParseOptionsSingle(t *testing.T) {
	_, paths, err := parseOptions("single", []string{"features.gob"})
	if err != nil {
		t.Fatalf("parseOptions() error = %v", err)
	}

	wantPaths := []string{"features.gob"}
	if !reflect.DeepEqual(paths, wantPaths) {
		t.Errorf("paths = %v, want %v", paths, wantPaths)
	}
}

func TestParseOptionsTree(t *testing.T) {
	opts, paths, err := parseOptions("tree", []string{
		"-o", "tree.gob",
		"-min_gap", "2.5",
		"features.gob",
	})
	if err != nil {
		t.Fatalf("parseOptions() error = %v", err)
	}

	if opts.Output != "tree.gob" {
		t.Errorf("Output = %q, want tree.gob", opts.Output)
	}
	if opts.MinGap != 0.025 {
		t.Errorf("MinGap = %v, want 0.025", opts.MinGap)
	}
	wantPaths := []string{"features.gob"}
	if !reflect.DeepEqual(paths, wantPaths) {
		t.Errorf("paths = %v, want %v", paths, wantPaths)
	}
}

func TestParseOptionsDumpTree(t *testing.T) {
	_, paths, err := parseOptions("dump_tree", []string{"tree.gob"})
	if err != nil {
		t.Fatalf("parseOptions() error = %v", err)
	}

	wantPaths := []string{"tree.gob"}
	if !reflect.DeepEqual(paths, wantPaths) {
		t.Errorf("paths = %v, want %v", paths, wantPaths)
	}
}

func TestParseOptionsDumpTreeJSON(t *testing.T) {
	opts, paths, err := parseOptions("dump_tree_json", []string{
		"-o", "danger_tree.all.json",
		"tree.gob",
	})
	if err != nil {
		t.Fatalf("parseOptions() error = %v", err)
	}

	if opts.Output != "danger_tree.all.json" {
		t.Errorf("Output = %q, want danger_tree.all.json", opts.Output)
	}
	wantPaths := []string{"tree.gob"}
	if !reflect.DeepEqual(paths, wantPaths) {
		t.Errorf("paths = %v, want %v", paths, wantPaths)
	}
}
