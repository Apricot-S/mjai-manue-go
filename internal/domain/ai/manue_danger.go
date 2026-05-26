package ai

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
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
	scene := newDangerScene(state, self, winner)
	if scene.isSafe(discard) {
		return 0, nil
	}
	node := e.root
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

type dangerScene struct {
	selfHand        []tile.Tile
	safeTiles       []tile.Tile
	visibleTiles    []tile.Tile
	doras           []tile.Tile
	roundWind       wind.Wind
	targetWind      wind.Wind
	prereachTiles   []tile.Tile
	earlyReachTiles []tile.Tile
	lateReachTiles  []tile.Tile
	reachTiles      []tile.Tile
}

func newDangerScene(state round.StateViewer, self seat.Seat, target seat.Seat) dangerScene {
	selfPlayer := state.Player(self)
	targetPlayer := state.Player(target)
	var selfHand []tile.Tile
	if h, ok := selfPlayer.Hand(); ok {
		selfHand = h.ToTiles()
	}
	prereachTiles := targetPlayer.DiscardedTiles()
	reachTiles := []tile.Tile(nil)
	if idx := targetPlayer.RiichiDiscardedTilesIndex(); idx >= 0 && idx < len(prereachTiles) {
		prereachTiles = prereachTiles[:idx+1]
		reachTiles = []tile.Tile{prereachTiles[idx]}
	}
	half := len(prereachTiles) / 2
	return dangerScene{
		selfHand:        selfHand,
		safeTiles:       state.SafeTiles(target),
		visibleTiles:    state.VisibleTiles(self),
		doras:           state.Doras(),
		roundWind:       state.RoundWind(),
		targetWind:      state.SeatWind(target),
		prereachTiles:   prereachTiles,
		earlyReachTiles: prereachTiles[:half],
		lateReachTiles:  prereachTiles[half:],
		reachTiles:      reachTiles,
	}
}

