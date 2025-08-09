package main

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/ai/core"
	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/game"
)

type evaluator func(*Scene, *base.Pai) (bool, error)
type evaluators map[string]evaluator

var defaultEvaluators = registerEvaluators()

type Scene struct {
	tehaiSet   *base.PaiSet
	anpaiSet   *base.PaiSet
	visibleSet *base.PaiSet
	doraSet    *base.PaiSet
	bakaze     *base.Pai
	targetKaze *base.Pai

	prereachSutehaiSet *base.PaiSet
	earlySutehaiSet    *base.PaiSet
	lateSutehaiSet     *base.PaiSet
	reachPaiSet        *base.PaiSet

	evaluators *evaluators
}

func NewScene(
	tehais, anpais, visibles, doras, prereachSutehais base.Pais,
	bakaze, targetKaze *base.Pai,
) (*Scene, error) {
	s := &Scene{
		bakaze:     bakaze,
		targetKaze: targetKaze,
		evaluators: defaultEvaluators,
	}

	var err error
	if s.tehaiSet, err = base.NewPaiSet(tehais); err != nil {
		return nil, err
	}
	if s.anpaiSet, err = base.NewPaiSet(anpais); err != nil {
		return nil, err
	}
	if s.visibleSet, err = base.NewPaiSet(visibles); err != nil {
		return nil, err
	}
	if s.doraSet, err = base.NewPaiSet(doras); err != nil {
		return nil, err
	}

	if s.prereachSutehaiSet, err = base.NewPaiSet(prereachSutehais); err != nil {
		return nil, err
	}

	halfLen := len(prereachSutehais) / 2
	if s.earlySutehaiSet, err = base.NewPaiSet(prereachSutehais[:halfLen]); err != nil {
		return nil, err
	}
	if s.lateSutehaiSet, err = base.NewPaiSet(prereachSutehais[halfLen:]); err != nil {
		return nil, err
	}

	var reachPais base.Pais = nil
	if len(prereachSutehais) != 0 {
		// prereach_sutehais can be empty in unit tests.
		reachPai := prereachSutehais[len(prereachSutehais)-1].RemoveRed()
		reachPais = base.Pais{*reachPai}
	}
	if s.reachPaiSet, err = base.NewPaiSet(reachPais); err != nil {
		return nil, err
	}

	return s, nil
}

func NewSceneWithState(gameState game.StateViewer, me *base.Player, target *base.Player) (*Scene, error) {
	s := &Scene{
		evaluators: defaultEvaluators,
	}

	var err error
	if s.tehaiSet, err = base.NewPaiSet(me.Tehais()); err != nil {
		return nil, err
	}
	if s.anpaiSet, err = base.NewPaiSet(gameState.Anpais(target)); err != nil {
		return nil, err
	}
	if s.visibleSet, err = base.NewPaiSet(gameState.VisiblePais(me)); err != nil {
		return nil, err
	}
	if s.doraSet, err = base.NewPaiSet(gameState.Doras()); err != nil {
		return nil, err
	}

	s.bakaze = gameState.Bakaze()
	s.targetKaze = gameState.Jikaze(target)

	var prereachSutehais base.Pais = nil
	var reachPais base.Pais = nil
	if idx := target.ReachSutehaiIndex(); idx != -1 {
		sutehais := target.Sutehais()
		prereachSutehais = sutehais[:idx+1]
		reachPai := sutehais[idx]
		reachPais = base.Pais{reachPai}
	}
	if s.prereachSutehaiSet, err = base.NewPaiSet(prereachSutehais); err != nil {
		return nil, err
	}
	if s.reachPaiSet, err = base.NewPaiSet(reachPais); err != nil {
		return nil, err
	}

	halfLen := len(prereachSutehais) / 2
	if s.earlySutehaiSet, err = base.NewPaiSet(prereachSutehais[:halfLen]); err != nil {
		return nil, err
	}
	if s.lateSutehaiSet, err = base.NewPaiSet(prereachSutehais[halfLen:]); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Scene) Evaluate(name string, pai *base.Pai) (bool, error) {
	if evaluator, ok := (*s.evaluators)[name]; ok {
		return evaluator(s, pai)
	}
	return false, fmt.Errorf("an unknown feature name was specified: %v", name)
}

