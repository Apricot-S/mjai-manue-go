package inbound

type Error struct{}

func NewError() *Error {
	return &Error{}
}

func (e *Error) isInboundEvent() {}
