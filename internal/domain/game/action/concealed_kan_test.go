package action_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestNewConcealedKan(t *testing.T) {
	actor := *seat.MustSeat(1)
	consumed := [4]tile.Tile{
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("5mr"),
	}

	got, err := action.NewConcealedKan(actor, consumed)
	if err != nil {
		t.Fatalf("NewConcealedKan() failed: %v", err)
	}

	if got.Actor() != actor {
		t.Errorf("Actor() = %v, want %v", got.Actor(), actor)
	}
	if got.Consumed() != consumed {
		t.Errorf("Consumed() = %v, want %v", got.Consumed(), consumed)
	}
}

func TestNewConcealedKan_UnknownTile(t *testing.T) {
	actor := *seat.MustSeat(1)
	unknown := tile.MustTileFromCode("?")

	tests := []struct {
		name     string
		consumed [4]tile.Tile
	}{
		{
			name: "consumed",
			consumed: [4]tile.Tile{
				unknown,
				tile.MustTileFromCode("5m"),
				tile.MustTileFromCode("5m"),
				tile.MustTileFromCode("5m"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := action.NewConcealedKan(actor, tt.consumed); err == nil {
				t.Error("NewConcealedKan() succeeded unexpectedly")
			}
		})
	}
}
