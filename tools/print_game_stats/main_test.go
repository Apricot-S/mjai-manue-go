package main

import (
	"bytes"
	"encoding/json/v2"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/configs"
)

func TestPrintStats(t *testing.T) {
	stats := configs.GameStats{
		NumTurnsDistribution: []float64{0.125, 0.5},
		YamitenStats: map[string]configs.YamitenStat{
			"0,0": {Total: 4, Tenpai: 1},
			"0,1": {Total: 2, Tenpai: 1},
		},
		RyukyokuTenpaiStat: configs.RyukyokuTenpaiStat{
			Total: 10,
			Noten: 6,
			TenpaiTurnDistribution: map[string]int{
				"0":    1,
				"0.25": 2,
				"17.5": 3,
			},
		},
	}

	var buf bytes.Buffer
	printStats(&buf, stats)
	got := buf.String()

	wantContains := []string{
		"numTurnsDistribution:\n   0: 0.125\n   1: 0.500\n",
		"yamitenStats:\n   0: 0.250(    1/    4)  0.500(    1/    2)    NaN(    0/    0)",
		"ryukyokuTenpaiStat:\n   0.00: 0.100 (1)\n   0.25: 0.200 (2)",
		"  17.50: 0.300 (3)\n  noten: 0.600 (6)\n",
	}
	for _, want := range wantContains {
		if !strings.Contains(got, want) {
			t.Errorf("printStats() missing %q in:\n%s", want, got)
		}
	}
}

func TestPrintStatsOmitsMissingSections(t *testing.T) {
	var buf bytes.Buffer
	printStats(&buf, configs.GameStats{})

	if got := buf.String(); got != "" {
		t.Errorf("printStats(empty) = %q, want empty output", got)
	}
}

func TestLoadStatsFromFile(t *testing.T) {
	path := writeGameStatsFile(t, configs.GameStats{
		NumHoras: 1,
		YamitenStats: map[string]configs.YamitenStat{
			"0,0": {Total: 1, Tenpai: 1},
		},
	})

	got, err := loadStatsFromFile(path)
	if err != nil {
		t.Fatalf("loadStatsFromFile() error = %v", err)
	}
	if got.NumHoras != 1 {
		t.Errorf("NumHoras = %d, want 1", got.NumHoras)
	}
	if got.YamitenStats["0,0"].Tenpai != 1 {
		t.Errorf("YamitenStats[0,0].Tenpai = %d, want 1", got.YamitenStats["0,0"].Tenpai)
	}
}

func writeGameStatsFile(t *testing.T, stats configs.GameStats) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "game_stats.json")
	data, err := json.Marshal(stats)
	if err != nil {
		t.Fatalf("failed to marshal stats: %v", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("failed to write stats: %v", err)
	}
	return path
}
