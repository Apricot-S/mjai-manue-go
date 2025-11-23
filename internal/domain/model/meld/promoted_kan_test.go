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
		{
			name:         "valid: C-CC-C",
			taken:        *tile.MustTileFromCode("C"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("C"), *tile.MustTileFromCode("C")},
			added:        *tile.MustTileFromCode("C"),
			target:       *playerid.MustPlayerID(3),
			wantTaken:    tile.MustTileFromCode("C"),
			wantConsumed: []tile.Tile{*tile.MustTileFromCode("C"), *tile.MustTileFromCode("C")},
			wantAdded:    tile.MustTileFromCode("C"),
			wantTarget:   playerid.MustPlayerID(3),
			wantErr:      false,
		},
		{
			name:         "valid: 5sr-5s5s-5s",
			taken:        *tile.MustTileFromCode("5sr"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("5s"), *tile.MustTileFromCode("5s")},
			added:        *tile.MustTileFromCode("5s"),
			target:       *playerid.MustPlayerID(0),
			wantTaken:    tile.MustTileFromCode("5sr"),
			wantConsumed: []tile.Tile{*tile.MustTileFromCode("5s"), *tile.MustTileFromCode("5s")},
			wantAdded:    tile.MustTileFromCode("5s"),
			wantTarget:   playerid.MustPlayerID(0),
			wantErr:      false,
		},
		{
			name:         "valid: 5s-5s5sr-5s",
			taken:        *tile.MustTileFromCode("5s"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("5s"), *tile.MustTileFromCode("5sr")},
			added:        *tile.MustTileFromCode("5s"),
			target:       *playerid.MustPlayerID(0),
			wantTaken:    tile.MustTileFromCode("5s"),
			wantConsumed: []tile.Tile{*tile.MustTileFromCode("5s"), *tile.MustTileFromCode("5sr")},
			wantAdded:    tile.MustTileFromCode("5s"),
			wantTarget:   playerid.MustPlayerID(0),
			wantErr:      false,
		},
		{
			name:         "valid: 5s-5s5s-5sr",
			taken:        *tile.MustTileFromCode("5s"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("5s"), *tile.MustTileFromCode("5s")},
			added:        *tile.MustTileFromCode("5sr"),
			target:       *playerid.MustPlayerID(0),
			wantTaken:    tile.MustTileFromCode("5s"),
			wantConsumed: []tile.Tile{*tile.MustTileFromCode("5s"), *tile.MustTileFromCode("5s")},
			wantAdded:    tile.MustTileFromCode("5sr"),
			wantTarget:   playerid.MustPlayerID(0),
			wantErr:      false,
		},
		{
			name:         "invalid: 5s-5s5s-5s",
			taken:        *tile.MustTileFromCode("5s"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("5s"), *tile.MustTileFromCode("5s")},
			added:        *tile.MustTileFromCode("5s"),
			target:       *playerid.MustPlayerID(0),
			wantTaken:    nil,
			wantConsumed: nil,
			wantAdded:    nil,
			wantTarget:   nil,
			wantErr:      true,
		},
		{
			name:         "invalid: ?",
			taken:        *tile.MustTileFromCode("?"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?")},
			added:        *tile.MustTileFromCode("?"),
			target:       *playerid.MustPlayerID(0),
			wantTaken:    nil,
			wantConsumed: nil,
			wantAdded:    nil,
			wantTarget:   nil,
			wantErr:      true,
		},
		{
			name:         "taken and the others do not match",
			taken:        *tile.MustTileFromCode("1m"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("2m"), *tile.MustTileFromCode("2m")},
			added:        *tile.MustTileFromCode("2m"),
			target:       *playerid.MustPlayerID(0),
			wantTaken:    nil,
			wantConsumed: nil,
			wantAdded:    nil,
			wantTarget:   nil,
			wantErr:      true,
		},
		{
			name:         "consumed tiles do not match",
			taken:        *tile.MustTileFromCode("1m"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m")},
			added:        *tile.MustTileFromCode("1m"),
			target:       *playerid.MustPlayerID(0),
			wantTaken:    nil,
			wantConsumed: nil,
			wantAdded:    nil,
			wantTarget:   nil,
			wantErr:      true,
		},
		{
			name:         "added and the others do not match",
			taken:        *tile.MustTileFromCode("2m"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("2m"), *tile.MustTileFromCode("2m")},
			added:        *tile.MustTileFromCode("1m"),
			target:       *playerid.MustPlayerID(0),
			wantTaken:    nil,
			wantConsumed: nil,
			wantAdded:    nil,
			wantTarget:   nil,
			wantErr:      true,
		},
		{
			name:         "sort tiles: 5s-5sr5s-5s to 5s-5s5sr-5s",
			taken:        *tile.MustTileFromCode("5s"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("5sr"), *tile.MustTileFromCode("5s")},
			added:        *tile.MustTileFromCode("5s"),
			target:       *playerid.MustPlayerID(0),
			wantTaken:    tile.MustTileFromCode("5s"),
			wantConsumed: []tile.Tile{*tile.MustTileFromCode("5s"), *tile.MustTileFromCode("5sr")},
			wantAdded:    tile.MustTileFromCode("5s"),
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
