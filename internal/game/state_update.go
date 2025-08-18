package game

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
)

var east, _ = base.NewPaiWithName("E")

func (s *StateImpl) Update(event inbound.Event) error {
	switch e := event.(type) {
	case *inbound.StartGame:
		s.currentEvent = e
		return s.onStartGame(e)
	case *inbound.EndKyoku, *inbound.EndGame:
		// For game records only
		return nil
	}

	// This is specially handled here because it's not an anpai if the dahai is followed by a hora.
	if _, isHora := event.(*inbound.Hora); !isHora {
		switch s.currentEvent.(type) {
		case *inbound.Dahai, *inbound.Kakan:
			for _, p := range s.players {
				if p.ID() != s.prevDahaiActor {
					p.AddExtraAnpais(*s.prevDahaiPai)
				}
			}
		}
	}

	s.currentEvent = event

	switch e := event.(type) {
	case *inbound.StartKyoku:
		return s.onStartKyoku(e)
	case *inbound.Tsumo:
		return s.onTsumo(e)
	case *inbound.Dahai:
		return s.onDahai(e)
	case *inbound.Chi:
		return s.onChi(e)
	case *inbound.Pon:
		return s.onPon(e)
	case *inbound.Daiminkan:
		return s.onDaiminkan(e)
	case *inbound.Ankan:
		return s.onAnkan(e)
	case *inbound.Kakan:
		return s.onKakan(e)
	case *inbound.Dora:
		return s.onDora(e)
	case *inbound.Reach:
		return s.onReach(e)
	case *inbound.ReachAccepted:
		return s.onReachAccepted(e)
	case *inbound.Hora:
		return s.onHora(e)
	case *inbound.Ryukyoku:
		return s.onRyukyoku(e)
	default:
		return fmt.Errorf("unknown event type: %v", event)
	}
}

func (s *StateImpl) onStartGame(event *inbound.StartGame) error {
	if event == nil {
		return fmt.Errorf("start_game event is nil")
	}

	var players [NumPlayers]base.Player
	for i, name := range event.Names {
		p, err := base.NewPlayer(i, name, InitScore)
		if err != nil {
			return err
		}
		players[i] = *p
	}

	s.players = players
	s.bakaze = *east
	s.kyokuNum = 1
	s.honba = 0
	s.oya = &s.players[0]
	s.chicha = &s.players[0]
	s.doraMarkers = make([]base.Pai, 0, MaxNumDoraMarkers)
	s.numPipais = NumInitPipais

	s.prevDahaiActor = noActor
	s.prevDahaiPai = nil
	s.currentEvent = nil

	s.playerID = event.ID
	s.lastActor = noActor
	s.lastAction = nil

	s.kuikaePais = make([]base.Pai, 0, 3)
	s.missedRon = false
	s.isFuriten = false
	s.isRinshanTsumo = false

	return nil
}

func (s *StateImpl) onStartKyoku(event *inbound.StartKyoku) error {
	if event == nil {
		return fmt.Errorf("start_kyoku event is nil")
	}

	s.bakaze = event.Bakaze
	s.kyokuNum = event.Kyoku
	s.honba = event.Honba
	s.oya = &s.players[event.Oya]
	s.doraMarkers = make([]base.Pai, 0, MaxNumDoraMarkers)
	s.doraMarkers = append(s.doraMarkers, event.DoraMarker)
	s.numPipais = NumInitPipais

	for i := range NumPlayers {
		var err error
		if event.Scores != nil {
			err = s.players[i].OnStartKyoku(event.Tehais[i], &event.Scores[i])
		} else {
			err = s.players[i].OnStartKyoku(event.Tehais[i], nil)
		}
		if err != nil {
			return err
		}
	}

	s.prevDahaiActor = noActor
	s.prevDahaiPai = nil

	s.lastActor = noActor
	s.lastAction = nil

	s.kuikaePais = make([]base.Pai, 0, 3)
	s.missedRon = false
	s.isFuriten = false
	s.isRinshanTsumo = false

	return nil
}

func (s *StateImpl) onTsumo(event *inbound.Tsumo) error {
	if event == nil {
		return fmt.Errorf("tsumo event is nil")
	}

	if s.numPipais <= 0 {
		return fmt.Errorf("tsumo is not possible if numPipais is 0 or negative: %d", s.numPipais)
	}
	s.numPipais--

	actor := event.Actor
	player := &s.players[actor]
	if err := player.OnTsumo(event.Pai); err != nil {
		return err
	}

	s.lastActor = actor
	s.lastAction = event

	return nil
}

