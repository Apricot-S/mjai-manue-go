package outbound

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
)

type Dahai struct {
	Type      string `json:"type"`
	Actor     int    `json:"actor"`
	Pai       string `json:"pai"`
	Tsumogiri bool   `json:"tsumogiri"`
	Log       string `json:"log,omitempty"`
}

func NewDahai(a *action.Discard, log string) *Dahai {
	return &Dahai{
		Type:      "dahai",
		Actor:     a.Actor().Index(),
		Pai:       a.Tile().String(),
		Tsumogiri: a.Tsumogiri(),
		Log:       log,
	}
}

func (*Dahai) outboundMessage() {}
