package inbound

type Hello struct{}

func NewHello() *Hello {
	return &Hello{}
}

func (h *Hello) isInboundEvent() {}
