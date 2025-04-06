package ai

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/configs"
	"github.com/Apricot-S/mjai-manue-go/internal/game"
)

type evaluator func(*Scene, *game.Pai) (bool, error)
type evaluators map[string]evaluator

var defaultEvaluators = registerEvaluators()

type Scene struct {
	gameState game.State
	me        *game.Player
	target    *game.Player

	tehaiSet   *game.PaiSet
	anpaiSet   *game.PaiSet
	visibleSet *game.PaiSet
	doraSet    *game.PaiSet
	bakaze     *game.Pai
	targetKaze *game.Pai

	prereachSutehaiSet *game.PaiSet
	earlySutehaiSet    *game.PaiSet
	lateSutehaiSet     *game.PaiSet
	reachPaiSet        *game.PaiSet

	evaluators *evaluators
}

func NewScene(gameState game.State, me *game.Player, target *game.Player) (*Scene, error) {
	s := &Scene{
		gameState:  gameState,
		me:         me,
		target:     target,
		evaluators: &defaultEvaluators,
	}

	var err error
	if s.tehaiSet, err = game.NewPaiSetWithPais(me.Tehais()); err != nil {
		return nil, err
	}
	if s.anpaiSet, err = game.NewPaiSetWithPais(gameState.Anpais(target)); err != nil {
		return nil, err
	}
	if s.visibleSet, err = game.NewPaiSetWithPais(gameState.VisiblePais(me)); err != nil {
		return nil, err
	}
	if s.doraSet, err = game.NewPaiSetWithPais(gameState.Doras()); err != nil {
		return nil, err
	}

	s.bakaze = gameState.Bakaze()
	s.targetKaze = gameState.Jikaze(target)

	var prereachSutehais game.Pais = nil
	var reachPais game.Pais = nil
	if idx := target.ReachSutehaiIndex(); idx != -1 {
		sutehais := target.Sutehais()
		prereachSutehais = sutehais[:idx+1]
		reachPai := sutehais[idx]
		reachPais = game.Pais{reachPai}
	}
	if s.prereachSutehaiSet, err = game.NewPaiSetWithPais(prereachSutehais); err != nil {
		return nil, err
	}
	if s.reachPaiSet, err = game.NewPaiSetWithPais(reachPais); err != nil {
		return nil, err
	}

	halfLen := len(prereachSutehais) / 2
	if s.earlySutehaiSet, err = game.NewPaiSetWithPais(prereachSutehais[:halfLen]); err != nil {
		return nil, err
	}
	if s.lateSutehaiSet, err = game.NewPaiSetWithPais(prereachSutehais[halfLen:]); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Scene) Evaluate(name string, pai *game.Pai) (bool, error) {
	if evaluator, ok := (*s.evaluators)[name]; ok {
		return evaluator(s, pai)
	}

	switch name {
	case "anpai":
		return s.isAnpai(pai)
	case "tsupai":
		return s.isTsupai(pai), nil
	case "suji":
		return s.isSuji(pai)
	case "weak_suji":
		return s.isWeakSuji(pai)
	case "reach_suji":
		return s.isReachSuji(pai)
	case "prereach_suji":
		return s.isPrereachSuji(pai)
	case "urasuji":
		return s.isUrasuji(pai)
	case "early_urasuji":
		return s.isEarlyUrasuji(pai)
	case "reach_urasuji":
		return s.isReachUrasuji(pai)
	case "aida4ken":
		return s.isAida4ken(pai)
	case "matagisuji":
		return s.isMatagisuji(pai)
	case "early_matagisuji":
		return s.isEarlyMatagisuji(pai)
	case "late_matagisuji":
		return s.isLateMatagisuji(pai)
	case "reach_matagisuji":
		return s.isReachMatagisuji(pai)
	case "senkisuji":
		return s.isSenkisuji(pai)
	case "early_senkisuji":
		return s.isEarlySenkisuji(pai)
	case "outer_prereach_sutehai":
		return s.isOuterPrereachSutehai(pai)
	case "outer_early_sutehai":
		return s.isOuterEarlySutehai(pai)
	case "dora":
		return s.isDora(pai)
	case "dora_suji":
		return s.isDoraSuji(pai)
	case "dora_matagi":
		return s.isDoraMatagi(pai)
	case "fanpai":
		return s.isFanpai(pai), nil
	case "ryenfonpai":
		return s.isRyenfonpai(pai), nil
	case "sangenpai":
		return s.isSangenpai(pai), nil
	case "fonpai":
		return s.isFonpai(pai), nil
	case "bakaze":
		return s.isBakaze(pai), nil
	case "jikaze":
		return s.isJikaze(pai), nil
	default:
		return false, fmt.Errorf("an unknown feature name was specified: %v", name)
	}
}

