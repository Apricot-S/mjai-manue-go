package event

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
)

type RiichiAccepted struct {
	actor  seat.Seat
	deltas *[common.NumPlayers]int
	scores *[common.NumPlayers]int
}

func NewRiichiAccepted(actor seat.Seat, deltas, scores *[common.NumPlayers]int) *RiichiAccepted {
	return &RiichiAccepted{
		actor:  actor,
		deltas: deltas,
		scores: scores,
	}
}

func (*RiichiAccepted) isEvent() {}

func (r *RiichiAccepted) Actor() seat.Seat {
	return r.actor
}

func (r *RiichiAccepted) Deltas() *[common.NumPlayers]int {
	return r.deltas
}

func (r *RiichiAccepted) Scores() *[common.NumPlayers]int {
	return r.scores
}
