package main

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/ai/core"
	"github.com/Apricot-S/mjai-manue-go/internal/game"
)

type evaluator func(*Scene, *game.Pai) (bool, error)
type evaluators map[string]evaluator

var defaultEvaluators = registerEvaluators()

type Scene struct {
	// gameState game.StateViewer
	// me        *game.Player
	// target    *game.Player

	// tehaiSet   *game.PaiSet
	anpaiSet *game.PaiSet
	// visibleSet *game.PaiSet
	// doraSet    *game.PaiSet
	// bakaze     *game.Pai
	// targetKaze *game.Pai

	prereachSutehaiSet *game.PaiSet
	// earlySutehaiSet    *game.PaiSet
	// lateSutehaiSet     *game.PaiSet
	reachPaiSet *game.PaiSet

	evaluators *evaluators
}

func NewScene(gameState game.StateViewer, me *game.Player, target *game.Player) (*Scene, error) {
	s := &Scene{
		// gameState:  gameState,
		// me:         me,
		// target:     target,
		evaluators: defaultEvaluators,
	}

	var err error
	// if s.tehaiSet, err = game.NewPaiSet(me.Tehais()); err != nil {
	// 	return nil, err
	// }
	if s.anpaiSet, err = game.NewPaiSet(gameState.Anpais(target)); err != nil {
		return nil, err
	}
	// if s.visibleSet, err = game.NewPaiSet(gameState.VisiblePais(me)); err != nil {
	// 	return nil, err
	// }
	// if s.doraSet, err = game.NewPaiSet(gameState.Doras()); err != nil {
	// 	return nil, err
	// }

	// s.bakaze = gameState.Bakaze()
	// s.targetKaze = gameState.Jikaze(target)

	var prereachSutehais game.Pais = nil
	var reachPais game.Pais = nil
	if idx := target.ReachSutehaiIndex(); idx != -1 {
		sutehais := target.Sutehais()
		prereachSutehais = sutehais[:idx+1]
		reachPai := sutehais[idx]
		reachPais = game.Pais{reachPai}
	}
	if s.prereachSutehaiSet, err = game.NewPaiSet(prereachSutehais); err != nil {
		return nil, err
	}
	if s.reachPaiSet, err = game.NewPaiSet(reachPais); err != nil {
		return nil, err
	}

	// halfLen := len(prereachSutehais) / 2
	// if s.earlySutehaiSet, err = game.NewPaiSet(prereachSutehais[:halfLen]); err != nil {
	// 	return nil, err
	// }
	// if s.lateSutehaiSet, err = game.NewPaiSet(prereachSutehais[halfLen:]); err != nil {
	// 	return nil, err
	// }

	return s, nil
}

func (s *Scene) Evaluate(name string, pai *game.Pai) (bool, error) {
	if evaluator, ok := (*s.evaluators)[name]; ok {
		return evaluator(s, pai)
	}
	return false, fmt.Errorf("an unknown feature name was specified: %v", name)
}

// func isAnpai(pai *game.Pai, anpaiSet *game.PaiSet) (bool, error) {
// 	return anpaiSet.Has(pai)
// }

func isTsupai(pai *game.Pai) bool {
	return pai.IsTsupai()
}

// Omotesuji (表筋) or Nakasuji (中筋)
func isSuji(pai *game.Pai, anpaiSet *game.PaiSet) (bool, error) {
	return isSujiOf(pai, anpaiSet)
}

// Katasuji (片筋) or Suji (筋)
func isWeakSuji(pai *game.Pai, anpaiSet *game.PaiSet) (bool, error) {
	return isWeakSujiOf(pai, anpaiSet)
}

// Suji for Riichi declaration tile. Including tiles like 4p against 1p Riichi.
func isReachSuji(pai *game.Pai, reachPaiSet *game.PaiSet) (bool, error) {
	return isWeakSujiOf(pai, reachPaiSet)
}

