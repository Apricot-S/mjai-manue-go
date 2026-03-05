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

func TestState_SeatWind(t *testing.T) {
	newInitStateForTest := func(roundNumber int, dealer id.ID) round.State {
		players := [4]player.Player{}
		return round.NewStateForTest(
			wind.East,
			roundNumber,
			0,
			0,
			[4]int{25000, 25000, 25000, 25000},
			dealer,
			*id.MustID(0),
			tile.Tiles{*tile.MustTileFromCode("1m")},
			round.NumInitWall,
			players,
		)
	}

	tests := []struct {
		name          string
		currentNumber int
		currentDealer id.ID
		playerID      id.ID
		want          wind.Wind
	}{
		{
			name:          "E1 p0 -> E",
			currentNumber: 1,
			currentDealer: *id.MustID(0),
			playerID:      *id.MustID(0),
			want:          wind.East,
		},
		{
			name:          "E1 p1 -> S",
			currentNumber: 1,
			currentDealer: *id.MustID(0),
			playerID:      *id.MustID(1),
			want:          wind.South,
		},
		{
			name:          "E1 p2 -> W",
			currentNumber: 1,
			currentDealer: *id.MustID(0),
			playerID:      *id.MustID(2),
			want:          wind.West,
		},
		{
			name:          "E1 p3 -> N",
			currentNumber: 1,
			currentDealer: *id.MustID(0),
			playerID:      *id.MustID(3),
			want:          wind.North,
		},
		{
			name:          "E2 p0 -> N",
			currentNumber: 2,
			currentDealer: *id.MustID(1),
			playerID:      *id.MustID(0),
			want:          wind.North,
		},
		{
			name:          "E2 p1 -> E",
			currentNumber: 2,
			currentDealer: *id.MustID(1),
			playerID:      *id.MustID(1),
			want:          wind.East,
		},
		{
			name:          "E2 p2 -> S",
			currentNumber: 2,
			currentDealer: *id.MustID(1),
			playerID:      *id.MustID(2),
			want:          wind.South,
		},
		{
			name:          "E2 p3 -> W",
			currentNumber: 2,
			currentDealer: *id.MustID(1),
			playerID:      *id.MustID(3),
			want:          wind.West,
		},
		{
			name:          "E3 p0 -> W",
			currentNumber: 3,
			currentDealer: *id.MustID(2),
			playerID:      *id.MustID(0),
			want:          wind.West,
		},
		{
			name:          "E3 p1 -> N",
			currentNumber: 3,
			currentDealer: *id.MustID(2),
			playerID:      *id.MustID(1),
			want:          wind.North,
		},
		{
			name:          "E3 p2 -> E",
			currentNumber: 3,
			currentDealer: *id.MustID(2),
			playerID:      *id.MustID(2),
			want:          wind.East,
		},
		{
			name:          "E3 p3 -> S",
			currentNumber: 3,
			currentDealer: *id.MustID(2),
			playerID:      *id.MustID(3),
			want:          wind.South,
		},
		{
			name:          "E4 p0 -> S",
			currentNumber: 4,
			currentDealer: *id.MustID(3),
			playerID:      *id.MustID(0),
			want:          wind.South,
		},
		{
			name:          "E4 p1 -> W",
			currentNumber: 4,
			currentDealer: *id.MustID(3),
			playerID:      *id.MustID(1),
			want:          wind.West,
		},
		{
			name:          "E4 p2 -> N",
			currentNumber: 4,
			currentDealer: *id.MustID(3),
			playerID:      *id.MustID(2),
			want:          wind.North,
		},
		{
			name:          "E4 p3 -> E",
			currentNumber: 4,
			currentDealer: *id.MustID(3),
			playerID:      *id.MustID(3),
			want:          wind.East,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newInitStateForTest(tt.currentNumber, tt.currentDealer)
			got := s.SeatWind(tt.playerID)
			if got != tt.want {
				t.Errorf("SeatWind() = %v, want %v", got, tt.want)
			}
		})
	}
}
