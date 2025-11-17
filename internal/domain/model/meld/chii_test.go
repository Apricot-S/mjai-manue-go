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
			name:         "valid tiles: 8-79m",
			taken:        *tile.MustTileFromCode("8m"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("7m"), *tile.MustTileFromCode("9m")},
			target:       0,
			wantTaken:    tile.MustTileFromCode("8m"),
			wantConsumed: []tile.Tile{*tile.MustTileFromCode("7m"), *tile.MustTileFromCode("9m")},
			wantTarget:   0,
			wantErr:      false,
		},
		{
			name:         "valid tiles: 5r-46m",
			taken:        *tile.MustTileFromCode("5mr"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("4m"), *tile.MustTileFromCode("6m")},
			target:       0,
			wantTaken:    tile.MustTileFromCode("5mr"),
			wantConsumed: []tile.Tile{*tile.MustTileFromCode("4m"), *tile.MustTileFromCode("6m")},
			wantTarget:   0,
			wantErr:      false,
		},
		{
			name:         "invalid tiles: ?-??",
			taken:        *tile.MustTileFromCode("?"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?")},
			target:       0,
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   0,
			wantErr:      true,
		},
		{
			name:         "invalid tiles: E-SW",
			taken:        *tile.MustTileFromCode("E"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("S"), *tile.MustTileFromCode("W")},
			target:       0,
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   0,
			wantErr:      true,
		},
		{
			name:         "invalid tiles: 8-9m1p",
			taken:        *tile.MustTileFromCode("8m"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("9m"), *tile.MustTileFromCode("1p")},
			target:       0,
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   0,
			wantErr:      true,
		},
		{
			name:         "invalid tiles: 6-89m",
			taken:        *tile.MustTileFromCode("6m"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("8m"), *tile.MustTileFromCode("9m")},
			target:       0,
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   0,
			wantErr:      true,
		},
		{
			name:         "invalid tiles: 6-79m",
			taken:        *tile.MustTileFromCode("6m"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("7m"), *tile.MustTileFromCode("9m")},
			target:       0,
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   0,
			wantErr:      true,
		},
		{
			name:         "invalid tiles: 6-7m8p",
			taken:        *tile.MustTileFromCode("6m"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("7m"), *tile.MustTileFromCode("8p")},
			target:       0,
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   0,
			wantErr:      true,
		},
		{
			name:         "sort tiles: 8-97s to 8-79s",
			taken:        *tile.MustTileFromCode("8s"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("9s"), *tile.MustTileFromCode("7s")},
			target:       3,
			wantTaken:    tile.MustTileFromCode("8s"),
			wantConsumed: []tile.Tile{*tile.MustTileFromCode("7s"), *tile.MustTileFromCode("9s")},
			wantTarget:   3,
			wantErr:      false,
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
