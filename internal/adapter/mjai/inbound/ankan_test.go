package inbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestParseEvent_Ankan(t *testing.T) {
	got := mustParseEventForTest(t, `{"type":"ankan","actor":2,"consumed":["3m","3m","3m","3m"]}`)
	kan, ok := got.(*event.ConcealedKan)
	if !ok {
		t.Fatalf("ParseEvent() = %T, want *event.ConcealedKan", got)
	}
	if kan.Actor() != *seat.MustSeat(2) {
		t.Errorf("Actor() = %v, want %v", kan.Actor(), *seat.MustSeat(2))
	}
	if kan.Consumed() != [4]tile.Tile{tile.MustTileFromCode("3m"), tile.MustTileFromCode("3m"), tile.MustTileFromCode("3m"), tile.MustTileFromCode("3m")} {
		t.Errorf("Consumed() = %v", kan.Consumed())
	}
}

func TestParseEvent_AnkanAdapterAllowsDomainInvalidMeld(t *testing.T) {
	got := mustParseEventForTest(t, `{"type":"ankan","actor":2,"consumed":["3m","3m","3m","4m"]}`)
	if _, ok := got.(*event.ConcealedKan); !ok {
		t.Fatalf("ParseEvent() = %T, want *event.ConcealedKan", got)
	}
}

func TestParseEvent_AnkanInvalidFields(t *testing.T) {
	tests := []string{
		`{"type":"ankan","actor":4,"consumed":["3m","3m","3m","3m"]}`,
		`{"type":"ankan","actor":2,"consumed":["3m","3m","3m"]}`,
		`{"type":"ankan","actor":2,"consumed":["3m","3m","3m","3m","3m"]}`,
		`{"type":"ankan","actor":2,"consumed":["3m","3m","3m","1z"]}`,
		`{"type":"ankan","actor":2,"consumed":["3m","3m","3m","?"]}`,
	}
	for _, payload := range tests {
		t.Run(payload, func(t *testing.T) {
			parseEventShouldFail(t, payload)
		})
	}
}
