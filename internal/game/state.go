package game

import (
	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
)

const (
	NumPlayers        = 4
	InitScore         = 25_000
	MaxNumDoraMarkers = 5
	NumInitPipais     = base.NumIDs*4 - 13*NumPlayers - 14
	FinalTurn         = float64(NumInitPipais) / float64(NumPlayers)

	// Indicates that no action has been taken by anyone.
	noActor = -1
)

func GetPlayerDistance(p1 *base.Player, p2 *base.Player) int {
	return (NumPlayers + p1.ID() - p2.ID()) % NumPlayers
}

func getNextKyoku(bakaze *base.Pai, kyokuNum int) (*base.Pai, int) {
	if kyokuNum == 4 {
		return bakaze.NextForDora(), 1
	}
	return bakaze, kyokuNum + 1
}

// StateViewer is an interface for referencing the game state.
type StateViewer interface {
	Players() *[NumPlayers]base.Player
	Bakaze() *base.Pai
	KyokuNum() int
	Honba() int
	Oya() *base.Player
	Chicha() *base.Player
	DoraMarkers() []base.Pai
	NumPipais() int

	Anpais(player *base.Player) []base.Pai
	VisiblePais(player *base.Player) []base.Pai
	Doras() []base.Pai
	Jikaze(player *base.Player) *base.Pai
	YakuhaiFan(pai *base.Pai, player *base.Player) int
	NextKyoku() (*base.Pai, int)
	Turn() float64
	RankedPlayers() [NumPlayers]base.Player

	Print()
}

// StateUpdater is an interface for updating the game state.
type StateUpdater interface {
	Update(event inbound.Event) error
}

// ActionCandidatesProvider is an interface for providing action candidates.
type ActionCandidatesProvider interface {
	DahaiCandidates() []base.Pai
	ReachDahaiCandidates() ([]base.Pai, error)
	IsTsumoPai(pai *base.Pai) bool
	FuroCandidates() ([]base.Furo, error)
	// mjai-manue does not consider Ankan and Kakan, so it is not necessary to implement them.
	HoraCandidate() (*base.Hora, error)
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
