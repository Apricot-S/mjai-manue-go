package outbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/outbound"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestMarshalMessage_Pon(t *testing.T) {
	pon, err := action.NewPon(
		*seat.MustSeat(2),
		*seat.MustSeat(0),
		tile.MustTileFromCode("E"),
		[2]tile.Tile{tile.MustTileFromCode("E"), tile.MustTileFromCode("E")},
	)
	if err != nil {
		t.Fatalf("NewPon() failed: %v", err)
	}
	msg, err := outbound.ToMessage(pon, "")
	if err != nil {
		t.Fatalf("ToMessage() failed: %v", err)
	}

	got, err := outbound.MarshalMessage(msg)
	if err != nil {
		t.Fatalf("MarshalMessage() failed: %v", err)
	}
	if want := `{"type":"pon","actor":2,"target":0,"pai":"E","consumed":["E","E"]}`; string(got) != want {
		t.Errorf("MarshalMessage() = %s, want %s", got, want)
	}
}

func TestMarshalMessage_Pon_Log(t *testing.T) {
	pon, err := action.NewPon(
		*seat.MustSeat(2),
		*seat.MustSeat(0),
		tile.MustTileFromCode("E"),
		[2]tile.Tile{tile.MustTileFromCode("E"), tile.MustTileFromCode("E")},
	)
	if err != nil {
		t.Fatalf("NewPon() failed: %v", err)
	}
	msg, err := outbound.ToMessage(pon, "call pon")
	if err != nil {
		t.Fatalf("ToMessage() failed: %v", err)
	}

	got, err := outbound.MarshalMessage(msg)
	if err != nil {
		t.Fatalf("MarshalMessage() failed: %v", err)
	}
	if want := `{"type":"pon","actor":2,"target":0,"pai":"E","consumed":["E","E"],"log":"call pon"}`; string(got) != want {
		t.Errorf("MarshalMessage() = %s, want %s", got, want)
	}
}
