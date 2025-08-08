package outbound

import (
	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

func mustPai(name string) *base.Pai {
	p, err := base.NewPaiWithName(name)
	if err != nil {
		panic(err)
	}
	return p
}

func mustPais(names ...string) []base.Pai {
	pais := make([]base.Pai, len(names))
	for i, n := range names {
		pais[i] = *mustPai(n)
	}
	return pais
}
