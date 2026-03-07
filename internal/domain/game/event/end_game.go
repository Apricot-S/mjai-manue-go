package event

type EndGame struct {
}

func (e EndGame) EventType() string {
	return "end_game"
}
