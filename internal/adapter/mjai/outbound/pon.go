package outbound

import "github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"

type Pon struct {
	Type     string   `json:"type"`
	Actor    int      `json:"actor"`
	Target   int      `json:"target"`
	Pai      string   `json:"pai"`
	Consumed []string `json:"consumed"`
	Log      string   `json:"log,omitempty"`
}

func NewPon(a *action.Pon, log string) *Pon {
	return &Pon{
		Type:     "pon",
		Actor:    a.Actor().Index(),
		Target:   a.Target().Index(),
		Pai:      a.Taken().String(),
		Consumed: tileCodes2(a.Consumed()),
		Log:      log,
	}
}

func (*Pon) outboundMessage() {}
