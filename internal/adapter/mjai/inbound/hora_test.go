package inbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestParseEvent_Hora(t *testing.T) {
	got := mustParseEventForTest(t, `{"type":"hora","actor":2,"target":3,"pai":"9m","uradora_markers":["6m"],"ura_markers":["6m"],"hora_tehais":["5m","5mr","7m","8m","1p","1p","1p","3p","4p","5pr","8s","8s","8s"],"yakus":[["reach",1]],"fu":50,"fan":4,"hora_points":8000,"deltas":[0,0,10300,-8300],"scores":[25000,30800,34700,9500]}`)
	win, ok := got.(*event.Win)
	if !ok {
		t.Fatalf("ParseEvent() = %T, want *event.Win", got)
	}
	if win.Actor() != *seat.MustSeat(2) {
		t.Errorf("Actor() = %v, want %v", win.Actor(), *seat.MustSeat(2))
	}
	if win.Target() != *seat.MustSeat(3) {
		t.Errorf("Target() = %v, want %v", win.Target(), *seat.MustSeat(3))
	}
	if win.WinningTile() == nil || *win.WinningTile() != *tile.MustTileFromCode("9m") {
		t.Errorf("WinningTile() = %v, want 9m", win.WinningTile())
	}
	if win.WinningPoints() != 8000 {
		t.Errorf("WinningPoints() = %v, want 8000", win.WinningPoints())
	}
	if win.Deltas() == nil || *win.Deltas() != [4]int{0, 0, 10300, -8300} {
		t.Errorf("Deltas() = %v", win.Deltas())
	}
	if win.Scores() == nil || *win.Scores() != [4]int{25000, 30800, 34700, 9500} {
		t.Errorf("Scores() = %v", win.Scores())
	}
}

func TestParseEvent_HoraOmitemptyFieldsAbsent(t *testing.T) {
	got := mustParseEventForTest(t, `{"type":"hora","actor":2,"target":3}`)
	win, ok := got.(*event.Win)
	if !ok {
		t.Fatalf("ParseEvent() = %T, want *event.Win", got)
	}
	if win.WinningTile() != nil {
		t.Errorf("WinningTile() = %v, want nil", win.WinningTile())
	}
	if win.WinningPoints() != 0 {
		t.Errorf("WinningPoints() = %v, want 0", win.WinningPoints())
	}
	if win.Deltas() != nil {
		t.Errorf("Deltas() = %v, want nil", win.Deltas())
	}
	if win.Scores() != nil {
		t.Errorf("Scores() = %v, want nil", win.Scores())
	}
}

func TestParseEvent_HoraDeltasOnly(t *testing.T) {
	got := mustParseEventForTest(t, `{"type":"hora","actor":2,"target":3,"pai":"9m","ura_markers":["6m"],"deltas":[0,0,10300,-8300],"hora_points":8000}`)
	win, ok := got.(*event.Win)
	if !ok {
		t.Fatalf("ParseEvent() = %T, want *event.Win", got)
	}
	if win.Scores() != nil {
		t.Errorf("Scores() = %v, want nil", win.Scores())
	}
	if win.Deltas() == nil || *win.Deltas() != [4]int{0, 0, 10300, -8300} {
		t.Errorf("Deltas() = %v", win.Deltas())
	}
}

func TestParseEvent_HoraInvalidFields(t *testing.T) {
	tests := []string{
		`{"type":"hora","actor":4,"target":3}`,
		`{"type":"hora","actor":2,"target":4}`,
		`{"type":"hora","actor":2,"target":3,"pai":"1z"}`,
		`{"type":"hora","actor":2,"target":3,"deltas":[0,0,10300]}`,
		`{"type":"hora","actor":2,"target":3,"scores":[25000,30800,34700,34700,34700]}`,
	}
	for _, payload := range tests {
		t.Run(payload, func(t *testing.T) {
			parseEventShouldFail(t, payload)
		})
	}
}
