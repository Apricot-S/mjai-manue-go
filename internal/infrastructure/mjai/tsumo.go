package mjai

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"io"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type Tsumo struct {
	Type  string `json:"type"`
	Actor int    `json:"actor"`
	Pai   string `json:"pai"`
}

func ParseTsumo(r io.Reader) (*event.Draw, error) {
	var msg Tsumo
	dec := jsontext.NewDecoder(r)
	if err := json.UnmarshalDecode(dec, &msg); err != nil {
		return nil, err
	}
	return msg.ToEvent()
}

func (m *Tsumo) ToEvent() (*event.Draw, error) {
	if m == nil {
		return nil, fmt.Errorf("tsumo message is nil")
	}
	if m.Type != "tsumo" {
		return nil, fmt.Errorf("unexpected message type: %q", m.Type)
	}

	actor, err := seat.NewSeat(m.Actor)
	if err != nil {
		return nil, fmt.Errorf("invalid actor: %w", err)
	}

	pai, err := tile.NewTileFromCode(m.Pai)
	if err != nil {
		return nil, fmt.Errorf("invalid pai: %w", err)
	}

	return event.NewDraw(*actor, *pai), nil
}
