package round

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

func newValidHands() [common.NumPlayers][initHandSize]tile.Tile {
	handCodes := [initHandSize]string{"1m", "1p", "1s", "2m", "2p", "2s", "3m", "3p", "3s", "4m", "4p", "4s", "5m"}
	var hands [common.NumPlayers][initHandSize]tile.Tile
	for player := range common.NumPlayers {
		for i, code := range handCodes {
			hands[player][i] = *tile.MustTileFromCode(code)
		}
	}
	return hands
}

func TestApplyStartRound(t *testing.T) {
	validDealer := *seat.MustSeat(1)
	validStartingDealer := *seat.MustSeat(0)
	validDora := *tile.MustTileFromCode("1m")
	validHands := newValidHands()
	validScores := &[common.NumPlayers]int{25000, 25000, 25000, 25000}

	ev, err := event.NewStartRound(
		wind.South,
		2,
		1,
		2,
		validDealer,
		validStartingDealer,
		validDora,
		validScores,
		validHands,
	)
	if err != nil {
		t.Fatalf("NewStartRound() failed: %v", err)
	}

	var s State
	if err := s.Apply(ev); err != nil {
		t.Fatalf("State.Apply() failed: %v", err)
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
	if got := s.StartingDealer().Index(); got != validStartingDealer.Index() {
		t.Fatalf("StartingDealer() = %d, want %d", got, validStartingDealer.Index())
	}
	if got := s.DoraIndicators(); got.Len() != 1 || got[0].ID() != validDora.ID() {
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
		if len(handTiles) != initHandSize {
			t.Fatalf("player %d hand size = %d, want %d", i, len(handTiles), initHandSize)
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

func TestApplyStartRoundWithNilScores(t *testing.T) {
	validDealer := *seat.MustSeat(1)
	validStartingDealer := *seat.MustSeat(0)
	validDora := *tile.MustTileFromCode("1m")
	validHands := newValidHands()

	ev, err := event.NewStartRound(
		wind.East,
		1,
		0,
		0,
		validDealer,
		validStartingDealer,
		validDora,
		nil,
		validHands,
	)
	if err != nil {
		t.Fatalf("NewStartRound() failed: %v", err)
	}

	scoresBefore := [common.NumPlayers]int{10000, 20000, 30000, 40000}
	var s State
	s.scores = scoresBefore
	if err := s.Apply(ev); err != nil {
		t.Fatalf("State.Apply() failed: %v", err)
	}

	if got := s.Scores(); got != scoresBefore {
		t.Fatalf("Scores() = %v, want %v", got, scoresBefore)
	}
}

func TestApplyStartRoundFallsBackToInvisiblePlayer(t *testing.T) {
	validDealer := *seat.MustSeat(1)
	validStartingDealer := *seat.MustSeat(0)
	validDora := *tile.MustTileFromCode("1m")
	unknownHands := newValidHands()
	for i := range unknownHands[0] {
		unknownHands[0][i] = *tile.MustTileFromCode("?")
	}

	ev, err := event.NewStartRound(
		wind.East,
		1,
		0,
		0,
		validDealer,
		validStartingDealer,
		validDora,
		&[common.NumPlayers]int{25000, 25000, 25000, 25000},
		unknownHands,
	)
	if err != nil {
		t.Fatalf("NewStartRound() failed: %v", err)
	}

	var s State
	if err := s.Apply(ev); err != nil {
		t.Fatalf("State.Apply() failed: %v", err)
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

func TestApplyStartRoundErrorsOnInvalidVisibleHand(t *testing.T) {
	validDealer := *seat.MustSeat(1)
	validStartingDealer := *seat.MustSeat(0)
	validDora := *tile.MustTileFromCode("1m")
	invalidHands := newValidHands()
	for i := range 5 {
		invalidHands[0][i] = *tile.MustTileFromCode("1m")
	}

	ev, err := event.NewStartRound(
		wind.East,
		1,
		0,
		0,
		validDealer,
		validStartingDealer,
		validDora,
		&[common.NumPlayers]int{25000, 25000, 25000, 25000},
		invalidHands,
	)
	if err != nil {
		t.Fatalf("NewStartRound() failed: %v", err)
	}

	var s State
	if err := s.Apply(ev); err == nil {
		t.Fatal("State.Apply() succeeded unexpectedly")
	}
}
