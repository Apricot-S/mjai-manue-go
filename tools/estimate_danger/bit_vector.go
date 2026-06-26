package main

import (
	"math/big"
	"slices"
)

type BitVector = big.Int

var one = big.NewInt(1)

func BoolArrayToBitVector(boolArray []bool) *BitVector {
	vector := big.NewInt(0)
	for _, b := range slices.Backward(boolArray) {
		vector.Lsh(vector, 1)
		if b {
			vector.Or(vector, one)
		}
	}
	return vector
}

func Matches(featureVector, positiveMask, negativeMask *BitVector) bool {
	temp1 := new(BitVector)
	temp1.And(featureVector, positiveMask)
	positiveMatch := temp1.Cmp(positiveMask) == 0

	temp2 := new(BitVector)
	temp2.Or(featureVector, negativeMask)
	negativeMatch := temp2.Cmp(negativeMask) == 0

	return positiveMatch && negativeMatch
}
