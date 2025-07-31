package base

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

func NewShuntsu(pai1, pai2, pai3 Pai) *Shuntsu {
	return &Shuntsu{pai1, pai2, pai3}
}

func (s *Shuntsu) ToString() string {
	return mentsuToString("shuntsu", s[:])
}

func (s *Shuntsu) Pais() []Pai {
	return s[:]
}

type Kotsu [3]Pai

func NewKotsu(pai1, pai2, pai3 Pai) *Kotsu {
	return &Kotsu{pai1, pai2, pai3}
}

func (k *Kotsu) ToString() string {
	return mentsuToString("kotsu", k[:])
}

func (k *Kotsu) Pais() []Pai {
	return k[:]
}

type Kantsu [4]Pai

func NewKantsu(pai1, pai2, pai3, pai4 Pai) *Kantsu {
	return &Kantsu{pai1, pai2, pai3, pai4}
}

func (k *Kantsu) ToString() string {
	return mentsuToString("kantsu", k[:])
}

func (k *Kantsu) Pais() []Pai {
	return k[:]
}

type Toitsu [2]Pai

func NewToitsu(pai1, pai2 Pai) *Toitsu {
	return &Toitsu{pai1, pai2}
}

func (t *Toitsu) ToString() string {
	return mentsuToString("toitsu", t[:])
}

func (t *Toitsu) Pais() []Pai {
	return t[:]
}