func (s *StateImpl) onDahai(event *inbound.Dahai) error {
	if event == nil {
		return fmt.Errorf("dahai event is nil")
	}

	pai := event.Pai
	actor := event.Actor
	player := &s.players[actor]
	if err := player.OnDahai(pai); err != nil {
		return err
	}

	s.prevDahaiActor = actor
	s.prevDahaiPai = &pai

	s.lastActor = actor
	s.lastAction = event

	tehaiCounts, err := base.NewPaiSet(s.players[s.playerID].Tehais())
	if err != nil {
		return err
	}

	if actor == s.playerID {
		s.kuikaePais = make([]base.Pai, 0, 3)
		if player.ReachState() != base.ReachAccepted {
			s.missedRon = false
			s.isFuriten = false
		}
		s.isRinshanTsumo = false

		for _, anpai := range s.Anpais(player) {
			if err := tehaiCounts.AddPai(&anpai, 1); err != nil {
				return fmt.Errorf("failed to add anpai %v to tehaiCounts: %w", anpai, err)
			}

			isHoraFrom, err := IsHoraForm(tehaiCounts)
			if err != nil {
				return fmt.Errorf("failed to check if tehaiCounts is hora form: %w", err)
			}

			if err := tehaiCounts.AddPai(&anpai, -1); err != nil {
				return fmt.Errorf("failed to remove anpai %v to tehaiCounts: %w", anpai, err)
			}

			if isHoraFrom {
				s.isFuriten = true
				break
			}
		}
	} else {
		if s.missedRon {
			// The previous ron-able tile was missed
			s.isFuriten = true
		}

		if err := tehaiCounts.AddPai(&pai, 1); err != nil {
			return fmt.Errorf("failed to add pai %v to tehaiCounts: %w", pai, err)
		}

		isHoraFrom, err := IsHoraForm(tehaiCounts)
		if err != nil {
			return fmt.Errorf("failed to check if tehaiCounts is hora form: %w", err)
		}
		if isHoraFrom && !s.missedRon && !s.isFuriten {
			s.missedRon = true
		}
	}

	return nil
}

func (s *StateImpl) onChi(event *inbound.Chi) error {
	if event == nil {
		return fmt.Errorf("chi event is nil")
	}

	if s.numPipais <= 0 {
		return fmt.Errorf("chi is not possible if numPipais is 0 or negative: %d", s.numPipais)
	}

	pai := event.Taken
	consumed := event.Consumed
	furo, err := base.NewChi(pai, consumed, event.Target)
	if err != nil {
		return err
	}

	actor := event.Actor
	if err := s.players[actor].OnChi(furo); err != nil {
		return err
	}

	target := event.Target
	if err := s.players[target].OnTargeted(furo); err != nil {
		return err
	}

	s.lastActor = actor
	s.lastAction = event

	if actor == s.playerID {
		s.kuikaePais = append(s.kuikaePais, *pai.RemoveRed())
		if pai.Number() == 5 {
			s.kuikaePais = append(s.kuikaePais, *pai.AddRed())
		}
		// Add suji kuikae for ryanmen chi (two-sided chi)
		n0 := int8(consumed[0].Number())
		n1 := int8(consumed[1].Number())
		diff := n0 - n1
		if diff == 1 || diff == -1 {
			// If consumed[0] and consumed[1] are consecutive, it's a ryanmen chi
			// Number of the taken tile
			nTaken := int8(pai.Number())
			// Find the smaller and larger of the consumed tiles
			nLow := n0
			nHigh := n1
			if nLow > nHigh {
				nLow, nHigh = nHigh, nLow
			}
			// For ryanmen chi, nTaken is nLow-1 or nHigh+1
			// Add the other end of the suji as a kuikae tile
			if nTaken < 7 && nTaken == nLow-1 {
				// Example: chi 4 with 5,6 -> also add 7 as kuikae
				sujiPai := pai.Next(3)
				s.kuikaePais = append(s.kuikaePais, *sujiPai)
				if sujiPai.Number() == 5 {
					s.kuikaePais = append(s.kuikaePais, *sujiPai.AddRed())
				}
			}
			if nTaken > 3 && nTaken == nHigh+1 {
				// Example: chi 5 with 3,4 -> also add 2 as kuikae
				sujiPai := pai.Next(-3)
				s.kuikaePais = append(s.kuikaePais, *sujiPai)
				if sujiPai.Number() == 5 {
					s.kuikaePais = append(s.kuikaePais, *sujiPai.AddRed())
				}
			}
		}
	}

	return nil
}

func (s *StateImpl) onPon(event *inbound.Pon) error {
	if event == nil {
		return fmt.Errorf("pon event is nil")
	}

	if s.numPipais <= 0 {
		return fmt.Errorf("pon is not possible if numPipais is 0 or negative: %d", s.numPipais)
	}

	pai := event.Taken
	furo, err := base.NewPon(pai, event.Consumed, event.Target)
	if err != nil {
		return err
	}

	actor := event.Actor
	if err := s.players[actor].OnPon(furo); err != nil {
		return err
	}

	target := event.Target
	if err := s.players[target].OnTargeted(furo); err != nil {
		return err
	}

	s.lastActor = actor
	s.lastAction = event

	if actor == s.playerID {
		s.kuikaePais = append(s.kuikaePais, *pai.RemoveRed())
		if !pai.IsTsupai() && pai.Number() == 5 {
			s.kuikaePais = append(s.kuikaePais, *pai.AddRed())
		}
	}

	return nil
}

