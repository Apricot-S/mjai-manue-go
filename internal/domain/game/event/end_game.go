package event

type EndGame struct {
}

func NewEndGame() *EndGame {
	return &EndGame{}
}

func (*EndGame) isEvent() {}
