package outbound

import "github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"

type Reach struct {
	Type  string `json:"type"`
	Actor int    `json:"actor"`
	Log   string `json:"log,omitempty"`
}

func NewReach(a *action.Riichi, log string) *Reach {
	return &Reach{
		Type:  "reach",
		Actor: a.Actor().Index(),
		Log:   log,
	}
}

func (*Reach) outboundMessage() {}
