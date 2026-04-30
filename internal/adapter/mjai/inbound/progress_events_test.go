package inbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func mustParseEventForTest(t *testing.T, payload string) event.Event {
	t.Helper()

	msg, err := inbound.ParseMessage([]byte(payload))
	if err != nil {
		t.Fatalf("ParseMessage() failed: %v", err)
	}
	ev, err := inbound.ParseEvent(msg)
	if err != nil {
		t.Fatalf("ParseEvent() failed: %v", err)
	}
	return ev
}

func TestParseEvent_Reach(t *testing.T) {
	got := mustParseEventForTest(t, `{"type":"reach","actor":2,"cannot_dahai":["1m"]}`)
	reach, ok := got.(*event.Riichi)
	if !ok {
		t.Fatalf("ParseEvent() = %T, want *event.Riichi", got)
	}
	if reach.Actor() != *seat.MustSeat(2) {
		t.Errorf("Actor() = %v, want %v", reach.Actor(), *seat.MustSeat(2))
	}
}

func TestParseEvent_ReachAccepted(t *testing.T) {
	got := mustParseEventForTest(t, `{"type":"reach_accepted","actor":1,"deltas":[0,-1000,0,0],"scores":[25000,24000,25000,25000]}`)
	accepted, ok := got.(*event.RiichiAccepted)
	if !ok {
		t.Fatalf("ParseEvent() = %T, want *event.RiichiAccepted", got)
	}
	if accepted.Actor() != *seat.MustSeat(1) {
		t.Errorf("Actor() = %v, want %v", accepted.Actor(), *seat.MustSeat(1))
	}
	if accepted.Scores() == nil || *accepted.Scores() != [4]int{25000, 24000, 25000, 25000} {
		t.Errorf("Scores() = %v", accepted.Scores())
	}
}

func TestParseEvent_ReachAcceptedActorOnly(t *testing.T) {
	got := mustParseEventForTest(t, `{"type":"reach_accepted","actor":1}`)
	accepted, ok := got.(*event.RiichiAccepted)
	if !ok {
		t.Fatalf("ParseEvent() = %T, want *event.RiichiAccepted", got)
	}
	if accepted.Actor() != *seat.MustSeat(1) {
		t.Errorf("Actor() = %v, want %v", accepted.Actor(), *seat.MustSeat(1))
	}
	if accepted.Scores() != nil {
		t.Errorf("Scores() = %v, want nil", accepted.Scores())
	}
	if accepted.Deltas() != nil {
		t.Errorf("Deltas() = %v, want nil", accepted.Deltas())
	}
}

func TestParseEvent_Calls(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		want    any
	}{
		{
			name:    "pon",
			payload: `{"type":"pon","actor":1,"target":3,"pai":"1s","consumed":["1s","1s"],"cannot_dahai":[]}`,
			want:    (*event.Pon)(nil),
		},
		{
			name:    "chi",
			payload: `{"type":"chi","actor":0,"target":3,"pai":"4s","consumed":["5sr","6s"]}`,
			want:    (*event.Chii)(nil),
		},
		{
			name:    "ankan",
			payload: `{"type":"ankan","actor":2,"consumed":["3m","3m","3m","3m"]}`,
			want:    (*event.ConcealedKan)(nil),
		},
		{
			name:    "kakan",
			payload: `{"type":"kakan","actor":3,"pai":"8m","consumed":["8m","8m","8m"]}`,
			want:    (*event.PromotedKan)(nil),
		},
		{
			name:    "daiminkan",
			payload: `{"type":"daiminkan","actor":1,"target":3,"pai":"4s","consumed":["4s","4s","4s"]}`,
			want:    (*event.CalledKan)(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mustParseEventForTest(t, tt.payload)
			switch tt.want.(type) {
			case *event.Pon:
				pon, ok := got.(*event.Pon)
				if !ok {
					t.Fatalf("ParseEvent() = %T, want *event.Pon", got)
				}
				if pon.Actor() != *seat.MustSeat(1) || pon.Target() != *seat.MustSeat(3) || pon.Taken() != *tile.MustTileFromCode("1s") {
					t.Errorf("Pon event mismatch: %+v", pon)
				}
			case *event.Chii:
				if _, ok := got.(*event.Chii); !ok {
					t.Fatalf("ParseEvent() = %T, want *event.Chii", got)
				}
			case *event.ConcealedKan:
				if _, ok := got.(*event.ConcealedKan); !ok {
					t.Fatalf("ParseEvent() = %T, want *event.ConcealedKan", got)
				}
			case *event.PromotedKan:
				if _, ok := got.(*event.PromotedKan); !ok {
					t.Fatalf("ParseEvent() = %T, want *event.PromotedKan", got)
				}
			case *event.CalledKan:
				if _, ok := got.(*event.CalledKan); !ok {
					t.Fatalf("ParseEvent() = %T, want *event.CalledKan", got)
				}
			}
		})
	}
}

