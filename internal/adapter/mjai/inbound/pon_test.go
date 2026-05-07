package inbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestParseEvent_Pon(t *testing.T) {
	got := mustParseEventForTest(t, `{"type":"pon","actor":1,"target":3,"pai":"1s","consumed":["1s","1s"],"cannot_dahai":[]}`)
	pon, ok := got.(*event.Pon)
	if !ok {
		t.Fatalf("ParseEvent() = %T, want *event.Pon", got)
	}
	if pon.Actor() != *seat.MustSeat(1) {
		t.Errorf("Actor() = %v, want %v", pon.Actor(), *seat.MustSeat(1))
	}
	if pon.Target() != *seat.MustSeat(3) {
		t.Errorf("Target() = %v, want %v", pon.Target(), *seat.MustSeat(3))
	}
	if pon.Taken() != tile.MustTileFromCode("1s") {
		t.Errorf("Taken() = %v, want 1s", pon.Taken())
	}
	if pon.Consumed() != [2]tile.Tile{tile.MustTileFromCode("1s"), tile.MustTileFromCode("1s")} {
		t.Errorf("Consumed() = %v", pon.Consumed())
	}
}

func TestParseEvent_PonAdapterAllowsDomainInvalidMeld(t *testing.T) {
	got := mustParseEventForTest(t, `{"type":"pon","actor":1,"target":3,"pai":"1s","consumed":["2s","2s"]}`)
	if _, ok := got.(*event.Pon); !ok {
		t.Fatalf("ParseEvent() = %T, want *event.Pon", got)
	}
}

func TestParseEvent_PonInvalidFields(t *testing.T) {
	tests := []string{
		`{"type":"pon","actor":4,"target":3,"pai":"1s","consumed":["1s","1s"]}`,
		`{"type":"pon","actor":1,"target":4,"pai":"1s","consumed":["1s","1s"]}`,
		`{"type":"pon","actor":1,"target":3,"pai":"1z","consumed":["1s","1s"]}`,
		`{"type":"pon","actor":1,"target":3,"pai":"?","consumed":["1s","1s"]}`,
		`{"type":"pon","actor":1,"target":3,"pai":"1s","consumed":["1s"]}`,
		`{"type":"pon","actor":1,"target":3,"pai":"1s","consumed":["1s","1s","1s"]}`,
		`{"type":"pon","actor":1,"target":3,"pai":"1s","consumed":["0s","1s"]}`,
		`{"type":"pon","actor":1,"target":3,"pai":"1s","consumed":["?","1s"]}`,
	}
	for _, payload := range tests {
		t.Run(payload, func(t *testing.T) {
			parseEventShouldFail(t, payload)
		})
	}
}
