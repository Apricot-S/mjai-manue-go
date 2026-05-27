package ai

import (
	"fmt"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func getOtherDiscardReactionCandidates(actions []action.Action, self player.PlayerViewer) ([]actionCandidate, error) {
	h, err := selfTurnHand(self)
	if err != nil {
		return nil, fmt.Errorf("cannot build reaction candidates: %w", err)
	}

	candidates := make([]actionCandidate, 0, len(actions))
	if pass := firstActionOfType[*action.Pass](actions); pass != nil {
		shanten, goals := service.AnalyzeShanten(h, service.AllowedExtraTiles(1))
		unknown := tile.MustTileFromCode("?")
		candidates = append(candidates, actionCandidate{
			traceKey:         "none",
			action:           pass,
			discardTile:      unknown,
			melds:            self.Melds(),
			turnHand:         h,
			afterDiscardHand: h,
			baseShanten:      shanten,
			shantenGoals:     goals,
			score:            scoreDiscardCandidate(unknown, shanten),
		})
	}

	callIndex := 0
	for _, a := range actions {
		callMeld, ok, err := actionCallMeld(a)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		callCandidates, err := getCallReactionCandidates(callIndex, a, callMeld, h, self.Melds())
		if err != nil {
			return nil, err
		}
		candidates = append(candidates, callCandidates...)
		callIndex++
	}
	return candidates, nil
}

func getCallReactionCandidates(
	callIndex int,
	callAction action.Action,
	callMeld meld.Meld,
	baseHand *hand.VisibleHand,
	baseMelds []meld.Meld,
) ([]actionCandidate, error) {
	turnHand, err := baseHand.Call(callMeld)
	if err != nil {
		return nil, fmt.Errorf("cannot build reaction candidates for call %d: %w", callIndex, err)
	}
	nextMelds := append(slices.Clone(baseMelds), callMeld)
	turnShanten, turnGoals := service.AnalyzeShanten(turnHand, service.AllowedExtraTiles(1))

	if _, ok := callMeld.(*meld.CalledKan); ok {
		unknown := tile.MustTileFromCode("?")
		shanten, goals := service.AnalyzeShanten(turnHand, service.AllowedExtraTiles(1))
		return []actionCandidate{{
			traceKey:         fmt.Sprintf("%d.none", callIndex),
			action:           callAction,
			discardTile:      unknown,
			melds:            nextMelds,
			turnHand:         turnHand,
			afterDiscardHand: turnHand,
			baseShanten:      shanten,
			shantenGoals:     goals,
			score:            scoreDiscardCandidate(unknown, shanten),
		}}, nil
	}

	swapCallTiles := callSwapTiles(callMeld)
	discardTiles := tile.Tiles(turnHand.ToTiles()).Distinct(func(t tile.Tile) bool {
		return isSwapCallTile(t, swapCallTiles)
	})
	candidates := make([]actionCandidate, 0, len(discardTiles))
	for _, discardTile := range discardTiles {
		afterDiscard, err := turnHand.Discard(discardTile)
		if err != nil {
			return nil, fmt.Errorf("cannot build reaction candidate %d.%s: %w", callIndex, discardTile, err)
		}
		shanten := candidateShanten(discardTile, turnShanten, turnGoals)
		candidates = append(candidates, actionCandidate{
			traceKey:         fmt.Sprintf("%d.%s", callIndex, discardTile),
			action:           callAction,
			discardTile:      discardTile,
			melds:            nextMelds,
			turnHand:         turnHand,
			afterDiscardHand: afterDiscard,
			baseShanten:      turnShanten,
			shantenGoals:     turnGoals,
			score:            scoreDiscardCandidate(discardTile, shanten),
		})
	}
	return candidates, nil
}

func actionCallMeld(a action.Action) (meld.Meld, bool, error) {
	switch call := a.(type) {
	case *action.Chii:
		m, err := meld.NewChii(call.Taken(), call.Consumed(), call.Target())
		if err != nil {
			return nil, false, err
		}
		return m, true, nil
	case *action.Pon:
		m, err := meld.NewPon(call.Taken(), call.Consumed(), call.Target())
		if err != nil {
			return nil, false, err
		}
		return m, true, nil
	case *action.CalledKan:
		m, err := meld.NewCalledKan(call.Taken(), call.Consumed(), call.Target())
		if err != nil {
			return nil, false, err
		}
		return m, true, nil
	default:
		return nil, false, nil
	}
}

func callSwapTiles(m meld.Meld) []tile.Tile {
	switch c := m.(type) {
	case *meld.Chii:
		return c.SwapCallTiles()
	case *meld.Pon:
		return c.SwapCallTiles()
	default:
		return nil
	}
}

func isSwapCallTile(t tile.Tile, swapCallTiles []tile.Tile) bool {
	return slices.ContainsFunc(swapCallTiles, func(s tile.Tile) bool {
		return t.HasSameSymbol(s)
	})
}
