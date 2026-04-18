package game

import "github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"

type StateViewer interface {
	Scores() [common.NumPlayers]int
}

type State struct {
	scores [common.NumPlayers]int
}