func (s dangerScene) evaluate(feature string, discard tile.Tile) (bool, error) {
	switch feature {
	case "anpai":
		return s.isSafe(discard), nil
	case "tsupai":
		return discard.IsHonors(), nil
	case "dora":
		return containsSameSymbol(s.doras, discard), nil
	case "dora_suji":
		return isSujiOf(discard, s.doras, false), nil
	case "dora_matagi":
		return isMatagisujiOf(discard, s.doras, s.safeTiles), nil
	case "fanpai":
		return fanpaiValue(discard, s.roundWind, s.targetWind) >= 1, nil
	case "ryenfonpai":
		return fanpaiValue(discard, s.roundWind, s.targetWind) >= 2, nil
	case "fonpai":
		return discard.IsHonors() && discard.ID() >= tile.MustTileFromCode("E").ID() && discard.ID() <= tile.MustTileFromCode("N").ID(), nil
	case "sangenpai":
		return discard == tile.MustTileFromCode("P") || discard == tile.MustTileFromCode("F") || discard == tile.MustTileFromCode("C"), nil
	case "bakaze":
		return windTile(s.roundWind).HasSameSymbol(discard), nil
	case "jikaze":
		return windTile(s.targetWind).HasSameSymbol(discard), nil
	case "suji":
		return isSujiOf(discard, s.safeTiles, false), nil
	case "weak_suji":
		return isSujiOf(discard, s.safeTiles, true), nil
	case "reach_suji":
		return isSujiOf(discard, s.reachTiles, true), nil
	case "prereach_suji":
		return isSujiOf(discard, s.prereachTiles, false), nil
	case "urasuji":
		return isUrasujiOf(discard, s.prereachTiles, s.safeTiles), nil
	case "early_urasuji":
		return isUrasujiOf(discard, s.earlyReachTiles, s.safeTiles), nil
	case "reach_urasuji":
		return isUrasujiOf(discard, s.reachTiles, s.safeTiles), nil
	case "matagisuji":
		return isMatagisujiOf(discard, s.prereachTiles, s.safeTiles), nil
	case "early_matagisuji":
		return isMatagisujiOf(discard, s.earlyReachTiles, s.safeTiles), nil
	case "late_matagisuji":
		return isMatagisujiOf(discard, s.lateReachTiles, s.safeTiles), nil
	case "reach_matagisuji":
		return isMatagisujiOf(discard, s.reachTiles, s.safeTiles), nil
	case "senkisuji":
		return isSenkisujiOf(discard, s.prereachTiles, s.safeTiles), nil
	case "early_senkisuji":
		return isSenkisujiOf(discard, s.earlyReachTiles, s.safeTiles), nil
	case "outer_prereach_sutehai":
		return isOuter(discard, s.prereachTiles), nil
	case "outer_early_sutehai":
		return isOuter(discard, s.earlyReachTiles), nil
	case "aida4ken":
		return isAida4Ken(discard, s.prereachTiles), nil
	}
	if strings.HasPrefix(feature, "chances<=") {
		n, ok := parseFeatureInt(feature, "chances<=")
		return ok && isNChanceOrLess(discard, n, s.visibleTiles), nil
	}
	if strings.HasPrefix(feature, "visible>=") {
		n, ok := parseFeatureInt(feature, "visible>=")
		return ok && countSameSymbol(s.visibleTiles, discard) >= n+1, nil
	}
	if strings.HasPrefix(feature, "suji_visible<=") {
		n, ok := parseFeatureInt(feature, "suji_visible<=")
		return ok && isSujiVisibleNoMoreThan(discard, n, s.visibleTiles), nil
	}
	if strings.HasPrefix(feature, "in_tehais>=") {
		n, ok := parseFeatureInt(feature, "in_tehais>=")
		return ok && countSameSymbol(s.selfHand, discard) >= n, nil
	}
	if strings.HasPrefix(feature, "suji_in_tehais>=") {
		n, ok := parseFeatureInt(feature, "suji_in_tehais>=")
		return ok && countSujiSymbols(discard, s.selfHand) >= n, nil
	}
	if strings.Contains(feature, "<=n<=") {
		return evalNumberRange(feature, discard), nil
	}
	if strings.Contains(feature, "_outer_prereach_sutehai") {
		n := leadingFeatureInt(feature)
		return isNOuterPrereachSutehai(discard, n, s.prereachTiles), nil
	}
	if strings.Contains(feature, "_inner_prereach_sutehai") {
		n := leadingFeatureInt(feature)
		return isNOuterPrereachSutehai(discard, -n, s.prereachTiles), nil
	}
	if strings.HasPrefix(feature, "same_type_in_prereach>=") {
		n, ok := parseFeatureInt(feature, "same_type_in_prereach>=")
		return ok && countSameColor(s.prereachTiles, discard)+1 >= n, nil
	}
	if strings.HasPrefix(feature, "+-") && strings.Contains(feature, "_in_prereach_sutehais>=") {
		return evalNeighborPrereach(feature, discard, s.prereachTiles), nil
	}
	return false, fmt.Errorf("cannot evaluate danger feature %q", feature)
}

func (s dangerScene) isSafe(discard tile.Tile) bool {
	return containsSameSymbol(s.safeTiles, discard)
}

func containsSameSymbol(tiles []tile.Tile, target tile.Tile) bool {
	for _, t := range tiles {
		if t.HasSameSymbol(target) {
			return true
		}
	}
	return false
}

func countSameSymbol(tiles []tile.Tile, target tile.Tile) int {
	count := 0
	for _, t := range tiles {
		if t.HasSameSymbol(target) {
			count++
		}
	}
	return count
}

func countSameColor(tiles []tile.Tile, target tile.Tile) int {
	if !target.IsSuits() {
		return 0
	}
	count := 0
	for _, t := range tiles {
		if t.IsSuits() && t.Color() == target.Color() {
			count++
		}
	}
	return count
}

func countSujiSymbols(target tile.Tile, tiles []tile.Tile) int {
	count := 0
	for _, s := range sujiTiles(target) {
		count += countSameSymbol(tiles, s)
	}
	return count
}

