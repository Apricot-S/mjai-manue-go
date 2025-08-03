package inbound

type EndKyoku struct{}

func NewEndKyoku() *EndKyoku {
	return &EndKyoku{}
}

func (n *EndKyoku) isInboundEvent() {}
