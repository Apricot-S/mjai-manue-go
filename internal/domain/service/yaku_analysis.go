package service

import (
	"maps"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/block"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/wind"
)

type WinEvent int

const (
	NoEvent WinEvent = iota + 1
	RobbingAKan
	AfterAKan
	LastTile
	FirstTile
)

// Note:
// This function calculates Fu and Han based on the original mjai-manue specification,
// not the exact Fu and Han.
func CalculateFuHan(
	hand *hand.VisibleHand,
	handBlocks []block.Block,
	melds []meld.Meld,
	prevalentWind wind.Wind,
	seatWind wind.Wind,
	doraIndicators []tile.Tile,
	riichi bool,
) (fu int, han int, yakus map[string]int) {
	meldBlocks := make([]block.Block, len(melds))
	for i, m := range melds {
		meldBlocks[i] = m.ToBlock()
	}

	allBlocks := slices.Concat(handBlocks, meldBlocks)

	allTiles := make(tile.Tiles, 0, 14)
	for _, b := range allBlocks {
		allTiles = append(allTiles, b.ToTiles()...)
	}

	isOpen := len(melds) > 0

	yakus = map[string]int{
		"reach": riichi_(riichi),
		"tyc":   tanyao(allTiles),
		"cty":   chantaiyao(allBlocks, isOpen),
		"pf":    pinfu(allBlocks, prevalentWind, seatWind, isOpen),
		"ykh":   yakuhai(allBlocks, prevalentWind, seatWind),
		"ipk":   iipeikou(allBlocks, isOpen),
		"ssj":   sanshokuDoujun(allBlocks, isOpen),
		"ikt":   ikkiTsuukan(allBlocks, isOpen),
		"tth":   toitoihou(allBlocks),
		"cis":   chiniisou(allBlocks, isOpen),
		"his":   honiisou(allBlocks, isOpen),
	}
	maps.DeleteFunc(yakus, func(k string, v int) bool {
		return v <= 0
	})

	if len(yakus) > 0 {
		numDoras := countDoras(doraIndicators, allTiles)
		if numDoras > 0 {
			yakus["dr"] = numDoras
		}

		numRedDoras := countRedDoras(hand, melds)
		if numRedDoras > 0 {
			yakus["adr"] = numRedDoras
		}
	}

	// TODO Calculate fu more accurately
	_, isPinfu := yakus["pf"]
	if isPinfu || isOpen {
		fu = 30
	} else {
		fu = 40
	}

	han = 0
	for _, h := range yakus {
		han += h
	}

	return fu, han, yakus
}

func countDoras(doraIndicators []tile.Tile, allTiles []tile.Tile) int {
	doras := make([]tile.Tile, len(doraIndicators))
	for i := range doras {
		doras[i] = *doraIndicators[i].NextForDora()
	}

	numDoras := 0
	for i := range allTiles {
		for j := range doras {
			if doras[j].HasSameSymbol(&allTiles[i]) {
				numDoras++
			}
		}
	}

	return numDoras
}

func countRedDoras(hand *hand.VisibleHand, melds []meld.Meld) int {
	numRedDoras := 0

	for _, ht := range hand.ToTiles() {
		if ht.IsRed() {
			numRedDoras++
		}
	}

	for _, m := range melds {
		for _, mt := range m.ToTiles() {
			if mt.IsRed() {
				numRedDoras++
			}
		}
	}

	return numRedDoras
}

func Has1Han(
	hand *hand.VisibleHand,
	melds []meld.Meld,
	winningTile *tile.Tile,
	prevalentWind wind.Wind,
	seatWind wind.Wind,
	tsumo bool,
	riichi bool,
	event WinEvent,
) bool {
	handWithWinningTile, err := hand.Draw(winningTile)
	if err != nil {
		panic(err)
	}

	if IsWinningFormChiitoitsu(handWithWinningTile) || IsWinningFormKokushimusou(handWithWinningTile) {
		return true
	}

	shanten, _ := AnalyzeShanten(handWithWinningTile, UpperBound(-1))
	if shanten > -1 {
		return false
	}

	if riichi {
		return true
	}

	return false
}

func riichi_(isRiichi bool) int {
	if isRiichi {
		return 1
	}
	return 0
}

func tanyao(allTiles []tile.Tile) int {
	for _, t := range allTiles {
		if t.IsYaochu() {
			return 0
		}
	}
	return 1
}

// Note:
// To reproduce the behavior of the original mjai-manue,
// Honroutou is also judged as Chantaiyao.
func chantaiyao(allBlocks []block.Block, isOpen bool) int {
	for _, b := range allBlocks {
		isYaochuBlock := false
		for _, p := range b.ToTiles() {
			if p.IsYaochu() {
				isYaochuBlock = true
				break
			}
		}
		if !isYaochuBlock {
			return 0
		}
	}

	if isOpen {
		return 1
	}
	return 2
}

