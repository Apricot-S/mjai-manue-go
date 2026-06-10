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
	seenNumbers := [10]bool{}
	count := 0
	for _, t := range tiles {
		if t.IsSuits() && t.Color() == target.Color() && !seenNumbers[t.Number()] {
			seenNumbers[t.Number()] = true
			count++
		}
	}
	return count
}

func hasSujiSymbolCount(target tile.Tile, threshold int, tiles []tile.Tile) bool {
	for _, s := range sujiTiles(target) {
		if countSameSymbol(tiles, s) >= threshold {
			return true
		}
	}
	return false
}

func fanpaiValue(target tile.Tile, roundWind wind.Wind, targetWind wind.Wind) int {
	if !target.IsHonors() {
		return 0
	}
	if target.IsDragon() {
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
	for _, offset := range []int{-3, 3} {
		if suji := target.Next(offset); suji != nil {
			result = append(result, *suji)
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
		offset := i
		if targetNumber >= 5 {
			offset = -i
		}
		kabe := target.Next(offset)
		if kabe == nil {
			continue
		}
		if countSameSymbol(visibleTiles, *kabe) >= 4-n {
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
		first := target.Next(n - target.Number())
		if first == nil {
			continue
		}
		second := first.Next(3)
		if second == nil {
			continue
		}
		if !containsSameSymbol(safeTiles, *first) && !containsSameSymbol(safeTiles, *second) {
			result = append(result, *first)
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

	number := target.Number()
	if number == 5 {
		return false
	}

	if number < 5 {
		for offset := 1; number+offset <= 5; offset++ {
			inner := target.Next(offset)
			if inner != nil && containsSameSymbol(tiles, *inner) {
				return true
			}
		}
	} else {
		for offset := -1; number+offset >= 5; offset-- {
			inner := target.Next(offset)
			if inner != nil && containsSameSymbol(tiles, *inner) {
				return true
			}
		}
	}

	return false
}

func isNOuterPrereachSutehai(target tile.Tile, n int, tiles []tile.Tile) bool {
	if !target.IsSuits() {
		return false
	}

	number := target.Number()
	if number == 5 {
		return false
	}

	if number >= 5 {
		n = -n
	}
	inner := target.Next(n)
	if inner == nil {
		return false
	}
	innerNumber := inner.Number()
	if (number >= 5 || innerNumber > 5) && (number <= 5 || innerNumber < 5) {
		return false
	}
	return containsSameSymbol(tiles, *inner)
}

func isAida4Ken(target tile.Tile, tiles []tile.Tile) bool {
	if !target.IsSuits() {
		return false
	}
	n := target.Number()
	matches := false
	if 2 <= n && n <= 5 {
		low := target.Next(-1)
		high := target.Next(4)
		matches = low != nil && high != nil && containsSameSymbol(tiles, *low) &&
			containsSameSymbol(tiles, *high)
	}
	if 5 <= n && n <= 8 {
		low := target.Next(-4)
		high := target.Next(1)
		matches = matches || low != nil && high != nil && containsSameSymbol(tiles, *low) &&
			containsSameSymbol(tiles, *high)
	}
	return matches
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
