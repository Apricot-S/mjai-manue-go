package meld_test

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestNewCalledKan(t *testing.T) {
	tests := []struct {
		name         string
		taken        tile.Tile
		consumed     [3]tile.Tile
		target       seat.Seat
		wantTaken    *tile.Tile
		wantConsumed []tile.Tile
		wantTarget   *seat.Seat
		wantErr      bool
	}{
		{
			name:         "valid tiles: 1m",
			taken:        tile.MustTileFromCode("1m"),
			consumed:     [3]tile.Tile{tile.MustTileFromCode("1m"), tile.MustTileFromCode("1m"), tile.MustTileFromCode("1m")},
			target:       seat.MustSeat(0),
			wantTaken:    new(tile.MustTileFromCode("1m")),
			wantConsumed: []tile.Tile{tile.MustTileFromCode("1m"), tile.MustTileFromCode("1m"), tile.MustTileFromCode("1m")},
			wantTarget:   new(seat.MustSeat(0)),
			wantErr:      false,
		},
		{
			name:         "valid tiles: C",
			taken:        tile.MustTileFromCode("C"),
			consumed:     [3]tile.Tile{tile.MustTileFromCode("C"), tile.MustTileFromCode("C"), tile.MustTileFromCode("C")},
			target:       seat.MustSeat(0),
			wantTaken:    new(tile.MustTileFromCode("C")),
			wantConsumed: []tile.Tile{tile.MustTileFromCode("C"), tile.MustTileFromCode("C"), tile.MustTileFromCode("C")},
			wantTarget:   new(seat.MustSeat(0)),
			wantErr:      false,
		},
		{
			name:         "valid tiles: 5sr-5s5s5s",
			taken:        tile.MustTileFromCode("5sr"),
			consumed:     [3]tile.Tile{tile.MustTileFromCode("5s"), tile.MustTileFromCode("5s"), tile.MustTileFromCode("5s")},
			target:       seat.MustSeat(0),
			wantTaken:    new(tile.MustTileFromCode("5sr")),
			wantConsumed: []tile.Tile{tile.MustTileFromCode("5s"), tile.MustTileFromCode("5s"), tile.MustTileFromCode("5s")},
			wantTarget:   new(seat.MustSeat(0)),
			wantErr:      false,
		},
		{
			name:         "valid tiles: 5p-5p5p5pr",
			taken:        tile.MustTileFromCode("5p"),
			consumed:     [3]tile.Tile{tile.MustTileFromCode("5p"), tile.MustTileFromCode("5p"), tile.MustTileFromCode("5pr")},
			target:       seat.MustSeat(0),
			wantTaken:    new(tile.MustTileFromCode("5p")),
			wantConsumed: []tile.Tile{tile.MustTileFromCode("5p"), tile.MustTileFromCode("5p"), tile.MustTileFromCode("5pr")},
			wantTarget:   new(seat.MustSeat(0)),
			wantErr:      false,
		},
		{
			name:         "invalid tiles: 5pr-5p5p5pr",
			taken:        tile.MustTileFromCode("5pr"),
			consumed:     [3]tile.Tile{tile.MustTileFromCode("5p"), tile.MustTileFromCode("5p"), tile.MustTileFromCode("5pr")},
			target:       seat.MustSeat(0),
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   new(seat.MustSeat(0)),
			wantErr:      true,
		},
		{
			name:         "invalid tiles: 5p-5p5pr5pr",
			taken:        tile.MustTileFromCode("5p"),
			consumed:     [3]tile.Tile{tile.MustTileFromCode("5p"), tile.MustTileFromCode("5pr"), tile.MustTileFromCode("5pr")},
			target:       seat.MustSeat(0),
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   new(seat.MustSeat(0)),
			wantErr:      true,
		},
		{
			name:         "invalid tiles: 5p-5p5p5p",
			taken:        tile.MustTileFromCode("5p"),
			consumed:     [3]tile.Tile{tile.MustTileFromCode("5p"), tile.MustTileFromCode("5p"), tile.MustTileFromCode("5p")},
			target:       seat.MustSeat(2),
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   new(seat.MustSeat(2)),
			wantErr:      true,
		},
		{
			name:         "invalid tiles: ?",
			taken:        tile.MustTileFromCode("?"),
			consumed:     [3]tile.Tile{tile.MustTileFromCode("?"), tile.MustTileFromCode("?"), tile.MustTileFromCode("?")},
			target:       seat.MustSeat(0),
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   new(seat.MustSeat(0)),
			wantErr:      true,
		},
		{
			name:         "taken and consumed do not match",
			taken:        tile.MustTileFromCode("1m"),
			consumed:     [3]tile.Tile{tile.MustTileFromCode("2m"), tile.MustTileFromCode("2m"), tile.MustTileFromCode("2m")},
			target:       seat.MustSeat(0),
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   new(seat.MustSeat(0)),
			wantErr:      true,
		},
		{
			name:         "consumed tiles do not match",
			taken:        tile.MustTileFromCode("1m"),
			consumed:     [3]tile.Tile{tile.MustTileFromCode("1m"), tile.MustTileFromCode("1m"), tile.MustTileFromCode("2m")},
			target:       seat.MustSeat(0),
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   new(seat.MustSeat(0)),
			wantErr:      true,
		},
		{
			name:         "sort tiles: 5m-5mr5m5m to 5m-5m5m5mr",
			taken:        tile.MustTileFromCode("5m"),
			consumed:     [3]tile.Tile{tile.MustTileFromCode("5mr"), tile.MustTileFromCode("5m"), tile.MustTileFromCode("5m")},
			target:       seat.MustSeat(0),
			wantTaken:    new(tile.MustTileFromCode("5m")),
			wantConsumed: []tile.Tile{tile.MustTileFromCode("5m"), tile.MustTileFromCode("5m"), tile.MustTileFromCode("5mr")},
			wantTarget:   new(seat.MustSeat(0)),
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := meld.NewCalledKan(tt.taken, tt.consumed, tt.target)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewCalledKan() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewCalledKan() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got.Taken(), tt.wantTaken) {
				t.Errorf("NewCalledKan().Taken() = %v, want %v", got.Taken(), tt.wantTaken)
			}
			if !reflect.DeepEqual(got.Consumed(), tt.wantConsumed) {
				t.Errorf("NewCalledKan().Consumed() = %v, want %v", got.Consumed(), tt.wantConsumed)
			}
			if !reflect.DeepEqual(got.Target(), tt.wantTarget) {
				t.Errorf("NewCalledKan().Target() = %v, want %v", got.Target(), tt.wantTarget)
			}
		})
	}
}

