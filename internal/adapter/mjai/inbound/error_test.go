package inbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/inbound"
)

func TestParseMessage_Error(t *testing.T) {
	msg, err := inbound.ParseMessage([]byte(`{"type":"error","message":"invalid join"}`))
	if err != nil {
		t.Fatalf("ParseMessage() failed: %v", err)
	}

	got, ok := msg.(*inbound.Error)
	if !ok {
		t.Fatalf("ParseMessage() = %T, want *inbound.Error", msg)
	}
	if got.Type != "error" {
		t.Errorf("Type = %q, want error", got.Type)
	}
	if got.Message != "invalid join" {
		t.Errorf("Message = %q, want invalid join", got.Message)
	}
}