// TODO Consider ryanmen criteria
func pinfu(allBlocks []block.Block, prevalentWind wind.Wind, seatWind wind.Wind, isOpen bool) int {
	if isOpen {
		return 0
	}

	for _, b := range allBlocks {
		switch b.(type) {
		case *block.Triplet, *block.Quad:
			return 0
		case *block.Pair:
			t := &b.ToTiles()[0]
			if t.IsSuits() {
				continue
			}
			if t.Number() > 4 {
				// dragons
				return 0
			}

			name := t.String()
			if name == prevalentWind.String() || name == seatWind.String() {
				return 0
			}
		}
	}

	return 1
}

func yakuhai(allBlocks []block.Block, prevalentWind wind.Wind, seatWind wind.Wind) int {
	han := 0

	for _, b := range allBlocks {
		switch b.(type) {
		case *block.Triplet, *block.Quad:
			t := b.ToTiles()[0]
			if t.IsSuits() {
				continue
			}
			if t.Number() > 4 {
				// dragons
				han++
				continue
			}

			name := t.String()
			if name == prevalentWind.String() {
				han++
			}
			if name == seatWind.String() {
				han++
			}
		}
	}

	return han
}

func iipeikou(allBlocks []block.Block, isOpen bool) int {
	if isOpen {
		return 0
	}

	for i, b1 := range allBlocks {
		if _, ok := b1.(*block.Sequence); !ok {
			continue
		}
		t1 := b1.ToTiles()[0]

		for _, b2 := range allBlocks[i+1:] {
			if _, ok := b2.(*block.Sequence); !ok {
				continue
			}
			if t1.HasSameSymbol(&b2.ToTiles()[0]) {
				return 1
			}
		}
	}
	return 0
}

func sanshokuDoujun(allBlocks []block.Block, isOpen bool) int {
	colorNumMap := map[rune]map[int]bool{
		tile.ManzuColor: {},
		tile.PinzuColor: {},
		tile.SouzuColor: {},
	}

	for _, b := range allBlocks {
		sequence, ok := b.(*block.Sequence)
		if !ok {
			continue
		}

		t := sequence.ToTiles()[0]
		colorNumMap[t.Color()][t.Number()] = true
	}

	for n := 1; n <= 7; n++ {
		if colorNumMap[tile.ManzuColor][n] &&
			colorNumMap[tile.PinzuColor][n] &&
			colorNumMap[tile.SouzuColor][n] {
			if isOpen {
				return 1
			}
			return 2
		}
	}

	return 0
}

func ikkiTsuukan(allBlocks []block.Block, isOpen bool) int {
	colorNumMap := map[rune]map[int]bool{
		tile.ManzuColor: {},
		tile.PinzuColor: {},
		tile.SouzuColor: {},
	}

	for _, b := range allBlocks {
		sequence, ok := b.(*block.Sequence)
		if !ok {
			continue
		}

		t := sequence.ToTiles()[0]
		colorNumMap[t.Color()][t.Number()] = true
	}

	for _, c := range []rune{tile.ManzuColor, tile.PinzuColor, tile.SouzuColor} {
		if colorNumMap[c][1] && colorNumMap[c][4] && colorNumMap[c][7] {
			if isOpen {
				return 1
			}
			return 2
		}
	}

	return 0
}

func toitoihou(allBlocks []block.Block) int {
	for _, b := range allBlocks {
		if _, ok := b.(*block.Sequence); ok {
			return 0
		}
	}
	return 2
}

func chiniisou(allBlocks []block.Block, isOpen bool) int {
	var color rune

	for i, b := range allBlocks {
		c := b.ToTiles()[0].Color()
		if c == tile.HonorsColor {
			return 0
		}

		if i == 0 {
			color = c
		} else if c != color {
			return 0
		}
	}

	if isOpen {
		return 5
	}
	return 6
}

// Note:
// To reproduce the behavior of the original mjai-manue,
// Tsuuiisou are also judged as Honiisou.
func honiisou(allBlocks []block.Block, isOpen bool) int {
	var color rune = 0
	hasHonors := false

	for _, b := range allBlocks {
		c := b.ToTiles()[0].Color()
		if c == tile.HonorsColor {
			hasHonors = true
			continue
		}

		if color == 0 {
			color = c
		} else if c != color {
			return 0
		}
	}

	if !hasHonors {
		return 0
	}

	if isOpen {
		return 2
	}
	return 3
}
