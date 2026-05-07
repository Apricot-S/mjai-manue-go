package round

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

func mustNewRoundStateForTest(t *testing.T, hands [common.NumPlayers][common.InitHandSize]tile.Tile) *State {
	t.Helper()

	validDealer := *seat.MustSeat(0)
	validDora := tile.MustTileFromCode("E")
	validScores := &[common.NumPlayers]int{25000, 25000, 25000, 25000}

	ev := event.NewStartRound(
		wind.East,
		1,
		0,
		0,
		validDealer,
		validDora,
		validScores,
		hands,
	)

	s, err := NewState(ev, *validScores)
	if err != nil {
		t.Fatalf("round.NewState() failed: %v", err)
	}
	return s
}

func newVisiblePlayersForTest(t *testing.T, hands [common.NumPlayers][common.InitHandSize]tile.Tile) [common.NumPlayers]player.Player {
	t.Helper()

	var players [common.NumPlayers]player.Player
	for i, handTiles := range hands {
		p, err := player.NewVisiblePlayer(handTiles)
		if err != nil {
			t.Fatalf("player.NewVisiblePlayer(%d) failed: %v", i, err)
		}
		players[i] = p
	}
	return players
}
