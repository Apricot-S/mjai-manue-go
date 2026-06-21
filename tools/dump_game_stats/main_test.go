package main

import (
	"encoding/json/v2"
	"os"
	"path/filepath"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/configs"
)

const horaLog = `{"type":"start_game","names":["a","b","c","d"]}
{"type":"start_kyoku","bakaze":"E","kyoku":1,"honba":0,"kyotaku":0,"oya":0,"dora_marker":"1m","tehais":[["1m","2m","3m","1p","2p","3p","1s","2s","3s","5m","5m","6m","7m"],["1m","1m","2m","2m","3m","3m","4m","4m","5m","5m","6m","6m","7m"],["1s","2s","3s","4s","5s","6s","7s","8s","9s","E","S","W","N"],["1p","1p","2p","2p","3p","3p","4p","4p","5p","5p","6p","6p","7p"]],"scores":[25000,25000,25000,25000]}
{"type":"tsumo","actor":0,"pai":"9m"}
{"type":"hora","actor":0,"target":0,"hora_points":4000,"scores":[29000,24000,23500,23500]}
{"type":"end_kyoku"}
{"type":"end_game"}
`

const ryukyokuLog = `{"type":"start_game","names":["a","b","c","d"]}
{"type":"start_kyoku","bakaze":"E","kyoku":1,"honba":0,"kyotaku":0,"oya":0,"dora_marker":"1m","tehais":[["1m","2m","3m","1p","2p","3p","1s","2s","3s","5m","5m","6m","7m"],["1m","1m","2m","2m","3m","3m","4m","4m","5m","5m","6m","6m","7m"],["1s","2s","3s","4s","5s","6s","7s","8s","9s","E","S","W","N"],["1p","1p","2p","2p","3p","3p","4p","4p","5p","5p","6p","6p","7p"]],"scores":[25000,25000,25000,25000]}
{"type":"tsumo","actor":0,"pai":"9m"}
{"type":"dahai","actor":0,"pai":"9m","tsumogiri":true}
{"type":"ryukyoku","tenpais":[true,false,false,false],"scores":[25000,25000,25000,25000]}
{"type":"end_kyoku"}
{"type":"end_game"}
`

func TestRunHoraStats(t *testing.T) {
	path := writeLogFile(t, horaLog)
	got, err := run([]string{path})
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}

	if got.NumHoras != 1 {
		t.Errorf("NumHoras = %d, want 1", got.NumHoras)
	}
	if got.NumTsumoHoras != 1 {
		t.Errorf("NumTsumoHoras = %d, want 1", got.NumTsumoHoras)
	}
	if got.AverageHoraPoints != 4000 {
		t.Errorf("AverageHoraPoints = %v, want 4000", got.AverageHoraPoints)
	}
	if got.OyaHoraPointsFreqs["total"] != 1 || got.OyaHoraPointsFreqs["4000"] != 1 {
		t.Errorf("OyaHoraPointsFreqs = %v, want total/4000 to be 1", got.OyaHoraPointsFreqs)
	}
	if got.KoHoraPointsFreqs["total"] != 0 {
		t.Errorf("KoHoraPointsFreqs[total] = %d, want 0", got.KoHoraPointsFreqs["total"])
	}
	if got.NumTurnsDistribution[0] != 1 {
		t.Errorf("NumTurnsDistribution[0] = %v, want 1", got.NumTurnsDistribution[0])
	}
}

func TestRunRyukyokuStats(t *testing.T) {
	path := writeLogFile(t, ryukyokuLog)
	got, err := run([]string{path})
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}

	if got.RyukyokuRatio != 1 {
		t.Errorf("RyukyokuRatio = %v, want 1", got.RyukyokuRatio)
	}
	if got.RyukyokuTenpaiStat.Total != 4 {
		t.Errorf("RyukyokuTenpaiStat.Total = %d, want 4", got.RyukyokuTenpaiStat.Total)
	}
	if got.RyukyokuTenpaiStat.Tenpai != 1 {
		t.Errorf("RyukyokuTenpaiStat.Tenpai = %d, want 1", got.RyukyokuTenpaiStat.Tenpai)
	}
	if got.RyukyokuTenpaiStat.Noten != 3 {
		t.Errorf("RyukyokuTenpaiStat.Noten = %d, want 3", got.RyukyokuTenpaiStat.Noten)
	}
	if got.RyukyokuTenpaiStat.TenpaiTurnDistribution["0.25"] != 1 {
		t.Errorf("TenpaiTurnDistribution[0.25] = %d, want 1", got.RyukyokuTenpaiStat.TenpaiTurnDistribution["0.25"])
	}
	if stat := got.YamitenStats["17,0"]; stat.Total != 1 || stat.Tenpai != 1 {
		t.Errorf("YamitenStats[17,0] = %+v, want total/tenpai to be 1", stat)
	}

	data, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("failed to marshal output: %v", err)
	}
	var decoded configs.GameStats
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal GameStats: %v", err)
	}
	if decoded.RyukyokuTenpaiStat.TenpaiTurnDistribution["0.25"] != 1 {
		t.Error("decoded GameStats lost tenpai turn distribution")
	}
}

func TestRunRejectsNoMatches(t *testing.T) {
	if _, err := run([]string{filepath.Join(t.TempDir(), "*.mjson")}); err == nil {
		t.Fatal("run() succeeded unexpectedly")
	}
}

func writeLogFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "game.mjson")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write log: %v", err)
	}
	return path
}
