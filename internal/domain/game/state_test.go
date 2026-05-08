package game_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
)

func TestNewDefaultState(t *testing.T) {
	state := game.NewDefaultState()

	want := [common.NumPlayers]int{25000, 25000, 25000, 25000}
	if got := state.Scores(); got != want {
		t.Errorf("Scores() = %v, want %v", got, want)
	}
}

func TestState_UpdateScores(t *testing.T) {
	state := game.NewState([common.NumPlayers]int{25000, 25000, 25000, 25000})

	want := [common.NumPlayers]int{26000, 24000, 25000, 25000}
	state.UpdateScores(want)

	if got := state.Scores(); got != want {
		t.Errorf("Scores() = %v, want %v", got, want)
	}
}
