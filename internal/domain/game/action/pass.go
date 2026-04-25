package action

type Pass struct {
}

func NewPass() *Pass {
	return &Pass{}
}

func (*Pass) isAction() {}
