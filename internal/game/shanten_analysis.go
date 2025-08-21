package game

import (
	"fmt"
	"math"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

// Goal represents a winning hand that can be transitioned from the current hand.
// It contains information about the current hand's shanten number,
// the winning hand's set composition, and the number of required and unrequired tiles.
type Goal struct {
	// Shanten is the shanten number from the current hand to the winning hand.
	// If the current hand and the winning hand are the same, it will be -1.
	Shanten int
	// Mentsus is a list of sets in the winning hand.
	Mentsus []base.Mentsu
	// CountVector is the number of each tile included in the winning hand.
	CountVector base.PaiSet
	// RequiredVector is the number of each tile required for the winning hand.
	RequiredVector base.PaiSet
	// ThrowableVector is the number of each tile not required for the winning hand.
	ThrowableVector base.PaiSet
}

const (
	InfinityShanten  = math.MaxInt
	MaxShantenNumber = 8

	numChows = 7 * 3
)

var (
	// chowStartIDs is a set of IDs of first tile in Chow.
	chowStartIDs = [numChows]uint8{
		0, 1, 2, 3, 4, 5, 6,
		9, 10, 11, 12, 13, 14, 15,
		18, 19, 20, 21, 22, 23, 24,
	}

	allPairs [base.NumIDs]base.Mentsu = func() [base.NumIDs]base.Mentsu {
		p := [base.NumIDs]base.Mentsu{}
		for i := range uint8(base.NumIDs) {
			pai, _ := base.NewPaiWithID(i)
			p[i] = base.NewToitsu(*pai, *pai)
		}
		return p
	}()

	allMelds [base.NumIDs + numChows]base.Mentsu = func() [base.NumIDs + numChows]base.Mentsu {
		m := [base.NumIDs + numChows]base.Mentsu{}
		for i := range uint8(base.NumIDs) {
			pai, _ := base.NewPaiWithID(i)
			m[i] = base.NewKotsu(*pai, *pai, *pai)
		}
		for chowId := range uint8(numChows) {
			i := chowStartIDs[chowId]
			pai, _ := base.NewPaiWithID(i)
			m[chowId+base.NumIDs] = base.NewShuntsu(*pai, *pai.Next(1), *pai.Next(2))
		}
		return m
	}()
)

func countPais(ps *base.PaiSet) (int, error) {
	sum := 0
	for _, c := range ps {
		if c < 0 {
			return 0, fmt.Errorf("negative number of tiles in the PaiSet")
		}
		if c > 4 {
			return 0, fmt.Errorf("more than 4 tiles of the same type in the PaiSet")
		}
		sum += c
	}

	return sum, nil
}

// AnalyzeShanten calculates the shanten number and the list of Goal for the given PaiSet.
// When the list of Goal is empty, [InfinityShanten] is returned as the shanten number.
// It does not consider Seven Pairs or Thirteen Orphans.
func AnalyzeShanten(ps *base.PaiSet) (int, []Goal, error) {
	return AnalyzeShantenWithOption(ps, 0, MaxShantenNumber)
}

func AnalyzeShantenWithOption(ps *base.PaiSet, allowedExtraPais int, upperbound int) (int, []Goal, error) {
	numPais, err := countPais(ps)
	if err != nil {
		return InfinityShanten, nil, err
	}

	targetVector := base.PaiSet{}
	numMentsus := min(numPais/3, 4)
	mentsus := make([]base.Mentsu, 0, numMentsus+1) // +1 for the pair
	allGoals := []Goal{}

	shanten := analyzeShantenInternal(
		ps,
		&targetVector,
		-1,
		numMentsus,
		0,
		upperbound,
		mentsus,
		&allGoals,
		allowedExtraPais,
	)
	newUpperbound := min(shanten+allowedExtraPais, upperbound)

	// Filter out the goals that exceed newUpperbound
	goals := make([]Goal, 0, len(allGoals))
	for _, goal := range allGoals {
		if goal.Shanten <= newUpperbound {
			for pid := range base.NumIDs {
				goal.RequiredVector[pid] = max(goal.CountVector[pid]-ps[pid], 0)
				goal.ThrowableVector[pid] = max(ps[pid]-goal.CountVector[pid], 0)
			}
			goals = append(goals, goal)
		}
	}
	goals = slices.Clip(goals)

	if len(goals) == 0 {
		return InfinityShanten, goals, nil
	}
	return shanten, goals, nil
}

// analyzeShantenInternal calculates the shanten number and
// the set of nearest winning hands using pruning DFS.
func analyzeShantenInternal(
	currentVector *base.PaiSet,
	targetVector *base.PaiSet,
	currentShanten int,
	numMeldsLeft int,
	minMeldId int,
	upperbound int,
	mentsus []base.Mentsu,
	goals *[]Goal,
	allowedExtraPais int,
) int {
	if numMeldsLeft == 0 {
		// Add a pair
		for i := range uint8(base.NumIDs) {
			if targetVector[i] > 2 {
				// Can't add a pair
				continue
			}

			pairDistance := max(targetVector[i]+2-currentVector[i], 0)
			newShanten := currentShanten + pairDistance

			if newShanten <= upperbound+allowedExtraPais {
				goalVector := *targetVector
				goalVector[i] += 2
				goal := Goal{
					Shanten:     newShanten,
					Mentsus:     makeNewMentsus(mentsus, allPairs[i]),
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
	for i := minMeldId; i < base.NumIDs; i++ {
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
		if pungDistance < 3 && newShanten <= upperbound+allowedExtraPais {
			targetVector[i] += 3
			upperbound = analyzeShantenInternal(
				currentVector,
				targetVector,
				newShanten,
				numMeldsLeft-1,
				i+1, // The same Pung can only be taken out once.
				upperbound,
				makeNewMentsus(mentsus, allMelds[i]),
				goals,
				allowedExtraPais,
			)
			targetVector[i] -= 3
		}
	}

	// Add Chows
	startChowId := max(minMeldId-base.NumIDs, 0)
	for chowId := startChowId; chowId < numChows; chowId++ {
		i := chowStartIDs[chowId]
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
		if chowDistance < 3 && newShanten <= upperbound+allowedExtraPais {
			targetVector[i]++
			targetVector[i+1]++
			targetVector[i+2]++
			upperbound = analyzeShantenInternal(
				currentVector,
				targetVector,
				newShanten,
				numMeldsLeft-1,
				chowId+base.NumIDs,
				upperbound,
				makeNewMentsus(mentsus, allMelds[chowId+base.NumIDs]),
				goals,
				allowedExtraPais,
			)
			targetVector[i]--
			targetVector[i+1]--
			targetVector[i+2]--
		}
	}

	return upperbound
}

func makeNewMentsus(mentsus []base.Mentsu, newMentsu base.Mentsu) []base.Mentsu {
	newMentsus := make([]base.Mentsu, len(mentsus), cap(mentsus))
	copy(newMentsus, mentsus)
	return append(newMentsus, newMentsu)
}

func AnalyzeShantenChitoitsu(ps *base.PaiSet) (int, error) {
	numPais, err := countPais(ps)
	if err != nil {
		return InfinityShanten, err
	}

	if numPais < 13 {
		return InfinityShanten, nil
	}

	numPairs := 0
	numKinds := 0
	for _, c := range ps {
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

	return shanten, nil
}

func AnalyzeShantenKokushimuso(ps *base.PaiSet) (int, error) {
	numPais, err := countPais(ps)
	if err != nil {
		return InfinityShanten, err
	}

	if numPais < 13 {
		return InfinityShanten, nil
	}

	numKinds := 0
	hasPair := 0
	for _, i := range yaochuhaiIndices {
		if ps[i] >= 1 {
			numKinds++
		}
		if ps[i] >= 2 {
			hasPair = 1
		}
	}

	return 13 - numKinds - hasPair, nil
}
