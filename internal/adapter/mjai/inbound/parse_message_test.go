package inbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/inbound"
)

func TestParseMessage_JoinUnsupported(t *testing.T) {
	if _, err := inbound.ParseMessage([]byte(`{"type":"join","name":"bot","room":"default"}`)); err == nil {
		t.Fatal("ParseMessage() succeeded unexpectedly")
	}
}
