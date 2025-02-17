package game_test

import (
	"math/rand/v2"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/game"
)

var choices = [...]int{1, 2, 4, 5, 7, 8, 10, 11, 13, 14}

func createRNG() *rand.Rand {
	src := rand.NewPCG(42, 42)
	rng := rand.New(src)
	return rng
}

func chooseHandLength(rng *rand.Rand) int {
	return choices[rng.IntN(len(choices))]
}

func fillHand(wall *[136]int, handLength int) *[34]int {
	hand := [34]int{}
	for _, tile := range wall[:handLength] {
		hand[tile]++
	}
	return &hand
}

func generateRandomPureHand(rng *rand.Rand) *[34]int {
	wall := [136]int{}
	for i := 0; i < 136; i++ {
		wall[i] = i / 4
	}
	rng.Shuffle(len(wall), func(i, j int) {
		wall[i], wall[j] = wall[j], wall[i]
	})
	handLength := chooseHandLength(rng)
	return fillHand(&wall, handLength)
}

func generateRandomFullPureHand(rng *rand.Rand) *[34]int {
	wall := [136]int{}
	for i := 0; i < 136; i++ {
		wall[i] = i / 4
	}
	rng.Shuffle(len(wall), func(i, j int) {
		wall[i], wall[j] = wall[j], wall[i]
	})
	return fillHand(&wall, 14)
}

func BenchmarkShantenAnalysis_Normal(b *testing.B) {
	rng := createRNG()
	b.ResetTimer()
	b.StopTimer()
	for range b.N {
		hand := generateRandomPureHand(rng)
		ps := game.NewPaiSet(*hand)
		b.StartTimer()
		_, _, _ = game.AnalyzeShanten(ps)
		b.StopTimer()
	}
}

func BenchmarkShantenAnalysis_FullNormal(b *testing.B) {
	rng := createRNG()
	b.ResetTimer()
	b.StopTimer()
	for range b.N {
		hand := generateRandomFullPureHand(rng)
		ps := game.NewPaiSet(*hand)
		b.StartTimer()
		_, _, _ = game.AnalyzeShanten(ps)
		b.StopTimer()
	}
}
