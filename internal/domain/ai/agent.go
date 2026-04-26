package ai

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
)

type Request struct {
	Self  seat.Seat
	Round round.StateViewer
}

type Decision struct {
	Action action.Action
	Log    string
}

type Agent interface {
	Decide(request Request) (Decision, error)
}