// func isAnpai(pai *base.Pai, anpaiSet *base.PaiSet) (bool, error) {
// 	return anpaiSet.Has(pai)
// }

func isTsupai(pai *base.Pai) bool {
	return pai.IsTsupai()
}

// Omotesuji (表筋) or Nakasuji (中筋)
func isSuji(pai *base.Pai, anpaiSet *base.PaiSet) (bool, error) {
	return isSujiOf(pai, anpaiSet)
}

// Katasuji (片筋) or Suji (筋)
func isWeakSuji(pai *base.Pai, anpaiSet *base.PaiSet) (bool, error) {
	return isWeakSujiOf(pai, anpaiSet)
}

// Suji for Riichi declaration tile. Including tiles like 4p against 1p Riichi.
func isReachSuji(pai *base.Pai, reachPaiSet *base.PaiSet) (bool, error) {
	return isWeakSujiOf(pai, reachPaiSet)
}

func isPrereachSuji(pai *base.Pai, prereachSutehaiSet *base.PaiSet) (bool, error) {
	return isSujiOf(pai, prereachSutehaiSet)
}

// Urasuji (裏筋)
// http://ja.wikipedia.org/wiki/%E7%AD%8B_(%E9%BA%BB%E9%9B%80)#.E8.A3.8F.E3.82.B9.E3.82.B8
func isUrasuji(pai *base.Pai, prereachSutehaiSet *base.PaiSet, anpaiSet *base.PaiSet) (bool, error) {
	return isUrasujiOf(pai, prereachSutehaiSet, anpaiSet)
}

func isEarlyUrasuji(pai *base.Pai, earlySutehaiSet *base.PaiSet, anpaiSet *base.PaiSet) (bool, error) {
	return isUrasujiOf(pai, earlySutehaiSet, anpaiSet)
}

func isReachUrasuji(pai *base.Pai, reachPaiSet *base.PaiSet, anpaiSet *base.PaiSet) (bool, error) {
	return isUrasujiOf(pai, reachPaiSet, anpaiSet)
}

func isUrasujiOf5(pai *base.Pai, prereachSutehaiSet *base.PaiSet, anpaiSet *base.PaiSet) (bool, error) {
	fiveSet := *prereachSutehaiSet
	for i := range base.NumIDs {
		isSuhai := (i / 9) < 3
		isFive := (i % 9) == 4
		if !(isSuhai && isFive) {
			fiveSet[i] = 0
		}
	}
	return isUrasujiOf(pai, &fiveSet, anpaiSet)
}

// Aidayonken (間四間)
// http://ja.wikipedia.org/wiki/%E7%AD%8B_(%E9%BA%BB%E9%9B%80)#.E9.96.93.E5.9B.9B.E9.96.93
func isAida4ken(pai *base.Pai, prereachSutehaiSet *base.PaiSet) (bool, error) {
	if pai.IsTsupai() {
		return false, nil
	}

	num := pai.Number()
	typ := pai.Type()

	var matchA, matchB bool
	var err error

	if 2 <= num && num <= 5 {
		if matchA, err = hasBoth(typ, num-1, num+4, prereachSutehaiSet); err != nil {
			return false, err
		}
	}

	if 5 <= num && num <= 8 {
		if matchB, err = hasBoth(typ, num-4, num+1, prereachSutehaiSet); err != nil {
			return false, err
		}
	}

	return matchA || matchB, nil
}

func hasBoth(paiType rune, n1, n2 uint8, set *base.PaiSet) (bool, error) {
	p1, err := base.NewPaiWithDetail(paiType, n1, false)
	if err != nil {
		return false, err
	}
	has1, err := set.Has(p1)
	if err != nil {
		return false, err
	}

	p2, err := base.NewPaiWithDetail(paiType, n2, false)
	if err != nil {
		return false, err
	}
	has2, err := set.Has(p2)
	if err != nil {
		return false, err
	}

	return has1 && has2, nil
}

