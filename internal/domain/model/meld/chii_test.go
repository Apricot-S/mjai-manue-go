package meld_test

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
)

func TestNewChii(t *testing.T) {
	tests := []struct {
		name         string
		taken        tile.Tile
		consumed     [2]tile.Tile
		target       int
		wantTaken    *tile.Tile
		wantConsumed []tile.Tile
		wantTarget   int
		wantErr      bool
	}{
		{
			name:         "valid target: 0",
			taken:        *tile.MustTileFromCode("1p"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("2p"), *tile.MustTileFromCode("3p")},
			target:       0,
			wantTaken:    tile.MustTileFromCode("1p"),
			wantConsumed: []tile.Tile{*tile.MustTileFromCode("2p"), *tile.MustTileFromCode("3p")},
			wantTarget:   0,
			wantErr:      false,
		},
		{
			name:         "valid target: 3",
			taken:        *tile.MustTileFromCode("7s"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("8s"), *tile.MustTileFromCode("9s")},
			target:       3,
			wantTaken:    tile.MustTileFromCode("7s"),
			wantConsumed: []tile.Tile{*tile.MustTileFromCode("8s"), *tile.MustTileFromCode("9s")},
			wantTarget:   3,
			wantErr:      false,
		},
		{
			name:         "invalid target: -1",
			taken:        *tile.MustTileFromCode("1m"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("2m"), *tile.MustTileFromCode("3m")},
			target:       -1,
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   -1,
			wantErr:      true,
		},
		{
			name:         "invalid target: 4",
			taken:        *tile.MustTileFromCode("1m"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("2m"), *tile.MustTileFromCode("3m")},
			target:       4,
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   4,
			wantErr:      true,
		},
		{
			name:         "invalid tiles: ?",
			taken:        *tile.MustTileFromCode("?"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?")},
			target:       0,
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   0,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := meld.NewChii(tt.taken, tt.consumed, tt.target)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewChii() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewChii() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got.Taken(), tt.wantTaken) {
				t.Errorf("NewChii().Taken() = %v, want %v", got.Taken(), tt.wantTaken)
			}
			if !reflect.DeepEqual(got.Consumed(), tt.wantConsumed) {
				t.Errorf("NewChii().Consumed() = %v, want %v", got.Consumed(), tt.wantConsumed)
			}
			if !reflect.DeepEqual(got.Target(), tt.wantTarget) {
				t.Errorf("NewChii().Target() = %v, want %v", got.Target(), tt.wantTarget)
			}
		})
	}
}
