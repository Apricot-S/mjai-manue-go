package outbound

import "github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"

type Hora struct {
	Type   string `json:"type"`
	Actor  int    `json:"actor"`
	Target int    `json:"target"`
	Pai    string `json:"pai"`
	Log    string `json:"log,omitempty"`
}

func NewHora(a *action.Win, log string) *Hora {
	return &Hora{
		Type:   "hora",
		Actor:  a.Actor().Index(),
		Target: a.Target().Index(),
		Pai:    a.WinningTile().String(),
		Log:    log,
	}
}

func (*Hora) outboundMessage() {}
