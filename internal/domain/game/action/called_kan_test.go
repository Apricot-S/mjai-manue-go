package action_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestNewCalledKan(t *testing.T) {
	actor := *seat.MustSeat(2)
	target := *seat.MustSeat(0)
	taken := *tile.MustTileFromCode("E")
	consumed := [3]tile.Tile{
		*tile.MustTileFromCode("E"),
		*tile.MustTileFromCode("E"),
		*tile.MustTileFromCode("E"),
	}

	got, err := action.NewCalledKan(actor, target, taken, consumed)
	if err != nil {
		t.Fatalf("NewCalledKan() failed: %v", err)
	}

	if got.Actor() != actor {
		t.Errorf("Actor() = %v, want %v", got.Actor(), actor)
	}
	if got.Target() != target {
		t.Errorf("Target() = %v, want %v", got.Target(), target)
	}
	if got.Taken() != taken {
		t.Errorf("Taken() = %v, want %v", got.Taken(), taken)
	}
	if got.Consumed() != consumed {
		t.Errorf("Consumed() = %v, want %v", got.Consumed(), consumed)
	}
}

func TestNewCalledKan_UnknownTile(t *testing.T) {
	actor := *seat.MustSeat(2)
	target := *seat.MustSeat(0)
	unknown := *tile.MustTileFromCode("?")

	tests := []struct {
		name     string
		taken    tile.Tile
		consumed [3]tile.Tile
	}{
		{
			name:  "taken",
			taken: unknown,
			consumed: [3]tile.Tile{
				*tile.MustTileFromCode("E"),
				*tile.MustTileFromCode("E"),
				*tile.MustTileFromCode("E"),
			},
		},
		{
			name:  "consumed",
			taken: *tile.MustTileFromCode("E"),
			consumed: [3]tile.Tile{
				unknown,
				*tile.MustTileFromCode("E"),
				*tile.MustTileFromCode("E"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := action.NewCalledKan(actor, target, tt.taken, tt.consumed); err == nil {
				t.Error("NewCalledKan() succeeded unexpectedly")
			}
		})
	}
}