// Matagisuji (跨ぎ筋)
// http://ja.wikipedia.org/wiki/%E7%AD%8B_(%E9%BA%BB%E9%9B%80)#.E3.81.BE.E3.81.9F.E3.81.8E.E3.82.B9.E3.82.B8
func isMatagisuji(pai *base.Pai, prereachSutehaiSet *base.PaiSet, anpaiSet *base.PaiSet) (bool, error) {
	return isMatagisujiOf(pai, prereachSutehaiSet, anpaiSet)
}

func isEarlyMatagisuji(pai *base.Pai, earlySutehaiSet *base.PaiSet, anpaiSet *base.PaiSet) (bool, error) {
	return isMatagisujiOf(pai, earlySutehaiSet, anpaiSet)
}

// func isLateMatagisuji(pai *base.Pai, lateSutehaiSet *base.PaiSet, anpaiSet *base.PaiSet) (bool, error) {
// 	return isMatagisujiOf(pai, lateSutehaiSet, anpaiSet)
// }

// func isReachMatagisuji(pai *base.Pai, reachPaiSet *base.PaiSet, anpaiSet *base.PaiSet) (bool, error) {
// 	return isMatagisujiOf(pai, reachPaiSet, anpaiSet)
// }

// // Senkisuji (疝気筋)
// // # http://ja.wikipedia.org/wiki/%E7%AD%8B_(%E9%BA%BB%E9%9B%80)#.E7.96.9D.E6.B0.97.E3.82.B9.E3.82.B8
// func isSenkisuji(pai *base.Pai, prereachSutehaiSet *base.PaiSet, anpaiSet *base.PaiSet) (bool, error) {
// 	return isSenkisujiOf(pai, prereachSutehaiSet, anpaiSet)
// }

// func isEarlySenkisuji(pai *base.Pai, earlySutehaiSet *base.PaiSet, anpaiSet *base.PaiSet) (bool, error) {
// 	return isSenkisujiOf(pai, earlySutehaiSet, anpaiSet)
// }

// func isOuterPrereachSutehai(pai *base.Pai, prereachSutehaiSet *base.PaiSet) (bool, error) {
// 	return isOuter(pai, prereachSutehaiSet)
// }

// func isOuterEarlySutehai(pai *base.Pai, earlySutehaiSet *base.PaiSet) (bool, error) {
// 	return isOuter(pai, earlySutehaiSet)
// }

// func isDora(pai *base.Pai, doraSet *base.PaiSet) (bool, error) {
// 	return doraSet.Has(pai)
// }

// func isDoraSuji(pai *base.Pai, doraSet *base.PaiSet) (bool, error) {
// 	return isWeakSujiOf(pai, doraSet)
// }

// func isDoraMatagi(pai *base.Pai, doraSet *base.PaiSet, anpaiSet *base.PaiSet) (bool, error) {
// 	return isMatagisujiOf(pai, doraSet, anpaiSet)
// }

// func isFanpai(pai *base.Pai, gameState game.StateViewer, target *base.Player) bool {
// 	return gameState.YakuhaiFan(pai, target) >= 1
// }

// func isRyenfonpai(pai *base.Pai, gameState game.StateViewer, target *base.Player) bool {
// 	return gameState.YakuhaiFan(pai, target) >= 2
// }

// func isSangenpai(pai *base.Pai) bool {
// 	return pai.IsTsupai() && pai.Number() >= 5
// }

// func isFonpai(pai *base.Pai) bool {
// 	return pai.IsTsupai() && pai.Number() < 5
// }

// func isBakaze(pai *base.Pai, bakaze *base.Pai) bool {
// 	return pai.HasSameSymbol(bakaze)
// }

// func isJikaze(pai *base.Pai, targetKaze *base.Pai) bool {
// 	return pai.HasSameSymbol(targetKaze)
// }

// func isNChanceOrLess(pai *base.Pai, n int, visibleSet *base.PaiSet) (bool, error) {
// 	if pai.IsTsupai() {
// 		return false, nil
// 	}

// 	paiNumber := pai.Number()
// 	if 4 <= paiNumber && paiNumber <= 6 {
// 		return false, nil
// 	}

// 	candidates := make([]uint8, 2)
// 	for i := uint8(1); i < 3; i++ {
// 		if paiNumber < 5 {
// 			candidates[i-1] = paiNumber + i
// 		} else {
// 			candidates[i-1] = paiNumber - i
// 		}
// 	}

