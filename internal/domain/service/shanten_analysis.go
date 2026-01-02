package service

import (
	"math"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/block"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tilecount"
)

// Goal represents a winning hand that can be transitioned from the current hand.
// It contains information about the current hand's shanten number,
// the winning hand's set composition, and the number of required and unrequired tiles.
type Goal struct {
	// Shanten is the shanten number from the current hand to the winning hand.
	// If the current hand and the winning hand are the same, it will be -1.
	Shanten int
	// Blocks is a list of blocks in the winning hand.
	Blocks []block.Block
	// CountVector is the number of each tile included in the winning hand.
	CountVector tilecount.TileCounts34
	// RequiredVector is the number of each tile required for the winning hand.
	RequiredVector tilecount.TileCounts34
	// ThrowableVector is the number of each tile not required for the winning hand.
	ThrowableVector tilecount.TileCounts34
}

const (
	InfinityShanten  = math.MaxInt
	MaxShantenNumber = 8

	numChowTypes = 7 * 3
	numMeldTypes = tile.NumTileType34 + numChowTypes
)

var (
	// chowStartIDs is a set of IDs of first tile in Chow.
	chowStartIDs = [numChowTypes]int{
		0, 1, 2, 3, 4, 5, 6,
		9, 10, 11, 12, 13, 14, 15,
		18, 19, 20, 21, 22, 23, 24,
	}

	allMelds [numMeldTypes]block.Block = func() [numMeldTypes]block.Block {
		ms := [numMeldTypes]block.Block{}
		for i := range tile.NumTileType34 {
			t := tile.MustTileFromID(i)
			ms[i] = block.MustTriplet(*t)
		}

		for chowID, i := range chowStartIDs {
			t := tile.MustTileFromID(i)
			ms[chowID+tile.NumTileType34] = block.MustSequence(*t)
		}
		return ms
	}()

	allPairs [tile.NumTileType34]block.Block = func() [tile.NumTileType34]block.Block {
		ps := [tile.NumTileType34]block.Block{}
		for i := range tile.NumTileType34 {
			t := tile.MustTileFromID(i)
			ps[i] = block.MustPair(*t)
		}
		return ps
	}()
)

type shantenConfig struct {
	allowedExtraTiles int
	upperBound        int
}

type shantenOption func(*shantenConfig)

func AllowedExtraTiles(n int) shantenOption {
	return func(cfg *shantenConfig) {
		cfg.allowedExtraTiles = n
	}
}

func UpperBound(n int) shantenOption {
	return func(cfg *shantenConfig) {
		cfg.upperBound = n
	}
}

