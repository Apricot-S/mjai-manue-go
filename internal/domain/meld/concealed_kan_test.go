package meld_test

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/tile"
)

func TestNewConcealedKan(t *testing.T) {
	tests := []struct {
		name         string
		consumed     [4]tile.Tile
		wantTaken    *tile.Tile
		wantConsumed []tile.Tile
		wantTarget   int
		wantErr      bool
	}{
		{
			name:         "valid tiles: 1p",
			consumed:     [4]tile.Tile{*tile.MustTileFromCode("1p"), *tile.MustTileFromCode("1p"), *tile.MustTileFromCode("1p"), *tile.MustTileFromCode("1p")},
			wantTaken:    nil,
			wantConsumed: []tile.Tile{*tile.MustTileFromCode("1p"), *tile.MustTileFromCode("1p"), *tile.MustTileFromCode("1p"), *tile.MustTileFromCode("1p")},
			wantTarget:   -1,
			wantErr:      false,
		},
		{
			name:         "valid tiles: C",
			consumed:     [4]tile.Tile{*tile.MustTileFromCode("C"), *tile.MustTileFromCode("C"), *tile.MustTileFromCode("C"), *tile.MustTileFromCode("C")},
			wantTaken:    nil,
			wantConsumed: []tile.Tile{*tile.MustTileFromCode("C"), *tile.MustTileFromCode("C"), *tile.MustTileFromCode("C"), *tile.MustTileFromCode("C")},
			wantTarget:   -1,
			wantErr:      false,
		},
		{
			name:         "valid tiles: 5m5m5m5mr",
			consumed:     [4]tile.Tile{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5mr")},
			wantTaken:    nil,
			wantConsumed: []tile.Tile{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5mr")},
			wantTarget:   -1,
			wantErr:      false,
		},
		{
			name:         "invalid tiles: 5m5m5m5m",
			consumed:     [4]tile.Tile{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m")},
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   -1,
			wantErr:      true,
		},
		{
			name:         "invalid tiles: 5m5m5mr5mr",
			consumed:     [4]tile.Tile{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5mr"), *tile.MustTileFromCode("5mr")},
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   -1,
			wantErr:      true,
		},
		{
			name:         "invalid tiles: ?",
			consumed:     [4]tile.Tile{*tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?"), *tile.MustTileFromCode("?")},
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   -1,
			wantErr:      true,
		},
		{
			name:         "consumed tiles do not match",
			consumed:     [4]tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("2m")},
			wantTaken:    nil,
			wantConsumed: nil,
			wantTarget:   -1,
			wantErr:      true,
		},
		{
			name:         "sort tiles: 5mr5m5m5m to 5m5m5m5mr",
			consumed:     [4]tile.Tile{*tile.MustTileFromCode("5mr"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m")},
			wantTaken:    nil,
			wantConsumed: []tile.Tile{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5mr")},
			wantTarget:   -1,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := meld.NewConcealedKan(tt.consumed)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewConcealedKan() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewConcealedKan() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got.Consumed(), tt.wantConsumed) {
				t.Errorf("NewConcealedKan().Consumed() = %v, want %v", got.Consumed(), tt.wantConsumed)
			}
		})
	}
}

func TestConcealedKan_ToTiles(t *testing.T) {
	tests := []struct {
		name     string
		consumed [4]tile.Tile
		want     []tile.Tile
	}{
		{
			name:     "1m1m1m1m",
			consumed: [4]tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
			want:     []tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
		},
		{
			name:     "5m5m5m5mr",
			consumed: [4]tile.Tile{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5mr")},
			want:     []tile.Tile{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5mr")},
		},
		{
			name:     "sort tiles: 5mr5m5m5m to 5m5m5m5mr",
			consumed: [4]tile.Tile{*tile.MustTileFromCode("5mr"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m")},
			want:     []tile.Tile{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5mr")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, err := meld.NewConcealedKan(tt.consumed)
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

func TestConcealedKan_String(t *testing.T) {
	tests := []struct {
		name     string
		consumed [4]tile.Tile
		want     string
	}{
		{
			name:     "1m1m1m1m",
			consumed: [4]tile.Tile{*tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m"), *tile.MustTileFromCode("1m")},
			want:     "[# 1m 1m #]",
		},
		{
			name:     "EEEE",
			consumed: [4]tile.Tile{*tile.MustTileFromCode("E"), *tile.MustTileFromCode("E"), *tile.MustTileFromCode("E"), *tile.MustTileFromCode("E")},
			want:     "[# E E #]",
		},
		{
			name:     "5s5s5s5sr",
			consumed: [4]tile.Tile{*tile.MustTileFromCode("5s"), *tile.MustTileFromCode("5s"), *tile.MustTileFromCode("5s"), *tile.MustTileFromCode("5sr")},
			want:     "[# 5s 5sr #]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, err := meld.NewConcealedKan(tt.consumed)
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