// func isPrereachSuji(pai *game.Pai, prereachSutehaiSet *game.PaiSet) (bool, error) {
// 	return isSujiOf(pai, prereachSutehaiSet)
// }

// // Urasuji (裏筋)
// // http://ja.wikipedia.org/wiki/%E7%AD%8B_(%E9%BA%BB%E9%9B%80)#.E8.A3.8F.E3.82.B9.E3.82.B8
// func isUrasuji(pai *game.Pai, prereachSutehaiSet *game.PaiSet, anpaiSet *game.PaiSet) (bool, error) {
// 	return isUrasujiOf(pai, prereachSutehaiSet, anpaiSet)
// }

// func isEarlyUrasuji(pai *game.Pai, earlySutehaiSet *game.PaiSet, anpaiSet *game.PaiSet) (bool, error) {
// 	return isUrasujiOf(pai, earlySutehaiSet, anpaiSet)
// }

// func isReachUrasuji(pai *game.Pai, reachPaiSet *game.PaiSet, anpaiSet *game.PaiSet) (bool, error) {
// 	return isUrasujiOf(pai, reachPaiSet, anpaiSet)
// }

// // Aidayonken (間四間)
// // http://ja.wikipedia.org/wiki/%E7%AD%8B_(%E9%BA%BB%E9%9B%80)#.E9.96.93.E5.9B.9B.E9.96.93
// func isAida4ken(pai *game.Pai, prereachSutehaiSet *game.PaiSet) (bool, error) {
// 	if pai.IsTsupai() {
// 		return false, nil
// 	}

// 	num := pai.Number()
// 	typ := pai.Type()

// 	if 2 <= num && num <= 5 {
// 		low, err := game.NewPaiWithDetail(typ, num-1, false)
// 		if err != nil {
// 			return false, err
// 		}
// 		hasLow, err := prereachSutehaiSet.Has(low)
// 		if err != nil {
// 			return false, err
// 		}

// 		high, err := game.NewPaiWithDetail(typ, num+4, false)
// 		if err != nil {
// 			return false, err
// 		}
// 		hasHigh, err := prereachSutehaiSet.Has(high)
// 		if err != nil {
// 			return false, err
// 		}

// 		if hasLow && hasHigh {
// 			return true, nil
// 		}
// 	}

// 	if 5 <= num && num <= 8 {
// 		low, err := game.NewPaiWithDetail(typ, num-4, false)
// 		if err != nil {
// 			return false, err
// 		}
// 		hasLow, err := prereachSutehaiSet.Has(low)
// 		if err != nil {
// 			return false, err
// 		}

// 		high, err := game.NewPaiWithDetail(typ, num+1, false)
// 		if err != nil {
// 			return false, err
// 		}
// 		hasHigh, err := prereachSutehaiSet.Has(high)
// 		if err != nil {
// 			return false, err
// 		}

// 		if hasLow && hasHigh {
// 			return true, nil
// 		}
// 	}

// 	return false, nil
// }

// // Matagisuji (跨ぎ筋)
// // http://ja.wikipedia.org/wiki/%E7%AD%8B_(%E9%BA%BB%E9%9B%80)#.E3.81.BE.E3.81.9F.E3.81.8E.E3.82.B9.E3.82.B8
// func isMatagisuji(pai *game.Pai, prereachSutehaiSet *game.PaiSet, anpaiSet *game.PaiSet) (bool, error) {
// 	return isMatagisujiOf(pai, prereachSutehaiSet, anpaiSet)
// }

// func isEarlyMatagisuji(pai *game.Pai, earlySutehaiSet *game.PaiSet, anpaiSet *game.PaiSet) (bool, error) {
// 	return isMatagisujiOf(pai, earlySutehaiSet, anpaiSet)
// }

// func isLateMatagisuji(pai *game.Pai, lateSutehaiSet *game.PaiSet, anpaiSet *game.PaiSet) (bool, error) {
// 	return isMatagisujiOf(pai, lateSutehaiSet, anpaiSet)
// }

