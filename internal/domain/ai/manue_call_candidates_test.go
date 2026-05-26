package ai

import (
	"strings"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestGetOtherDiscardReactionCandidates_BuildsPassAndCallDiscards(t *testing.T) {
	self := seat.MustSeat(0)
	target := seat.MustSeat(3)
	pon, err := action.NewPon(self, target, tile.MustTileFromCode("5p"), [2]tile.Tile{
		tile.MustTileFromCode("5p"),
		tile.MustTileFromCode("5p"),
	})
	if err != nil {
		t.Fatalf("NewPon() failed: %v", err)
	}
	pass := action.NewPass(self)

	got, err := getOtherDiscardReactionCandidates([]action.Action{pass, pon}, stubPlayerViewer{
		hand: hand.CodesToHand([]string{
			"1m", "2m", "3m", "4m", "5m", "6m", "7m",
			"1p", "2p", "5p", "5p", "5pr", "E",
		}),
		riichiState: player.NotRiichi,
	})
	if err != nil {
		t.Fatalf("getOtherDiscardReactionCandidates() failed: %v", err)
	}

	found := map[string]bool{"none": false, "0.4m": false}
	for _, candidate := range got {
		if _, ok := found[candidate.traceKey]; ok {
			found[candidate.traceKey] = true
		}
		if candidate.traceKey == "0.5pr" {
			t.Errorf("getOtherDiscardReactionCandidates() included kuikae discard %q", candidate.traceKey)
		}
		if strings.HasPrefix(candidate.traceKey, "0.") && candidate.action != pon {
			t.Errorf("call candidate action = %v, want pon", candidate.action)
		}
	}
	for traceKey, ok := range found {
		if !ok {
			t.Errorf("getOtherDiscardReactionCandidates() does not contain %q", traceKey)
		}
	}
}