func fanpaiValue(target tile.Tile, roundWind wind.Wind, targetWind wind.Wind) int {
	if !target.IsHonors() {
		return 0
	}
	if target == tile.MustTileFromCode("P") || target == tile.MustTileFromCode("F") || target == tile.MustTileFromCode("C") {
		return 1
	}
	value := 0
	if windTile(roundWind).HasSameSymbol(target) {
		value++
	}
	if windTile(targetWind).HasSameSymbol(target) {
		value++
	}
	return value
}

func sujiTiles(target tile.Tile) []tile.Tile {
	if !target.IsSuits() {
		return nil
	}
	result := make([]tile.Tile, 0, 2)
	for _, n := range []int{target.Number() - 3, target.Number() + 3} {
		if 1 <= n && n <= 9 {
			result = append(result, tile.MustTileFromID(target.RemoveRed().ID()+n-target.Number()))
		}
	}
	return result
}

func isSujiOf(target tile.Tile, tiles []tile.Tile, weak bool) bool {
	sujis := sujiTiles(target)
	if len(sujis) == 0 {
		return false
	}
	matches := 0
	for _, s := range sujis {
		if containsSameSymbol(tiles, s) {
			matches++
		}
	}
	if weak {
		return matches > 0
	}
	return matches == len(sujis)
}

func isSujiVisibleNoMoreThan(target tile.Tile, n int, visibleTiles []tile.Tile) bool {
	if !target.IsSuits() {
		return false
	}
	for _, s := range sujiTiles(target) {
		if countSameSymbol(visibleTiles, s) < n+1 {
			return true
		}
	}
	return false
}

func isNChanceOrLess(target tile.Tile, n int, visibleTiles []tile.Tile) bool {
	if !target.IsSuits() {
		return false
	}
	targetNumber := target.Number()
	if 4 <= targetNumber && targetNumber <= 6 {
		return false
	}
	for i := 1; i < 3; i++ {
		kabeNumber := targetNumber + i
		if targetNumber >= 5 {
			kabeNumber = targetNumber - i
		}
		kabe := tile.MustTileFromID(target.RemoveRed().ID() + kabeNumber - targetNumber)
		if countSameSymbol(visibleTiles, kabe) >= 4-n {
			return true
		}
	}
	return false
}

func possibleSujis(target tile.Tile, safeTiles []tile.Tile) []tile.Tile {
	if !target.IsSuits() {
		return nil
	}
	result := make([]tile.Tile, 0, 2)
	for _, n := range []int{target.Number() - 3, target.Number()} {
		if n < 1 || n+3 > 9 {
			continue
		}
		first := tile.MustTileFromID(target.RemoveRed().ID() + n - target.Number())
		second := tile.MustTileFromID(first.ID() + 3)
		if !containsSameSymbol(safeTiles, first) && !containsSameSymbol(safeTiles, second) {
			result = append(result, first)
		}
	}
	return result
}

func isUrasujiOf(target tile.Tile, tiles []tile.Tile, safeTiles []tile.Tile) bool {
	for _, s := range possibleSujis(target, safeTiles) {
		if low := s.Next(-1); low != nil && containsSameSymbol(tiles, *low) {
			return true
		}
		if high := s.Next(4); high != nil && containsSameSymbol(tiles, *high) {
			return true
		}
	}
	return false
}

func isSenkisujiOf(target tile.Tile, tiles []tile.Tile, safeTiles []tile.Tile) bool {
	for _, s := range possibleSujis(target, safeTiles) {
		if low := s.Next(-2); low != nil && containsSameSymbol(tiles, *low) {
			return true
		}
		if high := s.Next(5); high != nil && containsSameSymbol(tiles, *high) {
			return true
		}
	}
	return false
}