// func isReachMatagisuji(pai *game.Pai, reachPaiSet *game.PaiSet, anpaiSet *game.PaiSet) (bool, error) {
// 	return isMatagisujiOf(pai, reachPaiSet, anpaiSet)
// }

// // Senkisuji (疝気筋)
// // # http://ja.wikipedia.org/wiki/%E7%AD%8B_(%E9%BA%BB%E9%9B%80)#.E7.96.9D.E6.B0.97.E3.82.B9.E3.82.B8
// func isSenkisuji(pai *game.Pai, prereachSutehaiSet *game.PaiSet, anpaiSet *game.PaiSet) (bool, error) {
// 	return isSenkisujiOf(pai, prereachSutehaiSet, anpaiSet)
// }

// func isEarlySenkisuji(pai *game.Pai, earlySutehaiSet *game.PaiSet, anpaiSet *game.PaiSet) (bool, error) {
// 	return isSenkisujiOf(pai, earlySutehaiSet, anpaiSet)
// }

// func isOuterPrereachSutehai(pai *game.Pai, prereachSutehaiSet *game.PaiSet) (bool, error) {
// 	return isOuter(pai, prereachSutehaiSet)
// }

// func isOuterEarlySutehai(pai *game.Pai, earlySutehaiSet *game.PaiSet) (bool, error) {
// 	return isOuter(pai, earlySutehaiSet)
// }

// func isDora(pai *game.Pai, doraSet *game.PaiSet) (bool, error) {
// 	return doraSet.Has(pai)
// }

// func isDoraSuji(pai *game.Pai, doraSet *game.PaiSet) (bool, error) {
// 	return isWeakSujiOf(pai, doraSet)
// }

// func isDoraMatagi(pai *game.Pai, doraSet *game.PaiSet, anpaiSet *game.PaiSet) (bool, error) {
// 	return isMatagisujiOf(pai, doraSet, anpaiSet)
// }

// func isFanpai(pai *game.Pai, gameState game.StateViewer, target *game.Player) bool {
// 	return gameState.YakuhaiFan(pai, target) >= 1
// }

// func isRyenfonpai(pai *game.Pai, gameState game.StateViewer, target *game.Player) bool {
// 	return gameState.YakuhaiFan(pai, target) >= 2
// }

// func isSangenpai(pai *game.Pai) bool {
// 	return pai.IsTsupai() && pai.Number() >= 5
// }

// func isFonpai(pai *game.Pai) bool {
// 	return pai.IsTsupai() && pai.Number() < 5
// }

// func isBakaze(pai *game.Pai, bakaze *game.Pai) bool {
// 	return pai.HasSameSymbol(bakaze)
// }

// func isJikaze(pai *game.Pai, targetKaze *game.Pai) bool {
// 	return pai.HasSameSymbol(targetKaze)
// }

// func isNChanceOrLess(pai *game.Pai, n int, visibleSet *game.PaiSet) (bool, error) {
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
// 		kabePai, err := game.NewPaiWithDetail(pai.Type(), num, false)
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

// func isVisibleNOrMore(pai *game.Pai, n int, visibleSet *game.PaiSet) (bool, error) {
// 	c, err := visibleSet.Count(pai)
// 	if err != nil {
// 		return false, err
// 	}
// 	return c >= n, nil
// }

// func isSujiVisible(pai *game.Pai, n int, visibleSet *game.PaiSet) (bool, error) {
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

// func isNumNOrInner(pai *game.Pai, n uint8) bool {
// 	if pai.IsTsupai() {
// 		return false
// 	}

// 	paiNumber := pai.Number()
// 	if n <= paiNumber && paiNumber <= 10-n {
// 		return true
// 	}

// 	return false
// }

// func isInTehais(pai *game.Pai, n int, tehaiSet *game.PaiSet) (bool, error) {
// 	c, err := tehaiSet.Count(pai)
// 	return c >= n, err
// }

