package inbound

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
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
	Scores     *[]int     `json:"scores,omitempty"`
}

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
			hands[playerIndex][tileIndex] = *t
		}
	}

	var scoresPtr *[common.NumPlayers]int
	if m.Scores != nil {
		if len(*m.Scores) != common.NumPlayers {
			return nil, fmt.Errorf("invalid scores length: %d", len(*m.Scores))
		}
		var scores [common.NumPlayers]int
		copy(scores[:], *m.Scores)
		scoresPtr = &scores
	}

	roundWind, err := wind.NewWind(m.Bakaze)
	if err != nil {
		return nil, fmt.Errorf("invalid bakaze: %w", err)
	}

	dealerSeat, err := seat.NewSeat(m.Oya)
	if err != nil {
		return nil, fmt.Errorf("invalid oya: %w", err)
	}

	doraIndicator, err := tile.NewTileFromCode(m.DoraMarker)
	if err != nil {
		return nil, fmt.Errorf("invalid dora marker: %w", err)
	}

	return event.NewStartRound(
		roundWind,
		m.Kyoku,
		m.Honba,
		m.Kyotaku,
		*dealerSeat,
		*doraIndicator,
		scoresPtr,
		hands,
	)
}