// AnalyzeShanten calculates the shanten number and the list of Goal for the given hand.
// When the list of Goal is empty, `InfinityShanten` is returned as the shanten number.
// It does not consider Seven Pairs or Thirteen Orphans.
func AnalyzeShanten(hand *hand.VisibleHand, opts ...shantenOption) (int, []Goal) {
	cfg := &shantenConfig{
		allowedExtraTiles: 0,
		upperBound:        MaxShantenNumber,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	tc34 := hand.ToTileCounts34()

	targetVector := tilecount.TileCounts34{}
	numRequiredMelds := min(tc34.NumTiles()/3, 4)
	blocks := make([]block.Block, 0, numRequiredMelds+1) // +1 for the pair
	allGoals := []Goal{}

	shanten := analyzeShantenInternal(
		tc34,
		&targetVector,
		-1,
		numRequiredMelds,
		0,
		cfg.upperBound,
		blocks,
		&allGoals,
		cfg.allowedExtraTiles,
	)
	newUpperbound := min(shanten+cfg.allowedExtraTiles, cfg.upperBound)

	// Filter out the goals that exceed newUpperbound
	goals := make([]Goal, 0, len(allGoals))
	for _, goal := range allGoals {
		if goal.Shanten <= newUpperbound {
			for pid := range tile.NumTileType34 {
				goal.RequiredVector[pid] = max(goal.CountVector[pid]-tc34[pid], 0)
				goal.ThrowableVector[pid] = max(tc34[pid]-goal.CountVector[pid], 0)
			}
			goals = append(goals, goal)
		}
	}
	goals = slices.Clip(goals)

	if len(goals) == 0 {
		return InfinityShanten, goals
	}

	return shanten, goals
}

// analyzeShantenInternal calculates the shanten number and
// the set of nearest winning hands using pruning DFS.
func analyzeShantenInternal(
	currentVector *tilecount.TileCounts34,
	targetVector *tilecount.TileCounts34,
	currentShanten int,
	numMeldsLeft int,
	minMeldID int,
	upperbound int,
	blocks []block.Block,
	goals *[]Goal,
	allowedExtraTiles int,
) int {
	if numMeldsLeft == 0 {
		// Add a pair
		for i := range tile.NumTileType34 {
			if targetVector[i] > 2 {
				// Can't add a pair
				continue
			}

			pairDistance := max(targetVector[i]+2-currentVector[i], 0)
			newShanten := currentShanten + pairDistance

			if newShanten <= upperbound+allowedExtraTiles {
				goalVector := *targetVector
				goalVector[i] += 2
				goal := Goal{
					Shanten:     newShanten,
					Blocks:      makeNewBlocks(blocks, allPairs[i]),
					CountVector: goalVector,
				}
				*goals = append(*goals, goal)

				if newShanten < upperbound {
					upperbound = newShanten
				}
			}
		}

		return upperbound
	}

	// Add Pungs
	for i := minMeldID; i < tile.NumTileType34; i++ {
		if targetVector[i] >= 2 {
			// Can't add a Pung
			continue
		}

		pungDistance := 3
		if currentVector[i] > targetVector[i] {
			pungDistance = max(targetVector[i]+3-currentVector[i], 0)
		}
		newShanten := currentShanten + pungDistance

		// If pungDistance == 3:
		// There are no common tiles between currentVector and the target Pung.
		// Therefore, the winning hand containing the target Pung is not the nearest winning hand.
		// Consequently, there is no need to search for a winning hand that contains
		// the target Pung, so this branch is pruned.
		if pungDistance < 3 && newShanten <= upperbound+allowedExtraTiles {
			targetVector[i] += 3
			upperbound = analyzeShantenInternal(
				currentVector,
				targetVector,
				newShanten,
				numMeldsLeft-1,
				i+1, // The same Pung can only be taken out once.
				upperbound,
				makeNewBlocks(blocks, allMelds[i]),
				goals,
				allowedExtraTiles,
			)
			targetVector[i] -= 3
		}
	}

	// Add Chows
	startChowID := max(minMeldID-tile.NumTileType34, 0)
	for chowID := startChowID; chowID < numChowTypes; chowID++ {
		i := chowStartIDs[chowID]
		if targetVector[i] >= 4 || targetVector[i+1] >= 4 || targetVector[i+2] >= 4 {
			// Can't add a Chow
			continue
		}

		chowDistance := 3
		if currentVector[i] > targetVector[i] {
			chowDistance--
		}
		if currentVector[i+1] > targetVector[i+1] {
			chowDistance--
		}
		if currentVector[i+2] > targetVector[i+2] {
			chowDistance--
		}
		newShanten := currentShanten + chowDistance

		// If chowDistance == 3:
		// There are no common tiles between currentVector and the target Chow.
		// Therefore, the winning hand containing the target Chow is not the nearest winning hand.
		// Consequently, there is no need to search for a winning hand that contains
		// the target Chow, so this branch is pruned.
		if chowDistance < 3 && newShanten <= upperbound+allowedExtraTiles {
			targetVector[i]++
			targetVector[i+1]++
			targetVector[i+2]++
			upperbound = analyzeShantenInternal(
				currentVector,
				targetVector,
				newShanten,
				numMeldsLeft-1,
				chowID+tile.NumTileType34,
				upperbound,
				makeNewBlocks(blocks, allMelds[chowID+tile.NumTileType34]),
				goals,
				allowedExtraTiles,
			)
			targetVector[i]--
			targetVector[i+1]--
			targetVector[i+2]--
		}
	}

	return upperbound
}

func makeNewBlocks(blocks []block.Block, newBlock block.Block) []block.Block {
	newBlocks := make([]block.Block, len(blocks), cap(blocks))
	copy(newBlocks, blocks)
	return append(newBlocks, newBlock)
}

func AnalyzeShantenChiitoitsu(hand *hand.VisibleHand) int {
	tc34 := hand.ToTileCounts34()

	if tc34.NumTiles() < 13 {
		return InfinityShanten
	}

	numPairs := 0
	numKinds := 0
	for _, c := range tc34 {
		if c >= 2 {
			numPairs++
		}
		if c >= 1 {
			numKinds++
		}
	}

	shanten := 6 - numPairs
	if numKinds < 7 {
		shanten += 7 - numKinds
	}

	return shanten
}

func AnalyzeShantenKokushimuso(hand *hand.VisibleHand) int {
	tc34 := hand.ToTileCounts34()

	if tc34.NumTiles() < 13 {
		return InfinityShanten
	}

	numKinds := 0
	hasPair := 0
	for _, i := range tile.YaochuhaiIDs {
		if tc34[i] >= 1 {
			numKinds++
		}
		if tc34[i] >= 2 {
			hasPair = 1
		}
	}

	return 13 - numKinds - hasPair
}
