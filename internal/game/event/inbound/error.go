package inbound

type Error struct{}

func NewError() *Error {
	return &Error{}
}

func (n *Error) isInboundEvent() {}
