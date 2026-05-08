package outbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/outbound"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
)

func TestMarshalMessage_Kyushukyuhai(t *testing.T) {
	msg, err := outbound.ToMessage(action.NewKyushukyuhai(seat.MustSeat(1)), "")
	if err != nil {
		t.Fatalf("ToMessage() failed: %v", err)
	}

	got, err := outbound.MarshalMessage(msg)
	if err != nil {
		t.Fatalf("MarshalMessage() failed: %v", err)
	}
	if want := `{"type":"ryukyoku","reason":"kyushukyuhai","actor":1}`; string(got) != want {
		t.Errorf("MarshalMessage() = %s, want %s", got, want)
	}
}

func TestMarshalMessage_Kyushukyuhai_Log(t *testing.T) {
	msg, err := outbound.ToMessage(action.NewKyushukyuhai(seat.MustSeat(1)), "abortive draw")
	if err != nil {
		t.Fatalf("ToMessage() failed: %v", err)
	}

	got, err := outbound.MarshalMessage(msg)
	if err != nil {
		t.Fatalf("MarshalMessage() failed: %v", err)
	}
	if want := `{"type":"ryukyoku","reason":"kyushukyuhai","actor":1,"log":"abortive draw"}`; string(got) != want {
		t.Errorf("MarshalMessage() = %s, want %s", got, want)
	}
}

func TestToMessage_Kyushukyuhai(t *testing.T) {
	msg, err := outbound.ToMessage(action.NewKyushukyuhai(seat.MustSeat(2)), "")
	if err != nil {
		t.Fatalf("ToMessage() failed: %v", err)
	}

	got, ok := msg.(*outbound.Kyushukyuhai)
	if !ok {
		t.Fatalf("ToMessage() = %T, want *outbound.Kyushukyuhai", msg)
	}
	if got.Type != "ryukyoku" {
		t.Errorf("Type = %q, want ryukyoku", got.Type)
	}
	if got.Reason != "kyushukyuhai" {
		t.Errorf("Reason = %q, want kyushukyuhai", got.Reason)
	}
	if got.Actor != 2 {
		t.Errorf("Actor = %d, want 2", got.Actor)
	}
	if got.Log != "" {
		t.Errorf("Log = %q, want empty", got.Log)
	}
}
