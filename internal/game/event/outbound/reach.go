package outbound

type Reach struct {
	action
}

func NewReach(actor int, log string) (*Reach, error) {
	event := &Reach{
		action: action{
			Actor: actor,
			Log:   log,
		},
	}

	if err := eventValidator.Struct(event); err != nil {
		return nil, err
	}
	return event, nil
}
