package action_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestNewPromotedKan(t *testing.T) {
	actor := seat.MustSeat(1)
	added := tile.MustTileFromCode("5mr")
	consumed := [3]tile.Tile{
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("5m"),
	}

	got, err := action.NewPromotedKan(actor, added, consumed)
	if err != nil {
		t.Fatalf("NewPromotedKan() failed: %v", err)
	}

	if got.Actor() != actor {
		t.Errorf("Actor() = %v, want %v", got.Actor(), actor)
	}
	if got.Added() != added {
		t.Errorf("Added() = %v, want %v", got.Added(), added)
	}
	if got.Consumed() != consumed {
		t.Errorf("Consumed() = %v, want %v", got.Consumed(), consumed)
	}
}

func TestNewPromotedKan_UnknownTile(t *testing.T) {
	actor := seat.MustSeat(1)
	unknown := tile.MustTileFromCode("?")

	tests := []struct {
		name     string
		added    tile.Tile
		consumed [3]tile.Tile
	}{
		{
			name:  "added",
			added: unknown,
			consumed: [3]tile.Tile{
				tile.MustTileFromCode("5m"),
				tile.MustTileFromCode("5m"),
				tile.MustTileFromCode("5m"),
			},
		},
		{
			name:  "consumed",
			added: tile.MustTileFromCode("5mr"),
			consumed: [3]tile.Tile{
				unknown,
				tile.MustTileFromCode("5m"),
				tile.MustTileFromCode("5m"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := action.NewPromotedKan(actor, tt.added, tt.consumed); err == nil {
				t.Error("NewPromotedKan() succeeded unexpectedly")
			}
		})
	}
}
