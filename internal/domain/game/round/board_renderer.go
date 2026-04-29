package round

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type BoardRenderer interface {
	RenderBoard() string
}

func (s *State) RenderBoard() string {
	var b strings.Builder

	fmt.Fprintf(&b, "%s-%d kyoku %d honba  ", s.RoundWind(), s.RoundNumber(), s.Honba())
	fmt.Fprintf(&b, "pipai: %d  ", s.NumLeftTiles())
	fmt.Fprintf(&b, "dora_marker: %s  ", formatDoraIndicators(s.DoraIndicators()))
	b.WriteByte('\n')

	for i := range common.NumPlayers {
		playerSeat := *seat.MustSeat(i)
		p := s.Player(playerSeat)
		// TODO: Track the current actor in round.State and render "*" for that player.
		actorMarker := " "
		leftBracket, rightBracket := "[", "]"
		if playerSeat == s.Dealer() {
			leftBracket, rightBracket = "{", "}"
		}
		fmt.Fprintf(&b, "%s%s%d%s tehai: %s %s\n",
			actorMarker, leftBracket, i, rightBracket, formatHand(p.HandTiles(), p.DrawnTile()), formatMelds(p.Melds()))
		fmt.Fprintf(&b, "     ho:    %s\n", formatRiver(p.River(), p.RiichiRiverIndex()))
	}

	b.WriteString(strings.Repeat("-", 80))
	b.WriteByte('\n')
	return b.String()
}

func formatDoraIndicators(tiles []tile.Tile) string {
	parts := make([]string, len(tiles))
	for i, t := range tiles {
		parts[i] = t.String()
	}
	return strings.Join(parts, " ")
}

func formatTiles(tiles []tile.Tile) string {
	var b strings.Builder
	for _, t := range tiles {
		fmt.Fprintf(&b, "%-3s", t)
	}
	return b.String()
}

func formatHand(handTiles []tile.Tile, drawnTile *tile.Tile) string {
	if drawnTile == nil {
		return formatTiles(handTiles)
	}

	tiles := slices.Clone(handTiles)
	tiles = append(tiles, *drawnTile)
	return formatTiles(tiles)
}

func formatMelds(melds []meld.Meld) string {
	if len(melds) == 0 {
		return ""
	}

	parts := make([]string, len(melds))
	for i, m := range melds {
		parts[i] = m.String()
	}
	return strings.Join(parts, " ")
}

func formatRiver(river []tile.Tile, riichiRiverIndex int) string {
	if riichiRiverIndex < 0 {
		return formatTiles(river)
	}

	before := formatTiles(river[:riichiRiverIndex])
	after := formatTiles(river[riichiRiverIndex:])
	return before + "=" + after
}
