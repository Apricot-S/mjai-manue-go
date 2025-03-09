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

type Game struct {
	prevActionType message.Type
	// -1 if no dahai
	prevDahaiActor    int
	prevDahaiPai      *Pai
	currentActionType message.Type
	tenpais           [numPlayers]bool

	players     [numPlayers]Player
	bakaze      Pai
	kyokuNum    int
	honba       int
	oya         *Player
	chicha      *Player
	doraMarkers []Pai
	numPipais   int
}

func (g *Game) Players() *[numPlayers]Player {
	return &g.players
}

func (g *Game) Bakaze() *Pai {
	return &g.bakaze
}

func (g *Game) KyokuNum() int {
	return g.kyokuNum
}

func (g *Game) Honba() int {
	return g.honba
}

func (g *Game) Oya() *Player {
	return g.oya
}

func (g *Game) Chicha() *Player {
	return g.chicha
}

func (g *Game) DoraMarkers() []Pai {
	return g.doraMarkers
}

func (g *Game) NumPipais() int {
	return g.numPipais
}