// 	return core.AnyMatch(candidates, func(num uint8) (bool, error) {
// 		kabePai, err := base.NewPaiWithDetail(pai.Type(), num, false)
// 		if err != nil {
// 			return false, err
// 		}

// 		count, err := visibleSet.Count(kabePai)
// 		if err != nil {
// 			return false, err
// 		}

// 		return count >= 4-n, nil
// 	})
// }

// func isVisibleNOrMore(pai *base.Pai, n int, visibleSet *base.PaiSet) (bool, error) {
// 	c, err := visibleSet.Count(pai)
// 	if err != nil {
// 		return false, err
// 	}
// 	return c >= n, nil
// }

// func isSujiVisible(pai *base.Pai, n int, visibleSet *base.PaiSet) (bool, error) {
// 	if pai.IsTsupai() {
// 		return false, nil
// 	}

// 	suji, err := getSuji(pai)
// 	if err != nil {
// 		return false, err
// 	}

// 	return core.AnyMatch(suji, func(sujiPai game.Pai) (bool, error) {
// 		visible, err := isVisibleNOrMore(&sujiPai, n+1, visibleSet)
// 		if err != nil {
// 			return false, err
// 		}
// 		return !visible, nil
// 	})
// }

// func isNumNOrInner(pai *base.Pai, n uint8) bool {
// 	if pai.IsTsupai() {
// 		return false
// 	}

// 	paiNumber := pai.Number()
// 	if n <= paiNumber && paiNumber <= 10-n {
// 		return true
// 	}

// 	return false
// }

// func isInTehais(pai *base.Pai, n int, tehaiSet *base.PaiSet) (bool, error) {
// 	c, err := tehaiSet.Count(pai)
// 	return c >= n, err
// }

// func isSujiInTehais(pai *base.Pai, n int, tehaiSet *base.PaiSet) (bool, error) {
// 	if pai.IsTsupai() {
// 		return false, nil
// 	}

// 	suji, err := getSuji(pai)
// 	if err != nil {
// 		return false, err
// 	}

// 	return core.AnyMatch(suji, func(sujiPai game.Pai) (bool, error) {
// 		c, err := tehaiSet.Count(&sujiPai)
// 		return c >= n, err
// 	})
// }

// func isNOrMoreOfNeighborsInPrereachSutehais(
// 	pai *base.Pai,
// 	n int,
// 	neighborDistance int,
// 	prereachSutehaiSet *base.PaiSet,
// ) (bool, error) {
// 	if pai.IsTsupai() {
// 		return false, nil
// 	}

// 	paiNumber := int(pai.Number())
// 	numbers := make([]int, 0, 2*neighborDistance+1)
// 	for i := -neighborDistance; i <= neighborDistance; i++ {
// 		numbers = append(numbers, paiNumber+i)
// 	}

// 	numNeighbors, err := core.Count(numbers, func(num int) (bool, error) {
// 		if num < 1 || 9 < num {
// 			return false, nil
// 		}

// 		neighborPai, err := base.NewPaiWithDetail(pai.Type(), uint8(num), false)
// 		if err != nil {
// 			return false, err
// 		}

// 		count, err := prereachSutehaiSet.Count(neighborPai)
// 		if err != nil {
// 			return false, err
// 		}

// 		return count > 0, nil
// 	})
// 	if err != nil {
// 		return false, err
// 	}

// 	return numNeighbors >= n, nil
// }

// // n can be negative.
// func isNOuterPrereachSutehai(pai *base.Pai, n int, prereachSutehaiSet *base.PaiSet) (bool, error) {
// 	if pai.IsTsupai() {
// 		return false, nil
// 	}

// 	paiNumber := int(pai.Number())
// 	if paiNumber == 5 {
// 		return false, nil
// 	}

// 	nInnerNumber := 0
// 	if paiNumber < 5 {
// 		nInnerNumber = paiNumber + n
// 	} else {
// 		nInnerNumber = paiNumber - n
// 	}

// 	if nInnerNumber < 1 || 9 < nInnerNumber {
// 		return false, nil
// 	}

