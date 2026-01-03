package meld_test

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/tile"
)

func TestNewPon(t *testing.T) {
	tests := []struct {
		name         string
		taken        tile.Tile
		consumed     [2]tile.Tile
		target       player.PlayerID
		wantTaken    *tile.Tile
		wantConsumed []tile.Tile
		wantTarget   *player.PlayerID
		wantErr      bool
	}{
		{
			name:         "valid tiles: 1m",
			taken:        *tile.MustTileFromCode("1m"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
			target:       *player.MustPlayerID(0),
			wantTaken:    tile.MustTileFromCode("1m"),
			wantConsumed: []tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
			wantTarget:   player.MustPlayerID(0),
			wantErr:      false,
		},
		{
			name:         "valid tiles: C",
			taken:        *tile.MustTileFromCode("C"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("C"), *tile.MustTileFromCode("C")},
			target:       *player.MustPlayerID(0),
			wantTaken:    tile.MustTileFromCode("C"),
			wantConsumed: []tile.Tile{*tile.MustTileFromCode("C"), *tile.MustTileFromCode("C")},
			wantTarget:   player.MustPlayerID(0),
			wantErr:      false,
		},
		{
			name:         "valid tiles: 5s-5s5s",
			taken:        *tile.MustTileFromCode("5s"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("5s"), *tile.MustTileFromCode("5s")},
			target:       *player.MustPlayerID(0),
			wantTaken:    tile.MustTileFromCode("5s"),
			wantConsumed: []tile.Tile{*tile.MustTileFromCode("5s"), *tile.MustTileFromCode("5s")},
			wantTarget:   player.MustPlayerID(0),
			wantErr:      false,
		},
		{
			name:         "valid tiles: 5sr-5s5s",
			taken:        *tile.MustTileFromCode("5sr"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("5s"), *tile.MustTileFromCode("5s")},
			target:       *player.MustPlayerID(0),
			wantTaken:    tile.MustTileFromCode("5sr"),
			wantConsumed: []tile.Tile{*tile.MustTileFromCode("5s"), *tile.MustTileFromCode("5s")},
			wantTarget:   player.MustPlayerID(0),
			wantErr:      false,
		},
		{
			name:         "valid tiles: 5p-5p5pr",
			taken:        *tile.MustTileFromCode("5p"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("5p"), *tile.MustTileFromCode("5pr")},
			target:       *player.MustPlayerID(0),
			wantTaken:    tile.MustTileFromCode("5p"),
			wantConsumed: []tile.Tile{*tile.MustTileFromCode("5p"), *tile.MustTileFromCode("5pr")},
			wantTarget:   player.MustPlayerID(0),
			wantErr:      false,
		},
		{
			name:         "invalid tiles: 5pr-5p5pr",
			taken:        *tile.MustTileFromCode("5pr"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("5p"), *tile.MustTileFromCode("5pr")},
			target:       *player.MustPlayerID(0),
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   player.MustPlayerID(0),
			wantErr:      true,
		},
		{
			name:         "invalid tiles: 5p-5pr5pr",
			taken:        *tile.MustTileFromCode("5p"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("5pr"), *tile.MustTileFromCode("5pr")},
			target:       *player.MustPlayerID(0),
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   player.MustPlayerID(0),
			wantErr:      true,
		},
		{
			name:         "invalid tiles: ?",
			taken:        *tile.MustTileFromCode("?"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?")},
			target:       *player.MustPlayerID(0),
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   player.MustPlayerID(0),
			wantErr:      true,
		},
		{
			name:         "taken and consumed do not match",
			taken:        *tile.MustTileFromCode("1m"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("2m"), *tile.MustTileFromCode("2m")},
			target:       *player.MustPlayerID(0),
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   player.MustPlayerID(0),
			wantErr:      true,
		},
		{
			name:         "consumed tiles do not match",
			taken:        *tile.MustTileFromCode("1m"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m")},
			target:       *player.MustPlayerID(0),
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   player.MustPlayerID(0),
			wantErr:      true,
		},
		{
			name:         "sort tiles: 5m-5mr5m to 5m-5m5mr",
			taken:        *tile.MustTileFromCode("5m"),
			consumed:     [2]tile.Tile{*tile.MustTileFromCode("5mr"), *tile.MustTileFromCode("5m")},
			target:       *player.MustPlayerID(0),
			wantTaken:    tile.MustTileFromCode("5m"),
			wantConsumed: []tile.Tile{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5mr")},
			wantTarget:   player.MustPlayerID(0),
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := meld.NewPon(tt.taken, tt.consumed, tt.target)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewPon() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewPon() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got.Taken(), tt.wantTaken) {
				t.Errorf("NewPon().Taken() = %v, want %v", got.Taken(), tt.wantTaken)
			}
			if !reflect.DeepEqual(got.Consumed(), tt.wantConsumed) {
				t.Errorf("NewPon().Consumed() = %v, want %v", got.Consumed(), tt.wantConsumed)
			}
			if !reflect.DeepEqual(got.Target(), tt.wantTarget) {
				t.Errorf("NewPon().Target() = %v, want %v", got.Target(), tt.wantTarget)
			}
		})
	}
}

func TestPon_ToTiles(t *testing.T) {
	tests := []struct {
		name     string
		taken    tile.Tile
		consumed [2]tile.Tile
		target   player.PlayerID
		want     []tile.Tile
	}{
		{
			name:     "1m-1m1m",
			taken:    *tile.MustTileFromCode("1m"),
			consumed: [2]tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
			target:   *player.MustPlayerID(0),
			want:     []tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
		},
		{
			name:     "sort tiles: 5m-5mr5m to 5m5m5mr",
			taken:    *tile.MustTileFromCode("5m"),
			consumed: [2]tile.Tile{*tile.MustTileFromCode("5mr"), *tile.MustTileFromCode("5m")},
			target:   *player.MustPlayerID(2),
			want:     []tile.Tile{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5mr")},
		},
		{
			name:     "sort tiles: 5mr-5m5m to 5m5m5mr",
			taken:    *tile.MustTileFromCode("5mr"),
			consumed: [2]tile.Tile{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m")},
			target:   *player.MustPlayerID(2),
			want:     []tile.Tile{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5mr")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := meld.NewPon(tt.taken, tt.consumed, tt.target)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got := p.ToTiles()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToTiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPon_String(t *testing.T) {
	tests := []struct {
		name     string
		taken    tile.Tile
		consumed [2]tile.Tile
		target   player.PlayerID
		want     string
	}{
		{
			name:     "1m-1m1m from 1",
			taken:    *tile.MustTileFromCode("1m"),
			consumed: [2]tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
			target:   *player.MustPlayerID(1),
			want:     "[1m(1)/1m 1m]",
		},
		{
			name:     "5sr-5s5s from 3",
			taken:    *tile.MustTileFromCode("5sr"),
			consumed: [2]tile.Tile{*tile.MustTileFromCode("5s"), *tile.MustTileFromCode("5s")},
			target:   *player.MustPlayerID(3),
			want:     "[5sr(3)/5s 5s]",
		},
		{
			name:     "5s-5s5sr from 3",
			taken:    *tile.MustTileFromCode("5s"),
			consumed: [2]tile.Tile{*tile.MustTileFromCode("5s"), *tile.MustTileFromCode("5sr")},
			target:   *player.MustPlayerID(3),
			want:     "[5s(3)/5s 5sr]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := meld.NewPon(tt.taken, tt.consumed, tt.target)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got := p.String()
			if got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPon_SwapCallTiles(t *testing.T) {
	tests := []struct {
		name     string
		taken    tile.Tile
		consumed [2]tile.Tile
		target   player.PlayerID
		want     []tile.Tile
	}{
		{
			name:     "1m-1m1m",
			taken:    *tile.MustTileFromCode("1m"),
			consumed: [2]tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
			target:   *player.MustPlayerID(0),
			want:     []tile.Tile{*tile.MustTileFromCode("1m")},
		},
		{
			name:     "P-PP",
			taken:    *tile.MustTileFromCode("P"),
			consumed: [2]tile.Tile{*tile.MustTileFromCode("P"), *tile.MustTileFromCode("P")},
			target:   *player.MustPlayerID(0),
			want:     []tile.Tile{*tile.MustTileFromCode("P")},
		},
		{
			name:     "5m-5m5m",
			taken:    *tile.MustTileFromCode("5m"),
			consumed: [2]tile.Tile{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m")},
			target:   *player.MustPlayerID(2),
			want:     []tile.Tile{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5mr")},
		},
		{
			name:     "5m-5m5mr",
			taken:    *tile.MustTileFromCode("5m"),
			consumed: [2]tile.Tile{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5mr")},
			target:   *player.MustPlayerID(2),
			want:     []tile.Tile{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5mr")},
		},
		{
			name:     "5mr-5m5m",
			taken:    *tile.MustTileFromCode("5mr"),
			consumed: [2]tile.Tile{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m")},
			target:   *player.MustPlayerID(2),
			want:     []tile.Tile{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5mr")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := meld.NewPon(tt.taken, tt.consumed, tt.target)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got := p.SwapCallTiles()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SwapCallTiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
