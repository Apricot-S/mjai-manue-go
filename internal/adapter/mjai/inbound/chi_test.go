package inbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestParseEvent_Chi(t *testing.T) {
	got := mustParseEventForTest(t, `{"type":"chi","actor":0,"target":3,"pai":"4s","consumed":["5sr","6s"]}`)
	chii, ok := got.(*event.Chii)
	if !ok {
		t.Fatalf("ParseEvent() = %T, want *event.Chii", got)
	}
	if chii.Actor() != *seat.MustSeat(0) {
		t.Errorf("Actor() = %v, want %v", chii.Actor(), *seat.MustSeat(0))
	}
	if chii.Target() != *seat.MustSeat(3) {
		t.Errorf("Target() = %v, want %v", chii.Target(), *seat.MustSeat(3))
	}
	if chii.Taken() != tile.MustTileFromCode("4s") {
		t.Errorf("Taken() = %v, want 4s", chii.Taken())
	}
	if chii.Consumed() != [2]tile.Tile{tile.MustTileFromCode("5sr"), tile.MustTileFromCode("6s")} {
		t.Errorf("Consumed() = %v", chii.Consumed())
	}
}

func TestParseEvent_ChiAdapterAllowsDomainInvalidMeld(t *testing.T) {
	got := mustParseEventForTest(t, `{"type":"chi","actor":0,"target":3,"pai":"4s","consumed":["4s","4s"]}`)
	if _, ok := got.(*event.Chii); !ok {
		t.Fatalf("ParseEvent() = %T, want *event.Chii", got)
	}
}

func TestParseEvent_ChiInvalidFields(t *testing.T) {
	tests := []string{
		`{"type":"chi","actor":4,"target":3,"pai":"4s","consumed":["5s","6s"]}`,
		`{"type":"chi","actor":0,"target":4,"pai":"4s","consumed":["5s","6s"]}`,
		`{"type":"chi","actor":0,"target":3,"pai":"1z","consumed":["5s","6s"]}`,
		`{"type":"chi","actor":0,"target":3,"pai":"?","consumed":["5s","6s"]}`,
		`{"type":"chi","actor":0,"target":3,"pai":"4s","consumed":["5s"]}`,
		`{"type":"chi","actor":0,"target":3,"pai":"4s","consumed":["5s","6s","7s"]}`,
		`{"type":"chi","actor":0,"target":3,"pai":"4s","consumed":["0s","6s"]}`,
		`{"type":"chi","actor":0,"target":3,"pai":"4s","consumed":["?","6s"]}`,
	}
	for _, payload := range tests {
		t.Run(payload, func(t *testing.T) {
			parseEventShouldFail(t, payload)
		})
	}
}
