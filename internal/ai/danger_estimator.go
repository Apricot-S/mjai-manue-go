package ai

import (
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

	return s, nil
}

func (s *Scene) Evaluate(name string, pai *game.Pai) bool {
	switch name {
	case "anpai":
		return s.isAnpai(pai)
	case "tsupai":
		return s.isTsupai(pai)
	// case "suji":
	// 	return s.Suji(pai)
	// // ... 他のすべてのfeature判定メソッドをcase文で列挙
	default:
		return false
	}
}

func (s *Scene) isAnpai(pai *game.Pai) bool {
	ret, _ := s.anpaiSet.Has(pai)
	return ret
}

func (s *Scene) isTsupai(pai *game.Pai) bool {
	return pai.IsTsupai()
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

func (e *DangerEstimator) EstimateProb(scene *Scene, pai *game.Pai) (*ProbInfo, error) {
	pai = pai.RemoveRed()
	node := e.root
	features := make([]Feature, 0)

	for node.FeatureName != nil {
		value := scene.Evaluate(*node.FeatureName, pai)
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
