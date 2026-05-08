package service_test

import (
	"math/rand/v2"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
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

func fillHand(wall []int, handLength int) *hand.VisibleHand {
	tc34 := [34]int{}
	for _, tile := range wall[:handLength] {
		tc34[tile]++
	}

	tiles := make([]tile.Tile, 0, handLength)
	for id, count := range tc34 {
		for i := range count {
			t := tile.MustTileFromID(id)
			if t.Number() == 5 && t.IsSuits() && i == 3 {
				t = t.AddRed()
			}

			tiles = append(tiles, t)
		}
	}

	return hand.MustVisibleHand(tiles)
}

func generateRandomPureHandImpl(rng *rand.Rand, handLength int) *hand.VisibleHand {
	wall := [136]int{}
	for i := range wall {
		wall[i] = i / 4
	}
	shuffleWall(rng, wall[:])
	return fillHand(wall[:], handLength)
}

func generateRandomPureHand(rng *rand.Rand) *hand.VisibleHand {
	handLength := chooseHandLength(rng)
	return generateRandomPureHandImpl(rng, handLength)
}

func generateRandomPureHand14(rng *rand.Rand) *hand.VisibleHand {
	return generateRandomPureHandImpl(rng, 14)
}

func generateRandomHalfFlushPureHandImpl(rng *rand.Rand, handLength int) *hand.VisibleHand {
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

func generateRandomHalfFlushPureHand(rng *rand.Rand) *hand.VisibleHand {
	handLength := chooseHandLength(rng)
	return generateRandomHalfFlushPureHandImpl(rng, handLength)
}

func generateRandomHalfFlushPureHand14(rng *rand.Rand) *hand.VisibleHand {
	return generateRandomHalfFlushPureHandImpl(rng, 14)
}

func generateRandomFullFlushPureHandImpl(rng *rand.Rand, handLength int) *hand.VisibleHand {
	colorStartOptions := [...]int{0, 9, 18}
	colorStart := colorStartOptions[rng.IntN(len(colorStartOptions))]

	wall := [36]int{}
	for i := range wall {
		wall[i] = i/4 + colorStart
	}
	shuffleWall(rng, wall[:])
	return fillHand(wall[:], handLength)
}

func generateRandomFullFlushPureHand(rng *rand.Rand) *hand.VisibleHand {
	handLength := chooseHandLength(rng)
	return generateRandomFullFlushPureHandImpl(rng, handLength)
}

func generateRandomFullFlushPureHand14(rng *rand.Rand) *hand.VisibleHand {
	return generateRandomFullFlushPureHandImpl(rng, 14)
}

var nonSimples = [...]int{0, 8, 9, 17, 18, 26, 27, 28, 29, 30, 31, 32, 33}

func generateRandomNonSimplePureHandImpl(rng *rand.Rand, handLength int) *hand.VisibleHand {
	wall := [52]int{}
	for i := range wall {
		wall[i] = nonSimples[i%13]
	}
	shuffleWall(rng, wall[:])
	return fillHand(wall[:], handLength)
}

func generateRandomNonSimplePureHand(rng *rand.Rand) *hand.VisibleHand {
	handLength := chooseHandLength(rng)
	return generateRandomNonSimplePureHandImpl(rng, handLength)
}

func generateRandomNonSimplePureHand14(rng *rand.Rand) *hand.VisibleHand {
	return generateRandomNonSimplePureHandImpl(rng, 14)
}

func BenchmarkShantenAnalysis_Normal(b *testing.B) {
	rng := createRNG()
	for b.Loop() {
		b.StopTimer()
		hand := generateRandomPureHand(rng)
		b.StartTimer()
		_, _ = service.AnalyzeShanten(hand)
	}
}

func BenchmarkShantenAnalysis_HalfFlush(b *testing.B) {
	rng := createRNG()
	for b.Loop() {
		b.StopTimer()
		hand := generateRandomHalfFlushPureHand(rng)
		b.StartTimer()
		_, _ = service.AnalyzeShanten(hand)
	}
}

func BenchmarkShantenAnalysis_FullFlush(b *testing.B) {
	rng := createRNG()
	for b.Loop() {
		b.StopTimer()
		hand := generateRandomFullFlushPureHand(rng)
		b.StartTimer()
		_, _ = service.AnalyzeShanten(hand)
	}
}

func BenchmarkShantenAnalysis_NonSimple(b *testing.B) {
	rng := createRNG()
	for b.Loop() {
		b.StopTimer()
		hand := generateRandomNonSimplePureHand(rng)
		b.StartTimer()
		_, _ = service.AnalyzeShanten(hand)
	}
}

func BenchmarkShantenAnalysis_14_Normal(b *testing.B) {
	rng := createRNG()
	for b.Loop() {
		b.StopTimer()
		hand := generateRandomPureHand14(rng)
		b.StartTimer()
		_, _ = service.AnalyzeShanten(hand)
	}
}

func BenchmarkShantenAnalysis_14_HalfFlush(b *testing.B) {
	rng := createRNG()
	for b.Loop() {
		b.StopTimer()
		hand := generateRandomHalfFlushPureHand14(rng)
		b.StartTimer()
		_, _ = service.AnalyzeShanten(hand)
	}
}

func BenchmarkShantenAnalysis_14_FullFlush(b *testing.B) {
	rng := createRNG()
	for b.Loop() {
		b.StopTimer()
		hand := generateRandomFullFlushPureHand14(rng)
		b.StartTimer()
		_, _ = service.AnalyzeShanten(hand)
	}
}

func BenchmarkShantenAnalysis_14_NonSimple(b *testing.B) {
	rng := createRNG()
	for b.Loop() {
		b.StopTimer()
		hand := generateRandomNonSimplePureHand14(rng)
		b.StartTimer()
		_, _ = service.AnalyzeShanten(hand)
	}
}
