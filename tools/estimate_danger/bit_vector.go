package main

import "math/big"

type BitVector struct {
	*big.Int
}

var one = big.NewInt(1)

func boolArrayToBitVector(boolArray []bool) *BitVector {
	bitVector := big.NewInt(0)
	for i := len(boolArray) - 1; i >= 0; i-- {
		bitVector.Lsh(bitVector, 1)
		if boolArray[i] {
			bitVector.Or(bitVector, one)
		}
	}
	return &BitVector{bitVector}
}
