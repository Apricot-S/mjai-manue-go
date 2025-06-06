package core

import (
	"math/rand/v2"
)

func CreateRNG() *rand.Rand {
	src := rand.NewPCG(0, 0)
	rng := rand.New(src)
	return rng
}

func ShuffleWall[T any](rng *rand.Rand, wall *[]T) {
	rng.Shuffle(len(*wall), func(i, j int) {
		(*wall)[i], (*wall)[j] = (*wall)[j], (*wall)[i]
	})
}
