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

func shuffleWall(rng *rand.Rand, wall []int) {
	rng.Shuffle(len(wall), func(i, j int) {
		wall[i], wall[j] = wall[j], wall[i]
	})
}

func fillHand(wall []int, handLength int) *[34]int {
	hand := [34]int{}
	for _, tile := range wall[:handLength] {
		hand[tile]++
	}
	return &hand
}

func generateRandomPureHandImpl(rng *rand.Rand, handLength int) *[34]int {
	wall := [136]int{}
	for i := range wall {
		wall[i] = i / 4
	}
	shuffleWall(rng, wall[:])
	return fillHand(wall[:], handLength)
}

func generateRandomPureHand(rng *rand.Rand) *[34]int {
	handLength := chooseHandLength(rng)
	return generateRandomPureHandImpl(rng, handLength)
}

func generateRandomPureHand14(rng *rand.Rand) *[34]int {
	return generateRandomPureHandImpl(rng, 14)
}

func generateRandomHalfFlushPureHandImpl(rng *rand.Rand, handLength int) *[34]int {
	colorStartOptions := [...]int{0, 9, 18}
	colorStart := colorStartOptions[rng.IntN(len(colorStartOptions))]

	wall := [64]int{}
	for i := range wall {
		if i < 36 {
			wall[i] = i/4 + colorStart
		} else {
			wall[i] = (i-36)/4 + 27
		}
	}
	shuffleWall(rng, wall[:])
	return fillHand(wall[:], handLength)
}

func generateRandomHalfFlushPureHand(rng *rand.Rand) *[34]int {
	handLength := chooseHandLength(rng)
	return generateRandomHalfFlushPureHandImpl(rng, handLength)
}

func generateRandomHalfFlushPureHand14(rng *rand.Rand) *[34]int {
	return generateRandomHalfFlushPureHandImpl(rng, 14)
}

func generateRandomFullFlushPureHandImpl(rng *rand.Rand, handLength int) *[34]int {
	colorStartOptions := [...]int{0, 9, 18}
	colorStart := colorStartOptions[rng.IntN(len(colorStartOptions))]

	wall := [36]int{}
	for i := range wall {
		wall[i] = i/4 + colorStart
	}
	shuffleWall(rng, wall[:])
	return fillHand(wall[:], handLength)
}

func generateRandomFullFlushPureHand(rng *rand.Rand) *[34]int {
	handLength := chooseHandLength(rng)
	return generateRandomFullFlushPureHandImpl(rng, handLength)
}

func generateRandomFullFlushPureHand14(rng *rand.Rand) *[34]int {
	return generateRandomFullFlushPureHandImpl(rng, 14)
}

var nonSimples = [...]int{0, 8, 9, 17, 18, 26, 27, 28, 29, 30, 31, 32, 33}

func generateRandomNonSimplePureHandImpl(rng *rand.Rand, handLength int) *[34]int {
	wall := [52]int{}
	for i := range wall {
		wall[i] = nonSimples[i%13]
	}
	shuffleWall(rng, wall[:])
	return fillHand(wall[:], handLength)
}

func generateRandomNonSimplePureHand(rng *rand.Rand) *[34]int {
	handLength := chooseHandLength(rng)
	return generateRandomNonSimplePureHandImpl(rng, handLength)
}

func generateRandomNonSimplePureHand14(rng *rand.Rand) *[34]int {
	return generateRandomNonSimplePureHandImpl(rng, 14)
}

func BenchmarkShantenAnalysis_Normal(b *testing.B) {
	rng := createRNG()
	b.StopTimer()
	b.ResetTimer()
	for range b.N {
		hand := generateRandomPureHand(rng)
		var ps game.PaiSet = *hand
		b.StartTimer()
		_, _, _ = game.AnalyzeShanten(&ps)
		b.StopTimer()
	}
}

func BenchmarkShantenAnalysis_Normal14(b *testing.B) {
	rng := createRNG()
	b.StopTimer()
	b.ResetTimer()
	for range b.N {
		hand := generateRandomPureHand14(rng)
		var ps game.PaiSet = *hand
		b.StartTimer()
		_, _, _ = game.AnalyzeShanten(&ps)
		b.StopTimer()
	}
}

func BenchmarkShantenAnalysis_HalfFlush(b *testing.B) {
	rng := createRNG()
	b.StopTimer()
	b.ResetTimer()
	for range b.N {
		hand := generateRandomHalfFlushPureHand(rng)
		var ps game.PaiSet = *hand
		b.StartTimer()
		_, _, _ = game.AnalyzeShanten(&ps)
		b.StopTimer()
	}
}

func BenchmarkShantenAnalysis_HalfFlush14(b *testing.B) {
	rng := createRNG()
	b.StopTimer()
	b.ResetTimer()
	for range b.N {
		hand := generateRandomHalfFlushPureHand14(rng)
		var ps game.PaiSet = *hand
		b.StartTimer()
		_, _, _ = game.AnalyzeShanten(&ps)
		b.StopTimer()
	}
}

func BenchmarkShantenAnalysis_FullFlush(b *testing.B) {
	rng := createRNG()
	b.StopTimer()
	b.ResetTimer()
	for range b.N {
		hand := generateRandomFullFlushPureHand(rng)
		var ps game.PaiSet = *hand
		b.StartTimer()
		_, _, _ = game.AnalyzeShanten(&ps)
		b.StopTimer()
	}
}

func BenchmarkShantenAnalysis_FullFlush14(b *testing.B) {
	rng := createRNG()
	b.StopTimer()
	b.ResetTimer()
	for range b.N {
		hand := generateRandomFullFlushPureHand14(rng)
		var ps game.PaiSet = *hand
		b.StartTimer()
		_, _, _ = game.AnalyzeShanten(&ps)
		b.StopTimer()
	}
}

func BenchmarkShantenAnalysis_NonSimple(b *testing.B) {
	rng := createRNG()
	b.StopTimer()
	b.ResetTimer()
	for range b.N {
		hand := generateRandomNonSimplePureHand(rng)
		var ps game.PaiSet = *hand
		b.StartTimer()
		_, _, _ = game.AnalyzeShanten(&ps)
		b.StopTimer()
	}
}

func BenchmarkShantenAnalysis_NonSimple14(b *testing.B) {
	rng := createRNG()
	b.StopTimer()
	b.ResetTimer()
	for range b.N {
		hand := generateRandomNonSimplePureHand14(rng)
		var ps game.PaiSet = *hand
		b.StartTimer()
		_, _, _ = game.AnalyzeShanten(&ps)
		b.StopTimer()
	}
}
