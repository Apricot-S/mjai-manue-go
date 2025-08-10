package main

import "math/big"

type BitVector = big.Int

var one = big.NewInt(1)

func BoolArrayToBitVector(boolArray []bool) *BitVector {
	bitVector := big.NewInt(0)
	for i := len(boolArray) - 1; i >= 0; i-- {
		bitVector.Lsh(bitVector, 1)
		if boolArray[i] {
			bitVector.Or(bitVector, one)
		}
	}
	return bitVector
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
