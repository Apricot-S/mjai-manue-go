package action_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestNewChii(t *testing.T) {
	actor := *seat.MustSeat(1)
	target := *seat.MustSeat(0)
	taken := *tile.MustTileFromCode("3m")
	consumed := [2]tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m")}

	got, err := action.NewChii(actor, target, taken, consumed)
	if err != nil {
		t.Fatalf("NewChii() failed: %v", err)
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

func TestNewChii_UnknownTile(t *testing.T) {
	actor := *seat.MustSeat(1)
	target := *seat.MustSeat(0)
	unknown := *tile.MustTileFromCode("?")
	consumed := [2]tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m")}

	tests := []struct {
		name     string
		taken    tile.Tile
		consumed [2]tile.Tile
	}{
		{
			name:     "taken",
			taken:    unknown,
			consumed: consumed,
		},
		{
			name:     "consumed",
			taken:    *tile.MustTileFromCode("3m"),
			consumed: [2]tile.Tile{unknown, *tile.MustTileFromCode("2m")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := action.NewChii(actor, target, tt.taken, tt.consumed); err == nil {
				t.Error("NewChii() succeeded unexpectedly")
			}
		})
	}
}
