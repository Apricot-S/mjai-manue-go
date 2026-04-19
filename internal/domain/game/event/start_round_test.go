package event_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

const initHandSize = 13

func newValidHands() [common.NumPlayers][initHandSize]tile.Tile {
	var hands [common.NumPlayers][initHandSize]tile.Tile
	base := tile.MustTileFromCode("1m")
	for player := range common.NumPlayers {
		for i := range initHandSize {
			hands[player][i] = *base
		}
	}
	return hands
}

func TestNewStartRound(t *testing.T) {
	validDealer := *seat.MustSeat(1)
	validStartingDealer := *seat.MustSeat(0)
	validDora := *tile.MustTileFromCode("1m")
	validHands := newValidHands()
	validScores := &[common.NumPlayers]int{25000, 25000, 25000, 25000}

	tests := []struct {
		name           string
		roundWind      wind.Wind
		roundNumber    int
		honba          int
		riichiDeposit  int
		dealer         seat.Seat
		startingDealer seat.Seat
		doraIndicator  tile.Tile
		scores         *[common.NumPlayers]int
		hands          [common.NumPlayers][initHandSize]tile.Tile
		wantErr        bool
	}{
		{
			name:           "valid parameters",
			roundWind:      wind.East,
			roundNumber:    1,
			honba:          0,
			riichiDeposit:  0,
			dealer:         validDealer,
			startingDealer: validStartingDealer,
			doraIndicator:  validDora,
			scores:         validScores,
			hands:          validHands,
			wantErr:        false,
		},
		{
			name:           "nil scores allowed",
			roundWind:      wind.East,
			roundNumber:    1,
			honba:          0,
			riichiDeposit:  0,
			dealer:         validDealer,
			startingDealer: validStartingDealer,
			doraIndicator:  validDora,
			scores:         nil,
			hands:          validHands,
			wantErr:        false,
		},
		{
			name:           "invalid round wind",
			roundWind:      wind.Wind(0),
			roundNumber:    1,
			honba:          0,
			riichiDeposit:  0,
			dealer:         validDealer,
			startingDealer: validStartingDealer,
			doraIndicator:  validDora,
			scores:         validScores,
			hands:          validHands,
			wantErr:        true,
		},
		{
			name:           "invalid round number 0",
			roundWind:      wind.East,
			roundNumber:    0,
			honba:          0,
			riichiDeposit:  0,
			dealer:         validDealer,
			startingDealer: validStartingDealer,
			doraIndicator:  validDora,
			scores:         validScores,
			hands:          validHands,
			wantErr:        true,
		},
		{
			name:           "invalid round number 5",
			roundWind:      wind.East,
			roundNumber:    5,
			honba:          0,
			riichiDeposit:  0,
			dealer:         validDealer,
			startingDealer: validStartingDealer,
			doraIndicator:  validDora,
			scores:         validScores,
			hands:          validHands,
			wantErr:        true,
		},
		{
			name:           "negative honba",
			roundWind:      wind.East,
			roundNumber:    1,
			honba:          -1,
			riichiDeposit:  0,
			dealer:         validDealer,
			startingDealer: validStartingDealer,
			doraIndicator:  validDora,
			scores:         validScores,
			hands:          validHands,
			wantErr:        true,
		},
		{
			name:           "negative riichi deposit",
			roundWind:      wind.East,
			roundNumber:    1,
			honba:          0,
			riichiDeposit:  -1,
			dealer:         validDealer,
			startingDealer: validStartingDealer,
			doraIndicator:  validDora,
			scores:         validScores,
			hands:          validHands,
			wantErr:        true,
		},
		{
			name:           "unknown dora indicator",
			roundWind:      wind.East,
			roundNumber:    1,
			honba:          0,
			riichiDeposit:  0,
			dealer:         validDealer,
			startingDealer: validStartingDealer,
			doraIndicator:  *tile.MustTileFromCode("?"),
			scores:         validScores,
			hands:          validHands,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := event.NewStartRound(
				tt.roundWind,
				tt.roundNumber,
				tt.honba,
				tt.riichiDeposit,
				tt.dealer,
				tt.startingDealer,
				tt.doraIndicator,
				tt.scores,
				tt.hands,
			)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("NewStartRound() failed: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewStartRound() succeeded unexpectedly")
			}
			if got == nil {
				t.Fatal("NewStartRound() returned nil without error")
			}
		})
	}
}
