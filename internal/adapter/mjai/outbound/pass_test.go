package outbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/outbound"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
)

func TestMarshalMessage_Pass(t *testing.T) {
	msg, err := outbound.ToMessage(action.NewPass(seat.MustSeat(2)), "")
	if err != nil {
		t.Fatalf("ToMessage() failed: %v", err)
	}

	got, err := outbound.MarshalMessage(msg)
	if err != nil {
		t.Fatalf("MarshalMessage() failed: %v", err)
	}
	if want := `{"type":"none","actor":2}`; string(got) != want {
		t.Errorf("MarshalMessage() = %s, want %s", got, want)
	}
}

func TestMarshalMessage_Pass_Log(t *testing.T) {
	msg, err := outbound.ToMessage(action.NewPass(seat.MustSeat(2)), "decline call")
	if err != nil {
		t.Fatalf("ToMessage() failed: %v", err)
	}

	got, err := outbound.MarshalMessage(msg)
	if err != nil {
		t.Fatalf("MarshalMessage() failed: %v", err)
	}
	if want := `{"type":"none","actor":2,"log":"decline call"}`; string(got) != want {
		t.Errorf("MarshalMessage() = %s, want %s", got, want)
	}
}

func TestToMessage_Pass(t *testing.T) {
	msg, err := outbound.ToMessage(action.NewPass(seat.MustSeat(1)), "")
	if err != nil {
		t.Fatalf("ToMessage() failed: %v", err)
	}

	got, ok := msg.(*outbound.Pass)
	if !ok {
		t.Fatalf("ToMessage() = %T, want *outbound.Pass", msg)
	}
	if got.Type != "none" {
		t.Errorf("Type = %q, want none", got.Type)
	}
	if got.Actor != 1 {
		t.Errorf("Actor = %d, want 1", got.Actor)
	}
	if got.Log != "" {
		t.Errorf("Log = %q, want empty", got.Log)
	}
}
