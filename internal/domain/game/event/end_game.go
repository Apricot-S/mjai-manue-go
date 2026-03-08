package event

type EndGame struct {
}

func NewEndGame() *EndGame {
	return &EndGame{}
}

func (e EndGame) EventType() string {
	return "end_game"
}
