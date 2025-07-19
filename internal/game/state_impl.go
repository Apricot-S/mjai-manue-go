package game

import (
	"github.com/Apricot-S/mjai-manue-go/internal/message"
)

type StateImpl struct {
	players     [NumPlayers]Player
	bakaze      Pai
	kyokuNum    int
	honba       int
	oya         *Player
	chicha      *Player
	doraMarkers []Pai
	numPipais   int

	prevEventType message.Type
	// -1 if prev action is not dahai
	prevDahaiActor   int
	prevDahaiPai     *Pai
	currentEventType message.Type

	playerID int
	// -1 if there is no action
	lastActor      int
	lastActionType message.Type

	// The tiles that cannot be discarded because they would result in swap calling (喰い替え)
	kuikaePais     []Pai
	missedRon      bool
	isFuriten      bool
	isRinshanTsumo bool
}
