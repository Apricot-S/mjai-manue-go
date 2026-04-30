package inbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
)

func mustParseEventForTest(t *testing.T, payload string) event.Event {
	t.Helper()

	msg, err := inbound.ParseMessage([]byte(payload))
	if err != nil {
		t.Fatalf("ParseMessage() failed: %v", err)
	}
	ev, err := inbound.ParseEvent(msg)
	if err != nil {
		t.Fatalf("ParseEvent() failed: %v", err)
	}
	return ev
}

func parseEventShouldFail(t *testing.T, payload string) {
	t.Helper()

	msg, err := inbound.ParseMessage([]byte(payload))
	if err != nil {
		return
	}
	if _, err := inbound.ParseEvent(msg); err == nil {
		t.Fatal("ParseEvent() succeeded unexpectedly")
	}
}
