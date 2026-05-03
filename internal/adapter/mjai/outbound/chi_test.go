package outbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/outbound"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestMarshalMessage_Chi(t *testing.T) {
	chii, err := action.NewChii(
		*seat.MustSeat(1),
		*seat.MustSeat(0),
		*tile.MustTileFromCode("3m"),
		[2]tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m")},
	)
	if err != nil {
		t.Fatalf("NewChii() failed: %v", err)
	}
	msg, err := outbound.ToMessage(chii, "")
	if err != nil {
		t.Fatalf("ToMessage() failed: %v", err)
	}

	got, err := outbound.MarshalMessage(msg)
	if err != nil {
		t.Fatalf("MarshalMessage() failed: %v", err)
	}
	if want := `{"type":"chi","actor":1,"target":0,"pai":"3m","consumed":["1m","2m"]}`; string(got) != want {
		t.Errorf("MarshalMessage() = %s, want %s", got, want)
	}
}

func TestMarshalMessage_Chi_Log(t *testing.T) {
	chii, err := action.NewChii(
		*seat.MustSeat(1),
		*seat.MustSeat(0),
		*tile.MustTileFromCode("3m"),
		[2]tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m")},
	)
	if err != nil {
		t.Fatalf("NewChii() failed: %v", err)
	}
	msg, err := outbound.ToMessage(chii, "call chi")
	if err != nil {
		t.Fatalf("ToMessage() failed: %v", err)
	}

	got, err := outbound.MarshalMessage(msg)
	if err != nil {
		t.Fatalf("MarshalMessage() failed: %v", err)
	}
	if want := `{"type":"chi","actor":1,"target":0,"pai":"3m","consumed":["1m","2m"],"log":"call chi"}`; string(got) != want {
		t.Errorf("MarshalMessage() = %s, want %s", got, want)
	}
}
