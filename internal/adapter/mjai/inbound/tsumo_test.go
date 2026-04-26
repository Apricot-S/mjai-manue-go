package inbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestParseEvent_Tsumo(t *testing.T) {
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
			msg, gotErr := inbound.ParseMessage(tt.b)
			var got event.Event
			if gotErr == nil {
				got, gotErr = inbound.ParseEvent(msg)
			}
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ParseEvent() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ParseEvent() succeeded unexpectedly")
			}

			draw, ok := got.(*event.Draw)
			if !ok {
				t.Fatalf("ParseEvent() = %T, want *event.Draw", got)
			}
			if draw.Actor() != *tt.wantActor {
				t.Errorf("Actor() = %v, want %v", draw.Actor(), *tt.wantActor)
			}
			if draw.Tile() != *tt.wantTile {
				t.Errorf("Tile() = %v, want %v", draw.Tile(), *tt.wantTile)
			}
		})
	}
}
