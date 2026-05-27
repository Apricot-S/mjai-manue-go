package ai

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

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
