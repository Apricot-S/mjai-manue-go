package inbound

type EndKyoku struct{}

func NewEndKyoku() *EndKyoku {
	return &EndKyoku{}
}

func (e *EndKyoku) isInboundEvent() {}
