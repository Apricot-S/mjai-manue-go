package outbound

import "github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"

type Kyushukyuhai struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
	Actor  int    `json:"actor"`
	Log    string `json:"log,omitempty"`
}

func NewKyushukyuhai(a *action.Kyushukyuhai, log string) *Kyushukyuhai {
	return &Kyushukyuhai{
		Type:   "ryukyoku",
		Reason: "kyushukyuhai",
		Actor:  a.Actor().Index(),
		Log:    log,
	}
}

func (*Kyushukyuhai) outboundMessage() {}
