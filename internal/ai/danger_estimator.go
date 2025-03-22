package ai

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/configs"
	"github.com/Apricot-S/mjai-manue-go/internal/game"
)

type Scene struct {
	gameState *game.State
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

	evaluators map[string]func(*game.Pai) (bool, error)
}

func NewScene(gameState *game.State, me *game.Player, target *game.Player) (*Scene, error) {
	s := &Scene{
		gameState: gameState,
		me:        me,
		target:    target,
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
	if idx := target.ReachSutehaiIndex(); idx != nil {
		prereachSutehais = target.Sutehais()[:*idx+1]
		reachPai := target.Sutehais()[*idx]
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

	s.registerEvaluators()

	return s, nil
}

func (s *Scene) Evaluate(name string, pai *game.Pai) (bool, error) {
	if evaluator, ok := s.evaluators[name]; ok {
		return evaluator(pai)
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
		return false, nil
	}
}

func (s *Scene) isAnpai(pai *game.Pai) (bool, error) {
	return s.anpaiSet.Has(pai)
}

func (s *Scene) isTsupai(pai *game.Pai) bool {
	return pai.IsTsupai()
}

func (s *Scene) isSuji(pai *game.Pai) (bool, error) {
	return isSujiOf(pai, s.anpaiSet)
}

func (s *Scene) isWeakSuji(pai *game.Pai) (bool, error) {
	return isWeakSujiOf(pai, s.anpaiSet)
}

func (s *Scene) isReachSuji(pai *game.Pai) (bool, error) {
	return isWeakSujiOf(pai, s.reachPaiSet)
}

func (s *Scene) isPrereachSuji(pai *game.Pai) (bool, error) {
	return isSujiOf(pai, s.prereachSutehaiSet)
}

func (s *Scene) isUrasuji(pai *game.Pai) (bool, error) {
	return isUrasujiOf(pai, s.prereachSutehaiSet, s.anpaiSet)
}

func (s *Scene) isEarlyUrasuji(pai *game.Pai) (bool, error) {
	return isUrasujiOf(pai, s.earlySutehaiSet, s.anpaiSet)
}

func (s *Scene) isReachUrasuji(pai *game.Pai) (bool, error) {
	return isUrasujiOf(pai, s.reachPaiSet, s.anpaiSet)
}

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

func isSujiOf(pai *game.Pai, targetPaiSet *game.PaiSet) (bool, error) {
	if pai.IsTsupai() {
		return false, nil
	}

	suji, err := getSuji(pai)
	if err != nil {
		return false, err
	}

	for _, s := range suji {
		hasPai, err := targetPaiSet.Has(&s)
		if err != nil {
			return false, err
		}

		if !hasPai {
			return false, nil
		}
	}
	return true, nil
}

func isWeakSujiOf(pai *game.Pai, targetPaiSet *game.PaiSet) (bool, error) {
	if pai.IsTsupai() {
		return false, nil
	}

	suji, err := getSuji(pai)
	if err != nil {
		return false, err
	}

	for _, s := range suji {
		hasPai, err := targetPaiSet.Has(&s)
		if err != nil {
			return false, err
		}

		if hasPai {
			return true, nil
		}
	}
	return false, nil
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
		allAlive := true
		for _, m := range []uint8{n, n + 3} {
			if m < 1 || 9 < m {
				allAlive = false
				break
			}

			sujiPai, err := game.NewPaiWithDetail(pai.Type(), m, false)
			if err != nil {
				return nil, err
			}

			isAnpai, err := anpaiSet.Has(sujiPai)
			if err != nil {
				return nil, err
			}
			if isAnpai {
				allAlive = false
				break
			}
		}

		if allAlive {
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

	for i := uint8(1); i < 3; i++ {
		var num uint8 = 0
		if paiNumber < 5 {
			num = paiNumber + i
		} else {
			num = paiNumber - i
		}

		kabePai, err := game.NewPaiWithDetail(pai.Type(), num, false)
		if err != nil {
			return false, err
		}

		count, err := visibleSet.Count(kabePai)
		if err != nil {
			return false, err
		}

		if count >= 4-n {
			return true, nil
		}
	}

	return false, nil
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

func isUrasujiOf(pai *game.Pai, targetPaiSet *game.PaiSet, anpaiSet *game.PaiSet) (bool, error) {
	sujis, err := getPossibleSujis(pai, anpaiSet)
	if err != nil {
		return false, err
	}

	for _, s := range sujis {
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
	}

	return false, nil
}

// Senkisuji (疝気筋) : Urasuji (裏筋) of urasuji
func isSenkisujiOf(pai *game.Pai, targetPaiSet *game.PaiSet, anpaiSet *game.PaiSet) (bool, error) {
	sujis, err := getPossibleSujis(pai, anpaiSet)
	if err != nil {
		return false, err
	}

	for _, s := range sujis {
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
	}

	return false, nil
}

func isMatagisujiOf(pai *game.Pai, targetPaiSet *game.PaiSet, anpaiSet *game.PaiSet) (bool, error) {
	sujis, err := getPossibleSujis(pai, anpaiSet)
	if err != nil {
		return false, err
	}

	for _, s := range sujis {
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
	}

	return false, nil
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

	for _, n := range innerNumbers {
		innerPai, err := game.NewPaiWithDetail(pai.Type(), n, false)
		if err != nil {
			return false, err
		}
		has, err := targetPaiSet.Has(innerPai)
		if err != nil {
			return false, err
		}

		if has {
			return true, nil
		}
	}

	return false, nil
}

func (s *Scene) registerEvaluators() {
	for i := range 4 {
		featureName := fmt.Sprintf("chances<=%d", i)
		s.evaluators[featureName] = func(pai *game.Pai) (bool, error) {
			return isNChanceOrLess(pai, i, s.visibleSet)
		}
	}

	for i := uint8(2); i < 6; i++ {
		featureName := fmt.Sprintf("%d<=n<=%d", i, 10-i)
		s.evaluators[featureName] = func(pai *game.Pai) (bool, error) {
			return isNumNOrInner(pai, i), nil
		}
	}
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
	root *configs.DangerNode
}

func NewDangerEstimator(root *configs.DangerNode) *DangerEstimator {
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
