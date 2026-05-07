package outbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/outbound"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
)

func TestMarshalMessage_Reach(t *testing.T) {
	msg, err := outbound.ToMessage(action.NewRiichi(seat.MustSeat(1)), "")
	if err != nil {
		t.Fatalf("ToMessage() failed: %v", err)
	}

	got, err := outbound.MarshalMessage(msg)
	if err != nil {
		t.Fatalf("MarshalMessage() failed: %v", err)
	}
	if want := `{"type":"reach","actor":1}`; string(got) != want {
		t.Errorf("MarshalMessage() = %s, want %s", got, want)
	}
}

func TestMarshalMessage_Reach_Log(t *testing.T) {
	msg, err := outbound.ToMessage(action.NewRiichi(seat.MustSeat(1)), "declare riichi")
	if err != nil {
		t.Fatalf("ToMessage() failed: %v", err)
	}

	got, err := outbound.MarshalMessage(msg)
	if err != nil {
		t.Fatalf("MarshalMessage() failed: %v", err)
	}
	if want := `{"type":"reach","actor":1,"log":"declare riichi"}`; string(got) != want {
		t.Errorf("MarshalMessage() = %s, want %s", got, want)
	}
}

func TestToMessage_Reach(t *testing.T) {
	msg, err := outbound.ToMessage(action.NewRiichi(seat.MustSeat(2)), "")
	if err != nil {
		t.Fatalf("ToMessage() failed: %v", err)
	}

	got, ok := msg.(*outbound.Reach)
	if !ok {
		t.Fatalf("ToMessage() = %T, want *outbound.Reach", msg)
	}
	if got.Type != "reach" {
		t.Errorf("Type = %q, want reach", got.Type)
	}
	if got.Actor != 2 {
		t.Errorf("Actor = %d, want 2", got.Actor)
	}
	if got.Log != "" {
		t.Errorf("Log = %q, want empty", got.Log)
	}
}
