package outbound

import "github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"

type Ankan struct {
	Type     string   `json:"type"`
	Actor    int      `json:"actor"`
	Consumed []string `json:"consumed"`
	Log      string   `json:"log,omitempty"`
}

func NewAnkan(a *action.ConcealedKan, log string) *Ankan {
	return &Ankan{
		Type:     "ankan",
		Actor:    a.Actor().Index(),
		Consumed: tileCodes4(a.Consumed()),
		Log:      log,
	}
}

func (*Ankan) outboundMessage() {}
