package inbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestParseEvent_Daiminkan(t *testing.T) {
	got := mustParseEventForTest(t, `{"type":"daiminkan","actor":1,"target":3,"pai":"4s","consumed":["4s","4s","4s"]}`)
	kan, ok := got.(*event.CalledKan)
	if !ok {
		t.Fatalf("ParseEvent() = %T, want *event.CalledKan", got)
	}
	if kan.Actor() != *seat.MustSeat(1) {
		t.Errorf("Actor() = %v, want %v", kan.Actor(), *seat.MustSeat(1))
	}
	if kan.Target() != *seat.MustSeat(3) {
		t.Errorf("Target() = %v, want %v", kan.Target(), *seat.MustSeat(3))
	}
	if kan.Taken() != tile.MustTileFromCode("4s") {
		t.Errorf("Taken() = %v, want 4s", kan.Taken())
	}
	if kan.Consumed() != [3]tile.Tile{tile.MustTileFromCode("4s"), tile.MustTileFromCode("4s"), tile.MustTileFromCode("4s")} {
		t.Errorf("Consumed() = %v", kan.Consumed())
	}
}

func TestParseEvent_DaiminkanAdapterAllowsDomainInvalidMeld(t *testing.T) {
	got := mustParseEventForTest(t, `{"type":"daiminkan","actor":1,"target":3,"pai":"4s","consumed":["5s","5s","5s"]}`)
	if _, ok := got.(*event.CalledKan); !ok {
		t.Fatalf("ParseEvent() = %T, want *event.CalledKan", got)
	}
}

func TestParseEvent_DaiminkanInvalidFields(t *testing.T) {
	tests := []string{
		`{"type":"daiminkan","actor":4,"target":3,"pai":"4s","consumed":["4s","4s","4s"]}`,
		`{"type":"daiminkan","actor":1,"target":4,"pai":"4s","consumed":["4s","4s","4s"]}`,
		`{"type":"daiminkan","actor":1,"target":3,"pai":"1z","consumed":["4s","4s","4s"]}`,
		`{"type":"daiminkan","actor":1,"target":3,"pai":"?","consumed":["4s","4s","4s"]}`,
		`{"type":"daiminkan","actor":1,"target":3,"pai":"4s","consumed":["4s","4s"]}`,
		`{"type":"daiminkan","actor":1,"target":3,"pai":"4s","consumed":["4s","4s","4s","4s"]}`,
		`{"type":"daiminkan","actor":1,"target":3,"pai":"4s","consumed":["4s","4s","0s"]}`,
		`{"type":"daiminkan","actor":1,"target":3,"pai":"4s","consumed":["4s","4s","?"]}`,
	}
	for _, payload := range tests {
		t.Run(payload, func(t *testing.T) {
			parseEventShouldFail(t, payload)
		})
	}
}
