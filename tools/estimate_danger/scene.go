package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

type Scene struct {
	dangerScene
	candidates tile.Tiles
}

var defaultFeatureNames = []string{
	"tsupai",
	"suji",
	"weak_suji",
	"reach_suji",
	"prereach_suji",
	"urasuji",
	"early_urasuji",
	"reach_urasuji",
	"urasuji_of_5",
	"aida4ken",
	"matagisuji",
	"early_matagisuji",
	"late_matagisuji",
	"reach_matagisuji",
	"senkisuji",
	"early_senkisuji",
	"outer_prereach_sutehai",
	"outer_early_sutehai",
	"chances<=0",
	"chances<=1",
	"chances<=2",
	"chances<=3",
	"visible>=1",
	"visible>=2",
	"visible>=3",
	"suji_visible<=0",
	"suji_visible<=1",
	"suji_visible<=2",
	"suji_visible<=3",
	"2<=n<=8",
	"3<=n<=7",
	"4<=n<=6",
	"5<=n<=5",
	"dora",
	"dora_suji",
	"dora_matagi",
	"in_tehais>=2",
	"in_tehais>=3",
	"in_tehais>=4",
	"suji_in_tehais>=1",
	"suji_in_tehais>=2",
	"suji_in_tehais>=3",
	"suji_in_tehais>=4",
	"+-1_in_prereach_sutehais>=1",
	"+-1_in_prereach_sutehais>=2",
	"+-2_in_prereach_sutehais>=1",
	"+-2_in_prereach_sutehais>=2",
	"+-2_in_prereach_sutehais>=3",
	"+-2_in_prereach_sutehais>=4",
	"1_outer_prereach_sutehai",
	"2_outer_prereach_sutehai",
	"1_inner_prereach_sutehai",
	"2_inner_prereach_sutehai",
	"same_type_in_prereach>=1",
	"same_type_in_prereach>=2",
	"same_type_in_prereach>=3",
	"same_type_in_prereach>=4",
	"same_type_in_prereach>=5",
	"same_type_in_prereach>=6",
	"same_type_in_prereach>=7",
	"same_type_in_prereach>=8",
	"fanpai",
	"ryenfonpai",
	"sangenpai",
	"fonpai",
	"bakaze",
	"jikaze",
}

func FeatureNames() []string {
	return slices.Clone(defaultFeatureNames)
}

func FeatureVectorToStr(featureVector *BitVector) string {
	var features []string
	for i, name := range defaultFeatureNames {
		if featureVector.Bit(i) != 0 {
			features = append(features, name)
		}
	}
	return strings.Join(features, " ")
}

func GetFeatureValue(featureVector *BitVector, featureName string) bool {
	index := slices.Index(defaultFeatureNames, featureName)
	return index >= 0 && featureVector.Bit(index) != 0
}

func NewScene(state round.StateViewer, self seat.Seat, target seat.Seat, discard tile.Tile) *Scene {
	ds := newDangerScene(state, self, target)
	// Ruby's training tool receives the game after the discard and adds the
	// discarded tile back to the actor's hand before collecting candidates.
	// CoffeeScript/runtime estimates from a live hand and does not need this
	// training-only adjustment.
	ds.selfHand = append(ds.selfHand, discard)
	return &Scene{dangerScene: ds, candidates: candidateTiles(ds.selfHand, ds.safeTiles)}
}

func NewSceneFromParams(
	selfHand, safeTiles, visibleTiles, doras, preRiichiTiles tile.Tiles,
	roundWind, targetWind wind.Wind,
) *Scene {
	ds := dangerScene{
		selfHand:       slices.Clone(selfHand),
		safeTiles:      slices.Clone(safeTiles),
		visibleTiles:   slices.Clone(visibleTiles),
		doras:          slices.Clone(doras),
		roundWind:      roundWind,
		targetWind:     targetWind,
		preRiichiTiles: slices.Clone(preRiichiTiles),
	}
	half := len(ds.preRiichiTiles) / 2
	ds.earlyPreRiichiTiles = ds.preRiichiTiles[:half]
	ds.latePreRiichiTiles = ds.preRiichiTiles[half:]
	if len(ds.preRiichiTiles) > 0 {
		ds.riichiDeclarationTiles = []tile.Tile{ds.preRiichiTiles[len(ds.preRiichiTiles)-1]}
	}
	return &Scene{dangerScene: ds, candidates: candidateTiles(ds.selfHand, ds.safeTiles)}
}