func (s *Scene) isAnpai(pai *game.Pai) (bool, error) {
	return s.anpaiSet.Has(pai)
}

func (s *Scene) isTsupai(pai *game.Pai) bool {
	return pai.IsTsupai()
}

// Omotesuji (表筋) or Nakasuji (中筋)
func (s *Scene) isSuji(pai *game.Pai) (bool, error) {
	return isSujiOf(pai, s.anpaiSet)
}

// Katasuji (片筋) or Suji (筋)
func (s *Scene) isWeakSuji(pai *game.Pai) (bool, error) {
	return isWeakSujiOf(pai, s.anpaiSet)
}

// Suji for Riichi declaration tile. Including tiles like 4p against 1p Riichi.
func (s *Scene) isReachSuji(pai *game.Pai) (bool, error) {
	return isWeakSujiOf(pai, s.reachPaiSet)
}

func (s *Scene) isPrereachSuji(pai *game.Pai) (bool, error) {
	return isSujiOf(pai, s.prereachSutehaiSet)
}

// Urasuji (裏筋)
// http://ja.wikipedia.org/wiki/%E7%AD%8B_(%E9%BA%BB%E9%9B%80)#.E8.A3.8F.E3.82.B9.E3.82.B8
func (s *Scene) isUrasuji(pai *game.Pai) (bool, error) {
	return isUrasujiOf(pai, s.prereachSutehaiSet, s.anpaiSet)
}

func (s *Scene) isEarlyUrasuji(pai *game.Pai) (bool, error) {
	return isUrasujiOf(pai, s.earlySutehaiSet, s.anpaiSet)
}

func (s *Scene) isReachUrasuji(pai *game.Pai) (bool, error) {
	return isUrasujiOf(pai, s.reachPaiSet, s.anpaiSet)
}

// Aidayonken (間四間)
// http://ja.wikipedia.org/wiki/%E7%AD%8B_(%E9%BA%BB%E9%9B%80)#.E9.96.93.E5.9B.9B.E9.96.93
func (s *Scene) isAida4ken(pai *game.Pai) (bool, error) {
	if pai.IsTsupai() {
		return false, nil
	}

	num := pai.Number()
	typ := pai.Type()

	if 2 <= num && num <= 5 {
		low, err := game.NewPaiWithDetail(typ, num-1, false)
		if err != nil {
			return false, err
		}
		hasLow, err := s.prereachSutehaiSet.Has(low)
		if err != nil {
			return false, err
		}

		high, err := game.NewPaiWithDetail(typ, num+4, false)
		if err != nil {
			return false, err
		}
		hasHigh, err := s.prereachSutehaiSet.Has(high)
		if err != nil {
			return false, err
		}

		if hasLow && hasHigh {
			return true, nil
		}
	}

	if 5 <= num && num <= 8 {
		low, err := game.NewPaiWithDetail(typ, num-4, false)
		if err != nil {
			return false, err
		}
		hasLow, err := s.prereachSutehaiSet.Has(low)
		if err != nil {
			return false, err
		}

		high, err := game.NewPaiWithDetail(typ, num+1, false)
		if err != nil {
			return false, err
		}
		hasHigh, err := s.prereachSutehaiSet.Has(high)
		if err != nil {
			return false, err
		}

		if hasLow && hasHigh {
			return true, nil
		}
	}

	return false, nil
}

