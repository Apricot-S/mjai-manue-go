package game

import "github.com/Apricot-S/mjai-manue-go/internal/message"

const (
	numPlayers        = 4
	initScore         = 25_000
	maxNumDoraMarkers = 5
	numInitPipais     = NumIDs*4 - 13*numPlayers - 14
	finalTurn         = numInitPipais / numPlayers
)

func getDistance(p1 *Player, p2 *Player) int {
	return (numPlayers + p1.ID() - p2.ID()) % numPlayers
}

func getNextKyoku(bakaze *Pai, kyokuNum int) (*Pai, int) {
	if kyokuNum == 4 {
		return bakaze.NextForDora(), 1
	}
	return bakaze, kyokuNum + 1
}

type State struct {
	players     [numPlayers]Player
	bakaze      Pai
	kyokuNum    int
	honba       int
	oya         *Player
	chicha      *Player
	doraMarkers []Pai
	numPipais   int

	prevActionType message.Type
	// -1 if no dahai
	prevDahaiActor    int
	prevDahaiPai      *Pai
	currentActionType message.Type
	tenpais           [numPlayers]bool
}

func (s *State) Players() *[numPlayers]Player {
	return &s.players
}

func (s *State) Bakaze() *Pai {
	return &s.bakaze
}

func (s *State) KyokuNum() int {
	return s.kyokuNum
}

func (s *State) Honba() int {
	return s.honba
}

func (s *State) Oya() *Player {
	return s.oya
}

func (s *State) Chicha() *Player {
	return s.chicha
}

func (s *State) DoraMarkers() []Pai {
	return s.doraMarkers
}

func (s *State) NumPipais() int {
	return s.numPipais
}
