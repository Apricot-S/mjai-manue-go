package service

import (
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
	tsumo bool,
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

	hasMeld := len(melds) > 0

	han = 0
	yakus = make(map[string]int)

	addYaku := func(existsYaku bool, name string, menzenFan int, kuiFan int) {
		if !existsYaku {
			return
		}

		var yakuFan int
		if hasMeld {
			yakuFan = kuiFan
		} else {
			yakuFan = menzenFan
		}

		if yakuFan > 0 {
			yakus[name] = yakuFan
			han += yakuFan
		}
	}

	addYaku(riichi, "reach", 1, 0)

	pinfu := false

	// TODO Calculate fu more accurately
	if pinfu || hasMeld {
		fu = 30
	} else {
		fu = 40
	}

	return fu, han, yakus
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
	return false
}
