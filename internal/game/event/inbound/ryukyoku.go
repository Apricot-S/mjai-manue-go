package inbound

type Ryukyoku struct {
	Scores *[4]int
}

func NewRyukyoku(scores *[4]int) (*Ryukyoku, error) {
	event := &Ryukyoku{
		Scores: scores,
	}

	if err := eventValidator.Struct(event); err != nil {
		return nil, err
	}
	return event, nil
}

func (s *Ryukyoku) isInboundEvent() {}