// func isSujiInTehais(pai *game.Pai, n int, tehaiSet *game.PaiSet) (bool, error) {
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
// 	pai *game.Pai,
// 	n int,
// 	neighborDistance int,
// 	prereachSutehaiSet *game.PaiSet,
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

// 		neighborPai, err := game.NewPaiWithDetail(pai.Type(), uint8(num), false)
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
// func isNOuterPrereachSutehai(pai *game.Pai, n int, prereachSutehaiSet *game.PaiSet) (bool, error) {
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

// 	innerPai, err := game.NewPaiWithDetail(pai.Type(), uint8(nInnerNumber), false)
// 	if err != nil {
// 		return false, err
// 	}

// 	return prereachSutehaiSet.Has(innerPai)
// }

// func isSameTypeInPrereach(pai *game.Pai, n int, prereachSutehaiSet *game.PaiSet) (bool, error) {
// 	if pai.IsTsupai() {
// 		return false, nil
// 	}

// 	numbers := []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9}
// 	numSameType, err := core.Count(numbers, func(num uint8) (bool, error) {
// 		target, err := game.NewPaiWithDetail(pai.Type(), num, false)
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

func isSujiOf(pai *game.Pai, targetPaiSet *game.PaiSet) (bool, error) {
	if pai.IsTsupai() {
		return false, nil
	}

	suji, err := getSuji(pai)
	if err != nil {
		return false, err
	}

	return core.AllMatch(suji, func(s game.Pai) (bool, error) {
		return targetPaiSet.Has(&s)
	})
}

func isWeakSujiOf(pai *game.Pai, targetPaiSet *game.PaiSet) (bool, error) {
	if pai.IsTsupai() {
		return false, nil
	}

	suji, err := getSuji(pai)
	if err != nil {
		return false, err
	}

	return core.AnyMatch(suji, func(s game.Pai) (bool, error) {
		return targetPaiSet.Has(&s)
	})
}

func getSuji(pai *game.Pai) ([]game.Pai, error) {
	if pai.IsTsupai() {
		return []game.Pai{}, nil
	}

	result := make([]game.Pai, 0, 2)
	paiNumber := pai.Number()
	candidates := []uint8{paiNumber - 3, paiNumber + 3}
	for _, n := range candidates {
		if 1 <= n && n <= 9 {
			sujiPai, err := game.NewPaiWithDetail(pai.Type(), n, false)
			if err != nil {
				return nil, err
			}
			result = append(result, *sujiPai)
		}
	}

	return result, nil
}

// func isUrasujiOf(pai *game.Pai, targetPaiSet *game.PaiSet, anpaiSet *game.PaiSet) (bool, error) {
// 	sujis, err := getPossibleSujis(pai, anpaiSet)
// 	if err != nil {
// 		return false, err
// 	}

// 	return core.AnyMatch(sujis, func(s game.Pai) (bool, error) {
// 		if low := s.Next(-1); low != nil {
// 			hasLow, err := targetPaiSet.Has(low)
// 			if err != nil {
// 				return false, err
// 			}
// 			if hasLow {
// 				return true, nil
// 			}
// 		}

// 		if high := s.Next(4); high != nil {
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

// // Senkisuji (疝気筋) : Urasuji (裏筋) of urasuji
// func isSenkisujiOf(pai *game.Pai, targetPaiSet *game.PaiSet, anpaiSet *game.PaiSet) (bool, error) {
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

// func isMatagisujiOf(pai *game.Pai, targetPaiSet *game.PaiSet, anpaiSet *game.PaiSet) (bool, error) {
// 	sujis, err := getPossibleSujis(pai, anpaiSet)
// 	if err != nil {
// 		return false, err
// 	}

// 	return core.AnyMatch(sujis, func(s game.Pai) (bool, error) {
// 		if low := s.Next(1); low != nil {
// 			hasLow, err := targetPaiSet.Has(low)
// 			if err != nil {
// 				return false, err
// 			}
// 			if hasLow {
// 				return true, nil
// 			}
// 		}