// Matagisuji (跨ぎ筋)
// http://ja.wikipedia.org/wiki/%E7%AD%8B_(%E9%BA%BB%E9%9B%80)#.E3.81.BE.E3.81.9F.E3.81.8E.E3.82.B9.E3.82.B8
func (s *Scene) isMatagisuji(pai *game.Pai) (bool, error) {
	return isMatagisujiOf(pai, s.prereachSutehaiSet, s.anpaiSet)
}

func (s *Scene) isEarlyMatagisuji(pai *game.Pai) (bool, error) {
	return isMatagisujiOf(pai, s.earlySutehaiSet, s.anpaiSet)
}

func (s *Scene) isLateMatagisuji(pai *game.Pai) (bool, error) {
	return isMatagisujiOf(pai, s.lateSutehaiSet, s.anpaiSet)
}

func (s *Scene) isReachMatagisuji(pai *game.Pai) (bool, error) {
	return isMatagisujiOf(pai, s.reachPaiSet, s.anpaiSet)
}

// Senkisuji (疝気筋)
// # http://ja.wikipedia.org/wiki/%E7%AD%8B_(%E9%BA%BB%E9%9B%80)#.E7.96.9D.E6.B0.97.E3.82.B9.E3.82.B8
func (s *Scene) isSenkisuji(pai *game.Pai) (bool, error) {
	return isSenkisujiOf(pai, s.prereachSutehaiSet, s.anpaiSet)
}

func (s *Scene) isEarlySenkisuji(pai *game.Pai) (bool, error) {
	return isSenkisujiOf(pai, s.earlySutehaiSet, s.anpaiSet)
}

func (s *Scene) isOuterPrereachSutehai(pai *game.Pai) (bool, error) {
	return isOuter(pai, s.prereachSutehaiSet)
}

func (s *Scene) isOuterEarlySutehai(pai *game.Pai) (bool, error) {
	return isOuter(pai, s.earlySutehaiSet)
}

func (s *Scene) isDora(pai *game.Pai) (bool, error) {
	return s.doraSet.Has(pai)
}

func (s *Scene) isDoraSuji(pai *game.Pai) (bool, error) {
	return isWeakSujiOf(pai, s.doraSet)
}

func (s *Scene) isDoraMatagi(pai *game.Pai) (bool, error) {
	return isMatagisujiOf(pai, s.doraSet, s.anpaiSet)
}

func (s *Scene) isFanpai(pai *game.Pai) bool {
	return s.gameState.YakuhaiFan(pai, s.target) >= 1
}

func (s *Scene) isRyenfonpai(pai *game.Pai) bool {
	return s.gameState.YakuhaiFan(pai, s.target) >= 2
}

func (s *Scene) isSangenpai(pai *game.Pai) bool {
	return pai.IsTsupai() && pai.Number() >= 5
}

func (s *Scene) isFonpai(pai *game.Pai) bool {
	return pai.IsTsupai() && pai.Number() < 5
}

func (s *Scene) isBakaze(pai *game.Pai) bool {
	return pai.HasSameSymbol(s.bakaze)
}

func (s *Scene) isJikaze(pai *game.Pai) bool {
	return pai.HasSameSymbol(s.targetKaze)
}

// n can be negative.
func isNOuterPrereachSutehai(pai *game.Pai, n int, prereachSutehaiSet *game.PaiSet) (bool, error) {
	if pai.IsTsupai() {
		return false, nil
	}

	paiNumber := int(pai.Number())
	if paiNumber == 5 {
		return false, nil
	}

	nInnerNumber := 0
	if paiNumber < 5 {
		nInnerNumber = paiNumber + n
	} else {
		nInnerNumber = paiNumber - n
	}

	if nInnerNumber < 1 || 9 < nInnerNumber {
		return false, nil
	}

	if (paiNumber >= 5 || nInnerNumber > 5) && (paiNumber <= 5 || nInnerNumber < 5) {
		return false, nil
	}

	innerPai, err := game.NewPaiWithDetail(pai.Type(), uint8(nInnerNumber), false)
	if err != nil {
		return false, err
	}

	return prereachSutehaiSet.Has(innerPai)
}