func candidateTiles(selfHand []tile.Tile, safeTiles []tile.Tile) tile.Tiles {
	candidates := tile.Tiles(selfHand).Distinct(func(t tile.Tile) bool {
		return tile.Tiles(safeTiles).ContainsSameSymbol(t)
	})
	for i := range candidates {
		candidates[i] = candidates[i].RemoveRed()
	}
	return candidates.Distinct(nil)
}

func (s *Scene) Candidates() tile.Tiles {
	// Callers treat candidates as read-only; avoid cloning for every discard.
	return s.candidates
}

func (s *Scene) FeatureVector(discard tile.Tile) (*BitVector, error) {
	boolArray := make([]bool, len(defaultFeatureNames))
	for i, name := range defaultFeatureNames {
		value, err := s.evaluate(name, discard)
		if err != nil {
			return nil, err
		}
		boolArray[i] = value
	}
	return BoolArrayToBitVector(boolArray), nil
}

type dangerScene struct {
	selfHand               []tile.Tile
	safeTiles              []tile.Tile
	visibleTiles           []tile.Tile
	doras                  []tile.Tile
	roundWind              wind.Wind
	targetWind             wind.Wind
	preRiichiTiles         []tile.Tile
	earlyPreRiichiTiles    []tile.Tile
	latePreRiichiTiles     []tile.Tile
	riichiDeclarationTiles []tile.Tile
}

func newDangerScene(state round.StateViewer, self seat.Seat, target seat.Seat) dangerScene {
	var selfHand []tile.Tile
	selfPlayer := state.Player(self)
	if h, ok := selfPlayer.Hand(); ok {
		selfHand = h.ToTiles()
		if drawnTile := selfPlayer.DrawnTile(); drawnTile != nil {
			selfHand = append(selfHand, *drawnTile)
		}
	}

	var preRiichiTiles []tile.Tile
	var riichiDeclarationTiles []tile.Tile
	targetPlayer := state.Player(target)
	discardedTiles := targetPlayer.DiscardedTiles()
	if idx := targetPlayer.RiichiDiscardedTilesIndex(); idx >= 0 && idx < len(discardedTiles) {
		preRiichiTiles = discardedTiles[:idx+1]
		riichiDeclarationTiles = []tile.Tile{discardedTiles[idx]}
	}

	half := len(preRiichiTiles) / 2

	return dangerScene{
		selfHand:               selfHand,
		safeTiles:              state.SafeTiles(target),
		visibleTiles:           state.VisibleTiles(self),
		doras:                  state.Doras(),
		roundWind:              state.RoundWind(),
		targetWind:             state.SeatWind(target),
		preRiichiTiles:         preRiichiTiles,
		earlyPreRiichiTiles:    preRiichiTiles[:half],
		latePreRiichiTiles:     preRiichiTiles[half:],
		riichiDeclarationTiles: riichiDeclarationTiles,
	}
}

