package inbound

type EndGame struct{}

func NewEndGame() *EndGame {
	return &EndGame{}
}

func (n *EndGame) isInboundEvent() {}
