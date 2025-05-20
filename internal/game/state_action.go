package game

import (
	"slices"
	"sort"

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

	candidates := slices.Clone(player.tehais)
	sort.Sort(candidates)
	candidates = slices.CompactFunc(candidates, func(p1, p2 Pai) bool {
		return p1.ID() == p2.ID()
	})

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
		// If the player has already accepted the reach, return nil.
		return nil, nil
	}

	tehaiPais := slices.Clone(player.tehais)
	tehaiCounts, err := NewPaiSetWithPais(tehaiPais)
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

	sort.Sort(tehaiPais)
	tehaiPais = slices.CompactFunc(tehaiPais, func(p1, p2 Pai) bool {
		return p1.ID() == p2.ID()
	})

	// The number of candidates will be equal or less than 13.
	candidates := make([]Pai, 0, 13)
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
	if s.lastActionType != message.TypeTsumo {
		return false
	}

	tehais := s.players[s.playerID].tehais
	return pai.ID() == tehais[len(tehais)-1].ID()
}

func (s *StateImpl) ForbiddenDahais() []Pai {
	return s.kuikaePais
}

func (s *StateImpl) FuroCandidates() ([]Furo, error) {
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
	if s.lastActor == noActor || s.lastActor == s.playerID {
		// Chi is not possible if the last actor is the player itself or no actor.
		return nil, nil
	}
	if s.lastActionType != message.TypeDahai {
		// Chi is only possible after dahai.
		return nil, nil
	}
	if GetPlayerDistance(&s.players[s.playerID], &s.players[s.lastActor]) != 1 {
		// Chi is only possible for kamicha's discarded tile.
		return nil, nil
	}
	if s.NumPipais() == 0 {
		// Chi is not possible if discarded tile is a last tile.
		return nil, nil
	}

	// TODO: Implement logic.
	return nil, nil
}

func (s *StateImpl) ponCandidates() ([]Furo, error) {
	if s.lastActor == noActor || s.lastActor == s.playerID {
		// Pon is not possible if the last actor is the player itself or no actor.
		return nil, nil
	}
	if s.lastActionType != message.TypeDahai {
		// Pon is only possible after dahai.
		return nil, nil
	}
	if s.NumPipais() == 0 {
		// Pon is not possible if discarded tile is a last tile.
		return nil, nil
	}

	// TODO: Implement logic.
	return nil, nil
}

func (s *StateImpl) daiminkanCandidates() ([]Furo, error) {
	if s.lastActor == noActor || s.lastActor == s.playerID {
		// DaiminkanCandidates is not possible if the last actor is the player itself or no actor.
		return nil, nil
	}
	if s.lastActionType != message.TypeDahai {
		// DaiminkanCandidates is only possible after dahai.
		return nil, nil
	}
	if s.NumPipais() == 0 {
		// DaiminkanCandidates is not possible if discarded tile is a last tile.
		return nil, nil
	}

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