func (s *StateImpl) onDaiminkan(event *inbound.Daiminkan) error {
	if event == nil {
		return fmt.Errorf("daiminkan event is nil")
	}

	if s.numPipais <= 0 {
		return fmt.Errorf("daiminkan is not possible if numPipais is 0 or negative: %d", s.numPipais)
	}

	pai := event.Taken
	furo, err := base.NewDaiminkan(pai, event.Consumed, event.Target)
	if err != nil {
		return err
	}

	actor := event.Actor
	if err := s.players[actor].OnDaiminkan(furo); err != nil {
		return err
	}

	target := event.Target
	if err := s.players[target].OnTargeted(furo); err != nil {
		return err
	}

	s.lastActor = actor
	s.lastAction = event

	if actor == s.playerID {
		s.isRinshanTsumo = true
	}

	return nil
}

func (s *StateImpl) onAnkan(event *inbound.Ankan) error {
	if event == nil {
		return fmt.Errorf("ankan event is nil")
	}

	if s.numPipais <= 0 {
		return fmt.Errorf("ankan is not possible if numPipais is 0 or negative: %d", s.numPipais)
	}

	furo, err := base.NewAnkan(event.Consumed)
	if err != nil {
		return err
	}

	actor := event.Actor
	if err := s.players[actor].OnAnkan(furo); err != nil {
		return err
	}

	s.lastActor = actor
	s.lastAction = event

	if actor == s.playerID {
		s.isRinshanTsumo = true
	}

	return nil
}

func (s *StateImpl) onKakan(event *inbound.Kakan) error {
	if event == nil {
		return fmt.Errorf("kakan event is nil")
	}

	if s.numPipais <= 0 {
		return fmt.Errorf("kakan is not possible if numPipais is 0 or negative: %d", s.numPipais)
	}

	pai := event.Added
	furo, err := base.NewKakan(event.Taken, event.Consumed, event.Added, event.Target)
	if err != nil {
		return err
	}

	actor := event.Actor
	if err := s.players[actor].OnKakan(furo); err != nil {
		return err
	}

	// For chankan
	s.prevDahaiActor = actor
	s.prevDahaiPai = &pai

	s.lastActor = actor
	s.lastAction = event

	if actor == s.playerID {
		s.isRinshanTsumo = true
	} else {
		tehaiCounts, err := base.NewPaiSet(s.players[s.playerID].Tehais())
		if err != nil {
			return err
		}

		if s.missedRon {
			// The previous ron-able tile was missed
			s.isFuriten = true
		}

		if err := tehaiCounts.AddPai(&pai, 1); err != nil {
			return fmt.Errorf("failed to add pai %v to tehaiCounts: %w", pai, err)
		}

		isHoraFrom, err := IsHoraForm(tehaiCounts)
		if err != nil {
			return fmt.Errorf("failed to check if tehaiCounts is hora form: %w", err)
		}
		if isHoraFrom && !s.missedRon && !s.isFuriten {
			s.missedRon = true
		}
	}

	return nil
}

func (s *StateImpl) onDora(event *inbound.Dora) error {
	if event == nil {
		return fmt.Errorf("dora event is nil")
	}

	if len(s.doraMarkers) >= MaxNumDoraMarkers {
		return fmt.Errorf("a 6th dora cannot be added")
	}

	s.doraMarkers = append(s.doraMarkers, event.DoraMarker)

	return nil
}

func (s *StateImpl) onReach(event *inbound.Reach) error {
	if event == nil {
		return fmt.Errorf("reach event is nil")
	}

	if s.numPipais <= 0 {
		return fmt.Errorf("reach is not possible if numPipais is 0 or negative: %d", s.numPipais)
	}

	actor := event.Actor
	player := &s.players[actor]
	if err := player.OnReach(); err != nil {
		return err
	}

	s.lastActor = actor
	s.lastAction = event

	return nil
}

func (s *StateImpl) onReachAccepted(event *inbound.ReachAccepted) error {
	if event == nil {
		return fmt.Errorf("reach_accepted event is nil")
	}

	actor := event.Actor
	player := &s.players[actor]
	var err error
	if event.Scores != nil {
		err = player.OnReachAccepted(&event.Scores[actor])
	} else {
		err = player.OnReachAccepted(nil)
	}
	if err != nil {
		return err
	}

	return nil
}

func (s *StateImpl) onHora(event *inbound.Hora) error {
	if event == nil {
		return fmt.Errorf("hora event is nil")
	}

	if event.Scores != nil {
		for i, score := range event.Scores {
			s.players[i].SetScore(score)
		}
	}

	// After hora, only end_kyoku comes, so reset the last action.
	s.lastActor = noActor
	s.lastAction = nil

	return nil
}

func (s *StateImpl) onRyukyoku(event *inbound.Ryukyoku) error {
	if event == nil {
		return fmt.Errorf("ryukyoku event is nil")
	}

	if event.Scores != nil {
		for i, score := range event.Scores {
			s.players[i].SetScore(score)
		}
	}

	// After ryukyoku, only end_kyoku comes, so reset the last action.
	s.lastActor = noActor
	s.lastAction = nil

	return nil
}
