package event

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type Win struct {
	actor  seat.Seat
	target seat.Seat
	// winningTile is optional because some hora messages/logs omit pai.
	// When nil, round state validation skips tile equality checks.
	winningTile   *tile.Tile
	winningPoints int
	deltas        *[common.NumPlayers]int
	scores        *[common.NumPlayers]int
}

func NewWin(
	actor, target seat.Seat,
	winningTile *tile.Tile,
	winningPoints int,
	deltas *[common.NumPlayers]int,
	scores *[common.NumPlayers]int,
) *Win {
	return &Win{
		actor:         actor,
		target:        target,
		winningTile:   winningTile,
		winningPoints: winningPoints,
		deltas:        deltas,
		scores:        scores,
	}
}

func (*Win) isEvent() {}

func (w *Win) Actor() seat.Seat {
	return w.actor
}

func (w *Win) Target() seat.Seat {
	return w.target
}

func (w *Win) WinningTile() *tile.Tile {
	return w.winningTile
}

func (w *Win) WinningPoints() int {
	return w.winningPoints
}

func (w *Win) Deltas() *[common.NumPlayers]int {
	return w.deltas
}

func (w *Win) Scores() *[common.NumPlayers]int {
	return w.scores
}
