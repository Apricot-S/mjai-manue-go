package inbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
)

func TestParseEvent_Reach(t *testing.T) {
	got := mustParseEventForTest(t, `{"type":"reach","actor":2,"cannot_dahai":["1m"]}`)
	reach, ok := got.(*event.Riichi)
	if !ok {
		t.Fatalf("ParseEvent() = %T, want *event.Riichi", got)
	}
	if reach.Actor() != seat.MustSeat(2) {
		t.Errorf("Actor() = %v, want %v", reach.Actor(), seat.MustSeat(2))
	}
}

func TestParseEvent_ReachInvalidActor(t *testing.T) {
	parseEventShouldFail(t, `{"type":"reach","actor":4}`)
}
