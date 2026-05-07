package outbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/outbound"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestMarshalMessage_Kakan(t *testing.T) {
	kakan, err := action.NewPromotedKan(
		*seat.MustSeat(1),
		tile.MustTileFromCode("5mr"),
		[3]tile.Tile{
			tile.MustTileFromCode("5m"),
			tile.MustTileFromCode("5m"),
			tile.MustTileFromCode("5m"),
		},
	)
	if err != nil {
		t.Fatalf("NewPromotedKan() failed: %v", err)
	}
	msg, err := outbound.ToMessage(kakan, "")
	if err != nil {
		t.Fatalf("ToMessage() failed: %v", err)
	}

	got, err := outbound.MarshalMessage(msg)
	if err != nil {
		t.Fatalf("MarshalMessage() failed: %v", err)
	}
	if want := `{"type":"kakan","actor":1,"pai":"5mr","consumed":["5m","5m","5m"]}`; string(got) != want {
		t.Errorf("MarshalMessage() = %s, want %s", got, want)
	}
}

func TestMarshalMessage_Kakan_Log(t *testing.T) {
	kakan, err := action.NewPromotedKan(
		*seat.MustSeat(1),
		tile.MustTileFromCode("5mr"),
		[3]tile.Tile{
			tile.MustTileFromCode("5m"),
			tile.MustTileFromCode("5m"),
			tile.MustTileFromCode("5m"),
		},
	)
	if err != nil {
		t.Fatalf("NewPromotedKan() failed: %v", err)
	}
	msg, err := outbound.ToMessage(kakan, "call kakan")
	if err != nil {
		t.Fatalf("ToMessage() failed: %v", err)
	}

	got, err := outbound.MarshalMessage(msg)
	if err != nil {
		t.Fatalf("MarshalMessage() failed: %v", err)
	}
	if want := `{"type":"kakan","actor":1,"pai":"5mr","consumed":["5m","5m","5m"],"log":"call kakan"}`; string(got) != want {
		t.Errorf("MarshalMessage() = %s, want %s", got, want)
	}
}
