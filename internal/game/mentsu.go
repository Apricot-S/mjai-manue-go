package game

import (
	"fmt"
	"strings"
)

type MentsuType int

const (
	Shuntsu MentsuType = iota + 1
	Kotsu
	Kantsu
	Toitsu
)

type Mentsu struct {
	typ  MentsuType
	pais []Pai
}

func NewMentsu(t MentsuType, pais []Pai) (*Mentsu, error) {
	switch t {
	case Shuntsu, Kotsu:
		if len(pais) != 3 {
			return nil, fmt.Errorf("invalid %v", t)
		}
	case Kantsu:
		if len(pais) != 4 {
			return nil, fmt.Errorf("invalid kantsu")
		}
	case Toitsu:
		if len(pais) != 2 {
			return nil, fmt.Errorf("invalid toitsu")
		}
	default:
		return nil, fmt.Errorf("unknown type")
	}
	return &Mentsu{typ: t, pais: pais}, nil
}

func (m *Mentsu) ToString() string {
	str := ""
	switch m.typ {
	case Shuntsu:
		str = "shuntsu"
	case Kotsu:
		str = "kotsu"
	case Kantsu:
		str = "kantsu"
	case Toitsu:
		str = "toitsu"
	}

	paiStrs := make([]string, len(m.pais))
	for i, p := range m.pais {
		paiStrs[i] = p.ToString()
	}

	str += fmt.Sprintf("[%s]", strings.Join(paiStrs, " "))
	return str
}

func (m *Mentsu) Type() MentsuType {
	return m.typ
}

func (m *Mentsu) Pais() []Pai {
	return m.pais
}
