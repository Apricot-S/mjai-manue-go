package main

import (
	"encoding/json/v2"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/configs"
)

func TestComputeRatios(t *testing.T) {
	got, err := computeRatios(scoreStats{
		"E1,0": {
			-100: 1,
			100:  3,
		},
	})
	if err != nil {
		t.Fatalf("computeRatios() error = %v", err)
	}

	if got["E1,0"][-100] != 0.25 {
		t.Errorf("ratio for -100 = %v, want 0.25", got["E1,0"][-100])
	}
	if got["E1,0"][100] != 0.75 {
		t.Errorf("ratio for 100 = %v, want 0.75", got["E1,0"][100])
	}
}

func TestLoadStatsFromFileReadsOriginalJSONKeys(t *testing.T) {
	path := filepath.Join(t.TempDir(), "score_stats.json")
	data := []byte(`{"scoreStats":{"E1,0":{"-100":1,"100":3}}}`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("failed to write input: %v", err)
	}

	got, err := loadStatsFromFile(path)
	if err != nil {
		t.Fatalf("loadStatsFromFile() error = %v", err)
	}
	if got.ScoreStats["E1,0"][-100] != 1 {
		t.Errorf("scoreStats[E1,0][-100] = %d, want 1", got.ScoreStats["E1,0"][-100])
	}
	if got.ScoreStats["E1,0"][100] != 3 {
		t.Errorf("scoreStats[E1,0][100] = %d, want 3", got.ScoreStats["E1,0"][100])
	}
}

func TestLoadStatsFromFileRejectsInvalidScoreKey(t *testing.T) {
	path := filepath.Join(t.TempDir(), "score_stats.json")
	data := []byte(`{"scoreStats":{"E1,0":{"bad":1}}}`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("failed to write input: %v", err)
	}

	if _, err := loadStatsFromFile(path); err == nil {
		t.Fatal("loadStatsFromFile() succeeded unexpectedly")
	}
}

func TestBuildEntry(t *testing.T) {
	got := buildEntry(map[int]float64{
		100:  0.25,
		-100: 0.75,
	}, 100)

	if got["0"] != 0.25 {
		t.Errorf("entry[0] = %v, want 0.25", got["0"])
	}
	if got["200"] != 1.0 {
		t.Errorf("entry[200] = %v, want 1", got["200"])
	}
}

func TestRunOutputsLoadableLightGameStats(t *testing.T) {
	path := writeScoreStatsFile(t, input{
		ScoreStats: completeScoreStatsForTest(),
	})

	got, err := run(path)
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}

	assertNear(t, got.WinProbsMap["E1,0,1"]["-100"], 0.25)
	assertNear(t, got.WinProbsMap["E1,0,1"]["100"], 0.75)
	assertNear(t, got.WinProbsMap["E1,0,1"]["300"], 1.0)
	assertNear(t, got.WinProbsMap["E1,1,0"]["-200"], 0.25)
	assertNear(t, got.WinProbsMap["E1,1,0"]["0"], 0.75)
	assertNear(t, got.WinProbsMap["E1,1,0"]["200"], 1.0)

	encoded, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("failed to marshal LightGameStats: %v", err)
	}
	var decoded configs.LightGameStats
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("failed to unmarshal LightGameStats: %v", err)
	}
	if decoded.WinProbsMap["E1,0,1"]["-100"] == 0 {
		t.Error("decoded WinProbsMap[E1,0,1][-100] = 0, want non-zero")
	}
}

func completeScoreStatsForTest() scoreStats {
	stats := make(scoreStats)
	rounds := []string{"E1", "E2", "E3", "E4", "S1", "S2", "S3", "S4"}
	for _, roundName := range rounds {
		for position := range 4 {
			key := roundName + "," + strconv.Itoa(position)
			stats[key] = map[int]int{0: 1}
		}
	}
	stats["E1,0"] = map[int]int{100: 1, -100: 1}
	stats["E1,1"] = map[int]int{100: 1, -100: 1}
	return stats
}

func writeScoreStatsFile(t *testing.T, in input) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "score_stats.json")
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("failed to marshal input: %v", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("failed to write input: %v", err)
	}
	return path
}

func assertNear(t *testing.T, got float64, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-12 {
		t.Errorf("got %v, want %v", got, want)
	}
}