func isNOrMoreOfNeighborsInPrereachSutehais(pai *game.Pai, n int, neighborDistance int, prereachSutehaiSet *game.PaiSet) (bool, error) {
	if pai.IsTsupai() {
		return false, nil
	}

	paiNumber := int(pai.Number())
	numbers := make([]int, 0, 2*neighborDistance+1)
	for i := -neighborDistance; i <= neighborDistance; i++ {
		numbers = append(numbers, paiNumber+i)
	}

	numNeighbors, err := count(numbers, func(num int) (bool, error) {
		if num < 1 || 9 < num {
			return false, nil
		}

		neighborPai, err := game.NewPaiWithDetail(pai.Type(), uint8(num), false)
		if err != nil {
			return false, err
		}

		count, err := prereachSutehaiSet.Count(neighborPai)
		if err != nil {
			return false, err
		}

		return count > 0, nil
	})
	if err != nil {
		return false, err
	}

	return numNeighbors >= n, nil
}

func isSujiOf(pai *game.Pai, targetPaiSet *game.PaiSet) (bool, error) {
	if pai.IsTsupai() {
		return false, nil
	}

	suji, err := getSuji(pai)
	if err != nil {
		return false, err
	}

	return allMatch(suji, func(s game.Pai) (bool, error) {
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

	return anyMatch(suji, func(s game.Pai) (bool, error) {
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

// Returns sujis which contain the given pai and is alive i.e. none of pais in the suji are anpai.
// Uses the first pai to represent the suji. e.g. 1p for 14p suji
func getPossibleSujis(pai *game.Pai, anpaiSet *game.PaiSet) ([]game.Pai, error) {
	if pai.IsTsupai() {
		return []game.Pai{}, nil
	}

	sujis := make([]game.Pai, 0, 2)
	paiNumber := pai.Number()
	candidates := []uint8{paiNumber - 3, paiNumber}

	for _, n := range candidates {
		isAlive, err := allMatch([]uint8{n, n + 3}, func(m uint8) (bool, error) {
			if m < 1 || m > 9 {
				return false, nil
			}

			sujiPai, err := game.NewPaiWithDetail(pai.Type(), m, false)
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
			sujiPai, err := game.NewPaiWithDetail(pai.Type(), n, false)
			if err != nil {
				return nil, err
			}
			sujis = append(sujis, *sujiPai)
		}
	}

	return sujis, nil
}

func isNChanceOrLess(pai *game.Pai, n int, visibleSet *game.PaiSet) (bool, error) {
	if pai.IsTsupai() {
		return false, nil
	}

	paiNumber := pai.Number()
	if 4 <= paiNumber && paiNumber <= 6 {
		return false, nil
	}

	candidates := make([]uint8, 2)
	for i := uint8(1); i < 3; i++ {
		if paiNumber < 5 {
			candidates[i-1] = paiNumber + i
		} else {
			candidates[i-1] = paiNumber - i
		}
	}

	return anyMatch(candidates, func(num uint8) (bool, error) {
		kabePai, err := game.NewPaiWithDetail(pai.Type(), num, false)
		if err != nil {
			return false, err
		}

		count, err := visibleSet.Count(kabePai)
		if err != nil {
			return false, err
		}

		return count >= 4-n, nil
	})
}

func isNumNOrInner(pai *game.Pai, n uint8) bool {
	if pai.IsTsupai() {
		return false
	}

	paiNumber := pai.Number()
	if n <= paiNumber && paiNumber <= 10-n {
		return true
	}

	return false
}

func isVisibleNOrMore(pai *game.Pai, n int, visibleSet *game.PaiSet) (bool, error) {
	c, err := visibleSet.Count(pai)
	if err != nil {
		return false, err
	}
	return c >= n, nil
}

func isUrasujiOf(pai *game.Pai, targetPaiSet *game.PaiSet, anpaiSet *game.PaiSet) (bool, error) {
	sujis, err := getPossibleSujis(pai, anpaiSet)
	if err != nil {
		return false, err
	}

	return anyMatch(sujis, func(s game.Pai) (bool, error) {
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

// Senkisuji (疝気筋) : Urasuji (裏筋) of urasuji
func isSenkisujiOf(pai *game.Pai, targetPaiSet *game.PaiSet, anpaiSet *game.PaiSet) (bool, error) {
	sujis, err := getPossibleSujis(pai, anpaiSet)
	if err != nil {
		return false, err
	}

	return anyMatch(sujis, func(s game.Pai) (bool, error) {
		if low := s.Next(-2); low != nil {
			hasLow, err := targetPaiSet.Has(low)
			if err != nil {
				return false, err
			}
			if hasLow {
				return true, nil
			}
		}

		if high := s.Next(5); high != nil {
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

func isMatagisujiOf(pai *game.Pai, targetPaiSet *game.PaiSet, anpaiSet *game.PaiSet) (bool, error) {
	sujis, err := getPossibleSujis(pai, anpaiSet)
	if err != nil {
		return false, err
	}

	return anyMatch(sujis, func(s game.Pai) (bool, error) {
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

func isOuter(pai *game.Pai, targetPaiSet *game.PaiSet) (bool, error) {
	if pai.IsTsupai() {
		return false, nil
	}

	paiNumber := pai.Number()
	if paiNumber == 5 {
		return false, nil
	}

	var innerNumbers []uint8
	if paiNumber < 5 {
		for i := paiNumber + 1; i < 6; i++ {
			innerNumbers = append(innerNumbers, i)
		}
	} else {
		for i := uint8(5); i < paiNumber; i++ {
			innerNumbers = append(innerNumbers, i)
		}
	}

	return anyMatch(innerNumbers, func(n uint8) (bool, error) {
		innerPai, err := game.NewPaiWithDetail(pai.Type(), n, false)
		if err != nil {
			return false, err
		}
		has, err := targetPaiSet.Has(innerPai)
		if err != nil {
			return false, err
		}
		return has, nil
	})
}

func registerEvaluators() evaluators {
	ev := evaluators{}

	for i := range 4 {
		featureName := fmt.Sprintf("chances<=%d", i)
		n := i
		ev[featureName] = func(scene *Scene, pai *game.Pai) (bool, error) {
			return isNChanceOrLess(pai, n, scene.visibleSet)
		}
	}

	// Whether i tiles are visible from one's perspective.
	// Includes one's own hand. Excludes the tile one is about to discard.
	for i := 1; i < 4; i++ {
		featureName := fmt.Sprintf("visible>=%d", i)
		n := i
		ev[featureName] = func(scene *Scene, pai *game.Pai) (bool, error) {
			return isVisibleNOrMore(pai, n+1, scene.visibleSet)
		}
	}

	// Among the Suji of that tile, whether one is visible no more than i copies.
	// The tile itself should not be counted.
	// In the case of 5p, this means "either 2p or 8p is visible no more than i copies,"
	// not "the combined visibility of 2p and 8p is no more than i copies."
	for i := range 4 {
		featureName := fmt.Sprintf("suji_visible<=%d", i)
		n := i
		ev[featureName] = func(scene *Scene, pai *game.Pai) (bool, error) {
			if pai.IsTsupai() {
				return false, nil
			}

			suji, err := getSuji(pai)
			if err != nil {
				return false, err
			}

			return anyMatch(suji, func(sujiPai game.Pai) (bool, error) {
				visible, err := isVisibleNOrMore(&sujiPai, n+1, scene.visibleSet)
				if err != nil {
					return false, err
				}
				return !visible, nil
			})
		}
	}

	for i := uint8(2); i < 6; i++ {
		featureName := fmt.Sprintf("%d<=n<=%d", i, 10-i)
		n := i
		ev[featureName] = func(scene *Scene, pai *game.Pai) (bool, error) {
			return isNumNOrInner(pai, n), nil
		}
	}

	for i := 2; i < 5; i++ {
		featureName := fmt.Sprintf("in_tehais>=%d", i)
		n := i
		ev[featureName] = func(scene *Scene, pai *game.Pai) (bool, error) {
			c, err := scene.tehaiSet.Count(pai)
			return c >= n, err
		}
	}

	// Among the Suji of that tile, whether one is held at least i copies.
	// The tile itself should not be counted.
	// In the case of 5p, this means "either 2p or 8p is held at least i copies,"
	// not "the combined total of 2p and 8p is at least i copies."
	for i := 1; i < 5; i++ {
		featureName := fmt.Sprintf("suji_in_tehais>=%d", i)
		n := i
		ev[featureName] = func(scene *Scene, pai *game.Pai) (bool, error) {
			if pai.IsTsupai() {
				return false, nil
			}

			suji, err := getSuji(pai)
			if err != nil {
				return false, err
			}

			return anyMatch(suji, func(sujiPai game.Pai) (bool, error) {
				c, err := scene.tehaiSet.Count(&sujiPai)
				return c >= n, err
			})
		}
	}

	for i := 1; i < 3; i++ {
		for j := 1; j < (i*2 + 1); j++ {
			featureName := fmt.Sprintf("+-%d_in_prereach_sutehais>=%d", i, j)
			distance := i
			threshold := j
			ev[featureName] = func(scene *Scene, pai *game.Pai) (bool, error) {
				return isNOrMoreOfNeighborsInPrereachSutehais(pai, threshold, distance, scene.prereachSutehaiSet)
			}
		}
	}

	for i := 1; i < 3; i++ {
		featureName := fmt.Sprintf("%d_outer_prereach_sutehai", i)
		n := i
		ev[featureName] = func(scene *Scene, pai *game.Pai) (bool, error) {
			return isNOuterPrereachSutehai(pai, n, scene.prereachSutehaiSet)
		}
	}

	for i := 1; i < 3; i++ {
		featureName := fmt.Sprintf("%d_inner_prereach_sutehai", i)
		n := i
		ev[featureName] = func(scene *Scene, pai *game.Pai) (bool, error) {
			return isNOuterPrereachSutehai(pai, -n, scene.prereachSutehaiSet)
		}
	}

	for i := 1; i < 9; i++ {
		featureName := fmt.Sprintf("same_type_in_prereach>=%d", i)
		n := i
		ev[featureName] = func(scene *Scene, pai *game.Pai) (bool, error) {
			if pai.IsTsupai() {
				return false, nil
			}

			numbers := []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9}
			numSameType, err := count(numbers, func(num uint8) (bool, error) {
				target, err := game.NewPaiWithDetail(pai.Type(), num, false)
				if err != nil {
					return false, err
				}
				return scene.prereachSutehaiSet.Has(target)
			})
			if err != nil {
				return false, err
			}

			return numSameType+1 >= n, nil
		}
	}

	return ev
}

type Feature struct {
	Name  string
	Value bool
}

type ProbInfo struct {
	Prob     float64
	Features []Feature
}

type DangerEstimator struct {
	root *configs.DecisionNode
}

func NewDangerEstimator(root *configs.DecisionNode) *DangerEstimator {
	return &DangerEstimator{
		root: root,
	}
}

func (e *DangerEstimator) EstimateProb(scene *Scene, pai *game.Pai) (*ProbInfo, error) {
	pai = pai.RemoveRed()
	node := e.root
	features := make([]Feature, 0)

	for node.FeatureName != nil {
		value, err := scene.Evaluate(*node.FeatureName, pai)
		if err != nil {
			return nil, err
		}

		features = append(features, Feature{
			Name:  *node.FeatureName,
			Value: value,
		})

		if value {
			node = node.Positive
		} else {
			node = node.Negative
		}
	}

	return &ProbInfo{
		Prob:     node.AverageProb,
		Features: features,
	}, nil
}
