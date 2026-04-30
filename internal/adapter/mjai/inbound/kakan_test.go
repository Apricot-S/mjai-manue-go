package inbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestParseEvent_Kakan(t *testing.T) {
	got := mustParseEventForTest(t, `{"type":"kakan","actor":3,"pai":"8m","consumed":["8m","8m","8m"]}`)
	kan, ok := got.(*event.PromotedKan)
	if !ok {
		t.Fatalf("ParseEvent() = %T, want *event.PromotedKan", got)
	}
	if kan.Actor() != *seat.MustSeat(3) {
		t.Errorf("Actor() = %v, want %v", kan.Actor(), *seat.MustSeat(3))
	}
	if kan.Added() != *tile.MustTileFromCode("8m") {
		t.Errorf("Added() = %v, want 8m", kan.Added())
	}
	if kan.Consumed() != [3]tile.Tile{*tile.MustTileFromCode("8m"), *tile.MustTileFromCode("8m"), *tile.MustTileFromCode("8m")} {
		t.Errorf("Consumed() = %v", kan.Consumed())
	}
}

func TestParseEvent_KakanAdapterAllowsDomainInvalidMeld(t *testing.T) {
	got := mustParseEventForTest(t, `{"type":"kakan","actor":3,"pai":"8m","consumed":["7m","7m","7m"]}`)
	if _, ok := got.(*event.PromotedKan); !ok {
		t.Fatalf("ParseEvent() = %T, want *event.PromotedKan", got)
	}
}

func TestParseEvent_KakanInvalidFields(t *testing.T) {
	tests := []string{
		`{"type":"kakan","actor":4,"pai":"8m","consumed":["8m","8m","8m"]}`,
		`{"type":"kakan","actor":3,"pai":"1z","consumed":["8m","8m","8m"]}`,
		`{"type":"kakan","actor":3,"pai":"?","consumed":["8m","8m","8m"]}`,
		`{"type":"kakan","actor":3,"pai":"8m","consumed":["8m","8m"]}`,
		`{"type":"kakan","actor":3,"pai":"8m","consumed":["8m","8m","8m","8m"]}`,
		`{"type":"kakan","actor":3,"pai":"8m","consumed":["8m","8m","0m"]}`,
		`{"type":"kakan","actor":3,"pai":"8m","consumed":["8m","8m","?"]}`,
	}
	for _, payload := range tests {
		t.Run(payload, func(t *testing.T) {
			parseEventShouldFail(t, payload)
		})
	}
}
