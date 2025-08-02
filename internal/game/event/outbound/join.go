package outbound

type Join struct {
	Name string
	Room string
}

func NewJoin(name string, room string) *Join {
	return &Join{
		Name: name,
		Room: room,
	}
}

func (n *Join) isOutboundEvent() {}
