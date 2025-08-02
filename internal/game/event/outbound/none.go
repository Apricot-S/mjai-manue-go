package outbound

type None struct{}

func NewNone() *None {
	return &None{}
}

func (n *None) isOutboundEvent() {}
