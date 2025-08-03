package inbound

type Ryukyoku struct {
	Scores *[4]int
}

func NewRyukyoku(scores *[4]int) *Ryukyoku {
	event := &Ryukyoku{
		Scores: scores,
	}

	return event
}

func (s *Ryukyoku) isInboundEvent() {}
