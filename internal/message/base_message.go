package message

type Message interface {
	GetType() Type
}

type BaseMessage struct {
	Type Type `json:"type" validate:"required"`
}

func (m *BaseMessage) GetType() Type {
	return m.Type
}