// 		if high := s.Next(2); high != nil {
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

// // Returns sujis which contain the given pai and is alive i.e. none of pais in the suji are anpai.
// // Uses the first pai to represent the suji. e.g. 1p for 14p suji
// func getPossibleSujis(pai *game.Pai, anpaiSet *game.PaiSet) ([]game.Pai, error) {
// 	if pai.IsTsupai() {
// 		return []game.Pai{}, nil
// 	}

// 	sujis := make([]game.Pai, 0, 2)
// 	paiNumber := pai.Number()
// 	candidates := []uint8{paiNumber - 3, paiNumber}

// 	for _, n := range candidates {
// 		isAlive, err := core.AllMatch([]uint8{n, n + 3}, func(m uint8) (bool, error) {
// 			if m < 1 || m > 9 {
// 				return false, nil
// 			}

// 			sujiPai, err := game.NewPaiWithDetail(pai.Type(), m, false)
// 			if err != nil {
// 				return false, err
// 			}

// 			isAnpai, err := anpaiSet.Has(sujiPai)
// 			if err != nil {
// 				return false, err
// 			}
// 			return !isAnpai, nil
// 		})
// 		if err != nil {
// 			return nil, err
// 		}

// 		if isAlive {
// 			sujiPai, err := game.NewPaiWithDetail(pai.Type(), n, false)
// 			if err != nil {
// 				return nil, err
// 			}
// 			sujis = append(sujis, *sujiPai)
// 		}
// 	}

// 	return sujis, nil
// }

