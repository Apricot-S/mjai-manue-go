package outbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/outbound"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestMarshalMessage_Ankan(t *testing.T) {
	ankan, err := action.NewConcealedKan(
		seat.MustSeat(1),
		[4]tile.Tile{
			tile.MustTileFromCode("5m"),
			tile.MustTileFromCode("5m"),
			tile.MustTileFromCode("5m"),
			tile.MustTileFromCode("5mr"),
		},
	)
	if err != nil {
		t.Fatalf("NewConcealedKan() failed: %v", err)
	}
	msg, err := outbound.ToMessage(ankan, "")
	if err != nil {
		t.Fatalf("ToMessage() failed: %v", err)
	}

	got, err := outbound.MarshalMessage(msg)
	if err != nil {
		t.Fatalf("MarshalMessage() failed: %v", err)
	}
	if want := `{"type":"ankan","actor":1,"consumed":["5m","5m","5m","5mr"]}`; string(got) != want {
		t.Errorf("MarshalMessage() = %s, want %s", got, want)
	}
}

func TestMarshalMessage_Ankan_Log(t *testing.T) {
	ankan, err := action.NewConcealedKan(
		seat.MustSeat(1),
		[4]tile.Tile{
			tile.MustTileFromCode("5m"),
			tile.MustTileFromCode("5m"),
			tile.MustTileFromCode("5m"),
			tile.MustTileFromCode("5mr"),
		},
	)
	if err != nil {
		t.Fatalf("NewConcealedKan() failed: %v", err)
	}
	msg, err := outbound.ToMessage(ankan, "call ankan")
	if err != nil {
		t.Fatalf("ToMessage() failed: %v", err)
	}

	got, err := outbound.MarshalMessage(msg)
	if err != nil {
		t.Fatalf("MarshalMessage() failed: %v", err)
	}
	if want := `{"type":"ankan","actor":1,"consumed":["5m","5m","5m","5mr"],"log":"call ankan"}`; string(got) != want {
		t.Errorf("MarshalMessage() = %s, want %s", got, want)
	}
}
