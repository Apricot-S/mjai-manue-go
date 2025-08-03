package inbound

type Reach struct {
	Actor int `validate:"min=0,max=3"`
}

func NewReach(actor int) (*Reach, error) {
	event := &Reach{
		Actor: actor,
	}

	if err := eventValidator.Struct(event); err != nil {
		return nil, err
	}
	return event, nil
}

func (s *Reach) isInboundEvent() {}
