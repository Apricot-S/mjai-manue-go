package game

import (
	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
)

type StateImpl struct {
	players     [NumPlayers]base.Player
	bakaze      base.Pai
	kyokuNum    int
	honba       int
	oya         *base.Player
	chicha      *base.Player
	doraMarkers []base.Pai
	numPipais   int

	// -1 if prev action is not dahai
	prevDahaiActor int
	prevDahaiPai   *base.Pai
	currentEvent   inbound.Event
	// -1 if there is no action
	lastActor  int
	lastAction inbound.Event
	kanCount   int
	// Status of the player who call kan (-1: none, 0-3: single player, 4: multiple players)
	kanPlayerStatus int

	playerID int
	// The tiles that cannot be discarded because they would result in swap calling (喰い替え)
	kuikaePais     []base.Pai
	missedRon      bool
	isFuriten      bool
	isRinshanTsumo bool
}