func (s dangerScene) evaluate(feature string, discard tile.Tile) (bool, error) {
	switch feature {
	case "anpai":
		return tile.Tiles(s.safeTiles).ContainsSameSymbol(discard), nil
	case "tsupai":
		return discard.IsHonors(), nil
	case "dora":
		return tile.Tiles(s.doras).ContainsSameSymbol(discard), nil
	case "dora_suji":
		return isSujiOf(discard, s.doras, true), nil
	case "dora_matagi":
		return isMatagisujiOf(discard, s.doras, s.safeTiles), nil
	case "fanpai":
		return fanpaiValue(discard, s.roundWind, s.targetWind) >= 1, nil
	case "ryenfonpai":
		return fanpaiValue(discard, s.roundWind, s.targetWind) >= 2, nil
	case "fonpai":
		return discard.IsWind(), nil
	case "sangenpai":
		return discard.IsDragon(), nil
	case "bakaze":
		return windTile(s.roundWind).HasSameSymbol(discard), nil
	case "jikaze":
		return windTile(s.targetWind).HasSameSymbol(discard), nil
	case "suji":
		return isSujiOf(discard, s.safeTiles, false), nil
	case "weak_suji":
		return isSujiOf(discard, s.safeTiles, true), nil
	case "reach_suji":
		return isSujiOf(discard, s.riichiDeclarationTiles, true), nil
	case "prereach_suji":
		return isSujiOf(discard, s.preRiichiTiles, false), nil
	case "urasuji":
		return isUrasujiOf(discard, s.preRiichiTiles, s.safeTiles), nil
	case "early_urasuji":
		return isUrasujiOf(discard, s.earlyPreRiichiTiles, s.safeTiles), nil
	case "reach_urasuji":
		return isUrasujiOf(discard, s.riichiDeclarationTiles, s.safeTiles), nil
	case "urasuji_of_5":
		fives := slices.DeleteFunc(slices.Clone(s.preRiichiTiles), func(t tile.Tile) bool {
			return !t.IsSuits() || t.Number() != 5
		})
		// Ruby training tool defines urasuji_of_5; CoffeeScript/runtime does not.
		// Keep it here because extracted feature vectors must match Ruby.
		return isUrasujiOf(discard, fives, s.safeTiles), nil
	case "matagisuji":
		return isMatagisujiOf(discard, s.preRiichiTiles, s.safeTiles), nil
	case "early_matagisuji":
		return isMatagisujiOf(discard, s.earlyPreRiichiTiles, s.safeTiles), nil
	case "late_matagisuji":
		return isMatagisujiOf(discard, s.latePreRiichiTiles, s.safeTiles), nil
	case "reach_matagisuji":
		return isMatagisujiOf(discard, s.riichiDeclarationTiles, s.safeTiles), nil
	case "senkisuji":
		return isSenkisujiOf(discard, s.preRiichiTiles, s.safeTiles), nil
	case "early_senkisuji":
		return isSenkisujiOf(discard, s.earlyPreRiichiTiles, s.safeTiles), nil
	case "outer_prereach_sutehai":
		return isOuter(discard, s.preRiichiTiles), nil
	case "outer_early_sutehai":
		return isOuter(discard, s.earlyPreRiichiTiles), nil
	case "aida4ken":
		return isAida4Ken(discard, s.preRiichiTiles), nil
	}

	n, matched, err := parseFeatureIntPrefix(feature, "chances<=")
	if err != nil {
		return false, err
	}
	if matched {
		return isNChanceOrLess(discard, n, s.visibleTiles), nil
	}
	n, matched, err = parseFeatureIntPrefix(feature, "visible>=")
	if err != nil {
		return false, err
	}
	if matched {
		return tile.Tiles(s.visibleTiles).CountSameSymbol(discard) >= n+1, nil
	}
	n, matched, err = parseFeatureIntPrefix(feature, "suji_visible<=")
	if err != nil {
		return false, err
	}
	if matched {
		return isSujiVisibleNoMoreThan(discard, n, s.visibleTiles), nil
	}
	n, matched, err = parseFeatureIntPrefix(feature, "in_tehais>=")
	if err != nil {
		return false, err
	}
	if matched {
		return tile.Tiles(s.selfHand).CountSameSymbol(discard) >= n, nil
	}
	n, matched, err = parseFeatureIntPrefix(feature, "suji_in_tehais>=")
	if err != nil {
		return false, err
	}
	if matched {
		return hasSujiSymbolCount(discard, n, s.selfHand), nil
	}
	if strings.Contains(feature, "<=n<=") {
		return evalNumberRange(feature, discard)
	}
	n, matched, err = parseFeatureIntSuffix(feature, "_outer_prereach_sutehai")
	if err != nil {
		return false, err
	}
	if matched {
		return isNOuterPreRiichiSutehai(discard, n, s.preRiichiTiles), nil
	}
	n, matched, err = parseFeatureIntSuffix(feature, "_inner_prereach_sutehai")
	if err != nil {
		return false, err
	}
	if matched {
		return isNOuterPreRiichiSutehai(discard, -n, s.preRiichiTiles), nil
	}
	n, matched, err = parseFeatureIntPrefix(feature, "same_type_in_prereach>=")
	if err != nil {
		return false, err
	}
	if matched {
		// Ruby training tool counts only same-suit tile numbers in pre-riichi
		// discards. CoffeeScript/runtime adds one for the candidate tile.
		return discard.IsSuits() && countSameColor(s.preRiichiTiles, discard) >= n, nil
	}
	if strings.HasPrefix(feature, "+-") && strings.Contains(feature, "_in_prereach_sutehais>=") {
		return evalNeighborPreRiichi(feature, discard, s.preRiichiTiles)
	}

	return false, fmt.Errorf("cannot evaluate danger feature %q", feature)
}

