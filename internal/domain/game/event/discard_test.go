package event

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestNewDiscard(t *testing.T) {
	validActor := *seat.MustSeat(1)
	validTile := *tile.MustTileFromCode("1m")

	tests := []struct {
		name      string
		actor     seat.Seat
		tile      tile.Tile
		tsumogiri bool
		wantErr   bool
	}{
		{
			name:      "valid parameters",
			actor:     validActor,
			tile:      validTile,
			tsumogiri: false,
			wantErr:   false,
		},
		{
			name:      "tsumogiri true allowed",
			actor:     validActor,
			tile:      validTile,
			tsumogiri: true,
			wantErr:   false,
		},
		{
			name:      "unknown tile not allowed",
			actor:     validActor,
			tile:      *tile.MustTileFromCode("?"),
			tsumogiri: false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDiscard(tt.actor, tt.tile, tt.tsumogiri)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NewDiscard() failed: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewDiscard() succeeded unexpectedly")
			}
			if got == nil {
				t.Fatal("NewDiscard() returned nil without error")
			}
		})
	}
}

func TestDiscardAccessors(t *testing.T) {
	actor := *seat.MustSeat(2)
	discardedTile := *tile.MustTileFromCode("E")

	got, err := NewDiscard(actor, discardedTile, true)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}

	if got.Actor().Index() != actor.Index() {
		t.Errorf("Actor() = %v, want %v", got.Actor(), actor)
	}
	if got.Tile().ID() != discardedTile.ID() {
		t.Errorf("Tile() = %v, want %v", got.Tile(), discardedTile)
	}
	if got.Tsumogiri() != true {
		t.Errorf("Tsumogiri() = %v, want %v", got.Tsumogiri(), true)
	}
}
