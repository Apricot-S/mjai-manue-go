package game

import (
	"fmt"
	"strings"
)

type Mentsu interface {
	ToString() string
	Pais() []Pai
}

func mentsuToString(name string, pais []Pai) string {
	paiStrs := make([]string, len(pais))
	for i, p := range pais {
		paiStrs[i] = p.ToString()
	}
	return fmt.Sprintf("%s: [%s]", name, strings.Join(paiStrs, " "))
}

type Shuntsu [3]Pai

func NewShuntsu(pais [3]Pai) *Shuntsu {
	s := Shuntsu(pais)
	return &s
}

func (s *Shuntsu) ToString() string {
	return mentsuToString("shuntsu", s[:])
}

func (s *Shuntsu) Pais() []Pai {
	return s[:]
}

type Kotsu [3]Pai

func NewKotsu(pais [3]Pai) *Kotsu {
	k := Kotsu(pais)
	return &k
}

func (k *Kotsu) ToString() string {
	return mentsuToString("kotsu", k[:])
}

func (k *Kotsu) Pais() []Pai {
	return k[:]
}

type Kantsu [4]Pai

func NewKantsu(pais [4]Pai) *Kantsu {
	k := Kantsu(pais)
	return &k
}

func (k *Kantsu) ToString() string {
	return mentsuToString("kantsu", k[:])
}

func (k *Kantsu) Pais() []Pai {
	return k[:]
}

type Toitsu [2]Pai

func NewToitsu(pais [2]Pai) *Toitsu {
	t := Toitsu(pais)
	return &t
}

func (t *Toitsu) ToString() string {
	return mentsuToString("toitsu", t[:])
}

func (t *Toitsu) Pais() []Pai {
	return t[:]
}
