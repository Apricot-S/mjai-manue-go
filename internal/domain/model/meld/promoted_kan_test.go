package meld_test

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/playerid"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
)

func TestNewPromotedKan(t *testing.T) {
	tests := []struct {
		name         string
		taken        tile.Tile
		consumed     [2]tile.Tile
		added        tile.Tile
		target       playerid.PlayerID
		wantTaken    *tile.Tile
		wantConsumed []tile.Tile
		wantAdded    *tile.Tile
		wantTarget   *playerid.PlayerID
		wantErr      bool
	}{
		{
			name:         "valid: 1p-1p1p-1p",
			taken:        *tile.MustTileFromCode("1p"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("1p"), *tile.MustTileFromCode("1p")},
			added:        *tile.MustTileFromCode("1p"),
			target:       *playerid.MustPlayerID(0),
			wantTaken:    tile.MustTileFromCode("1p"),
			wantConsumed: []tile.Tile{*tile.MustTileFromCode("1p"), *tile.MustTileFromCode("1p")},
			wantAdded:    tile.MustTileFromCode("1p"),
			wantTarget:   playerid.MustPlayerID(0),
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := meld.NewPromotedKan(tt.taken, tt.consumed, tt.added, tt.target)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewPromotedKan() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewPromotedKan() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got.Taken(), tt.wantTaken) {
				t.Errorf("NewPromotedKan().Taken() = %v, want %v", got.Taken(), tt.wantTaken)
			}
			if !reflect.DeepEqual(got.Consumed(), tt.wantConsumed) {
				t.Errorf("NewPromotedKan().Consumed() = %v, want %v", got.Consumed(), tt.wantConsumed)
			}
			if !reflect.DeepEqual(got.Added(), tt.wantAdded) {
				t.Errorf("NewPromotedKan().Added() = %v, want %v", got.Added(), tt.wantAdded)
			}
			if !reflect.DeepEqual(got.Target(), tt.wantTarget) {
				t.Errorf("NewPromotedKan().Target() = %v, want %v", got.Target(), tt.wantTarget)
			}
		})
	}
}
