package inbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/infrastructure/mjai/inbound"
)

func TestParseMessage_EndKyoku(t *testing.T) {
	msg, err := inbound.ParseMessage([]byte(`{"type":"end_kyoku"}`))
	if err != nil {
		t.Fatalf("ParseMessage() failed: %v", err)
	}

	got, ok := msg.(*inbound.EndKyoku)
	if !ok {
		t.Fatalf("ParseMessage() = %T, want *inbound.EndKyoku", msg)
	}
	if got.Type != "end_kyoku" {
		t.Errorf("Type = %q, want end_kyoku", got.Type)
	}
}

func TestParseEvent_EndKyoku(t *testing.T) {
	got, err := (&inbound.EndKyoku{Type: "end_kyoku"}).ToEvent()
	if err != nil {
		t.Fatalf("ToEvent() failed: %v", err)
	}
	if got == nil {
		t.Fatal("ToEvent() = nil, want *event.EndRound")
	}
}

func TestParseEvent_EndKyoku_UnexpectedType(t *testing.T) {
	if _, err := (&inbound.EndKyoku{Type: "end_game"}).ToEvent(); err == nil {
		t.Fatal("ToEvent() succeeded unexpectedly")
	}
}
