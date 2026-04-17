package event

type EndRound struct {
}

func NewEndRound() *EndRound {
	return &EndRound{}
}

func (*EndRound) isEvent() {}