func isMatagisujiOf(target tile.Tile, tiles []tile.Tile, safeTiles []tile.Tile) bool {
	for _, s := range possibleSujis(target, safeTiles) {
		if low := s.Next(1); low != nil && containsSameSymbol(tiles, *low) {
			return true
		}
		if high := s.Next(2); high != nil && containsSameSymbol(tiles, *high) {
			return true
		}
	}
	return false
}

func isOuter(target tile.Tile, tiles []tile.Tile) bool {
	if !target.IsSuits() {
		return false
	}
	if target.Number() == 5 {
		return false
	}
	var innerNumbers []int
	if target.Number() < 5 {
		for n := target.Number() + 1; n < 6; n++ {
			innerNumbers = append(innerNumbers, n)
		}
	} else {
		for n := 5; n < target.Number(); n++ {
			innerNumbers = append(innerNumbers, n)
		}
	}
	for _, n := range innerNumbers {
		if containsSameSymbol(tiles, tile.MustTileFromID(target.RemoveRed().ID()+n-target.Number())) {
			return true
		}
	}
	return false
}

func isNOuterPrereachSutehai(target tile.Tile, n int, tiles []tile.Tile) bool {
	if !target.IsSuits() || target.Number() == 5 {
		return false
	}
	innerNumber := target.Number() + n
	if target.Number() >= 5 {
		innerNumber = target.Number() - n
	}
	if innerNumber < 1 || innerNumber > 9 {
		return false
	}
	if (target.Number() >= 5 || innerNumber > 5) && (target.Number() <= 5 || innerNumber < 5) {
		return false
	}
	return containsSameSymbol(tiles, tile.MustTileFromID(target.RemoveRed().ID()+innerNumber-target.Number()))
}

func isAida4Ken(target tile.Tile, tiles []tile.Tile) bool {
	if !target.IsSuits() {
		return false
	}
	n := target.Number()
	if 2 <= n && n <= 5 {
		return containsSameSymbol(tiles, tile.MustTileFromID(target.RemoveRed().ID()-1)) &&
			containsSameSymbol(tiles, tile.MustTileFromID(target.RemoveRed().ID()+4))
	}
	if 5 <= n && n <= 8 {
		return containsSameSymbol(tiles, tile.MustTileFromID(target.RemoveRed().ID()-4)) &&
			containsSameSymbol(tiles, tile.MustTileFromID(target.RemoveRed().ID()+1))
	}
	return false
}

func evalNumberRange(feature string, target tile.Tile) bool {
	if !target.IsSuits() {
		return false
	}
	parts := strings.Split(feature, "<=n<=")
	if len(parts) != 2 {
		return false
	}
	minN, err1 := strconv.Atoi(parts[0])
	maxN, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return false
	}
	return minN <= target.Number() && target.Number() <= maxN
}

func evalNeighborPrereach(feature string, target tile.Tile, tiles []tile.Tile) bool {
	parts := strings.Split(strings.TrimPrefix(feature, "+-"), "_in_prereach_sutehais>=")
	if len(parts) != 2 {
		return false
	}
	distance, err1 := strconv.Atoi(parts[0])
	threshold, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil || !target.IsSuits() {
		return false
	}
	count := 0
	if low := target.Next(-distance); low != nil {
		count += countSameSymbol(tiles, *low)
	}
	if high := target.Next(distance); high != nil {
		count += countSameSymbol(tiles, *high)
	}
	return count >= threshold
}

func parseFeatureInt(feature string, prefix string) (int, bool) {
	if !strings.HasPrefix(feature, prefix) {
		return 0, false
	}
	n, err := strconv.Atoi(strings.TrimPrefix(feature, prefix))
	return n, err == nil
}

func leadingFeatureInt(feature string) int {
	part := strings.SplitN(feature, "_", 2)[0]
	n, _ := strconv.Atoi(part)
	return n
}

func windTile(w wind.Wind) tile.Tile {
	switch w {
	case wind.East:
		return tile.MustTileFromCode("E")
	case wind.South:
		return tile.MustTileFromCode("S")
	case wind.West:
		return tile.MustTileFromCode("W")
	default:
		return tile.MustTileFromCode("N")
	}
}
