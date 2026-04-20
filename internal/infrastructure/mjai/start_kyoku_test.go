package mjai_test

import (
	"encoding/json/v2"
	"strings"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
	"github.com/Apricot-S/mjai-manue-go/internal/infrastructure/mjai"
)

func unknownHand() []string {
	hand := make([]string, 13)
	for i := range hand {
		hand[i] = "?"
	}
	return hand
}

func toJSONHand(hand []string) string {
	data, _ := json.Marshal(hand)
	return string(data)
}

func TestParseStartKyoku_Valid(t *testing.T) {
	payload := `{
		"type":"start_kyoku",
		"bakaze":"E",
		"kyoku":1,
		"honba":0,
		"kyotaku":1,
		"oya":2,
		"dora_marker":"5mr",
		"tehais":[` +
		toJSONHand(unknownHand()) + "," +
		toJSONHand(unknownHand()) + "," +
		toJSONHand(unknownHand()) + "," +
		toJSONHand(unknownHand()) + `],
		"scores":[25000,25000,25000,25000]
	}`

	got, err := mjai.ParseStartKyoku(strings.NewReader(payload))
	if err != nil {
		t.Fatalf("ParseStartKyoku() failed: %v", err)
	}
	if got.RoundWind() != wind.East {
		t.Fatalf("RoundWind() = %v, want %v", got.RoundWind(), wind.East)
	}
	if got.RoundNumber() != 1 {
		t.Fatalf("RoundNumber() = %d, want 1", got.RoundNumber())
	}
	if got.Honba() != 0 {
		t.Fatalf("Honba() = %d, want 0", got.Honba())
	}
	if got.RiichiDeposit() != 1 {
		t.Fatalf("RiichiDeposit() = %d, want 1", got.RiichiDeposit())
	}
	if got.Dealer().Index() != 2 {
		t.Fatalf("Dealer() = %d, want 2", got.Dealer().Index())
	}
	if got.StartingDealer().Index() != 2 {
		t.Fatalf("StartingDealer() = %d, want 2", got.StartingDealer().Index())
	}
	if got.DoraIndicator().String() != "5mr" {
		t.Fatalf("DoraIndicator() = %v, want 5mr", got.DoraIndicator())
	}
	if got.Scores() == nil {
		t.Fatal("scores must not be nil")
	}
	if got.Scores()[0] != 25000 {
		t.Fatalf("Scores()[0] = %d, want 25000", got.Scores()[0])
	}
}

func TestParseStartKyoku_NoScores(t *testing.T) {
	payload := `{
		"type":"start_kyoku",
		"bakaze":"E",
		"kyoku":1,
		"honba":0,
		"kyotaku":0,
		"oya":0,
		"dora_marker":"5mr",
		"tehais":[` +
		toJSONHand(unknownHand()) + "," +
		toJSONHand(unknownHand()) + "," +
		toJSONHand(unknownHand()) + "," +
		toJSONHand(unknownHand()) + `]
	}`

	got, err := mjai.ParseStartKyoku(strings.NewReader(payload))
	if err != nil {
		t.Fatalf("ParseStartKyoku() failed: %v", err)
	}
	if got.Scores() != nil {
		t.Fatalf("Scores() = %v, want nil", got.Scores())
	}
}

func TestParseStartKyoku_InvalidType(t *testing.T) {
	payload := `{
		"type":"start_game",
		"bakaze":"E",
		"kyoku":1,
		"honba":0,
		"kyotaku":0,
		"oya":0,
		"dora_marker":"5mr",
		"tehais":[` +
		toJSONHand(unknownHand()) + "," +
		toJSONHand(unknownHand()) + "," +
		toJSONHand(unknownHand()) + "," +
		toJSONHand(unknownHand()) + `]
	}`

	if _, err := mjai.ParseStartKyoku(strings.NewReader(payload)); err == nil {
		t.Fatal("ParseStartKyoku() succeeded unexpectedly")
	}
}

func TestParseStartKyoku_InvalidDoraMarker(t *testing.T) {
	payload := `{
		"type":"start_kyoku",
		"bakaze":"E",
		"kyoku":1,
		"honba":0,
		"kyotaku":0,
		"oya":0,
		"dora_marker":"invalid",
		"tehais":[` +
		toJSONHand(unknownHand()) + "," +
		toJSONHand(unknownHand()) + "," +
		toJSONHand(unknownHand()) + "," +
		toJSONHand(unknownHand()) + `]
	}`

	if _, err := mjai.ParseStartKyoku(strings.NewReader(payload)); err == nil {
		t.Fatal("ParseStartKyoku() succeeded unexpectedly")
	}
}

func TestParseStartKyoku_InvalidTehaisLength(t *testing.T) {
	payload := `{
		"type":"start_kyoku",
		"bakaze":"E",
		"kyoku":1,
		"honba":0,
		"kyotaku":0,
		"oya":0,
		"dora_marker":"5mr",
		"tehais":[` +
		toJSONHand(unknownHand()) + "," +
		toJSONHand(unknownHand()) + "," +
		toJSONHand(unknownHand()) + `]
	}`

	if _, err := mjai.ParseStartKyoku(strings.NewReader(payload)); err == nil {
		t.Fatal("ParseStartKyoku() succeeded unexpectedly")
	}
}
