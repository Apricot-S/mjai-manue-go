package inbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestParseEvent_Dora(t *testing.T) {
	got := mustParseEventForTest(t, `{"type":"dora","dora_marker":"6p"}`)
	dora, ok := got.(*event.Dora)
	if !ok {
		t.Fatalf("ParseEvent() = %T, want *event.Dora", got)
	}
	if dora.Indicator() != tile.MustTileFromCode("6p") {
		t.Errorf("Indicator() = %v, want 6p", dora.Indicator())
	}
}

func TestParseEvent_DoraInvalidFields(t *testing.T) {
	tests := []string{
		`{"type":"dora"}`,
		`{"type":"dora","dora_marker":"1z"}`,
		`{"type":"dora","dora_marker":"?"}`,
	}
	for _, payload := range tests {
		t.Run(payload, func(t *testing.T) {
			parseEventShouldFail(t, payload)
		})
	}
}
