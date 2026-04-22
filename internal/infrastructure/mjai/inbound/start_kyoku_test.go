package inbound_test

import (
	"encoding/json/v2"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
	"github.com/Apricot-S/mjai-manue-go/internal/infrastructure/mjai/inbound"
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

	got, err := inbound.ParseStartKyoku([]byte(payload))
	if err != nil {
		t.Fatalf("ParseStartKyoku() failed: %v", err)
	}
	if got.RoundWind() != wind.East {
		t.Errorf("RoundWind() = %v, want %v", got.RoundWind(), wind.East)
	}
	if got.RoundNumber() != 1 {
		t.Errorf("RoundNumber() = %d, want 1", got.RoundNumber())
	}
	if got.Honba() != 0 {
		t.Errorf("Honba() = %d, want 0", got.Honba())
	}
	if got.RiichiDeposit() != 1 {
		t.Errorf("RiichiDeposit() = %d, want 1", got.RiichiDeposit())
	}
	if got.Dealer().Index() != 2 {
		t.Errorf("Dealer() = %d, want 2", got.Dealer().Index())
	}
	if got.DoraIndicator().String() != "5mr" {
		t.Errorf("DoraIndicator() = %v, want 5mr", got.DoraIndicator())
	}
	if got.Scores() == nil {
		t.Fatal("scores must not be nil")
	}
	if got.Scores()[0] != 25000 {
		t.Errorf("Scores()[0] = %d, want 25000", got.Scores()[0])
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

	got, err := inbound.ParseStartKyoku([]byte(payload))
	if err != nil {
		t.Fatalf("ParseStartKyoku() failed: %v", err)
	}
	if got.Scores() != nil {
		t.Errorf("Scores() = %v, want nil", got.Scores())
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

	if _, err := inbound.ParseStartKyoku([]byte(payload)); err == nil {
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

	if _, err := inbound.ParseStartKyoku([]byte(payload)); err == nil {
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

	if _, err := inbound.ParseStartKyoku([]byte(payload)); err == nil {
		t.Fatal("ParseStartKyoku() succeeded unexpectedly")
	}
}
