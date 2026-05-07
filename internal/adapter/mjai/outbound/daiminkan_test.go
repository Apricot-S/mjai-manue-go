package outbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/outbound"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestMarshalMessage_Daiminkan(t *testing.T) {
	daiminkan, err := action.NewCalledKan(
		seat.MustSeat(2),
		seat.MustSeat(0),
		tile.MustTileFromCode("E"),
		[3]tile.Tile{
			tile.MustTileFromCode("E"),
			tile.MustTileFromCode("E"),
			tile.MustTileFromCode("E"),
		},
	)
	if err != nil {
		t.Fatalf("NewCalledKan() failed: %v", err)
	}
	msg, err := outbound.ToMessage(daiminkan, "")
	if err != nil {
		t.Fatalf("ToMessage() failed: %v", err)
	}

	got, err := outbound.MarshalMessage(msg)
	if err != nil {
		t.Fatalf("MarshalMessage() failed: %v", err)
	}
	if want := `{"type":"daiminkan","actor":2,"target":0,"pai":"E","consumed":["E","E","E"]}`; string(got) != want {
		t.Errorf("MarshalMessage() = %s, want %s", got, want)
	}
}

func TestMarshalMessage_Daiminkan_Log(t *testing.T) {
	daiminkan, err := action.NewCalledKan(
		seat.MustSeat(2),
		seat.MustSeat(0),
		tile.MustTileFromCode("E"),
		[3]tile.Tile{
			tile.MustTileFromCode("E"),
			tile.MustTileFromCode("E"),
			tile.MustTileFromCode("E"),
		},
	)
	if err != nil {
		t.Fatalf("NewCalledKan() failed: %v", err)
	}
	msg, err := outbound.ToMessage(daiminkan, "call daiminkan")
	if err != nil {
		t.Fatalf("ToMessage() failed: %v", err)
	}

	got, err := outbound.MarshalMessage(msg)
	if err != nil {
		t.Fatalf("MarshalMessage() failed: %v", err)
	}
	if want := `{"type":"daiminkan","actor":2,"target":0,"pai":"E","consumed":["E","E","E"],"log":"call daiminkan"}`; string(got) != want {
		t.Errorf("MarshalMessage() = %s, want %s", got, want)
	}
}