// 	if (paiNumber >= 5 || nInnerNumber > 5) && (paiNumber <= 5 || nInnerNumber < 5) {
// 		return false, nil
// 	}

// 	innerPai, err := base.NewPaiWithDetail(pai.Type(), uint8(nInnerNumber), false)
// 	if err != nil {
// 		return false, err
// 	}

// 	return prereachSutehaiSet.Has(innerPai)
// }

// func isSameTypeInPrereach(pai *base.Pai, n int, prereachSutehaiSet *base.PaiSet) (bool, error) {
// 	if pai.IsTsupai() {
// 		return false, nil
// 	}

// 	numbers := []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9}
// 	numSameType, err := core.Count(numbers, func(num uint8) (bool, error) {
// 		target, err := base.NewPaiWithDetail(pai.Type(), num, false)
// 		if err != nil {
// 			return false, err
// 		}
// 		return prereachSutehaiSet.Has(target)
// 	})
// 	if err != nil {
// 		return false, err
// 	}

// 	return numSameType+1 >= n, nil
// }

func isSujiOf(pai *base.Pai, targetPaiSet *base.PaiSet) (bool, error) {
	if pai.IsTsupai() {
		return false, nil
	}

	suji, err := getSuji(pai)
	if err != nil {
		return false, err
	}

	return core.AllMatch(suji, func(s base.Pai) (bool, error) {
		return targetPaiSet.Has(&s)
	})
}

func isWeakSujiOf(pai *base.Pai, targetPaiSet *base.PaiSet) (bool, error) {
	if pai.IsTsupai() {
		return false, nil
	}

	suji, err := getSuji(pai)
	if err != nil {
		return false, err
	}

	return core.AnyMatch(suji, func(s base.Pai) (bool, error) {
		return targetPaiSet.Has(&s)
	})
}

func getSuji(pai *base.Pai) ([]base.Pai, error) {
	if pai.IsTsupai() {
		return []base.Pai{}, nil
	}

	result := make([]base.Pai, 0, 2)
	paiNumber := pai.Number()
	candidates := []uint8{paiNumber - 3, paiNumber + 3}
	for _, n := range candidates {
		if 1 <= n && n <= 9 {
			sujiPai, err := base.NewPaiWithDetail(pai.Type(), n, false)
			if err != nil {
				return nil, err
			}
			result = append(result, *sujiPai)
		}
	}

	return result, nil
}

func isUrasujiOf(pai *base.Pai, targetPaiSet *base.PaiSet, anpaiSet *base.PaiSet) (bool, error) {
	sujis, err := getPossibleSujis(pai, anpaiSet)
	if err != nil {
		return false, err
	}

	return core.AnyMatch(sujis, func(s base.Pai) (bool, error) {
		if low := s.Next(-1); low != nil {
			hasLow, err := targetPaiSet.Has(low)
			if err != nil {
				return false, err
			}
			if hasLow {
				return true, nil
			}
		}

		if high := s.Next(4); high != nil {
			hasHigh, err := targetPaiSet.Has(high)
			if err != nil {
				return false, err
			}
			if hasHigh {
				return true, nil
			}
		}

		return false, nil
	})
}

// // Senkisuji (疝気筋) : Urasuji (裏筋) of urasuji
// func isSenkisujiOf(pai *base.Pai, targetPaiSet *base.PaiSet, anpaiSet *base.PaiSet) (bool, error) {
// 	sujis, err := getPossibleSujis(pai, anpaiSet)
// 	if err != nil {
// 		return false, err
// 	}

// 	return core.AnyMatch(sujis, func(s game.Pai) (bool, error) {
// 		if low := s.Next(-2); low != nil {
// 			hasLow, err := targetPaiSet.Has(low)
// 			if err != nil {
// 				return false, err
// 			}
// 			if hasLow {
// 				return true, nil
// 			}
// 		}

// 		if high := s.Next(5); high != nil {
// 			hasHigh, err := targetPaiSet.Has(high)
// 			if err != nil {
// 				return false, err
// 			}
// 			if hasHigh {
// 				return true, nil
// 			}
// 		}

// 		return false, nil
// 	})
// }

