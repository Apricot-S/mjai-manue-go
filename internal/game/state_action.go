package game

import (
	"fmt"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
)

func isEndOfRound(currentEvent inbound.Event) bool {
	_, isHora := currentEvent.(*inbound.Hora)
	_, isRyukyoku := currentEvent.(*inbound.Ryukyoku)
	return isHora || isRyukyoku
}

func (s *StateImpl) DahaiCandidates() []base.Pai {
	if isEndOfRound(s.currentEvent) {
		return nil
	}

	player := &s.players[s.playerID]
	if !player.CanDahai() {
		return nil
	}
	if player.ReachState() == base.ReachDeclared {
		// If the player has already declared reach, return nil.
		return nil
	}
	if player.ReachState() == base.ReachAccepted {
		// If the player has already accepted the reach, only the drawn tile is a candidate.
		tehais := player.Tehais()
		candidates := []base.Pai{tehais[len(tehais)-1]}
		return candidates
	}

	candidates := base.GetUniquePais(player.Tehais(), nil)

	// Remove the kuikae tiles from the candidates.
	kuikaeSet := make(map[uint8]struct{}, len(s.kuikaePais))
	for _, k := range s.kuikaePais {
		kuikaeSet[k.ID()] = struct{}{}
	}
	candidates = slices.DeleteFunc(candidates, func(p base.Pai) bool {
		_, found := kuikaeSet[p.ID()]
		return found
	})

	return candidates
}

