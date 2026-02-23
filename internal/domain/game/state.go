package game

const NumPlayers = 4

type StateViewer interface {
	Scores() [NumPlayers]int
}

type State struct {
	scores [NumPlayers]int
}
