package inbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
)

func TestParseEvent_ReachAccepted(t *testing.T) {
	tests := []struct {
		name       string
		payload    string
		wantDeltas *[4]int
		wantScores *[4]int
		wantErr    bool
	}{
		{
			name:    "actor only omitempty fields absent",
			payload: `{"type":"reach_accepted","actor":1}`,
		},
		{
			name:       "deltas and scores present",
			payload:    `{"type":"reach_accepted","actor":1,"deltas":[0,-1000,0,0],"scores":[25000,24000,25000,25000]}`,
			wantDeltas: &[4]int{0, -1000, 0, 0},
			wantScores: &[4]int{25000, 24000, 25000, 25000},
		},
		{
			name:       "deltas only",
			payload:    `{"type":"reach_accepted","actor":1,"deltas":[0,-1000,0,0]}`,
			wantDeltas: &[4]int{0, -1000, 0, 0},
		},
		{
			name:       "scores only",
			payload:    `{"type":"reach_accepted","actor":1,"scores":[25000,24000,25000,25000]}`,
			wantScores: &[4]int{25000, 24000, 25000, 25000},
		},
		{
			name:    "invalid actor",
			payload: `{"type":"reach_accepted","actor":4}`,
			wantErr: true,
		},
		{
			name:    "invalid deltas length",
			payload: `{"type":"reach_accepted","actor":1,"deltas":[0,-1000,0]}`,
			wantErr: true,
		},
		{
			name:    "invalid scores length",
			payload: `{"type":"reach_accepted","actor":1,"scores":[25000,24000,25000,25000,25000]}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				parseEventShouldFail(t, tt.payload)
				return
			}
			got := mustParseEventForTest(t, tt.payload)
			accepted, ok := got.(*event.RiichiAccepted)
			if !ok {
				t.Fatalf("ParseEvent() = %T, want *event.RiichiAccepted", got)
			}
			if accepted.Actor() != seat.MustSeat(1) {
				t.Errorf("Actor() = %v, want %v", accepted.Actor(), seat.MustSeat(1))
			}
			if got, want := accepted.Deltas(), tt.wantDeltas; (got == nil) != (want == nil) || got != nil && *got != *want {
				t.Errorf("Deltas() = %v, want %v", got, want)
			}
			if got, want := accepted.Scores(), tt.wantScores; (got == nil) != (want == nil) || got != nil && *got != *want {
				t.Errorf("Scores() = %v, want %v", got, want)
			}
		})
	}
}