func isMatagisujiOf(pai *base.Pai, targetPaiSet *base.PaiSet, anpaiSet *base.PaiSet) (bool, error) {
	sujis, err := getPossibleSujis(pai, anpaiSet)
	if err != nil {
		return false, err
	}

	return core.AnyMatch(sujis, func(s base.Pai) (bool, error) {
		if low := s.Next(1); low != nil {
			hasLow, err := targetPaiSet.Has(low)
			if err != nil {
				return false, err
			}
			if hasLow {
				return true, nil
			}
		}

		if high := s.Next(2); high != nil {
			hasHigh, err := targetPaiSet.Has(high)
			if err != nil {
				return false, err
			}
			if hasHigh {
				return true, nil
			}
		}

		return false, nil
	})
}

// Returns sujis which contain the given pai and is alive i.e. none of pais in the suji are anpai.
// Uses the first pai to represent the suji. e.g. 1p for 14p suji
func getPossibleSujis(pai *base.Pai, anpaiSet *base.PaiSet) ([]base.Pai, error) {
	if pai.IsTsupai() {
		return []base.Pai{}, nil
	}

	sujis := make([]base.Pai, 0, 2)
	paiNumber := pai.Number()
	candidates := []uint8{paiNumber - 3, paiNumber}

	for _, n := range candidates {
		isAlive, err := core.AllMatch([]uint8{n, n + 3}, func(m uint8) (bool, error) {
			if m < 1 || m > 9 {
				return false, nil
			}

			sujiPai, err := base.NewPaiWithDetail(pai.Type(), m, false)
			if err != nil {
				return false, err
			}

			isAnpai, err := anpaiSet.Has(sujiPai)
			if err != nil {
				return false, err
			}
			return !isAnpai, nil
		})
		if err != nil {
			return nil, err
		}

		if isAlive {
			sujiPai, err := base.NewPaiWithDetail(pai.Type(), n, false)
			if err != nil {
				return nil, err
			}
			sujis = append(sujis, *sujiPai)
		}
	}

	return sujis, nil
}

// func isOuter(pai *base.Pai, targetPaiSet *base.PaiSet) (bool, error) {
// 	if pai.IsTsupai() {
// 		return false, nil
// 	}

// 	paiNumber := pai.Number()
// 	if paiNumber == 5 {
// 		return false, nil
// 	}

// 	var innerNumbers []uint8
// 	if paiNumber < 5 {
// 		for i := paiNumber + 1; i < 6; i++ {
// 			innerNumbers = append(innerNumbers, i)
// 		}
// 	} else {
// 		for i := uint8(5); i < paiNumber; i++ {
// 			innerNumbers = append(innerNumbers, i)
// 		}
// 	}

// 	return core.AnyMatch(innerNumbers, func(n uint8) (bool, error) {
// 		innerPai, err := base.NewPaiWithDetail(pai.Type(), n, false)
// 		if err != nil {
// 			return false, err
// 		}
// 		has, err := targetPaiSet.Has(innerPai)
// 		if err != nil {
// 			return false, err
// 		}
// 		return has, nil
// 	})
// }

