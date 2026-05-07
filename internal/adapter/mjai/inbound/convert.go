package inbound

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func parseSeatField(name string, value int) (*seat.Seat, error) {
	s, err := seat.NewSeat(value)
	if err != nil {
		return nil, fmt.Errorf("invalid %s: %w", name, err)
	}
	return &s, nil
}

func parseTileField(name string, value string) (*tile.Tile, error) {
	if value == "" {
		return nil, fmt.Errorf("missing %s", name)
	}
	t, err := tile.NewTileFromCode(value)
	if err != nil {
		return nil, fmt.Errorf("invalid %s: %w", name, err)
	}
	return &t, nil
}

func parseKnownTileField(name string, value string) (*tile.Tile, error) {
	t, err := parseTileField(name, value)
	if err != nil {
		return nil, err
	}
	if t.IsUnknown() {
		return nil, fmt.Errorf("%s must not be unknown", name)
	}
	return t, nil
}

func parseConsumed2(values []string) ([2]tile.Tile, error) {
	if len(values) != 2 {
		return [2]tile.Tile{}, fmt.Errorf("consumed must contain 2 tiles, got %d", len(values))
	}
	var consumed [2]tile.Tile
	for i, value := range values {
		t, err := parseKnownTileField(fmt.Sprintf("consumed[%d]", i), value)
		if err != nil {
			return [2]tile.Tile{}, err
		}
		consumed[i] = *t
	}
	return consumed, nil
}

func parseConsumed3(values []string) ([3]tile.Tile, error) {
	if len(values) != 3 {
		return [3]tile.Tile{}, fmt.Errorf("consumed must contain 3 tiles, got %d", len(values))
	}
	var consumed [3]tile.Tile
	for i, value := range values {
		t, err := parseKnownTileField(fmt.Sprintf("consumed[%d]", i), value)
		if err != nil {
			return [3]tile.Tile{}, err
		}
		consumed[i] = *t
	}
	return consumed, nil
}

func parseConsumed4(values []string) ([4]tile.Tile, error) {
	if len(values) != 4 {
		return [4]tile.Tile{}, fmt.Errorf("consumed must contain 4 tiles, got %d", len(values))
	}
	var consumed [4]tile.Tile
	for i, value := range values {
		t, err := parseKnownTileField(fmt.Sprintf("consumed[%d]", i), value)
		if err != nil {
			return [4]tile.Tile{}, err
		}
		consumed[i] = *t
	}
	return consumed, nil
}

func parseScoresField(name string, values []int) ([common.NumPlayers]int, error) {
	if len(values) != common.NumPlayers {
		return [common.NumPlayers]int{}, fmt.Errorf("%s must contain %d values, got %d", name, common.NumPlayers, len(values))
	}
	var parsed [common.NumPlayers]int
	copy(parsed[:], values)
	return parsed, nil
}

func parseOptionalScoresField(name string, values []int) (*[common.NumPlayers]int, error) {
	if values == nil {
		return nil, nil
	}
	parsed, err := parseScoresField(name, values)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func parseTenpaisField(values []bool) ([common.NumPlayers]bool, error) {
	if len(values) != common.NumPlayers {
		return [common.NumPlayers]bool{}, fmt.Errorf("tenpais must contain %d values, got %d", common.NumPlayers, len(values))
	}
	var parsed [common.NumPlayers]bool
	copy(parsed[:], values)
	return parsed, nil
}

func parseOptionalTenpaisField(values []bool) (*[common.NumPlayers]bool, error) {
	if values == nil {
		return nil, nil
	}
	parsed, err := parseTenpaisField(values)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
