package round_test

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
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
			seat.MustSeat(0),
			seat.MustSeat(0),
			tile.Tiles{tile.MustTileFromCode("1m")},
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
			seat.MustSeat(0),
			seat.MustSeat(0),
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
			doraIndicators: tile.Tiles{tile.MustTileFromCode("1m")},
			want:           tile.Tiles{tile.MustTileFromCode("2m")},
		},
		{
			name:           "double dora",
			doraIndicators: tile.Tiles{tile.MustTileFromCode("5mr"), tile.MustTileFromCode("S")},
			want:           tile.Tiles{tile.MustTileFromCode("6m"), tile.MustTileFromCode("W")},
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

func TestState_DoraIndicators_ReturnsCopy(t *testing.T) {
	players := [4]player.Player{}
	s := round.NewStateForTest(
		wind.East,
		1,
		0,
		0,
		[4]int{25000, 25000, 25000, 25000},
		seat.MustSeat(0),
		seat.MustSeat(0),
		tile.Tiles{tile.MustTileFromCode("1m")},
		round.NumInitWall,
		players,
	)

	got := s.DoraIndicators()
	got[0] = tile.MustTileFromCode("9m")

	if want := tile.MustTileFromCode("1m"); s.DoraIndicators()[0] != want {
		t.Errorf("DoraIndicators() exposed internal slice; got %v, want %v", s.DoraIndicators()[0], want)
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
			seat.MustSeat(0),
			seat.MustSeat(0),
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
	newInitStateForTest := func(roundNumber int, dealer seat.Seat) round.State {
		players := [4]player.Player{}
		return round.NewStateForTest(
			wind.East,
			roundNumber,
			0,
			0,
			[4]int{25000, 25000, 25000, 25000},
			dealer,
			seat.MustSeat(0),
			tile.Tiles{tile.MustTileFromCode("1m")},
			round.NumInitWall,
			players,
		)
	}

	tests := []struct {
		name          string
		currentNumber int
		currentDealer seat.Seat
		playerID      seat.Seat
		want          wind.Wind
	}{
		{
			name:          "E1 p0 -> E",
			currentNumber: 1,
			currentDealer: seat.MustSeat(0),
			playerID:      seat.MustSeat(0),
			want:          wind.East,
		},
		{
			name:          "E1 p1 -> S",
			currentNumber: 1,
			currentDealer: seat.MustSeat(0),
			playerID:      seat.MustSeat(1),
			want:          wind.South,
		},
		{
			name:          "E1 p2 -> W",
			currentNumber: 1,
			currentDealer: seat.MustSeat(0),
			playerID:      seat.MustSeat(2),
			want:          wind.West,
		},
		{
			name:          "E1 p3 -> N",
			currentNumber: 1,
			currentDealer: seat.MustSeat(0),
			playerID:      seat.MustSeat(3),
			want:          wind.North,
		},
		{
			name:          "E2 p0 -> N",
			currentNumber: 2,
			currentDealer: seat.MustSeat(1),
			playerID:      seat.MustSeat(0),
			want:          wind.North,
		},
		{
			name:          "E2 p1 -> E",
			currentNumber: 2,
			currentDealer: seat.MustSeat(1),
			playerID:      seat.MustSeat(1),
			want:          wind.East,
		},
		{
			name:          "E2 p2 -> S",
			currentNumber: 2,
			currentDealer: seat.MustSeat(1),
			playerID:      seat.MustSeat(2),
			want:          wind.South,
		},
		{
			name:          "E2 p3 -> W",
			currentNumber: 2,
			currentDealer: seat.MustSeat(1),
			playerID:      seat.MustSeat(3),
			want:          wind.West,
		},
		{
			name:          "E3 p0 -> W",
			currentNumber: 3,
			currentDealer: seat.MustSeat(2),
			playerID:      seat.MustSeat(0),
			want:          wind.West,
		},
		{
			name:          "E3 p1 -> N",
			currentNumber: 3,
			currentDealer: seat.MustSeat(2),
			playerID:      seat.MustSeat(1),
			want:          wind.North,
		},
		{
			name:          "E3 p2 -> E",
			currentNumber: 3,
			currentDealer: seat.MustSeat(2),
			playerID:      seat.MustSeat(2),
			want:          wind.East,
		},
		{
			name:          "E3 p3 -> S",
			currentNumber: 3,
			currentDealer: seat.MustSeat(2),
			playerID:      seat.MustSeat(3),
			want:          wind.South,
		},
		{
			name:          "E4 p0 -> S",
			currentNumber: 4,
			currentDealer: seat.MustSeat(3),
			playerID:      seat.MustSeat(0),
			want:          wind.South,
		},
		{
			name:          "E4 p1 -> W",
			currentNumber: 4,
			currentDealer: seat.MustSeat(3),
			playerID:      seat.MustSeat(1),
			want:          wind.West,
		},
		{
			name:          "E4 p2 -> N",
			currentNumber: 4,
			currentDealer: seat.MustSeat(3),
			playerID:      seat.MustSeat(2),
			want:          wind.North,
		},
		{
			name:          "E4 p3 -> E",
			currentNumber: 4,
			currentDealer: seat.MustSeat(3),
			playerID:      seat.MustSeat(3),
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

func TestState_VisibleTiles(t *testing.T) {
	newInitStateForTest := func(players [4]player.Player) round.State {
		return round.NewStateForTest(
			wind.East,
			1,
			0,
			0,
			[4]int{25000, 25000, 25000, 25000},
			seat.MustSeat(0),
			seat.MustSeat(0),
			tile.Tiles{tile.MustTileFromCode("5mr")},
			round.NumInitWall,
			players,
		)
	}

	player0, _ := player.NewVisiblePlayer(
		[13]tile.Tile{
			tile.MustTileFromCode("1m"), tile.MustTileFromCode("2m"), tile.MustTileFromCode("3m"), tile.MustTileFromCode("4m"),
			tile.MustTileFromCode("1p"), tile.MustTileFromCode("2p"), tile.MustTileFromCode("3p"), tile.MustTileFromCode("4p"),
			tile.MustTileFromCode("1s"), tile.MustTileFromCode("2s"), tile.MustTileFromCode("3s"), tile.MustTileFromCode("4s"),
			tile.MustTileFromCode("E"),
		},
	)
	player1 := player.NewInvisiblePlayer()
	player2 := player.NewInvisiblePlayer()
	player3 := player.NewInvisiblePlayer()

	player0.Draw(tile.MustTileFromCode("5m"))
	player0.Discard(tile.MustTileFromCode("1m"), false)
	player1.Chii(*meld.MustChii(
		tile.MustTileFromCode("1m"),
		[2]tile.Tile{tile.MustTileFromCode("2m"), tile.MustTileFromCode("3m")},
		seat.MustSeat(0),
	))
	player0.TakeFromRiver(tile.MustTileFromCode("1m"))
	player1.Discard(tile.MustTileFromCode("5p"), false)
	player0.Draw(tile.MustTileFromCode("6m"))

	tests := []struct {
		name     string
		players  [4]player.Player
		playerID seat.Seat
		want     tile.Tiles
	}{
		{
			name:     "player0",
			players:  [4]player.Player{player0, player1, player2, player3},
			playerID: seat.MustSeat(0),
			want: []tile.Tile{
				tile.MustTileFromCode("5p"),
				tile.MustTileFromCode("1m"), tile.MustTileFromCode("2m"), tile.MustTileFromCode("3m"),
				tile.MustTileFromCode("5mr"),
				tile.MustTileFromCode("2m"), tile.MustTileFromCode("3m"), tile.MustTileFromCode("4m"), tile.MustTileFromCode("5m"),
				tile.MustTileFromCode("1p"), tile.MustTileFromCode("2p"), tile.MustTileFromCode("3p"), tile.MustTileFromCode("4p"),
				tile.MustTileFromCode("1s"), tile.MustTileFromCode("2s"), tile.MustTileFromCode("3s"), tile.MustTileFromCode("4s"),
				tile.MustTileFromCode("E"),
				tile.MustTileFromCode("6m"),
			},
		},
		{
			name:     "player1",
			players:  [4]player.Player{player0, player1, player2, player3},
			playerID: seat.MustSeat(1),
			want: []tile.Tile{
				tile.MustTileFromCode("5p"),
				tile.MustTileFromCode("1m"), tile.MustTileFromCode("2m"), tile.MustTileFromCode("3m"),
				tile.MustTileFromCode("5mr"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newInitStateForTest(tt.players)
			got := s.VisibleTiles(tt.playerID)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VisibleTiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestState_SafeTiles(t *testing.T) {
	newInitStateForTest := func(players [4]player.Player) round.State {
		return round.NewStateForTest(
			wind.East,
			1,
			0,
			0,
			[4]int{25000, 25000, 25000, 25000},
			seat.MustSeat(0),
			seat.MustSeat(0),
			tile.Tiles{tile.MustTileFromCode("1m")},
			round.NumInitWall,
			players,
		)
	}

	type testCase struct {
		name     string
		players  [4]player.Player
		playerID seat.Seat
		want     tile.Tiles
	}
	var tests []testCase

	tests = append(
		tests,
		testCase{
			name: "initial state",
			players: [4]player.Player{
				player.NewInvisiblePlayer(),
				player.NewInvisiblePlayer(),
				player.NewInvisiblePlayer(),
				player.NewInvisiblePlayer(),
			},
			playerID: seat.MustSeat(0),
			want:     nil,
		},
	)

	player1 := player.NewInvisiblePlayer()
	player1.AddExtraSafeTiles(tile.MustTileFromCode("1p"))

	tests = append(
		tests,
		testCase{
			name: "only extra safe tiles",
			players: [4]player.Player{
				player.NewInvisiblePlayer(),
				player1,
				player.NewInvisiblePlayer(),
				player.NewInvisiblePlayer(),
			},
			playerID: seat.MustSeat(1),
			want:     tile.Tiles{tile.MustTileFromCode("1p")},
		},
	)
	tests = append(
		tests,
		testCase{
			name: "other player",
			players: [4]player.Player{
				player.NewInvisiblePlayer(),
				player1,
				player.NewInvisiblePlayer(),
				player.NewInvisiblePlayer(),
			},
			playerID: seat.MustSeat(2),
			want:     nil,
		},
	)

	player0 := player.NewInvisiblePlayer()
	player0.Draw(tile.MustTileFromCode("1p"))
	player0.Discard(tile.MustTileFromCode("1p"), false)
	player0.AddExtraSafeTiles(tile.MustTileFromCode("1m"))
	tests = append(
		tests,
		testCase{
			name: "discarded tiles and extra safe tiles",
			players: [4]player.Player{
				player0,
				player.NewInvisiblePlayer(),
				player.NewInvisiblePlayer(),
				player.NewInvisiblePlayer(),
			},
			playerID: seat.MustSeat(0),
			want:     tile.Tiles{tile.MustTileFromCode("1p"), tile.MustTileFromCode("1m")},
		},
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newInitStateForTest(tt.players)
			got := s.SafeTiles(tt.playerID)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SafeTiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