func registerEvaluators() *evaluators {
	ev := evaluators{}

	// ev["anpai"] = func(scene *Scene, pai *base.Pai) (bool, error) {
	// 	return isAnpai(pai, scene.anpaiSet)
	// }
	ev["tsupai"] = func(scene *Scene, pai *base.Pai) (bool, error) {
		return isTsupai(pai), nil
	}
	ev["suji"] = func(scene *Scene, pai *base.Pai) (bool, error) {
		return isSuji(pai, scene.anpaiSet)
	}
	ev["weak_suji"] = func(scene *Scene, pai *base.Pai) (bool, error) {
		return isWeakSuji(pai, scene.anpaiSet)
	}
	ev["reach_suji"] = func(scene *Scene, pai *base.Pai) (bool, error) {
		return isReachSuji(pai, scene.reachPaiSet)
	}
	ev["prereach_suji"] = func(scene *Scene, pai *base.Pai) (bool, error) {
		return isPrereachSuji(pai, scene.prereachSutehaiSet)
	}
	ev["urasuji"] = func(scene *Scene, pai *base.Pai) (bool, error) {
		return isUrasuji(pai, scene.prereachSutehaiSet, scene.anpaiSet)
	}
	ev["early_urasuji"] = func(scene *Scene, pai *base.Pai) (bool, error) {
		return isEarlyUrasuji(pai, scene.earlySutehaiSet, scene.anpaiSet)
	}
	ev["reach_urasuji"] = func(scene *Scene, pai *base.Pai) (bool, error) {
		return isReachUrasuji(pai, scene.reachPaiSet, scene.anpaiSet)
	}
	ev["urasuji_of_5"] = func(scene *Scene, pai *base.Pai) (bool, error) {
		return isUrasujiOf5(pai, scene.prereachSutehaiSet, scene.anpaiSet)
	}
	ev["aida4ken"] = func(scene *Scene, pai *base.Pai) (bool, error) {
		return isAida4ken(pai, scene.prereachSutehaiSet)
	}
	ev["matagisuji"] = func(scene *Scene, pai *base.Pai) (bool, error) {
		return isMatagisuji(pai, scene.prereachSutehaiSet, scene.anpaiSet)
	}
	ev["early_matagisuji"] = func(scene *Scene, pai *base.Pai) (bool, error) {
		return isEarlyMatagisuji(pai, scene.earlySutehaiSet, scene.anpaiSet)
	}
	// ev["late_matagisuji"] = func(scene *Scene, pai *base.Pai) (bool, error) {
	// 	return isLateMatagisuji(pai, scene.lateSutehaiSet, scene.anpaiSet)
	// }
	// ev["reach_matagisuji"] = func(scene *Scene, pai *base.Pai) (bool, error) {
	// 	return isReachMatagisuji(pai, scene.reachPaiSet, scene.anpaiSet)
	// }
	// ev["senkisuji"] = func(scene *Scene, pai *base.Pai) (bool, error) {
	// 	return isSenkisuji(pai, scene.prereachSutehaiSet, scene.anpaiSet)
	// }
	// ev["early_senkisuji"] = func(scene *Scene, pai *base.Pai) (bool, error) {
	// 	return isEarlySenkisuji(pai, scene.earlySutehaiSet, scene.anpaiSet)
	// }
	// ev["outer_prereach_sutehai"] = func(scene *Scene, pai *base.Pai) (bool, error) {
	// 	return isOuterPrereachSutehai(pai, scene.prereachSutehaiSet)
	// }
	// ev["outer_early_sutehai"] = func(scene *Scene, pai *base.Pai) (bool, error) {
	// 	return isOuterEarlySutehai(pai, scene.earlySutehaiSet)
	// }
	// ev["dora"] = func(scene *Scene, pai *base.Pai) (bool, error) {
	// 	return isDora(pai, scene.doraSet)
	// }
	// ev["dora_suji"] = func(scene *Scene, pai *base.Pai) (bool, error) {
	// 	return isDoraSuji(pai, scene.doraSet)
	// }
	// ev["dora_matagi"] = func(scene *Scene, pai *base.Pai) (bool, error) {
	// 	return isDoraMatagi(pai, scene.doraSet, scene.anpaiSet)
	// }
	// ev["fanpai"] = func(scene *Scene, pai *base.Pai) (bool, error) {
	// 	return isFanpai(pai, scene.gameState, scene.target), nil
	// }
	// ev["ryenfonpai"] = func(scene *Scene, pai *base.Pai) (bool, error) {
	// 	return isRyenfonpai(pai, scene.gameState, scene.target), nil
	// }
	// ev["sangenpai"] = func(scene *Scene, pai *base.Pai) (bool, error) {
	// 	return isSangenpai(pai), nil
	// }
	// ev["fonpai"] = func(scene *Scene, pai *base.Pai) (bool, error) {
	// 	return isFonpai(pai), nil
	// }
	// ev["bakaze"] = func(scene *Scene, pai *base.Pai) (bool, error) {
	// 	return isBakaze(pai, scene.bakaze), nil
	// }
	// ev["jikaze"] = func(scene *Scene, pai *base.Pai) (bool, error) {
	// 	return isJikaze(pai, scene.targetKaze), nil
	// }

	// for i := range 4 {
	// 	featureName := fmt.Sprintf("chances<=%d", i)
	// 	n := i
	// 	ev[featureName] = func(scene *Scene, pai *base.Pai) (bool, error) {
	// 		return isNChanceOrLess(pai, n, scene.visibleSet)
	// 	}
	// }

	// // Whether i tiles are visible from one's perspective.
	// // Includes one's own hand. Excludes the tile one is about to discard.
	// for i := 1; i < 4; i++ {
	// 	featureName := fmt.Sprintf("visible>=%d", i)
	// 	n := i
	// 	ev[featureName] = func(scene *Scene, pai *base.Pai) (bool, error) {
	// 		return isVisibleNOrMore(pai, n+1, scene.visibleSet)
	// 	}
	// }

	// // Among the Suji of that tile, whether one is visible no more than i copies.
	// // The tile itself should not be counted.
	// // In the case of 5p, this means "either 2p or 8p is visible no more than i copies,"
	// // not "the combined visibility of 2p and 8p is no more than i copies."
	// for i := range 4 {
	// 	featureName := fmt.Sprintf("suji_visible<=%d", i)
	// 	n := i
	// 	ev[featureName] = func(scene *Scene, pai *base.Pai) (bool, error) {
	// 		return isSujiVisible(pai, n, scene.visibleSet)
	// 	}
	// }

	// for i := uint8(2); i < 6; i++ {
	// 	featureName := fmt.Sprintf("%d<=n<=%d", i, 10-i)
	// 	n := i
	// 	ev[featureName] = func(scene *Scene, pai *base.Pai) (bool, error) {
	// 		return isNumNOrInner(pai, n), nil
	// 	}
	// }

	// for i := 2; i < 5; i++ {
	// 	featureName := fmt.Sprintf("in_tehais>=%d", i)
	// 	n := i
	// 	ev[featureName] = func(scene *Scene, pai *base.Pai) (bool, error) {
	// 		return isInTehais(pai, n, scene.tehaiSet)
	// 	}
	// }

	// // Among the Suji of that tile, whether one is held at least i copies.
	// // The tile itself should not be counted.
	// // In the case of 5p, this means "either 2p or 8p is held at least i copies,"
	// // not "the combined total of 2p and 8p is at least i copies."
	// for i := 1; i < 5; i++ {
	// 	featureName := fmt.Sprintf("suji_in_tehais>=%d", i)
	// 	n := i
	// 	ev[featureName] = func(scene *Scene, pai *base.Pai) (bool, error) {
	// 		return isSujiInTehais(pai, n, scene.tehaiSet)
	// 	}
	// }

	// for i := 1; i < 3; i++ {
	// 	for j := 1; j < (i*2 + 1); j++ {
	// 		featureName := fmt.Sprintf("+-%d_in_prereach_sutehais>=%d", i, j)
	// 		distance := i
	// 		threshold := j
	// 		ev[featureName] = func(scene *Scene, pai *base.Pai) (bool, error) {
	// 			return isNOrMoreOfNeighborsInPrereachSutehais(
	// 				pai,
	// 				threshold,
	// 				distance,
	// 				scene.prereachSutehaiSet,
	// 			)
	// 		}
	// 	}
	// }

	// for i := 1; i < 3; i++ {
	// 	featureName := fmt.Sprintf("%d_outer_prereach_sutehai", i)
	// 	n := i
	// 	ev[featureName] = func(scene *Scene, pai *base.Pai) (bool, error) {
	// 		return isNOuterPrereachSutehai(pai, n, scene.prereachSutehaiSet)
	// 	}
	// }

	// for i := 1; i < 3; i++ {
	// 	featureName := fmt.Sprintf("%d_inner_prereach_sutehai", i)
	// 	n := i
	// 	ev[featureName] = func(scene *Scene, pai *base.Pai) (bool, error) {
	// 		return isNOuterPrereachSutehai(pai, -n, scene.prereachSutehaiSet)
	// 	}
	// }

	// for i := 1; i < 9; i++ {
	// 	featureName := fmt.Sprintf("same_type_in_prereach>=%d", i)
	// 	n := i
	// 	ev[featureName] = func(scene *Scene, pai *base.Pai) (bool, error) {
	// 		return isSameTypeInPrereach(pai, n, scene.prereachSutehaiSet)
	// 	}
	// }

	return &ev
}
