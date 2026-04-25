package inbound_test

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/infrastructure/mjai/inbound"
)

func TestParseMessage_EndGame(t *testing.T) {
	msg, err := inbound.ParseMessage([]byte(`{"type":"end_game","scores":[25000,26000,24000,25000]}`))
	if err != nil {
		t.Fatalf("ParseMessage() failed: %v", err)
	}

	got, ok := msg.(*inbound.EndGame)
	if !ok {
		t.Fatalf("ParseMessage() = %T, want *inbound.EndGame", msg)
	}
	if got.Type != "end_game" {
		t.Errorf("Type = %q, want end_game", got.Type)
	}
	if got.Scores == nil {
		t.Fatal("Scores must not be nil")
	}
	if want := []int{25000, 26000, 24000, 25000}; !reflect.DeepEqual(got.Scores, want) {
		t.Errorf("Scores = %v, want %v", got.Scores, want)
	}
}

func TestParseMessage_EndGame_NoScores(t *testing.T) {
	msg, err := inbound.ParseMessage([]byte(`{"type":"end_game"}`))
	if err != nil {
		t.Fatalf("ParseMessage() failed: %v", err)
	}

	got, ok := msg.(*inbound.EndGame)
	if !ok {
		t.Fatalf("ParseMessage() = %T, want *inbound.EndGame", msg)
	}
	if got.Scores != nil {
		t.Errorf("Scores = %v, want nil", got.Scores)
	}
}
