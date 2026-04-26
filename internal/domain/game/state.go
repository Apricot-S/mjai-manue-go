package game

import "github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"

type State struct {
	scores [common.NumPlayers]int
}

func NewState(scores [common.NumPlayers]int) *State {
	return &State{
		scores: scores,
	}
}

func NewDefaultState() *State {
	return NewState([common.NumPlayers]int{25000, 25000, 25000, 25000})
}

func (s *State) Scores() [common.NumPlayers]int {
	return s.scores
}

func (s *State) UpdateScores(scores [common.NumPlayers]int) {
	s.scores = scores
}
