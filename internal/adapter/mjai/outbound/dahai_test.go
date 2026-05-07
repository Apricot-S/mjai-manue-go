package outbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/outbound"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestMarshalMessage_Dahai(t *testing.T) {
	discard, err := action.NewDiscard(seat.MustSeat(1), tile.MustTileFromCode("5mr"), false)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}
	msg, err := outbound.ToMessage(discard, "")
	if err != nil {
		t.Fatalf("ToMessage() failed: %v", err)
	}

	got, err := outbound.MarshalMessage(msg)
	if err != nil {
		t.Fatalf("MarshalMessage() failed: %v", err)
	}
	if want := `{"type":"dahai","actor":1,"pai":"5mr","tsumogiri":false}`; string(got) != want {
		t.Errorf("MarshalMessage() = %s, want %s", got, want)
	}
}

func TestMarshalMessage_Dahai_Log(t *testing.T) {
	discard, err := action.NewDiscard(seat.MustSeat(1), tile.MustTileFromCode("5mr"), false)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}
	msg, err := outbound.ToMessage(discard, "discard drawn tile")
	if err != nil {
		t.Fatalf("ToMessage() failed: %v", err)
	}

	got, err := outbound.MarshalMessage(msg)
	if err != nil {
		t.Fatalf("MarshalMessage() failed: %v", err)
	}
	if want := `{"type":"dahai","actor":1,"pai":"5mr","tsumogiri":false,"log":"discard drawn tile"}`; string(got) != want {
		t.Errorf("MarshalMessage() = %s, want %s", got, want)
	}
}

func TestToMessage_Dahai(t *testing.T) {
	discard, err := action.NewDiscard(seat.MustSeat(2), tile.MustTileFromCode("E"), true)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}

	msg, err := outbound.ToMessage(discard, "")
	if err != nil {
		t.Fatalf("ToMessage() failed: %v", err)
	}
	got, ok := msg.(*outbound.Dahai)
	if !ok {
		t.Fatalf("ToMessage() = %T, want *outbound.Dahai", msg)
	}
	if got.Type != "dahai" {
		t.Errorf("Type = %q, want dahai", got.Type)
	}
	if got.Actor != 2 {
		t.Errorf("Actor = %d, want 2", got.Actor)
	}
	if got.Pai != "E" {
		t.Errorf("Pai = %q, want E", got.Pai)
	}
	if !got.Tsumogiri {
		t.Error("Tsumogiri = false, want true")
	}
	if got.Log != "" {
		t.Errorf("Log = %q, want empty", got.Log)
	}
}

func TestToMessage_Dahai_Log(t *testing.T) {
	discard, err := action.NewDiscard(seat.MustSeat(2), tile.MustTileFromCode("E"), true)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}

	msg, err := outbound.ToMessage(discard, "selected by agent")
	if err != nil {
		t.Fatalf("ToMessage() failed: %v", err)
	}
	got, ok := msg.(*outbound.Dahai)
	if !ok {
		t.Fatalf("ToMessage() = %T, want *outbound.Dahai", msg)
	}
	if got.Log != "selected by agent" {
		t.Errorf("Log = %q, want selected by agent", got.Log)
	}
}
