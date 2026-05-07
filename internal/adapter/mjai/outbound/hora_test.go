package outbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/outbound"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestMarshalMessage_Hora(t *testing.T) {
	win, err := action.NewWin(*seat.MustSeat(1), *seat.MustSeat(0), tile.MustTileFromCode("5mr"))
	if err != nil {
		t.Fatalf("NewWin() failed: %v", err)
	}
	msg, err := outbound.ToMessage(win, "")
	if err != nil {
		t.Fatalf("ToMessage() failed: %v", err)
	}

	got, err := outbound.MarshalMessage(msg)
	if err != nil {
		t.Fatalf("MarshalMessage() failed: %v", err)
	}
	if want := `{"type":"hora","actor":1,"target":0,"pai":"5mr"}`; string(got) != want {
		t.Errorf("MarshalMessage() = %s, want %s", got, want)
	}
}

func TestMarshalMessage_Hora_Log(t *testing.T) {
	win, err := action.NewWin(*seat.MustSeat(1), *seat.MustSeat(0), tile.MustTileFromCode("5mr"))
	if err != nil {
		t.Fatalf("NewWin() failed: %v", err)
	}
	msg, err := outbound.ToMessage(win, "win")
	if err != nil {
		t.Fatalf("ToMessage() failed: %v", err)
	}

	got, err := outbound.MarshalMessage(msg)
	if err != nil {
		t.Fatalf("MarshalMessage() failed: %v", err)
	}
	if want := `{"type":"hora","actor":1,"target":0,"pai":"5mr","log":"win"}`; string(got) != want {
		t.Errorf("MarshalMessage() = %s, want %s", got, want)
	}
}
