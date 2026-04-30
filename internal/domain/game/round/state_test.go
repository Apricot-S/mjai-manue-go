package round

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

func NewStateForTest(
	roundWind wind.Wind,
	roundNumber int,
	honba int,
	riichiDeposit int,
	scores [common.NumPlayers]int,
	dealer seat.Seat,
	startingDealer seat.Seat,
	doraIndicators tile.Tiles,
	numLeftTiles int,
	players [common.NumPlayers]player.Player,
) State {
	return State{
		roundWind,
		roundNumber,
		honba,
		riichiDeposit,
		scores,
		dealer,
		startingDealer,
		doraIndicators,
		numLeftTiles,
		players,
	}
}

func newValidHands() [common.NumPlayers][common.InitHandSize]tile.Tile {
	handCodes := [common.InitHandSize]string{"1m", "1p", "1s", "2m", "2p", "2s", "3m", "3p", "3s", "4m", "4p", "4s", "5m"}
	var hands [common.NumPlayers][common.InitHandSize]tile.Tile
	for player := range common.NumPlayers {
		for i, code := range handCodes {
			hands[player][i] = *tile.MustTileFromCode(code)
		}
	}
	return hands
}

func TestNewState(t *testing.T) {
	validDealer := *seat.MustSeat(1)
	validDora := *tile.MustTileFromCode("1m")
	validHands := newValidHands()
	validScores := &[common.NumPlayers]int{25000, 25000, 25000, 25000}

	ev := event.NewStartRound(
		wind.South,
		2,
		1,
		2,
		validDealer,
		validDora,
		validScores,
		validHands,
	)

	s, err := NewState(ev, *validScores)
	if err != nil {
		t.Fatalf("NewState() failed: %v", err)
	}

	if got := s.RoundWind(); got != wind.South {
		t.Fatalf("RoundWind() = %v, want %v", got, wind.South)
	}
	if got := s.RoundNumber(); got != 2 {
		t.Fatalf("RoundNumber() = %d, want %d", got, 2)
	}
	if got := s.Honba(); got != 1 {
		t.Fatalf("Honba() = %d, want %d", got, 1)
	}
	if got := s.RiichiDeposit(); got != 2 {
		t.Fatalf("RiichiDeposit() = %d, want %d", got, 2)
	}
	if got := s.Dealer().Index(); got != validDealer.Index() {
		t.Fatalf("Dealer() = %d, want %d", got, validDealer.Index())
	}
	if got := s.StartingDealer().Index(); got != 0 {
		t.Fatalf("StartingDealer() = %d, want %d", got, 0)
	}
	if got := s.DoraIndicators(); len(got) != 1 || got[0].ID() != validDora.ID() {
		t.Fatalf("DoraIndicators() = %v, want [%v]", got, validDora)
	}
	if got := s.NumLeftTiles(); got != NumInitWall {
		t.Fatalf("NumLeftTiles() = %d, want %d", got, NumInitWall)
	}
	if got := s.Scores(); got != *validScores {
		t.Fatalf("Scores() = %v, want %v", got, *validScores)
	}

	for i := range common.NumPlayers {
		playerSeat := *seat.MustSeat(i)
		handTiles := s.Player(playerSeat).HandTiles()
		if len(handTiles) != common.InitHandSize {
			t.Fatalf("player %d hand size = %d, want %d", i, len(handTiles), common.InitHandSize)
		}

		expectedHand, err := hand.NewVisibleHand(validHands[i][:])
		if err != nil {
			t.Fatalf("NewVisibleHand() failed for expected hand: %v", err)
		}

		actualHand, ok := s.Player(playerSeat).Hand()
		if !ok {
			t.Fatalf("player %d should be visible", i)
		}

		if !reflect.DeepEqual(actualHand.ToTiles(), expectedHand.ToTiles()) {
			t.Fatalf("player %d hand = %v, want %v", i, actualHand.ToTiles(), expectedHand.ToTiles())
		}
	}
}

