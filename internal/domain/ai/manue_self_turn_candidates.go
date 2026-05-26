package ai

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func getSelfTurnCandidates(actions []action.Action, self player.PlayerViewer) ([]actionCandidate, error) {
	h, err := selfTurnHand(self)
	if err != nil {
		return nil, fmt.Errorf("cannot build self-turn candidates: %w", err)
	}
	turnShanten, turnGoals := service.AnalyzeShanten(h, service.AllowedExtraTiles(1))

	// Self-turn candidates currently cover discard and riichi+discard only.
	// Concealed kan, promoted kan, and kyushukyuhai are intentionally not
	// selected.
	riichi := firstActionOfType[*action.Riichi](actions)
	riichiDeclared := self.RiichiState() == player.RiichiDeclared
	var candidates []actionCandidate
	for _, discard := range normalizedSelfTurnDiscards(actions) {
		afterDiscard, err := h.Discard(discard.Tile())
		if err != nil {
			return nil, fmt.Errorf("cannot build self-turn candidate for %s: %w", discard.Tile(), err)
		}
		shanten := candidateShanten(discard.Tile(), turnShanten, turnGoals)
		if riichi != nil && shanten <= 0 {
			// Match Manue's riichi candidate filtering: only regular-hand shanten
			// from AnalyzeShanten is considered here, so Seven Pairs and
			// Thirteen Orphans do not create riichi candidates.
			candidates = append(candidates, buildSelfTurnCandidate(riichi, discard.Tile(), h, afterDiscard, turnShanten, shanten, turnGoals, true, true))
		}
		if riichiDeclared && shanten > 0 {
			continue
		}
		candidates = append(candidates, buildSelfTurnCandidate(discard, discard.Tile(), h, afterDiscard, turnShanten, shanten, turnGoals, false, riichiDeclared))
	}
	return candidates, nil
}

// normalizedSelfTurnDiscards preserves CoffeeScript output behavior: when the
// same exact tile can be discarded from hand or as tsumogiri, keep tsumogiri.
func normalizedSelfTurnDiscards(actions []action.Action) []*action.Discard {
	discards := make([]*action.Discard, 0, len(actions))
	indexByTile := make(map[tile.Tile]int, len(actions))
	for _, a := range actions {
		discard, ok := a.(*action.Discard)
		if !ok {
			continue
		}
		if i, ok := indexByTile[discard.Tile()]; ok {
			if discard.Tsumogiri() {
				discards[i] = discard
			}
			continue
		}
		indexByTile[discard.Tile()] = len(discards)
		discards = append(discards, discard)
	}
	return discards
}

func buildSelfTurnCandidate(
	immediateAction action.Action,
	discardTile tile.Tile,
	turnHand *hand.VisibleHand,
	afterDiscardHand *hand.VisibleHand,
	baseShanten int,
	shanten int,
	goals []service.Goal,
	riichi bool,
	scoreAsRiichi bool,
) actionCandidate {
	return actionCandidate{
		traceKey:         formatDiscardTraceKey(riichi, discardTile),
		action:           immediateAction,
		riichi:           riichi,
		scoreAsRiichi:    scoreAsRiichi,
		discardTile:      discardTile,
		turnHand:         turnHand,
		afterDiscardHand: afterDiscardHand,
		baseShanten:      baseShanten,
		shantenGoals:     goals,
		score:            scoreDiscardCandidate(discardTile, shanten),
	}
}
