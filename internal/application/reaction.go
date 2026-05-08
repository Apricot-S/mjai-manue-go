package application

import "github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"

type ReactionKind int

const (
	ReactionNone ReactionKind = iota + 1
	ReactionAction
)

type Reaction struct {
	kind   ReactionKind
	action action.Action
	log    string
}

func NewNoReaction() Reaction {
	return Reaction{kind: ReactionNone}
}

func NewActionReaction(a action.Action, log string) Reaction {
	return Reaction{
		kind:   ReactionAction,
		action: a,
		log:    log,
	}
}

func (r Reaction) Kind() ReactionKind {
	return r.kind
}

func (r Reaction) Action() action.Action {
	return r.action
}

func (r Reaction) Log() string {
	return r.log
}
