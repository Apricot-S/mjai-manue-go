package message

// "type" is "none"
type Skip struct {
	Action
}

func NewSkip(actor int, log string) *Skip {
	return &Skip{
		Action: Action{
			Message: Message{TypeNone},
			Actor:   actor,
			Log:     log,
		},
	}
}
