package inbound

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

type StartKyoku struct {
	Type       string     `json:"type"`
	Bakaze     string     `json:"bakaze"`
	Kyoku      int        `json:"kyoku"`
	Honba      int        `json:"honba"`
	Kyotaku    int        `json:"kyotaku"`
	Oya        int        `json:"oya"`
	DoraMarker string     `json:"dora_marker"`
	Tehais     [][]string `json:"tehais"`
	Scores     []int      `json:"scores,omitempty"`
}

func (*StartKyoku) inboundMessage() {}

func (m *StartKyoku) ToEvent() (*event.StartRound, error) {
	if m == nil {
		return nil, fmt.Errorf("start kyoku message is nil")
	}
	if m.Type != "start_kyoku" {
		return nil, fmt.Errorf("unexpected message type: %q", m.Type)
	}
	if len(m.Tehais) != common.NumPlayers {
		return nil, fmt.Errorf("invalid tehais length: %d", len(m.Tehais))
	}

	var hands [common.NumPlayers][common.InitHandSize]tile.Tile
	for playerIndex, hand := range m.Tehais {
		if len(hand) != common.InitHandSize {
			return nil, fmt.Errorf("invalid hand length for player %d: %d", playerIndex, len(hand))
		}
		for tileIndex, code := range hand {
			t, err := tile.NewTileFromCode(code)
			if err != nil {
				return nil, fmt.Errorf("invalid tile code for player %d index %d: %w", playerIndex, tileIndex, err)
			}
			hands[playerIndex][tileIndex] = t
		}
	}

	scores, err := parseOptionalScoresField("scores", m.Scores)
	if err != nil {
		return nil, err
	}

	roundWind, err := wind.NewWind(m.Bakaze)
	if err != nil {
		return nil, fmt.Errorf("invalid bakaze: %w", err)
	}

	dealerSeat, err := parseSeatField("oya", m.Oya)
	if err != nil {
		return nil, err
	}

	doraIndicator, err := parseKnownTileField("dora_marker", m.DoraMarker)
	if err != nil {
		return nil, err
	}

	return event.NewStartRound(
		roundWind,
		m.Kyoku,
		m.Honba,
		m.Kyotaku,
		*dealerSeat,
		*doraIndicator,
		scores,
		hands,
	), nil
}
