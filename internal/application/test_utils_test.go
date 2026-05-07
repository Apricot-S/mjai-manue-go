package application_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/application"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/ai"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

func mustNewBotForTest(t *testing.T, self seat.Seat) *application.Bot {
	t.Helper()

	return application.NewBot(self, newTsumogiriAgentForTest(), nil)
}

func newTsumogiriAgentForTest() ai.Agent {
	return ai.NewTsumogiriAgent()
}

func mustNewStartRoundForTest(t *testing.T, hands [common.NumPlayers][common.InitHandSize]tile.Tile) *event.StartRound {
	t.Helper()

	return event.NewStartRound(
		wind.East,
		1,
		0,
		0,
		*seat.MustSeat(0),
		tile.MustTileFromCode("E"),
		&[common.NumPlayers]int{25000, 25000, 25000, 25000},
		hands,
	)
}

func newValidHands() [common.NumPlayers][common.InitHandSize]tile.Tile {
	return [common.NumPlayers][common.InitHandSize]tile.Tile{
		{
			tile.MustTileFromCode("1m"),
			tile.MustTileFromCode("2m"),
			tile.MustTileFromCode("3m"),
			tile.MustTileFromCode("4m"),
			tile.MustTileFromCode("5m"),
			tile.MustTileFromCode("6m"),
			tile.MustTileFromCode("7m"),
			tile.MustTileFromCode("8m"),
			tile.MustTileFromCode("9m"),
			tile.MustTileFromCode("1p"),
			tile.MustTileFromCode("2p"),
			tile.MustTileFromCode("3p"),
			tile.MustTileFromCode("4p"),
		},
		unknownHand(),
		unknownHand(),
		unknownHand(),
	}
}

func unknownHand() [common.InitHandSize]tile.Tile {
	var hand [common.InitHandSize]tile.Tile
	for i := range common.InitHandSize {
		hand[i] = tile.MustTileFromCode("?")
	}
	return hand
}
