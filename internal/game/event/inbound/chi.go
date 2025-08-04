package inbound

import (
	"fmt"
	"slices"
	"sort"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

type Chi struct {
	Actor    int `validate:"min=0,max=3"`
	Target   int `validate:"min=0,max=3"`
	Taken    base.Pai
	Consumed [2]base.Pai
}

func NewChi(actor int, target int, taken base.Pai, consumed [2]base.Pai) (*Chi, error) {
	event := &Chi{
		Actor:    actor,
		Target:   target,
		Taken:    taken,
		Consumed: consumed,
	}

	if getPlayerDistance(event.Actor, event.Target) != 1 {
		return nil, fmt.Errorf("target must be the kamicha of actor: actor=%d, target=%d", actor, target)
	}

	isSameColor := !slices.ContainsFunc(event.Consumed[:], func(p base.Pai) bool {
		return event.Taken.Type() != p.Type()
	})
	if !isSameColor {
		return nil, fmt.Errorf("taken tile must be the same color as the consumed tile: %v", event)
	}

	pais := slices.Concat(base.Pais{event.Taken}, event.Consumed[:])
	sort.Sort(pais)

	isSuhai := !slices.ContainsFunc(pais, func(p base.Pai) bool {
		return p.IsTsupai() || p.IsUnknown()
	})
	if !isSuhai {
		return nil, fmt.Errorf("chi tiles must be suhai: %v", event)
	}

	isSequence := pais[0].Number()+1 == pais[1].Number() && pais[1].Number()+1 == pais[2].Number()
	if !isSequence {
		return nil, fmt.Errorf("consumed tiles must be a sequence with the taken tile: %v", event)
	}

	if err := eventValidator.Struct(event); err != nil {
		return nil, err
	}
	return event, nil
}

func getPlayerDistance(p1 int, p2 int) int {
	const NumPlayers = 4
	return (NumPlayers + p1 - p2) % NumPlayers
}

func (c *Chi) isInboundEvent() {}
