package outbound

import "github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"

type Chi struct {
	Type     string   `json:"type"`
	Actor    int      `json:"actor"`
	Target   int      `json:"target"`
	Pai      string   `json:"pai"`
	Consumed []string `json:"consumed"`
	Log      string   `json:"log,omitempty"`
}

func NewChi(a *action.Chii, log string) *Chi {
	return &Chi{
		Type:     "chi",
		Actor:    a.Actor().Index(),
		Target:   a.Target().Index(),
		Pai:      a.Taken().String(),
		Consumed: tileCodes2(a.Consumed()),
		Log:      log,
	}
}

func (*Chi) outboundMessage() {}
