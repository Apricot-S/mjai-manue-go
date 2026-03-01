package round_test

import (
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
