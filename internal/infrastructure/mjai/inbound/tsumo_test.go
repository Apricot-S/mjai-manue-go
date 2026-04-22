package inbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/infrastructure/mjai/inbound"
)

func TestParseTsumo(t *testing.T) {
	tests := []struct {
		name      string
		b         []byte
		wantActor *seat.Seat
		wantTile  *tile.Tile
		wantErr   bool
	}{
		{
			name:      "valid",
			b:         []byte(`{"type":"tsumo","actor":1,"pai":"E","possible_actions":[]}`),
			wantActor: seat.MustSeat(1),
			wantTile:  tile.MustTileFromCode("E"),
			wantErr:   false,
		},
		{
			name:      "allow unknown",
			b:         []byte(`{"type":"tsumo","actor":0,"pai":"?"}`),
			wantActor: seat.MustSeat(0),
			wantTile:  tile.MustTileFromCode("?"),
			wantErr:   false,
		},
		{
			name:      "invalid actor",
			b:         []byte(`{"type":"tsumo","actor":5,"pai":"?"}`),
			wantActor: nil,
			wantTile:  nil,
			wantErr:   true,
		},
		{
			name:      "invalid tile",
			b:         []byte(`{"type":"tsumo","actor":0,"pai":"1z"}`),
			wantActor: nil,
			wantTile:  nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := inbound.ParseTsumo(tt.b)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ParseTsumo() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ParseTsumo() succeeded unexpectedly")
			}
			if got.Actor() != *tt.wantActor {
				t.Errorf("Actor() = %v, want %v", got, tt.wantActor)
			}
			if got.Tile() != *tt.wantTile {
				t.Errorf("Tile() = %v, want %v", got, tt.wantTile)
			}
		})
	}
}
