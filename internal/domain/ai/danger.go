package ai

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type DangerTreeNode interface {
	LeafProb() (float64, bool)
	Feature() (string, bool)
	NegativeNode() DangerTreeNode
	PositiveNode() DangerTreeNode
}

type DecisionTreeDangerEstimator struct {
	root DangerTreeNode
}

func NewDangerEstimator(root DangerTreeNode) *DecisionTreeDangerEstimator {
	return &DecisionTreeDangerEstimator{root: root}
}

func (e *DecisionTreeDangerEstimator) EstimateDealInProb(
	state round.StateViewer,
	self seat.Seat,
	winner seat.Seat,
	discard tile.Tile,
) (float64, error) {
	if e == nil || e.root == nil {
		return 0, fmt.Errorf("cannot estimate deal-in probability: danger tree is nil")
	}
	discard = discard.RemoveRed()
	if containsSameSymbol(state.SafeTiles(winner), discard) {
		return 0, nil
	}
	scene := newDangerScene(state, self, winner)
	return estimateDangerTreeProb(e.root, scene, discard)
}

func estimateDangerTreeProb(root DangerTreeNode, scene dangerScene, discard tile.Tile) (float64, error) {
	node := root
	for node != nil {
		if prob, ok := node.LeafProb(); ok {
			return prob, nil
		}
		feature, ok := node.Feature()
		if !ok {
			return 0, fmt.Errorf("cannot estimate deal-in probability: non-leaf node has no feature")
		}
		value, err := scene.evaluate(feature, discard)
		if err != nil {
			return 0, err
		}
		if value {
			node = node.PositiveNode()
		} else {
			node = node.NegativeNode()
		}
	}
	return 0, fmt.Errorf("cannot estimate deal-in probability: danger tree branch is nil")
}
