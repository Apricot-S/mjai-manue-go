package event

type EndRound struct {
}

func (e EndRound) EventType() string {
	return "end_round"
}
