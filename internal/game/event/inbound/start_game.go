package inbound

type StartGame struct {
	ID    int `validate:"min=0,max=3"`
	Names [4]string
}

func NewStartGame(id int, names [4]string) (*StartGame, error) {
	event := &StartGame{
		ID:    id,
		Names: names,
	}

	if err := eventValidator.Struct(event); err != nil {
		return nil, err
	}
	return event, nil
}

func (s *StartGame) isInboundEvent() {}
