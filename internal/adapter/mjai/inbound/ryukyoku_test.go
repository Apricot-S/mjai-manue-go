package inbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
)

func TestParseEvent_Ryukyoku(t *testing.T) {
	got := mustParseEventForTest(t, `{"type":"ryukyoku","reason":"fanpai","tehais":[["E"],["?"],["?"],["?"]],"tenpais":[false,true,false,true],"deltas":[-1500,1500,-1500,1500],"scores":[23500,26500,23500,26500]}`)
	drawRound, ok := got.(*event.DrawRound)
	if !ok {
		t.Fatalf("ParseEvent() = %T, want *event.DrawRound", got)
	}
	if drawRound.Reason() != "fanpai" {
		t.Errorf("Reason() = %q, want fanpai", drawRound.Reason())
	}
	if drawRound.Tenpais() == nil || *drawRound.Tenpais() != [4]bool{false, true, false, true} {
		t.Errorf("Tenpais() = %v", drawRound.Tenpais())
	}
	if drawRound.Deltas() == nil || *drawRound.Deltas() != [4]int{-1500, 1500, -1500, 1500} {
		t.Errorf("Deltas() = %v", drawRound.Deltas())
	}
	if drawRound.Scores() == nil || *drawRound.Scores() != [4]int{23500, 26500, 23500, 26500} {
		t.Errorf("Scores() = %v", drawRound.Scores())
	}
}

func TestParseEvent_RyukyokuOmitemptyFieldsAbsent(t *testing.T) {
	got := mustParseEventForTest(t, `{"type":"ryukyoku"}`)
	drawRound, ok := got.(*event.DrawRound)
	if !ok {
		t.Fatalf("ParseEvent() = %T, want *event.DrawRound", got)
	}
	if drawRound.Reason() != "" {
		t.Errorf("Reason() = %q, want empty", drawRound.Reason())
	}
	if drawRound.Tenpais() != nil {
		t.Errorf("Tenpais() = %v, want nil", drawRound.Tenpais())
	}
	if drawRound.Deltas() != nil {
		t.Errorf("Deltas() = %v, want nil", drawRound.Deltas())
	}
	if drawRound.Scores() != nil {
		t.Errorf("Scores() = %v, want nil", drawRound.Scores())
	}
}

func TestParseEvent_RyukyokuDeltasOnly(t *testing.T) {
	got := mustParseEventForTest(t, `{"type":"ryukyoku","deltas":[-1500,1500,-1500,1500]}`)
	drawRound, ok := got.(*event.DrawRound)
	if !ok {
		t.Fatalf("ParseEvent() = %T, want *event.DrawRound", got)
	}
	if drawRound.Scores() != nil {
		t.Errorf("Scores() = %v, want nil", drawRound.Scores())
	}
	if drawRound.Deltas() == nil || *drawRound.Deltas() != [4]int{-1500, 1500, -1500, 1500} {
		t.Errorf("Deltas() = %v", drawRound.Deltas())
	}
}

func TestParseEvent_RyukyokuInvalidFields(t *testing.T) {
	tests := []string{
		`{"type":"ryukyoku","tenpais":[false,true,false]}`,
		`{"type":"ryukyoku","deltas":[-1500,1500,-1500,-1500,1500]}`,
		`{"type":"ryukyoku","scores":[23500,26500,23500]}`,
	}
	for _, payload := range tests {
		t.Run(payload, func(t *testing.T) {
			parseEventShouldFail(t, payload)
		})
	}
}