func TestParseEvent_Dora(t *testing.T) {
	got := mustParseEventForTest(t, `{"type":"dora","dora_marker":"6p"}`)
	dora, ok := got.(*event.Dora)
	if !ok {
		t.Fatalf("ParseEvent() = %T, want *event.Dora", got)
	}
	if dora.Indicator() != *tile.MustTileFromCode("6p") {
		t.Errorf("Indicator() = %v, want 6p", dora.Indicator())
	}
}

func TestParseEvent_Hora(t *testing.T) {
	got := mustParseEventForTest(t, `{"type":"hora","actor":2,"target":3,"pai":"9m","uradora_markers":["6m"],"hora_tehais":["5m","5mr","7m","8m","1p","1p","1p","3p","4p","5pr","8s","8s","8s"],"yakus":[["reach",1]],"fu":50,"fan":4,"hora_points":8000,"deltas":[0,0,10300,-8300],"scores":[25000,30800,34700,9500]}`)
	win, ok := got.(*event.Win)
	if !ok {
		t.Fatalf("ParseEvent() = %T, want *event.Win", got)
	}
	if win.Actor() != *seat.MustSeat(2) || win.Target() != *seat.MustSeat(3) {
		t.Errorf("Win seats mismatch: actor=%v target=%v", win.Actor(), win.Target())
	}
	if win.WinningTile() == nil || *win.WinningTile() != *tile.MustTileFromCode("9m") {
		t.Errorf("WinningTile() = %v, want 9m", win.WinningTile())
	}
	if win.WinningPoints() != 8000 {
		t.Errorf("WinningPoints() = %v, want 8000", win.WinningPoints())
	}
	if win.Scores() == nil || *win.Scores() != [4]int{25000, 30800, 34700, 9500} {
		t.Errorf("Scores() = %v", win.Scores())
	}
}

func TestParseEvent_HoraOptionalScoresAndUraMarkersAlias(t *testing.T) {
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

func TestParseEvent_HoraOptionalPai(t *testing.T) {
	got := mustParseEventForTest(t, `{"type":"hora","actor":2,"target":3,"deltas":[0,0,10300,-8300],"hora_points":8000}`)
	win, ok := got.(*event.Win)
	if !ok {
		t.Fatalf("ParseEvent() = %T, want *event.Win", got)
	}
	if win.WinningTile() != nil {
		t.Errorf("WinningTile() = %v, want nil", win.WinningTile())
	}
	if win.WinningPoints() != 8000 {
		t.Errorf("WinningPoints() = %v, want 8000", win.WinningPoints())
	}
}

func TestParseEvent_Ryukyoku(t *testing.T) {
	got := mustParseEventForTest(t, `{"type":"ryukyoku","reason":"fanpai","tehais":[["?"]],"tenpais":[false,true,false,true],"deltas":[-1500,1500,-1500,1500],"scores":[23500,26500,23500,26500]}`)
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
}

func TestParseEvent_RyukyokuOptionalScoresAndDeltas(t *testing.T) {
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

func TestParseEvent_InvalidProgressEventField(t *testing.T) {
	msg, err := inbound.ParseMessage([]byte(`{"type":"pon","actor":1,"target":3,"pai":"1s","consumed":["1s"]}`))
	if err != nil {
		t.Fatalf("ParseMessage() failed: %v", err)
	}
	if _, err := inbound.ParseEvent(msg); err == nil {
		t.Fatal("ParseEvent() succeeded unexpectedly")
	}
}
