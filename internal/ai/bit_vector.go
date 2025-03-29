package ai

import "github.com/Apricot-S/mjai-manue-go/internal/game"

type BitVector uint64

func NewBitVector(countVector *game.PaiSet, threshold int) BitVector {
	var bv BitVector
	for i, c := range countVector {
		if c >= threshold {
			bv |= BitVector(1) << i
		}
	}
	return bv
}

func (bv BitVector) IsSubsetOf(other BitVector) bool {
	return (bv & other) == bv
}

func (bv BitVector) HasIntersectionWith(other BitVector) bool {
	return (bv & other) != 0
}
