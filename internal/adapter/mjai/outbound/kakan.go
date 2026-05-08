package outbound

import "github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"

type Kakan struct {
	Type     string   `json:"type"`
	Actor    int      `json:"actor"`
	Pai      string   `json:"pai"`
	Consumed []string `json:"consumed"`
	Log      string   `json:"log,omitempty"`
}

func NewKakan(a *action.PromotedKan, log string) *Kakan {
	return &Kakan{
		Type:     "kakan",
		Actor:    a.Actor().Index(),
		Pai:      a.Added().String(),
		Consumed: tileCodes3(a.Consumed()),
		Log:      log,
	}
}

func (*Kakan) outboundMessage() {}
