package inbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/infrastructure/mjai/inbound"
)

func TestParseMessage_Hello(t *testing.T) {
	msg, err := inbound.ParseMessage([]byte(`{"type":"hello","protocol":"mjsonp","protocol_version":3}`))
	if err != nil {
		t.Fatalf("ParseMessage() failed: %v", err)
	}

	got, ok := msg.(*inbound.Hello)
	if !ok {
		t.Fatalf("ParseMessage() = %T, want *inbound.Hello", msg)
	}
	if got.Type != "hello" {
		t.Errorf("Type = %q, want hello", got.Type)
	}
	if got.Protocol != "mjsonp" {
		t.Errorf("Protocol = %q, want mjsonp", got.Protocol)
	}
	if got.ProtocolVersion != 3 {
		t.Errorf("ProtocolVersion = %d, want 3", got.ProtocolVersion)
	}
}
