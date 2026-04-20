package event_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

func newValidHands() [common.NumPlayers][common.InitHandSize]tile.Tile {
	var hands [common.NumPlayers][common.InitHandSize]tile.Tile
	base := tile.MustTileFromCode("1m")
	for player := range common.NumPlayers {
		for i := range common.InitHandSize {
			hands[player][i] = *base
		}
	}
	return hands
}

func TestNewStartRound(t *testing.T) {
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
		dealer        seat.Seat
		doraIndicator tile.Tile
		scores        *[common.NumPlayers]int
		hands         [common.NumPlayers][common.InitHandSize]tile.Tile
		wantErr       bool
	}{
		{
			name:          "valid parameters",
			roundWind:     wind.East,
			roundNumber:   1,
			honba:         0,
			riichiDeposit: 0,
			dealer:        validDealer,
			doraIndicator: validDora,
			scores:        validScores,
			hands:         validHands,
			wantErr:       false,
		},
		{
			name:          "nil scores allowed",
			roundWind:     wind.East,
			roundNumber:   1,
			honba:         0,
			riichiDeposit: 0,
			dealer:        validDealer,
			doraIndicator: validDora,
			scores:        nil,
			hands:         validHands,
			wantErr:       false,
		},
		{
			name:          "invalid round wind",
			roundWind:     wind.Wind(0),
			roundNumber:   1,
			honba:         0,
			riichiDeposit: 0,
			dealer:        validDealer,
			doraIndicator: validDora,
			scores:        validScores,
			hands:         validHands,
			wantErr:       true,
		},
		{
			name:          "invalid round number 0",
			roundWind:     wind.East,
			roundNumber:   0,
			honba:         0,
			riichiDeposit: 0,
			dealer:        validDealer,
			doraIndicator: validDora,
			scores:        validScores,
			hands:         validHands,
			wantErr:       true,
		},
		{
			name:          "invalid round number 5",
			roundWind:     wind.East,
			roundNumber:   5,
			honba:         0,
			riichiDeposit: 0,
			dealer:        validDealer,
			doraIndicator: validDora,
			scores:        validScores,
			hands:         validHands,
			wantErr:       true,
		},
		{
			name:          "negative honba",
			roundWind:     wind.East,
			roundNumber:   1,
			honba:         -1,
			riichiDeposit: 0,
			dealer:        validDealer,
			doraIndicator: validDora,
			scores:        validScores,
			hands:         validHands,
			wantErr:       true,
		},
		{
			name:          "negative riichi deposit",
			roundWind:     wind.East,
			roundNumber:   1,
			honba:         0,
			riichiDeposit: -1,
			dealer:        validDealer,
			doraIndicator: validDora,
			scores:        validScores,
			hands:         validHands,
			wantErr:       true,
		},
		{
			name:          "unknown dora indicator",
			roundWind:     wind.East,
			roundNumber:   1,
			honba:         0,
			riichiDeposit: 0,
			dealer:        validDealer,
			doraIndicator: *tile.MustTileFromCode("?"),
			scores:        validScores,
			hands:         validHands,
			wantErr:       true,
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

func TestStartRoundAccessors(t *testing.T) {
	validDealer := *seat.MustSeat(1)
	validDora := *tile.MustTileFromCode("1m")
	validHands := newValidHands()
	validScores := &[common.NumPlayers]int{25000, 25000, 25000, 25000}

	got, err := event.NewStartRound(
		wind.South,
		2,
		1,
		2,
		validDealer,
		validDora,
		validScores,
		validHands,
	)
	if err != nil {
		t.Fatalf("NewStartRound() failed: %v", err)
	}
	if got.RoundWind() != wind.South {
		t.Fatalf("RoundWind() = %v, want %v", got.RoundWind(), wind.South)
	}
	if got.RoundNumber() != 2 {
		t.Fatalf("RoundNumber() = %d, want %d", got.RoundNumber(), 2)
	}
	if got.Honba() != 1 {
		t.Fatalf("Honba() = %d, want %d", got.Honba(), 1)
	}
	if got.RiichiDeposit() != 2 {
		t.Fatalf("RiichiDeposit() = %d, want %d", got.RiichiDeposit(), 2)
	}
	if got.Dealer().Index() != validDealer.Index() {
		t.Fatalf("Dealer() = %v, want %v", got.Dealer(), validDealer)
	}
	if got.DoraIndicator().ID() != validDora.ID() {
		t.Fatalf("DoraIndicator() = %v, want %v", got.DoraIndicator(), validDora)
	}
	if got.Scores() != validScores {
		t.Fatalf("Scores() = %v, want %v", got.Scores(), validScores)
	}
	if got.Hands() != validHands {
		t.Fatalf("Hands() = %v, want %v", got.Hands(), validHands)
	}
}
