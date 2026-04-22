package inbound_test

import (
	"io"
	"strings"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/infrastructure/mjai/inbound"
)

func TestParseTsumo(t *testing.T) {
	tests := []struct {
		name      string
		r         io.Reader
		wantActor *seat.Seat
		wantTile  *tile.Tile
		wantErr   bool
	}{
		{
			name:      "valid",
			r:         strings.NewReader(`{"type":"tsumo","actor":1,"pai":"E","possible_actions":[]}`),
			wantActor: seat.MustSeat(1),
			wantTile:  tile.MustTileFromCode("E"),
			wantErr:   false,
		},
		{
			name:      "allow unknown",
			r:         strings.NewReader(`{"type":"tsumo","actor":0,"pai":"?"}`),
			wantActor: seat.MustSeat(0),
			wantTile:  tile.MustTileFromCode("?"),
			wantErr:   false,
		},
		{
			name:      "invalid actor",
			r:         strings.NewReader(`{"type":"tsumo","actor":5,"pai":"?"}`),
			wantActor: nil,
			wantTile:  nil,
			wantErr:   true,
		},
		{
			name:      "invalid tile",
			r:         strings.NewReader(`{"type":"tsumo","actor":0,"pai":"1z"}`),
			wantActor: nil,
			wantTile:  nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := inbound.ParseTsumo(tt.r)
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
