package message

type None struct {
	Message
}

func NewNone() *None {
	return &None{
		Message: Message{Type: TypeNone},
	}
}
