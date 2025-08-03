package inbound

type ReachAccepted struct {
	Actor  int `validate:"min=0,max=3"`
	Scores *[4]int
}

func NewReachAccepted(actor int, scores *[4]int) (*ReachAccepted, error) {
	event := &ReachAccepted{
		Actor:  actor,
		Scores: scores,
	}

	if err := eventValidator.Struct(event); err != nil {
		return nil, err
	}
	return event, nil
}

func (s *ReachAccepted) isInboundEvent() {}
