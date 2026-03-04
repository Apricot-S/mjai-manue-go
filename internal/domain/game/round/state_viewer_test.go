package round_test

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/id"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

func TestState_NextRound(t *testing.T) {
	newInitStateForTest := func(roundWind wind.Wind, roundNumber int) round.State {
		players := [4]player.Player{}
		return round.NewStateForTest(
			roundWind,
			roundNumber,
			0,
			0,
			[4]int{25000, 25000, 25000, 25000},
			*id.MustID(0),
			*id.MustID(0),
			tile.Tiles{*tile.MustTileFromCode("1m")},
			round.NumInitWall,
			players,
		)
	}

	tests := []struct {
		name          string
		currentWind   wind.Wind
		currentNumber int
		wantWind      wind.Wind
		wantNumber    int
	}{
		{
			name:          "E1 -> E2",
			currentWind:   wind.East,
			currentNumber: 1,
			wantWind:      wind.East,
			wantNumber:    2,
		},
		{
			name:          "E4 -> S1",
			currentWind:   wind.East,
			currentNumber: 4,
			wantWind:      wind.South,
			wantNumber:    1,
		},
		{
			name:          "S2 -> S3",
			currentWind:   wind.South,
			currentNumber: 2,
			wantWind:      wind.South,
			wantNumber:    3,
		},
		{
			name:          "S4 -> W1",
			currentWind:   wind.South,
			currentNumber: 4,
			wantWind:      wind.West,
			wantNumber:    1,
		},
		{
			name:          "W3 -> W4",
			currentWind:   wind.West,
			currentNumber: 3,
			wantWind:      wind.West,
			wantNumber:    4,
		},
		{
			name:          "W4 -> N1",
			currentWind:   wind.West,
			currentNumber: 4,
			wantWind:      wind.North,
			wantNumber:    1,
		},
		{
			name:          "N4 -> E1",
			currentWind:   wind.North,
			currentNumber: 4,
			wantWind:      wind.East,
			wantNumber:    1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newInitStateForTest(tt.currentWind, tt.currentNumber)
			gotWind, gotNumber := s.NextRound()
			if gotWind != tt.wantWind {
				t.Errorf("NextRound() = %v, want %v", gotWind, tt.wantWind)
			}
			if gotNumber != tt.wantNumber {
				t.Errorf("NextRound() = %v, want %v", gotNumber, tt.wantNumber)
			}
		})
	}
}

func TestState_Doras(t *testing.T) {
	newInitStateForTest := func(doraIndicators tile.Tiles) round.State {
		players := [4]player.Player{}
		return round.NewStateForTest(
			wind.East,
			1,
			0,
			0,
			[4]int{25000, 25000, 25000, 25000},
			*id.MustID(0),
			*id.MustID(0),
			doraIndicators,
			round.NumInitWall,
			players,
		)
	}

	tests := []struct {
		name           string
		doraIndicators tile.Tiles
		want           tile.Tiles
	}{
		{
			name:           "empty",
			doraIndicators: tile.Tiles{},
			want:           tile.Tiles{},
		},
		{
			name:           "single dora",
			doraIndicators: tile.Tiles{*tile.MustTileFromCode("1m")},
			want:           tile.Tiles{*tile.MustTileFromCode("2m")},
		},
		{
			name:           "double dora",
			doraIndicators: tile.Tiles{*tile.MustTileFromCode("5mr"), *tile.MustTileFromCode("S")},
			want:           tile.Tiles{*tile.MustTileFromCode("6m"), *tile.MustTileFromCode("W")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newInitStateForTest(tt.doraIndicators)
			got := s.Doras()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Doras() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestState_Turn(t *testing.T) {
	newInitStateForTest := func(numLeftTiles int) round.State {
		players := [4]player.Player{}
		return round.NewStateForTest(
			wind.East,
			1,
			0,
			0,
			[4]int{25000, 25000, 25000, 25000},
			*id.MustID(0),
			*id.MustID(0),
			nil,
			numLeftTiles,
			players,
		)
	}

	tests := []struct {
		name         string
		numLeftTiles int
		want         float64
	}{
		{
			name:         "before first draw",
			numLeftTiles: round.NumInitWall,
			want:         0.0,
		},
		{
			name:         "after first draw",
			numLeftTiles: round.NumInitWall - 1,
			want:         0.25,
		},
		{
			name:         "after last draw",
			numLeftTiles: 0,
			want:         17.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newInitStateForTest(tt.numLeftTiles)
			got := s.Turn()
			if got != tt.want {
				t.Errorf("Turn() = %v, want %v", got, tt.want)
			}
		})
	}
}
