package game

import (
	"github.com/go-json-experiment/json/jsontext"
)

const (
	NumPlayers        = 4
	InitScore         = 25_000
	MaxNumDoraMarkers = 5
	NumInitPipais     = NumIDs*4 - 13*NumPlayers - 14
	FinalTurn         = NumInitPipais / NumPlayers

	// Indicates that no event has been triggered.
	noEvent = ""
	// Indicates that no action has been taken by anyone.
	noActor = -1
)

func GetPlayerDistance(p1 *Player, p2 *Player) int {
	return (NumPlayers + p1.ID() - p2.ID()) % NumPlayers
}

func getNextKyoku(bakaze *Pai, kyokuNum int) (*Pai, int) {
	if kyokuNum == 4 {
		return bakaze.NextForDora(), 1
	}
	return bakaze, kyokuNum + 1
}

// StateViewer is an interface for referencing the game state.
type StateViewer interface {
	Players() *[NumPlayers]Player
	Bakaze() *Pai
	KyokuNum() int
	Honba() int
	Oya() *Player
	Chicha() *Player
	DoraMarkers() []Pai
	NumPipais() int

	Anpais(player *Player) []Pai
	VisiblePais(player *Player) []Pai
	Doras() []Pai
	Jikaze(player *Player) *Pai
	YakuhaiFan(pai *Pai, player *Player) int
	NextKyoku() (*Pai, int)
	Turn() int
	RankedPlayers() [NumPlayers]Player

	Print()
}

// StateUpdater is an interface for updating the game state.
type StateUpdater interface {
	OnStartGame(event jsontext.Value) error
	Update(event jsontext.Value) error
}

// ActionCandidatesProvider is an interface for providing action candidates.
type ActionCandidatesProvider interface {
	DahaiCandidates() []Pai
	ReachDahaiCandidates() ([]Pai, error)
	IsTsumoPai(pai *Pai) bool
	FuroCandidates() ([]Furo, error)
	// mjai-manue does not consider Ankan and Kakan, so it is not necessary to implement them.
	HoraCandidate() (*Hora, error)
}

type StateAnalyzer interface {
	StateViewer
	ActionCandidatesProvider
}

type State interface {
	StateViewer
	StateUpdater
	ActionCandidatesProvider
}