// ReachDahaiCandidates returns the candidates for the reach declaration tile.
func (s *StateImpl) ReachDahaiCandidates() ([]base.Pai, error) {
	if isEndOfRound(s.currentEvent) {
		return nil, nil
	}

	player := &s.players[s.playerID]
	if !player.CanDahai() {
		return nil, nil
	}
	if !player.IsMenzen() {
		return nil, nil
	}
	if player.ReachState() == base.ReachAccepted {
		// If the player has already accepted the reach, the player cannot declare reach.
		return nil, nil
	}
	if s.NumPipais() < 4 {
		// If there are no remaining tiles to draw, the player cannot declare reach.
		return nil, nil
	}
	if player.Score() < base.KyotakuPoint {
		// If the player does not have enough points to declare reach, return nil.
		return nil, nil
	}

	tehaiCounts, err := base.NewPaiSet(player.Tehais())
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

	tehaiPais := base.GetUniquePais(player.Tehais(), nil)
	candidates := make([]base.Pai, 0, len(tehaiPais))
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

func (s *StateImpl) IsTsumoPai(pai *base.Pai) bool {
	if s.lastActor == noActor {
		return false
	}
	if s.lastActor != s.playerID {
		return false
	}
	_, isTsumo := s.lastAction.(*inbound.Tsumo)
	_, isReach := s.lastAction.(*inbound.Reach)
	if !isTsumo && !isReach {
		return false
	}

	tehais := s.players[s.playerID].Tehais()
	return pai.ID() == tehais[len(tehais)-1].ID()
}

func (s *StateImpl) FuroCandidates() ([]base.Furo, error) {
	if isEndOfRound(s.currentEvent) {
		return nil, nil
	}

	if s.lastActor == noActor || s.lastActor == s.playerID {
		// Furo is not possible if the last actor is the player itself or no actor.
		return nil, nil
	}
	if _, ok := s.lastAction.(*inbound.Dahai); !ok {
		// Furo is only possible after dahai.
		return nil, nil
	}
	if _, ok := s.currentEvent.(*inbound.ReachAccepted); ok {
		// If a call is accepted to the riichi declaration tile,
		// the client sends a call message ("chi", "pon" or "kan") for the "dahai" first.
		// Afterwards, the server responds with "reach_accepted".
		return nil, nil
	}
	if s.NumPipais() == 0 {
		// Furo is not possible if discarded tile is a last tile.
		return nil, nil
	}

	player := &s.players[s.playerID]
	if player.ReachState() != base.NotReach {
		// If the player has already declared the reach, the player cannot furo.
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

func (s *StateImpl) chiCandidates() ([]base.Furo, error) {
	if GetPlayerDistance(&s.players[s.playerID], &s.players[s.lastActor]) != 1 {
		// Chi is only possible for kamicha's discarded tile.
		return nil, nil
	}

	taken := s.prevDahaiPai
	if taken == nil {
		return nil, nil
	}
	if taken.IsTsupai() {
		// Chi is not possible for honors.
		return nil, nil
	}

	player := &s.players[s.playerID]
	tehais := player.Tehais()

	furos := make([]base.Furo, 0, 5)
	used := make(map[[2]uint8]struct{})
	takenNum := int(taken.Number())
	for d := -2; d <= 0; d++ {
		// taken can be the lowest, middle, or highest tile in a sequence.
		n1 := takenNum + d
		n2 := takenNum + d + 1
		n3 := takenNum + d + 2
		if n1 < 1 || n3 > 9 {
			continue
		}

		// Find two tiles that can form a sequence with the taken tile.
		var cands []base.Pai
		for _, p := range tehais {
			if p.Type() != taken.Type() {
				continue
			}
			if p.HasSameSymbol(taken) {
				continue
			}
			if n := int(p.Number()); n == n1 || n == n2 || n == n3 {
				cands = append(cands, p)
			}
		}

		for i := range cands {
			for j := i + 1; j < len(cands); j++ {
				// Two candidates other than the taken tile
				p1 := cands[i]
				p2 := cands[j]
				// Check if taken, p1, and p2 form a sequence.
				nums := []int{takenNum, int(p1.Number()), int(p2.Number())}
				slices.Sort(nums)
				if nums[0] != n1 || nums[1] != n2 || nums[2] != n3 {
					continue
				}

				// Preventing duplicate
				id1 := p1.ID()
				id2 := p2.ID()
				if id1 > id2 {
					id1, id2 = id2, id1
				}
				key := [2]uint8{id1, id2}
				if _, ok := used[key]; ok {
					// If the pair of tiles has already been used, skip it.
					continue
				}
				used[key] = struct{}{}

				// Create a chi
				consumed := [2]base.Pai{p1, p2}
				furo, err := base.NewChi(*taken, consumed, s.prevDahaiActor)
				if err != nil {
					return nil, err
				}

				// Kuikae check
				// Create a hand without the two tiles consumed by chi.
				rest := make([]base.Pai, 0, len(tehais)-2)
				usedIdx := map[int]bool{}
				p1used := false
				p2used := false
				for idx, tp := range tehais {
					if !p1used && tp.ID() == p1.ID() {
						usedIdx[idx] = true
						p1used = true
						continue
					}
					if !p2used && tp.ID() == p2.ID() && !usedIdx[idx] {
						usedIdx[idx] = true
						p2used = true
						continue
					}
				}
				for idx, tp := range tehais {
					if !usedIdx[idx] {
						rest = append(rest, tp)
					}
				}
				// If all remaining tiles are kuikae candidates, skip.
				allKuikae := true
				for _, tp := range rest {
					if !base.IsKuikae(furo, &tp) {
						allKuikae = false
						break
					}
				}
				if allKuikae {
					continue
				}

				furos = append(furos, furo)
			}
		}
	}
	return furos, nil
}

func (s *StateImpl) ponCandidates() ([]base.Furo, error) {
	taken := s.prevDahaiPai
	if taken == nil {
		return nil, nil
	}

	consumedPais := make([]base.Pai, 0, 3)
	player := &s.players[s.playerID]
	for _, p := range player.Tehais() {
		if p.HasSameSymbol(taken) {
			consumedPais = append(consumedPais, p)
		}
	}
	numConsumedPais := len(consumedPais)
	if numConsumedPais < 2 {
		return nil, nil
	}

	furos := make([]base.Furo, 0, 2)
	used := make(map[[2]uint8]struct{})
	for i := range numConsumedPais {
		for j := i + 1; j < numConsumedPais; j++ {
			// Preventing duplicate
			id1 := consumedPais[i].ID()
			id2 := consumedPais[j].ID()
			if id1 > id2 {
				id1, id2 = id2, id1
			}
			key := [2]uint8{id1, id2}
			if _, ok := used[key]; ok {
				// If the pair of tiles has already been used, skip it.
				continue
			}
			used[key] = struct{}{}

			consumed := [2]base.Pai{consumedPais[i], consumedPais[j]}
			furo, err := base.NewPon(*taken, consumed, s.prevDahaiActor)
			if err != nil {
				return nil, err
			}
			furos = append(furos, furo)
		}
	}
	return furos, nil
}

func (s *StateImpl) daiminkanCandidates() ([]base.Furo, error) {
	taken := s.prevDahaiPai
	if taken == nil {
		return nil, nil
	}

	consumedPais := make([]base.Pai, 0, 3)
	player := &s.players[s.playerID]
	for _, p := range player.Tehais() {
		if p.HasSameSymbol(taken) {
			consumedPais = append(consumedPais, p)
		}
	}
	numConsumedPais := len(consumedPais)
	if numConsumedPais != 3 {
		return nil, nil
	}

	consumed := [3]base.Pai(consumedPais)
	furo, err := base.NewDaiminkan(*taken, consumed, s.prevDahaiActor)
	if err != nil {
		return nil, err
	}
	furos := []base.Furo{furo}
	return furos, nil
}

func (s *StateImpl) HoraCandidate() (*base.Hora, error) {
	if isEndOfRound(s.currentEvent) {
		return nil, nil
	}

	if s.lastActor == noActor {
		return nil, nil
	}

	_, isTsumo := s.lastAction.(*inbound.Tsumo)
	_, isDahai := s.lastAction.(*inbound.Dahai)
	_, isKakan := s.lastAction.(*inbound.Kakan)
	isTsumoSituation := s.lastActor == s.playerID && isTsumo
	isRonSituation := s.lastActor != s.playerID && (isDahai || isKakan)
	if !isTsumoSituation && !isRonSituation {
		// The last action is not a valid situation for hora.
		return nil, nil
	}
	if s.isFuriten && !isTsumoSituation {
		// If the player is in furiten, the player cannot ron.
		return nil, nil
	}

	player := &s.players[s.playerID]
	tehais := player.Tehais()

	tehaiCounts, err := base.NewPaiSet(tehais)
	if err != nil {
		return nil, err
	}
	if isRonSituation {
		dahaiPai := s.prevDahaiPai
		if dahaiPai == nil {
			return nil, fmt.Errorf("dahaiPai is nil, but canRon is true")
		}
		if err := tehaiCounts.AddPai(dahaiPai, 1); err != nil {
			return nil, fmt.Errorf("failed to add dahaiPai %v to tehaiCounts: %w", dahaiPai, err)
		}
	}
	isHoraFrom, err := IsHoraForm(tehaiCounts)
	if err != nil {
		return nil, fmt.Errorf("failed to check if tehaiCounts is hora form: %w", err)
	}
	if !isHoraFrom {
		return nil, nil
	}

	var horaPai base.Pai
	if isTsumoSituation {
		horaPai = tehais[len(tehais)-1]
	} else {
		horaPai = *s.prevDahaiPai
	}

	// Situation Yaku
	hasMenzenchinTsumoho := isTsumoSituation && player.IsMenzen()
	hasReach := player.ReachState() == base.ReachAccepted
	hasChankan := isKakan
	hasRinshankaiho := s.isRinshanTsumo
	hasHaiteimoyueOrHoteiraoyui := s.NumPipais() == 0

	has1Fan := hasMenzenchinTsumoho || hasReach || hasChankan || hasRinshankaiho || hasHaiteimoyueOrHoteiraoyui
	if !has1Fan {
		has1Fan, err = Has1Fan(s, s.playerID, tehais, player.Furos(), &horaPai, isTsumoSituation)
		if err != nil {
			return nil, fmt.Errorf("failed to check if has 1 fan: %w", err)
		}
	}

	if has1Fan {
		hora, err := base.NewHora(horaPai, s.lastActor)
		if err != nil {
			return nil, fmt.Errorf("failed to create hora: %w", err)
		}
		return hora, nil
	}

	return nil, nil
}
