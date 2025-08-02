package outbound

type Skip struct {
	action
}

func NewSkip(actor int, log string) (*Skip, error) {
	event := &Skip{
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
