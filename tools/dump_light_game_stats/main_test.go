package main

import (
	"encoding/json/v2"
	"os"
	"path/filepath"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/inbound"
)

const testLog = `{"type":"start_game","names":["a","b","c","d"]}
{"type":"start_kyoku","bakaze":"E","kyoku":1,"honba":0,"kyotaku":0,"oya":0,"dora_marker":"1m","tehais":[["2m","3m","4m","5m","6m","7m","8m","9m","1p","2p","3p","4p","5p"],["1m","1m","2m","2m","3m","3m","4m","4m","5m","5m","6m","6m","7m"],["1s","2s","3s","4s","5s","6s","7s","8s","9s","E","S","W","N"],["1p","1p","2p","2p","3p","3p","4p","4p","5p","5p","6p","6p","7p"]],"scores":[25000,25000,25000,25000]}
{"type":"ryukyoku","scores":[26000,24000,25000,25000]}
{"type":"end_kyoku"}
{"type":"start_kyoku","bakaze":"E","kyoku":2,"honba":0,"kyotaku":0,"oya":1,"dora_marker":"2m","tehais":[["2m","3m","4m","5m","6m","7m","8m","9m","1p","2p","3p","4p","5p"],["1m","1m","2m","2m","3m","3m","4m","4m","5m","5m","6m","6m","7m"],["1s","2s","3s","4s","5s","6s","7s","8s","9s","E","S","W","N"],["1p","1p","2p","2p","3p","3p","4p","4p","5p","5p","6p","6p","7p"]],"scores":[26000,24000,25000,25000]}
{"type":"ryukyoku","scores":[27000,23000,25000,25000]}
{"type":"end_kyoku"}
{"type":"end_game"}
`

func TestScoreCounter(t *testing.T) {
	counter := newScoreCounter()
	for _, msg := range []string{
		`{"type":"start_game","names":["a","b","c","d"]}`,
		`{"type":"start_kyoku","bakaze":"E","kyoku":1,"scores":[25000,25000,25000,25000],"honba":0,"kyotaku":0,"oya":0,"dora_marker":"1m","tehais":[["2m","3m","4m","5m","6m","7m","8m","9m","1p","2p","3p","4p","5p"],["1m","1m","2m","2m","3m","3m","4m","4m","5m","5m","6m","6m","7m"],["1s","2s","3s","4s","5s","6s","7s","8s","9s","E","S","W","N"],["1p","1p","2p","2p","3p","3p","4p","4p","5p","5p","6p","6p","7p"]]}`,
		`{"type":"ryukyoku","scores":[26000,24000,25000,25000]}`,
		`{"type":"end_game"}`,
	} {
		parsed := mustParseMessage(t, msg)
		if err := counter.onMessage(parsed); err != nil {
			t.Fatalf("onMessage(%s) error = %v", msg, err)
		}
	}

	assertFreq(t, counter.stats, "E1,0", "1000", 1)
	assertFreq(t, counter.stats, "E1,1", "-1000", 1)
	assertFreq(t, counter.stats, "E1,2", "0", 1)
	assertFreq(t, counter.stats, "E1,3", "0", 1)
}

func TestRun(t *testing.T) {
	path := writeLogFile(t, testLog)
	got, err := run([]string{path})
	if err != nil {
		t.Fatalf("run() error = %v", err)
	}

	assertFreq(t, got.ScoreStats, "E1,0", "2000", 1)
	assertFreq(t, got.ScoreStats, "E1,1", "-2000", 1)
	assertFreq(t, got.ScoreStats, "E2,0", "1000", 1)
	assertFreq(t, got.ScoreStats, "E2,1", "-1000", 1)

	data, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("failed to marshal output: %v", err)
	}
	var decoded output
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}
	assertFreq(t, decoded.ScoreStats, "E1,0", "2000", 1)
}

func TestRunRejectsNoMatches(t *testing.T) {
	if _, err := run([]string{filepath.Join(t.TempDir(), "*.mjson")}); err == nil {
		t.Fatal("run() succeeded unexpectedly")
	}
}

func assertFreq(t *testing.T, stats scoreStats, key string, scoreDiff string, want int) {
	t.Helper()
	if got := stats[key][scoreDiff]; got != want {
		t.Errorf("stats[%q][%q] = %d, want %d", key, scoreDiff, got, want)
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

func mustParseMessage(t *testing.T, payload string) inbound.Message {
	t.Helper()
	msg, err := inbound.ParseMessage([]byte(payload))
	if err != nil {
		t.Fatalf("ParseMessage(%s) error = %v", payload, err)
	}
	return msg
}
