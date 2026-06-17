package outbound

import "github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"

type Pass struct {
	Type  string `json:"type"`
	Actor int    `json:"actor"`
	Log   string `json:"log,omitempty"`
}

func NewPass(a *action.Pass, log string) *Pass {
	return &Pass{
		Type:  "none",
		Actor: a.Actor().Index(),
		Log:   log,
	}
}

func (*Pass) outboundMessage() {}