func evalNumberRange(feature string, target tile.Tile) (bool, error) {
	if !target.IsSuits() {
		return false, nil
	}
	parts := strings.Split(feature, "<=n<=")
	if len(parts) != 2 {
		return false, fmt.Errorf("parse danger feature %q: invalid number range", feature)
	}
	minN, err := strconv.Atoi(parts[0])
	if err != nil {
		return false, fmt.Errorf("parse danger feature %q minimum: %w", feature, err)
	}
	maxN, err := strconv.Atoi(parts[1])
	if err != nil {
		return false, fmt.Errorf("parse danger feature %q maximum: %w", feature, err)
	}
	return minN <= target.Number() && target.Number() <= maxN, nil
}

func evalNeighborPreRiichi(feature string, target tile.Tile, tiles []tile.Tile) (bool, error) {
	if !target.IsSuits() {
		return false, nil
	}
	parts := strings.Split(strings.TrimPrefix(feature, "+-"), "_in_prereach_sutehais>=")
	if len(parts) != 2 {
		return false, fmt.Errorf("parse danger feature %q: invalid neighbor pre-riichi feature", feature)
	}
	distance, err := strconv.Atoi(parts[0])
	if err != nil {
		return false, fmt.Errorf("parse danger feature %q distance: %w", feature, err)
	}
	threshold, err := strconv.Atoi(parts[1])
	if err != nil {
		return false, fmt.Errorf("parse danger feature %q threshold: %w", feature, err)
	}

	count := 0
	for offset := -distance; offset <= distance; offset++ {
		if neighbor := target.Next(offset); neighbor != nil {
			if tile.Tiles(tiles).ContainsSameSymbol(*neighbor) {
				count++
			}
		}
	}
	return count >= threshold, nil
}

func parseFeatureIntPrefix(feature string, prefix string) (value int, matched bool, err error) {
	if !strings.HasPrefix(feature, prefix) {
		return 0, false, nil
	}
	value, err = strconv.Atoi(strings.TrimPrefix(feature, prefix))
	if err != nil {
		return 0, true, fmt.Errorf("parse danger feature %q: %w", feature, err)
	}
	return value, true, nil
}

func parseFeatureIntSuffix(feature string, suffix string) (value int, matched bool, err error) {
	if !strings.HasSuffix(feature, suffix) {
		return 0, false, nil
	}
	value, err = strconv.Atoi(strings.TrimSuffix(feature, suffix))
	if err != nil {
		return 0, true, fmt.Errorf("parse danger feature %q: %w", feature, err)
	}
	return value, true, nil
}
