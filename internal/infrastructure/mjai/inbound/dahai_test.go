package inbound_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/infrastructure/mjai/inbound"
)

func TestParseDahai(t *testing.T) {
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
			wantActor:     seat.MustSeat(1),
			wantTile:      tile.MustTileFromCode("W"),
			wantTsumogiri: false,
			wantErr:       false,
		},
		{
			name:          "valid tsumogiri",
			b:             []byte(`{"type":"dahai","actor":3,"pai":"9m","tsumogiri":true}`),
			wantActor:     seat.MustSeat(3),
			wantTile:      tile.MustTileFromCode("9m"),
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
			got, err := inbound.ParseDahai(tt.b)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("ParseDahai() failed: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Error("ParseDahai() succeeded unexpectedly")
				return
			}

			if got.Actor() != *tt.wantActor {
				t.Errorf("Actor() = %v, want %v", got.Actor(), *tt.wantActor)
			}
			if got.Tile() != *tt.wantTile {
				t.Errorf("Tile() = %v, want %v", got.Tile(), *tt.wantTile)
			}
			if got.Tsumogiri() != tt.wantTsumogiri {
				t.Errorf("Tsumogiri() = %v, want %v", got.Tsumogiri(), tt.wantTsumogiri)
			}
		})
	}
}
