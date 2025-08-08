package inbound

type EndGame struct{}

func NewEndGame() *EndGame {
	return &EndGame{}
}

func (e *EndGame) isInboundEvent() {}
