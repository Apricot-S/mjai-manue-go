package event

type EndRound struct {
}

func NewEndRound() *EndRound {
	return &EndRound{}
}

func (e EndRound) EventType() string {
	return "end_round"
}
