package base

import (
	"fmt"
)

type Hora struct {
	pai    Pai
	target int
}

func NewHora(pai Pai, target int) (*Hora, error) {
	if target < 0 || target > 3 {
		return nil, fmt.Errorf("hora: invalid target player index (must be 0-3, got: %d)", target)
	}

	return &Hora{
		pai:    pai,
		target: target,
	}, nil
}

func (h *Hora) Pai() *Pai {
	return &h.pai
}

func (h *Hora) Target() int {
	return h.target
}
