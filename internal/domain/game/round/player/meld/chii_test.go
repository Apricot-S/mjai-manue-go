package meld_test

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestNewChii(t *testing.T) {
	tests := []struct {
		name         string
		taken        tile.Tile
		consumed     [2]tile.Tile
		target       seat.Seat
		wantTaken    *tile.Tile
		wantConsumed []tile.Tile
		wantTarget   *seat.Seat
		wantErr      bool
	}{
		{
			name:         "valid tiles: 8-79m",
			taken:        tile.MustTileFromCode("8m"),
			consumed:     [2]tile.Tile{tile.MustTileFromCode("7m"), tile.MustTileFromCode("9m")},
			target:       *seat.MustSeat(0),
			wantTaken:    new(tile.MustTileFromCode("8m")),
			wantConsumed: []tile.Tile{tile.MustTileFromCode("7m"), tile.MustTileFromCode("9m")},
			wantTarget:   seat.MustSeat(0),
			wantErr:      false,
		},
		{
			name:         "valid tiles: 5r-46m",
			taken:        tile.MustTileFromCode("5mr"),
			consumed:     [2]tile.Tile{tile.MustTileFromCode("4m"), tile.MustTileFromCode("6m")},
			target:       *seat.MustSeat(0),
			wantTaken:    new(tile.MustTileFromCode("5mr")),
			wantConsumed: []tile.Tile{tile.MustTileFromCode("4m"), tile.MustTileFromCode("6m")},
			wantTarget:   seat.MustSeat(0),
			wantErr:      false,
		},
		{
			name:         "invalid tiles: ?-??",
			taken:        tile.MustTileFromCode("?"),
			consumed:     [2]tile.Tile{tile.MustTileFromCode("?"), tile.MustTileFromCode("?")},
			target:       *seat.MustSeat(0),
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   seat.MustSeat(0),
			wantErr:      true,
		},
		{
			name:         "invalid tiles: E-SW",
			taken:        tile.MustTileFromCode("E"),
			consumed:     [2]tile.Tile{tile.MustTileFromCode("S"), tile.MustTileFromCode("W")},
			target:       *seat.MustSeat(0),
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   seat.MustSeat(0),
			wantErr:      true,
		},
		{
			name:         "invalid tiles: 8-9m1p",
			taken:        tile.MustTileFromCode("8m"),
			consumed:     [2]tile.Tile{tile.MustTileFromCode("9m"), tile.MustTileFromCode("1p")},
			target:       *seat.MustSeat(0),
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   seat.MustSeat(0),
			wantErr:      true,
		},
		{
			name:         "invalid tiles: 6-89m",
			taken:        tile.MustTileFromCode("6m"),
			consumed:     [2]tile.Tile{tile.MustTileFromCode("8m"), tile.MustTileFromCode("9m")},
			target:       *seat.MustSeat(0),
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   seat.MustSeat(0),
			wantErr:      true,
		},
		{
			name:         "invalid tiles: 6-79m",
			taken:        tile.MustTileFromCode("6m"),
			consumed:     [2]tile.Tile{tile.MustTileFromCode("7m"), tile.MustTileFromCode("9m")},
			target:       *seat.MustSeat(0),
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   seat.MustSeat(0),
			wantErr:      true,
		},
		{
			name:         "invalid tiles: 6-7m8p",
			taken:        tile.MustTileFromCode("6m"),
			consumed:     [2]tile.Tile{tile.MustTileFromCode("7m"), tile.MustTileFromCode("8p")},
			target:       *seat.MustSeat(0),
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   seat.MustSeat(0),
			wantErr:      true,
		},
		{
			name:         "sort tiles: 8-97s to 8-79s",
			taken:        tile.MustTileFromCode("8s"),
			consumed:     [2]tile.Tile{tile.MustTileFromCode("9s"), tile.MustTileFromCode("7s")},
			target:       *seat.MustSeat(3),
			wantTaken:    new(tile.MustTileFromCode("8s")),
			wantConsumed: []tile.Tile{tile.MustTileFromCode("7s"), tile.MustTileFromCode("9s")},
			wantTarget:   seat.MustSeat(3),
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

func TestChii_ToTiles(t *testing.T) {
	tests := []struct {
		name     string
		taken    tile.Tile
		consumed [2]tile.Tile
		target   seat.Seat
		want     []tile.Tile
	}{
		{
			name:     "1-23m",
			taken:    tile.MustTileFromCode("1m"),
			consumed: [2]tile.Tile{tile.MustTileFromCode("2m"), tile.MustTileFromCode("3m")},
			target:   *seat.MustSeat(0),
			want:     []tile.Tile{tile.MustTileFromCode("1m"), tile.MustTileFromCode("2m"), tile.MustTileFromCode("3m")},
		},
		{
			name:     "5r-46p",
			taken:    tile.MustTileFromCode("5pr"),
			consumed: [2]tile.Tile{tile.MustTileFromCode("4p"), tile.MustTileFromCode("6p")},
			target:   *seat.MustSeat(0),
			want:     []tile.Tile{tile.MustTileFromCode("4p"), tile.MustTileFromCode("5pr"), tile.MustTileFromCode("6p")},
		},
		{
			name:     "6-5r4p",
			taken:    tile.MustTileFromCode("6p"),
			consumed: [2]tile.Tile{tile.MustTileFromCode("5pr"), tile.MustTileFromCode("4p")},
			target:   *seat.MustSeat(0),
			want:     []tile.Tile{tile.MustTileFromCode("4p"), tile.MustTileFromCode("5pr"), tile.MustTileFromCode("6p")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := meld.NewChii(tt.taken, tt.consumed, tt.target)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got := c.ToTiles()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToTiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChii_String(t *testing.T) {
	tests := []struct {
		name     string
		taken    tile.Tile
		consumed [2]tile.Tile
		target   seat.Seat
		want     string
	}{
		{
			name:     "1-23m from 0",
			taken:    tile.MustTileFromCode("1m"),
			consumed: [2]tile.Tile{tile.MustTileFromCode("2m"), tile.MustTileFromCode("3m")},
			target:   *seat.MustSeat(0),
			want:     "[1m(0)/2m 3m]",
		},
		{
			name:     "5r-46s from 1",
			taken:    tile.MustTileFromCode("5sr"),
			consumed: [2]tile.Tile{tile.MustTileFromCode("4s"), tile.MustTileFromCode("6s")},
			target:   *seat.MustSeat(1),
			want:     "[5sr(1)/4s 6s]",
		},
		{
			name:     "6-45pr from 2",
			taken:    tile.MustTileFromCode("6p"),
			consumed: [2]tile.Tile{tile.MustTileFromCode("4p"), tile.MustTileFromCode("5pr")},
			target:   *seat.MustSeat(2),
			want:     "[6p(2)/4p 5pr]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := meld.NewChii(tt.taken, tt.consumed, tt.target)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got := c.String()
			if got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChii_SwapCallTiles(t *testing.T) {
	tests := []struct {
		name     string
		taken    tile.Tile
		consumed [2]tile.Tile
		target   seat.Seat
		want     []tile.Tile
	}{
		{
			name:     "8-79m",
			taken:    tile.MustTileFromCode("8m"),
			consumed: [2]tile.Tile{tile.MustTileFromCode("7m"), tile.MustTileFromCode("9m")},
			target:   *seat.MustSeat(0),
			want:     []tile.Tile{tile.MustTileFromCode("8m")},
		},
		{
			name:     "5-46m",
			taken:    tile.MustTileFromCode("5m"),
			consumed: [2]tile.Tile{tile.MustTileFromCode("4m"), tile.MustTileFromCode("6m")},
			target:   *seat.MustSeat(0),
			want:     []tile.Tile{tile.MustTileFromCode("5m"), tile.MustTileFromCode("5mr")},
		},
		{
			name:     "5r-46m",
			taken:    tile.MustTileFromCode("5mr"),
			consumed: [2]tile.Tile{tile.MustTileFromCode("4m"), tile.MustTileFromCode("6m")},
			target:   *seat.MustSeat(0),
			want:     []tile.Tile{tile.MustTileFromCode("5m"), tile.MustTileFromCode("5mr")},
		},
		{
			name:     "9-78m",
			taken:    tile.MustTileFromCode("9m"),
			consumed: [2]tile.Tile{tile.MustTileFromCode("7m"), tile.MustTileFromCode("8m")},
			target:   *seat.MustSeat(0),
			want:     []tile.Tile{tile.MustTileFromCode("6m"), tile.MustTileFromCode("9m")},
		},
		{
			name:     "5-34m",
			taken:    tile.MustTileFromCode("5m"),
			consumed: [2]tile.Tile{tile.MustTileFromCode("3m"), tile.MustTileFromCode("4m")},
			target:   *seat.MustSeat(0),
			want:     []tile.Tile{tile.MustTileFromCode("2m"), tile.MustTileFromCode("5m"), tile.MustTileFromCode("5mr")},
		},
		{
			name:     "5r-34m",
			taken:    tile.MustTileFromCode("5mr"),
			consumed: [2]tile.Tile{tile.MustTileFromCode("3m"), tile.MustTileFromCode("4m")},
			target:   *seat.MustSeat(0),
			want:     []tile.Tile{tile.MustTileFromCode("2m"), tile.MustTileFromCode("5m"), tile.MustTileFromCode("5mr")},
		},
		{
			name:     "8-67m",
			taken:    tile.MustTileFromCode("8m"),
			consumed: [2]tile.Tile{tile.MustTileFromCode("6m"), tile.MustTileFromCode("7m")},
			target:   *seat.MustSeat(0),
			want:     []tile.Tile{tile.MustTileFromCode("5m"), tile.MustTileFromCode("5mr"), tile.MustTileFromCode("8m")},
		},
		{
			name:     "3-12m",
			taken:    tile.MustTileFromCode("3m"),
			consumed: [2]tile.Tile{tile.MustTileFromCode("1m"), tile.MustTileFromCode("2m")},
			target:   *seat.MustSeat(0),
			want:     []tile.Tile{tile.MustTileFromCode("3m")},
		},
		{
			name:     "1-23m",
			taken:    tile.MustTileFromCode("1m"),
			consumed: [2]tile.Tile{tile.MustTileFromCode("2m"), tile.MustTileFromCode("3m")},
			target:   *seat.MustSeat(0),
			want:     []tile.Tile{tile.MustTileFromCode("1m"), tile.MustTileFromCode("4m")},
		},
		{
			name:     "2-34m",
			taken:    tile.MustTileFromCode("2m"),
			consumed: [2]tile.Tile{tile.MustTileFromCode("3m"), tile.MustTileFromCode("4m")},
			target:   *seat.MustSeat(0),
			want:     []tile.Tile{tile.MustTileFromCode("2m"), tile.MustTileFromCode("5m"), tile.MustTileFromCode("5mr")},
		},
		{
			name:     "5-67m",
			taken:    tile.MustTileFromCode("5m"),
			consumed: [2]tile.Tile{tile.MustTileFromCode("6m"), tile.MustTileFromCode("7m")},
			target:   *seat.MustSeat(0),
			want:     []tile.Tile{tile.MustTileFromCode("5m"), tile.MustTileFromCode("5mr"), tile.MustTileFromCode("8m")},
		},
		{
			name:     "5r-67m",
			taken:    tile.MustTileFromCode("5mr"),
			consumed: [2]tile.Tile{tile.MustTileFromCode("6m"), tile.MustTileFromCode("7m")},
			target:   *seat.MustSeat(0),
			want:     []tile.Tile{tile.MustTileFromCode("5m"), tile.MustTileFromCode("5mr"), tile.MustTileFromCode("8m")},
		},
		{
			name:     "7-89m",
			taken:    tile.MustTileFromCode("7m"),
			consumed: [2]tile.Tile{tile.MustTileFromCode("8m"), tile.MustTileFromCode("9m")},
			target:   *seat.MustSeat(0),
			want:     []tile.Tile{tile.MustTileFromCode("7m")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := meld.NewChii(tt.taken, tt.consumed, tt.target)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got := c.SwapCallTiles()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SwapCallTiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
