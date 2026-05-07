package inbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestParseEvent_Dahai(t *testing.T) {
	tests := []struct {
		name          string
		b             []byte
		wantActor     *seat.Seat
		wantTile      *tile.Tile
		wantTsumogiri bool
		wantErr       bool
	}{
		{
			name:          "valid",
			b:             []byte(`{"type":"dahai","actor":1,"pai":"W","tsumogiri":false}`),
			wantActor:     new(seat.MustSeat(1)),
			wantTile:      new(tile.MustTileFromCode("W")),
			wantTsumogiri: false,
			wantErr:       false,
		},
		{
			name:          "valid tsumogiri",
			b:             []byte(`{"type":"dahai","actor":3,"pai":"9m","tsumogiri":true}`),
			wantActor:     new(seat.MustSeat(3)),
			wantTile:      new(tile.MustTileFromCode("9m")),
			wantTsumogiri: true,
			wantErr:       false,
		},
		{
			name:          "unknown tile not allowed",
			b:             []byte(`{"type":"dahai","actor":0,"pai":"?","tsumogiri":false}`),
			wantActor:     nil,
			wantTile:      nil,
			wantTsumogiri: false,
			wantErr:       true,
		},
		{
			name:          "invalid actor",
			b:             []byte(`{"type":"dahai","actor":5,"pai":"E","tsumogiri":false}`),
			wantActor:     nil,
			wantTile:      nil,
			wantTsumogiri: false,
			wantErr:       true,
		},
		{
			name:          "invalid tile",
			b:             []byte(`{"type":"dahai","actor":0,"pai":"1z","tsumogiri":false}`),
			wantActor:     nil,
			wantTile:      nil,
			wantTsumogiri: false,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := inbound.ParseMessage(tt.b)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("ParseMessage() failed: %v", err)
				}
				return
			}

			got, err := inbound.ParseEvent(msg)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("ParseEvent() failed: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Error("ParseEvent() succeeded unexpectedly")
				return
			}

			discard, ok := got.(*event.Discard)
			if !ok {
				t.Fatalf("ParseEvent() = %T, want *event.Discard", got)
			}

			if discard.Actor() != *tt.wantActor {
				t.Errorf("Actor() = %v, want %v", discard.Actor(), *tt.wantActor)
			}
			if discard.Tile() != *tt.wantTile {
				t.Errorf("Tile() = %v, want %v", discard.Tile(), *tt.wantTile)
			}
			if discard.Tsumogiri() != tt.wantTsumogiri {
				t.Errorf("Tsumogiri() = %v, want %v", discard.Tsumogiri(), tt.wantTsumogiri)
			}
		})
	}
}
