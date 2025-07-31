package base

import (
	"errors"
	"fmt"
)

// PaiSet represents the number of each tile.
// Red 5 are not distinguished from normal 5.
type PaiSet [NumIDs]int

func NewPaiSet(pais []Pai) (*PaiSet, error) {
	ps := &PaiSet{}
	err := ps.AddPais(pais)
	if err != nil {
		return nil, fmt.Errorf("failed to create PaiSet with pais: %w", err)
	}
	return ps, nil
}

func GetAll() *PaiSet {
	var ps PaiSet
	for i := range ps {
		ps[i] = 4
	}
	return &ps
}

func (ps *PaiSet) ToPais() []Pai {
	pais := []Pai{}
	for id, count := range ps {
		for range count {
			pai, _ := NewPaiWithID(uint8(id))
			pais = append(pais, *pai)
		}
	}
	return pais
}

func (ps *PaiSet) Count(pai *Pai) (int, error) {
	if pai.IsUnknown() {
		return 0, errors.New("PaiSet does not contain unknowns")
	}

	id := pai.RemoveRed().ID()

	return ps[id], nil
}

func (ps *PaiSet) Has(pai *Pai) (bool, error) {
	c, err := ps.Count(pai)
	if err != nil {
		return false, err
	}
	return c > 0, nil
}

func (ps *PaiSet) AddPai(pai *Pai, n int) error {
	if pai.IsUnknown() {
		return errors.New("PaiSet cannot contain unknowns")
	}

	id := pai.RemoveRed().ID()
	ps[id] += n

	return nil
}

func (ps *PaiSet) AddPais(pais []Pai) error {
	for _, pai := range pais {
		err := ps.AddPai(&pai, 1)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ps *PaiSet) RemovePaiSet(paiSet *PaiSet) {
	for i, c := range paiSet {
		ps[i] -= c
	}
}

func (ps *PaiSet) ToString() string {
	return PaisToStr(ps.ToPais())
}
