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
