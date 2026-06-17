package outbound

import "github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"

type Daiminkan struct {
	Type     string   `json:"type"`
	Actor    int      `json:"actor"`
	Target   int      `json:"target"`
	Pai      string   `json:"pai"`
	Consumed []string `json:"consumed"`
	Log      string   `json:"log,omitempty"`
}

func NewDaiminkan(a *action.CalledKan, log string) *Daiminkan {
	return &Daiminkan{
		Type:     "daiminkan",
		Actor:    a.Actor().Index(),
		Target:   a.Target().Index(),
		Pai:      a.Taken().String(),
		Consumed: tileCodes3(a.Consumed()),
		Log:      log,
	}
}

func (*Daiminkan) outboundMessage() {}
