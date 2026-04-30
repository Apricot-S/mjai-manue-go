package inbound

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type Hora struct {
	Type           string   `json:"type"`
	Actor          int      `json:"actor"`
	Target         int      `json:"target"`
	Pai            string   `json:"pai,omitempty"`
	UradoraMarkers []string `json:"uradora_markers,omitempty"`
	UraMarkers     []string `json:"ura_markers,omitempty"`
	HoraTehais     []string `json:"hora_tehais,omitempty"`
	Yakus          [][]any  `json:"yakus,omitempty"`
	Fu             int      `json:"fu,omitempty"`
	Fan            int      `json:"fan,omitempty"`
	HoraPoints     int      `json:"hora_points,omitempty"`
	Deltas         []int    `json:"deltas,omitempty"`
	Scores         []int    `json:"scores,omitempty"`
}

func (*Hora) inboundMessage() {}

func (m *Hora) ToEvent() (*event.Win, error) {
	if m == nil {
		return nil, fmt.Errorf("hora message is nil")
	}
	if m.Type != "hora" {
		return nil, fmt.Errorf("unexpected message type: %q", m.Type)
	}

	actor, err := parseSeatField("actor", m.Actor)
	if err != nil {
		return nil, err
	}
	target, err := parseSeatField("target", m.Target)
	if err != nil {
		return nil, err
	}
	var winningTile *tile.Tile
	if m.Pai != "" {
		winningTile, err = parseKnownTileField("pai", m.Pai)
		if err != nil {
			return nil, err
		}
	}

	deltas, err := parseOptionalScoresField("deltas", m.Deltas)
	if err != nil {
		return nil, err
	}
	scores, err := parseOptionalScoresField("scores", m.Scores)
	if err != nil {
		return nil, err
	}

	return event.NewWin(
		*actor,
		*target,
		winningTile,
		m.HoraPoints,
		deltas,
		scores,
	), nil
}