func TestNewStateRejectsInvalidStartRound(t *testing.T) {
	validDealer := *seat.MustSeat(1)
	validDora := *tile.MustTileFromCode("1m")
	validHands := newValidHands()
	validScores := &[common.NumPlayers]int{25000, 25000, 25000, 25000}

	tests := []struct {
		name          string
		roundWind     wind.Wind
		roundNumber   int
		honba         int
		riichiDeposit int
		doraIndicator tile.Tile
	}{
		{
			name:          "invalid round wind",
			roundWind:     wind.Wind(0),
			roundNumber:   1,
			honba:         0,
			riichiDeposit: 0,
			doraIndicator: validDora,
		},
		{
			name:          "invalid round number 0",
			roundWind:     wind.East,
			roundNumber:   0,
			honba:         0,
			riichiDeposit: 0,
			doraIndicator: validDora,
		},
		{
			name:          "invalid round number 5",
			roundWind:     wind.East,
			roundNumber:   5,
			honba:         0,
			riichiDeposit: 0,
			doraIndicator: validDora,
		},
		{
			name:          "negative honba",
			roundWind:     wind.East,
			roundNumber:   1,
			honba:         -1,
			riichiDeposit: 0,
			doraIndicator: validDora,
		},
		{
			name:          "negative riichi deposit",
			roundWind:     wind.East,
			roundNumber:   1,
			honba:         0,
			riichiDeposit: -1,
			doraIndicator: validDora,
		},
		{
			name:          "unknown dora indicator",
			roundWind:     wind.East,
			roundNumber:   1,
			honba:         0,
			riichiDeposit: 0,
			doraIndicator: *tile.MustTileFromCode("?"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev := event.NewStartRound(
				tt.roundWind,
				tt.roundNumber,
				tt.honba,
				tt.riichiDeposit,
				validDealer,
				tt.doraIndicator,
				validScores,
				validHands,
			)

			if _, err := NewState(ev, *validScores); err == nil {
				t.Fatal("NewState() succeeded unexpectedly")
			}
		})
	}
}

func TestNewStateWithNilScores(t *testing.T) {
	validDealer := *seat.MustSeat(1)
	validDora := *tile.MustTileFromCode("1m")
	validHands := newValidHands()

	ev := event.NewStartRound(
		wind.East,
		1,
		0,
		0,
		validDealer,
		validDora,
		nil,
		validHands,
	)

	scoresBefore := [common.NumPlayers]int{10000, 20000, 30000, 40000}
	s, err := NewState(ev, scoresBefore)
	if err != nil {
		t.Fatalf("NewState() failed: %v", err)
	}

	if got := s.Scores(); got != scoresBefore {
		t.Fatalf("Scores() = %v, want %v", got, scoresBefore)
	}
}

func TestNewStateFallsBackToInvisiblePlayer(t *testing.T) {
	validDealer := *seat.MustSeat(1)
	validDora := *tile.MustTileFromCode("1m")
	unknownHands := newValidHands()
	for i := range unknownHands[0] {
		unknownHands[0][i] = *tile.MustTileFromCode("?")
	}

	ev := event.NewStartRound(
		wind.East,
		1,
		0,
		0,
		validDealer,
		validDora,
		&[common.NumPlayers]int{25000, 25000, 25000, 25000},
		unknownHands,
	)

	s, err := NewState(ev, [common.NumPlayers]int{25000, 25000, 25000, 25000})
	if err != nil {
		t.Fatalf("NewState() failed: %v", err)
	}

	firstSeat := *seat.MustSeat(0)
	if _, ok := s.Player(firstSeat).Hand(); ok {
		t.Fatalf("expected seat 0 to be invisible when visible hand initialization fails")
	}

	for i := 1; i < common.NumPlayers; i++ {
		playerSeat := *seat.MustSeat(i)
		if _, ok := s.Player(playerSeat).Hand(); !ok {
			t.Fatalf("expected seat %d to be visible", i)
		}
	}
}

func TestNewStateErrorsOnInvalidVisibleHand(t *testing.T) {
	validDealer := *seat.MustSeat(1)
	validDora := *tile.MustTileFromCode("1m")
	invalidHands := newValidHands()
	for i := range 5 {
		invalidHands[0][i] = *tile.MustTileFromCode("1m")
	}

	ev := event.NewStartRound(
		wind.East,
		1,
		0,
		0,
		validDealer,
		validDora,
		&[common.NumPlayers]int{25000, 25000, 25000, 25000},
		invalidHands,
	)

	_, err := NewState(ev, [common.NumPlayers]int{25000, 25000, 25000, 25000})
	if err == nil {
		t.Fatal("NewState() succeeded unexpectedly")
	}
}
