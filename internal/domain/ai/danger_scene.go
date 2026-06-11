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
	selfPlayer := state.Player(self)
	targetPlayer := state.Player(target)
	var selfHand []tile.Tile
	if h, ok := selfPlayer.Hand(); ok {
		selfHand = h.ToTiles()
	}
	preRiichiTiles := []tile.Tile(nil)
	riichiDeclarationTiles := []tile.Tile(nil)
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

	if strings.HasPrefix(feature, "chances<=") {
		n, ok := parseFeatureInt(feature, "chances<=")
		return ok && isNChanceOrLess(discard, n, s.visibleTiles), nil
	}
	if strings.HasPrefix(feature, "visible>=") {
		n, ok := parseFeatureInt(feature, "visible>=")
		return ok && tile.Tiles(s.visibleTiles).CountSameSymbol(discard) >= n+1, nil
	}
	if strings.HasPrefix(feature, "suji_visible<=") {
		n, ok := parseFeatureInt(feature, "suji_visible<=")
		return ok && isSujiVisibleNoMoreThan(discard, n, s.visibleTiles), nil
	}
	if strings.HasPrefix(feature, "in_tehais>=") {
		n, ok := parseFeatureInt(feature, "in_tehais>=")
		return ok && tile.Tiles(s.selfHand).CountSameSymbol(discard) >= n, nil
	}
	if strings.HasPrefix(feature, "suji_in_tehais>=") {
		n, ok := parseFeatureInt(feature, "suji_in_tehais>=")
		return ok && hasSujiSymbolCount(discard, n, s.selfHand), nil
	}
	if strings.Contains(feature, "<=n<=") {
		return evalNumberRange(feature, discard), nil
	}
	if strings.Contains(feature, "_outer_prereach_sutehai") {
		n := leadingFeatureInt(feature)
		return isNOuterPreRiichiSutehai(discard, n, s.preRiichiTiles), nil
	}
	if strings.Contains(feature, "_inner_prereach_sutehai") {
		n := leadingFeatureInt(feature)
		return isNOuterPreRiichiSutehai(discard, -n, s.preRiichiTiles), nil
	}
	if strings.HasPrefix(feature, "same_type_in_prereach>=") {
		n, ok := parseFeatureInt(feature, "same_type_in_prereach>=")
		return ok && discard.IsSuits() && countSameColor(s.preRiichiTiles, discard)+1 >= n, nil
	}
	if strings.HasPrefix(feature, "+-") && strings.Contains(feature, "_in_prereach_sutehais>=") {
		return evalNeighborPreRiichi(feature, discard, s.preRiichiTiles), nil
	}

	return false, fmt.Errorf("cannot evaluate danger feature %q", feature)
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

func evalNeighborPreRiichi(feature string, target tile.Tile, tiles []tile.Tile) bool {
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
	for offset := -distance; offset <= distance; offset++ {
		if neighbor := target.Next(offset); neighbor != nil {
			if tile.Tiles(tiles).ContainsSameSymbol(*neighbor) {
				count++
			}
		}
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
	part, _, _ := strings.Cut(feature, "_")
	n, _ := strconv.Atoi(part)
	return n
}
