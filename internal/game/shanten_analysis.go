package game

import (
	"fmt"
	"math"
	"slices"
)

// Goal represents a winning form that can be transitioned from the current hand.
// It contains information about the current hand's shanten number,
// the winning form's set composition, and the number of required and unrequired tiles.
type Goal struct {
	// Shanten is the shanten number of the current hand.
	Shanten int
	// Mentsus is a list of sets in the winning form.
	Mentsus []Mentsu
	// CountVector is the number of each tile included in the winning form.
	CountVector [NumIDs]int
	// RequiredVector is the number of each tile required for the winning form.
	RequiredVector [NumIDs]int
	// ThrowableVector is the number of each tile not required for the winning form.
	ThrowableVector [NumIDs]int
}

const (
	maxShantenNumber = 8
	numChows         = 21
)

var (
	// 1-7{m,p,s}
	chows = [...]uint8{
		0, 1, 2, 3, 4, 5, 6,
		9, 10, 11, 12, 13, 14, 15,
		18, 19, 20, 21, 22, 23, 24,
	}
)

// AnalyzeShanten calculates the shanten number of the given PaiSet.
// It returns the shanten number and a list of [Goal].
// It does not consider Seven Pairs or Thirteen Orphans.
func AnalyzeShanten(ps *PaiSet) (int, []Goal, error) {
	return AnalyzeShantenWithOption(ps, 0, maxShantenNumber)
}

func AnalyzeShantenWithOption(ps *PaiSet, allowedExtraPais int, upperbound int) (int, []Goal, error) {
	currentVector := ps.Array()
	if slices.Min(currentVector[:]) < 0 {
		return -1, nil, fmt.Errorf("negative number of tiles in the PaiSet")
	}

	targetVector := [NumIDs]int{}
	goals := []Goal{}
	numMentsus := sum(currentVector) / 3
	if numMentsus > 4 {
		numMentsus = 4
	}

	shanten :=
		analyzeShantenInternal(
			&currentVector,
			&targetVector,
			-1,
			numMentsus,
			0,
			upperbound,
			[]Mentsu{},
			goals,
			allowedExtraPais)
	newUpperbound := min(shanten+allowedExtraPais, upperbound)

	for i, goal := range goals {
		if goal.Shanten <= newUpperbound {
			requiredVector := [NumIDs]int{}
			for pid := 0; pid < int(NumIDs); pid++ {
				requiredVector[pid] = max(goal.CountVector[pid]-currentVector[pid], 0)
			}
			goals[i].RequiredVector = requiredVector

			throwableVector := [NumIDs]int{}
			for pid := 0; pid < int(NumIDs); pid++ {
				throwableVector[pid] = max(currentVector[pid]-goal.CountVector[pid], 0)
			}
			goals[i].ThrowableVector = throwableVector
		}
	}

	if len(goals) == 0 {
		return math.MaxInt, []Goal{}, nil
	}
	return shanten, goals, nil
}

func sum(arr [NumIDs]int) int {
	sum := 0
	for _, v := range arr {
		sum += v
	}
	return sum
}

func analyzeShantenInternal(
	currentVector *[NumIDs]int,
	targetVector *[NumIDs]int,
	currentShanten int,
	numMeldsLeft int,
	minMeldId uint8,
	upperbound int,
	mentsus []Mentsu,
	goals []Goal,
	allowedExtraPais int,
) int {
	if numMeldsLeft == 0 {
		// Add a pair
		for i := uint8(0); i < NumIDs; i++ {
			if targetVector[i] > 2 {
				// Can't add a pair
				continue
			}

			pairDistance := max(targetVector[i]+2-currentVector[i], 0)
			newShanten := currentShanten + pairDistance

			if newShanten <= upperbound+allowedExtraPais {
				pai, _ := NewPaiWithID(i)
				pais := []Pai{*pai, *pai}
				toitsu, _ := NewMentsu(Toitsu, pais)
				newMentsus := append(mentsus, *toitsu)

				goalVector := *targetVector
				goalVector[i] += 2
				goal := Goal{
					Shanten:     newShanten,
					Mentsus:     newMentsus,
					CountVector: goalVector,
				}
				goals = append(goals, goal)

				if newShanten < upperbound {
					upperbound = newShanten
				}
			}
		}

		return upperbound
	}

	// Add Pungs
	if minMeldId < NumIDs {
		for i := minMeldId; i < NumIDs; i++ {
			if targetVector[i] >= 2 {
				// Can't add a Pung
				continue
			}

			pungDistance := 0
			if currentVector[i] <= targetVector[i] {
				pungDistance = 3
			} else if currentVector[i] < targetVector[i]+3 {
				pungDistance = (targetVector[i] + 3) - currentVector[i]
			}
			newShanten := currentShanten + pungDistance

			if pungDistance < 3 && newShanten <= upperbound+allowedExtraPais {
				pai, _ := NewPaiWithID(i)
				pais := []Pai{*pai, *pai, *pai}
				kotsu, _ := NewMentsu(Kotsu, pais)
				newMentsus := append(mentsus, *kotsu)

				targetVector[i] += 3
				upperbound = analyzeShantenInternal(
					currentVector,
					targetVector,
					newShanten,
					numMeldsLeft-1,
					i,
					upperbound,
					newMentsus,
					goals,
					allowedExtraPais,
				)
				targetVector[i] -= 3
			}
		}
	}

	// Add Chows
	startChowId := uint8(0)
	if minMeldId < NumIDs {
		startChowId = 0
	} else {
		startChowId = minMeldId - NumIDs
	}

	for chowId := startChowId; chowId < numChows; chowId++ {
		i := chows[chowId]
		if targetVector[i] >= 4 || targetVector[i+1] >= 4 || targetVector[i+2] >= 4 {
			// Can't add a Chow
			continue
		}

		chowDistance := 0
		if currentVector[i] <= targetVector[i] {
			chowDistance++
		}
		if currentVector[i+1] <= targetVector[i+1] {
			chowDistance++
		}
		if currentVector[i+2] <= targetVector[i+2] {
			chowDistance++
		}
		newShanten := currentShanten + chowDistance

		if chowDistance < 3 && newShanten <= upperbound+allowedExtraPais {
			pai, _ := NewPaiWithID(i)
			pais := []Pai{*pai, *pai.Next(1), *pai.Next(2)}
			shuntsu, _ := NewMentsu(Shuntsu, pais)
			newMentsus := append(mentsus, *shuntsu)

			targetVector[i]++
			targetVector[i+1]++
			targetVector[i+2]++
			upperbound = analyzeShantenInternal(
				currentVector,
				targetVector,
				newShanten,
				numMeldsLeft-1,
				chowId+NumIDs,
				upperbound,
				newMentsus,
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
