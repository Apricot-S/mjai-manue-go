package game

import (
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/message"
)

func (s *StateImpl) DahaiCandidates() []Pai {
	player := &s.players[s.playerID]
	if !player.CanDahai() {
		return nil
	}
	if player.ReachState() == Declared {
		// If the player has already declared reach, return nil.
		return nil
	}
	if player.ReachState() == Accepted {
		// If the player has already accepted the reach, only the drawn tile is a candidate.
		candidates := []Pai{player.tehais[len(player.tehais)-1]}
		return candidates
	}

	candidates := GetUniquePais(player.tehais, nil)

	// Remove the kuikae tiles from the candidates.
	kuikaeSet := make(map[uint8]struct{}, len(s.kuikaePais))
	for _, k := range s.kuikaePais {
		kuikaeSet[k.ID()] = struct{}{}
	}
	candidates = slices.DeleteFunc(candidates, func(p Pai) bool {
		_, found := kuikaeSet[p.ID()]
		return found
	})

	return candidates
}

// ReachDahaiCandidates returns the candidates for the reach declaration tile.
func (s *StateImpl) ReachDahaiCandidates() ([]Pai, error) {
	player := &s.players[s.playerID]
	if !player.CanDahai() {
		return nil, nil
	}
	if !player.IsMenzen() {
		return nil, nil
	}
	if player.ReachState() == Accepted {
		// If the player has already accepted the reach, the player cannot declare reach.
		return nil, nil
	}
	if s.NumPipais() < 4 {
		// If there are no remaining tiles to draw, the player cannot declare reach.
		return nil, nil
	}
	if player.Score() < kyotakuPoint {
		// If the player does not have enough points to declare reach, return nil.
		return nil, nil
	}

	tehaiCounts, err := NewPaiSetWithPais(player.tehais)
	if err != nil {
		return nil, err
	}
	shanten, _, err := AnalyzeShantenWithOption(tehaiCounts, 0, 0)
	if err != nil {
		return nil, err
	}
	if shanten > 0 {
		// If the hand is not tenpai, return nil.
		return nil, nil
	}

	tehaiPais := GetUniquePais(player.tehais, nil)
	candidates := make([]Pai, 0, len(tehaiPais))
	for _, p := range tehaiPais {
		i := p.RemoveRed().ID()
		tehaiCounts[i] -= 1
		shanten, _, err := AnalyzeShantenWithOption(tehaiCounts, 0, 0)
		if err != nil {
			return nil, err
		}
		if shanten <= 0 {
			candidates = append(candidates, p)
		}
		tehaiCounts[i] += 1
	}

	return candidates, nil
}

func (s *StateImpl) IsTsumoPai(pai *Pai) bool {
	if s.lastActor == noActor {
		return false
	}
	if s.lastActor != s.playerID {
		return false
	}
	if s.lastActionType != message.TypeTsumo && s.lastActionType != message.TypeReach {
		return false
	}

	tehais := s.players[s.playerID].tehais
	return pai.ID() == tehais[len(tehais)-1].ID()
}

func (s *StateImpl) FuroCandidates() ([]Furo, error) {
	if s.lastActor == noActor || s.lastActor == s.playerID {
		// Furo is not possible if the last actor is the player itself or no actor.
		return nil, nil
	}
	if s.lastActionType != message.TypeDahai {
		// Furo is only possible after dahai.
		return nil, nil
	}
	if s.NumPipais() == 0 {
		// Furo is not possible if discarded tile is a last tile.
		return nil, nil
	}

	cc, err := s.chiCandidates()
	if err != nil {
		return nil, err
	}
	pc, err := s.ponCandidates()
	if err != nil {
		return nil, err
	}
	dc, err := s.daiminkanCandidates()
	if err != nil {
		return nil, err
	}
	return slices.Concat(cc, pc, dc), nil
}

func (s *StateImpl) chiCandidates() ([]Furo, error) {
	if GetPlayerDistance(&s.players[s.playerID], &s.players[s.lastActor]) != 1 {
		// Chi is only possible for kamicha's discarded tile.
		return nil, nil
	}

	// TODO: Implement logic.
	return nil, nil
}

func (s *StateImpl) ponCandidates() ([]Furo, error) {
	// TODO: Implement logic.
	return nil, nil
}

func (s *StateImpl) daiminkanCandidates() ([]Furo, error) {
	// TODO: Implement logic.
	return nil, nil
}

func (s *StateImpl) HoraCandidate() (*Hora, error) {
	if s.lastActor == noActor {
		return nil, nil
	}
	if s.lastActor == s.playerID && s.lastActionType != message.TypeTsumo {
		// If the last actor is the player itself, it cannot be a ron hora.
		return nil, nil
	}
	if s.lastActor != s.playerID &&
		s.lastActionType != message.TypeDahai && s.lastActionType != message.TypeKakan {
		// If the last actor is not the player itself, it cannot be a tsumo hora.
		return nil, nil
	}

	// TODO: Implement logic.
	return nil, nil
}
