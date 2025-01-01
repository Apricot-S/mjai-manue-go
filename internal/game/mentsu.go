package game

import (
	"fmt"
	"strings"
)

type MentsuType int

const (
	Shuntsu MentsuType = iota
	Kotsu
	Kantsu
	Toitsu
)

type Mentsu struct {
	Type MentsuType
	Pais []Pai
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
	return &Mentsu{Type: t, Pais: pais}, nil
}

func (m *Mentsu) ToString() string {
	str := ""
	switch m.Type {
	case Shuntsu:
		str = "shuntsu"
	case Kotsu:
		str = "kotsu"
	case Kantsu:
		str = "kantsu"
	case Toitsu:
		str = "toitsu"
	}

	paiStrs := make([]string, len(m.Pais))
	for i, p := range m.Pais {
		paiStrs[i] = p.ToString()
	}

	str += fmt.Sprintf("[%s]", strings.Join(paiStrs, " "))
	return str
}
