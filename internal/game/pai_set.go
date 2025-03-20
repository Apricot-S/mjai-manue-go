package game

import (
	"errors"
	"fmt"
)

type PaiSet struct {
	array [NumIDs]int
}

func NewPaiSet(array [NumIDs]int) *PaiSet {
	return &PaiSet{array}
}

func NewPaiSetWithPais(pais []Pai) (*PaiSet, error) {
	ps := &PaiSet{}
	err := ps.AddPais(pais)
	if err != nil {
		return nil, fmt.Errorf("failed to create PaiSet with pais: %w", err)
	}
	return ps, nil
}

func GetAll() *PaiSet {
	var array [NumIDs]int
	for i := range array {
		array[i] = 4
	}
	return &PaiSet{array}
}

func (ps *PaiSet) Array() [NumIDs]int {
	return ps.array
}

func (ps *PaiSet) ToPais() []Pai {
	pais := []Pai{}
	for id, count := range ps.array {
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

	var id uint8
	if pai.IsRed() {
		id = pai.RemoveRed().ID()
	} else {
		id = pai.ID()
	}

	return ps.array[id], nil
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

	var id uint8
	if pai.IsRed() {
		id = pai.RemoveRed().ID()
	} else {
		id = pai.ID()
	}
	ps.array[id] += n

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
	for i, c := range paiSet.array {
		ps.array[i] -= c
	}
}

func (ps *PaiSet) ToString() string {
	return PaisToStr(ps.ToPais())
}
