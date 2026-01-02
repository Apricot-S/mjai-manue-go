package service

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/hand"
)

func IsTenpaiGeneral(hand *hand.VisibleHand) bool {
	shanten, _ := AnalyzeShanten(hand, UpperBound(0))
	return shanten <= 0
}

func IsTenpaiAll(hand *hand.VisibleHand) bool {
	shanten, _ := AnalyzeShanten(hand, UpperBound(0))
	if shanten <= 0 {
		return true
	}

	shanten = AnalyzeShantenChiitoitsu(hand)
	if shanten <= 0 {
		return true
	}

	shanten = AnalyzeShantenKokushimuso(hand)
	return shanten <= 0
}