func TestCalledKan_ToTiles(t *testing.T) {
	tests := []struct {
		name     string
		taken    tile.Tile
		consumed [3]tile.Tile
		target   seat.Seat
		want     []tile.Tile
	}{
		{
			name:     "1m-1m1m1m",
			taken:    tile.MustTileFromCode("1m"),
			consumed: [3]tile.Tile{tile.MustTileFromCode("1m"), tile.MustTileFromCode("1m"), tile.MustTileFromCode("1m")},
			target:   seat.MustSeat(0),
			want:     []tile.Tile{tile.MustTileFromCode("1m"), tile.MustTileFromCode("1m"), tile.MustTileFromCode("1m"), tile.MustTileFromCode("1m")},
		},
		{
			name:     "sort tiles: 5m-5mr5m5m to 5m5m5m5mr",
			taken:    tile.MustTileFromCode("5m"),
			consumed: [3]tile.Tile{tile.MustTileFromCode("5mr"), tile.MustTileFromCode("5m"), tile.MustTileFromCode("5m")},
			target:   seat.MustSeat(2),
			want:     []tile.Tile{tile.MustTileFromCode("5m"), tile.MustTileFromCode("5m"), tile.MustTileFromCode("5m"), tile.MustTileFromCode("5mr")},
		},
		{
			name:     "sort tiles: 5mr-5m5m5m to 5m5m5m5mr",
			taken:    tile.MustTileFromCode("5mr"),
			consumed: [3]tile.Tile{tile.MustTileFromCode("5m"), tile.MustTileFromCode("5m"), tile.MustTileFromCode("5m")},
			target:   seat.MustSeat(2),
			want:     []tile.Tile{tile.MustTileFromCode("5m"), tile.MustTileFromCode("5m"), tile.MustTileFromCode("5m"), tile.MustTileFromCode("5mr")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, err := meld.NewCalledKan(tt.taken, tt.consumed, tt.target)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got := k.ToTiles()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToTiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalledKan_String(t *testing.T) {
	tests := []struct {
		name     string
		taken    tile.Tile
		consumed [3]tile.Tile
		target   seat.Seat
		want     string
	}{
		{
			name:     "1m-1m1m1m from 1",
			taken:    tile.MustTileFromCode("1m"),
			consumed: [3]tile.Tile{tile.MustTileFromCode("1m"), tile.MustTileFromCode("1m"), tile.MustTileFromCode("1m")},
			target:   seat.MustSeat(1),
			want:     "[1m(1)/1m 1m 1m]",
		},
		{
			name:     "5sr-5s5s5s from 3",
			taken:    tile.MustTileFromCode("5sr"),
			consumed: [3]tile.Tile{tile.MustTileFromCode("5s"), tile.MustTileFromCode("5s"), tile.MustTileFromCode("5s")},
			target:   seat.MustSeat(3),
			want:     "[5sr(3)/5s 5s 5s]",
		},
		{
			name:     "5s-5s5s5sr from 3",
			taken:    tile.MustTileFromCode("5s"),
			consumed: [3]tile.Tile{tile.MustTileFromCode("5s"), tile.MustTileFromCode("5s"), tile.MustTileFromCode("5sr")},
			target:   seat.MustSeat(3),
			want:     "[5s(3)/5s 5s 5sr]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, err := meld.NewCalledKan(tt.taken, tt.consumed, tt.target)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got := k.String()
			if got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