// func isOuter(pai *game.Pai, targetPaiSet *game.PaiSet) (bool, error) {
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
// 		innerPai, err := game.NewPaiWithDetail(pai.Type(), n, false)
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

	// ev["anpai"] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 	return isAnpai(pai, scene.anpaiSet)
	// }
	ev["tsupai"] = func(scene *Scene, pai *game.Pai) (bool, error) {
		return isTsupai(pai), nil
	}
	ev["suji"] = func(scene *Scene, pai *game.Pai) (bool, error) {
		return isSuji(pai, scene.anpaiSet)
	}
	ev["weak_suji"] = func(scene *Scene, pai *game.Pai) (bool, error) {
		return isWeakSuji(pai, scene.anpaiSet)
	}
	ev["reach_suji"] = func(scene *Scene, pai *game.Pai) (bool, error) {
		return isReachSuji(pai, scene.reachPaiSet)
	}
	// ev["prereach_suji"] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 	return isPrereachSuji(pai, scene.prereachSutehaiSet)
	// }
	// ev["urasuji"] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 	return isUrasuji(pai, scene.prereachSutehaiSet, scene.anpaiSet)
	// }
	// ev["early_urasuji"] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 	return isEarlyUrasuji(pai, scene.earlySutehaiSet, scene.anpaiSet)
	// }
	// ev["reach_urasuji"] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 	return isReachUrasuji(pai, scene.reachPaiSet, scene.anpaiSet)
	// }
	// ev["aida4ken"] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 	return isAida4ken(pai, scene.prereachSutehaiSet)
	// }
	// ev["matagisuji"] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 	return isMatagisuji(pai, scene.prereachSutehaiSet, scene.anpaiSet)
	// }
	// ev["early_matagisuji"] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 	return isEarlyMatagisuji(pai, scene.earlySutehaiSet, scene.anpaiSet)
	// }
	// ev["late_matagisuji"] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 	return isLateMatagisuji(pai, scene.lateSutehaiSet, scene.anpaiSet)
	// }
	// ev["reach_matagisuji"] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 	return isReachMatagisuji(pai, scene.reachPaiSet, scene.anpaiSet)
	// }
	// ev["senkisuji"] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 	return isSenkisuji(pai, scene.prereachSutehaiSet, scene.anpaiSet)
	// }
	// ev["early_senkisuji"] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 	return isEarlySenkisuji(pai, scene.earlySutehaiSet, scene.anpaiSet)
	// }
	// ev["outer_prereach_sutehai"] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 	return isOuterPrereachSutehai(pai, scene.prereachSutehaiSet)
	// }
	// ev["outer_early_sutehai"] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 	return isOuterEarlySutehai(pai, scene.earlySutehaiSet)
	// }
	// ev["dora"] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 	return isDora(pai, scene.doraSet)
	// }
	// ev["dora_suji"] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 	return isDoraSuji(pai, scene.doraSet)
	// }
	// ev["dora_matagi"] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 	return isDoraMatagi(pai, scene.doraSet, scene.anpaiSet)
	// }
	// ev["fanpai"] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 	return isFanpai(pai, scene.gameState, scene.target), nil
	// }
	// ev["ryenfonpai"] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 	return isRyenfonpai(pai, scene.gameState, scene.target), nil
	// }
	// ev["sangenpai"] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 	return isSangenpai(pai), nil
	// }
	// ev["fonpai"] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 	return isFonpai(pai), nil
	// }
	// ev["bakaze"] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 	return isBakaze(pai, scene.bakaze), nil
	// }
	// ev["jikaze"] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 	return isJikaze(pai, scene.targetKaze), nil
	// }

	// for i := range 4 {
	// 	featureName := fmt.Sprintf("chances<=%d", i)
	// 	n := i
	// 	ev[featureName] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 		return isNChanceOrLess(pai, n, scene.visibleSet)
	// 	}
	// }

	// // Whether i tiles are visible from one's perspective.
	// // Includes one's own hand. Excludes the tile one is about to discard.
	// for i := 1; i < 4; i++ {
	// 	featureName := fmt.Sprintf("visible>=%d", i)
	// 	n := i
	// 	ev[featureName] = func(scene *Scene, pai *game.Pai) (bool, error) {
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
	// 	ev[featureName] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 		return isSujiVisible(pai, n, scene.visibleSet)
	// 	}
	// }

	// for i := uint8(2); i < 6; i++ {
	// 	featureName := fmt.Sprintf("%d<=n<=%d", i, 10-i)
	// 	n := i
	// 	ev[featureName] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 		return isNumNOrInner(pai, n), nil
	// 	}
	// }

	// for i := 2; i < 5; i++ {
	// 	featureName := fmt.Sprintf("in_tehais>=%d", i)
	// 	n := i
	// 	ev[featureName] = func(scene *Scene, pai *game.Pai) (bool, error) {
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
	// 	ev[featureName] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 		return isSujiInTehais(pai, n, scene.tehaiSet)
	// 	}
	// }

	// for i := 1; i < 3; i++ {
	// 	for j := 1; j < (i*2 + 1); j++ {
	// 		featureName := fmt.Sprintf("+-%d_in_prereach_sutehais>=%d", i, j)
	// 		distance := i
	// 		threshold := j
	// 		ev[featureName] = func(scene *Scene, pai *game.Pai) (bool, error) {
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
	// 	ev[featureName] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 		return isNOuterPrereachSutehai(pai, n, scene.prereachSutehaiSet)
	// 	}
	// }

	// for i := 1; i < 3; i++ {
	// 	featureName := fmt.Sprintf("%d_inner_prereach_sutehai", i)
	// 	n := i
	// 	ev[featureName] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 		return isNOuterPrereachSutehai(pai, -n, scene.prereachSutehaiSet)
	// 	}
	// }

	// for i := 1; i < 9; i++ {
	// 	featureName := fmt.Sprintf("same_type_in_prereach>=%d", i)
	// 	n := i
	// 	ev[featureName] = func(scene *Scene, pai *game.Pai) (bool, error) {
	// 		return isSameTypeInPrereach(pai, n, scene.prereachSutehaiSet)
	// 	}
	// }

	return &ev
}
