package game

import (
	"fmt"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/message"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

func (s *StateImpl) OnStartGame(event jsontext.Value) error {
	if event == nil {
		return fmt.Errorf("start_game message is nil")
	}

	var e message.StartGame
	if err := json.Unmarshal(event, &e); err != nil {
		return fmt.Errorf("failed to unmarshal start_game: %w", err)
	}

	id := e.ID
	if id < 0 || id >= NumPlayers {
		return fmt.Errorf("invalid player ID: %d", id)
	}

	names := []string{"", "", "", ""}
	if e.Names != nil {
		names = slices.Clone(e.Names)
	}
	if len(names) != NumPlayers {
		return fmt.Errorf("number of players must be 4, but got %d", len(names))
	}

	east, err := NewPaiWithName("E")
	if err != nil {
		return err
	}

	var players [NumPlayers]Player
	for i, name := range names {
		p, err := NewPlayer(i, name, InitScore)
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
	s.doraMarkers = make([]Pai, 0, MaxNumDoraMarkers)
	s.numPipais = NumInitPipais

	s.prevEventType = noEvent
	s.prevDahaiActor = noActor
	s.prevDahaiPai = nil
	s.currentEventType = noEvent

	s.playerID = id
	s.lastActor = noActor
	s.lastActionType = noEvent

	s.kuikaePais = make([]Pai, 0, 3)
	s.missedRon = false
	s.isFuriten = false
	s.isRinshanTsumo = false

	return nil
}

func (s *StateImpl) Update(event jsontext.Value) error {
	s.prevEventType = s.currentEventType

	var msg message.Message
	if err := json.Unmarshal(event, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	// This is specially handled here because it's not an anpai if the dahai is followed by a hora.
	if msg.Type != message.TypeHora &&
		(s.prevEventType == message.TypeDahai || s.prevEventType == message.TypeKakan) {
		for _, p := range s.players {
			if p.ID() != s.prevDahaiActor {
				p.AddExtraAnpais(*s.prevDahaiPai)
			}
		}
	}

	switch msg.Type {
	case message.TypeStartKyoku:
		var e message.StartKyoku
		if err := json.Unmarshal(event, &e); err != nil {
			return fmt.Errorf("failed to unmarshal start_kyoku: %w", err)
		}
		s.currentEventType = message.TypeStartKyoku
		s.onStartKyoku(&e)
	case message.TypeTsumo:
		var e message.Tsumo
		if err := json.Unmarshal(event, &e); err != nil {
			return fmt.Errorf("failed to unmarshal tsumo: %w", err)
		}
		s.currentEventType = message.TypeTsumo
		s.onTsumo(&e)
	case message.TypeDahai:
		var e message.Dahai
		if err := json.Unmarshal(event, &e); err != nil {
			return fmt.Errorf("failed to unmarshal dahai: %w", err)
		}
		s.currentEventType = message.TypeDahai
		s.onDahai(&e)
	case message.TypeChi:
		var e message.Chi
		if err := json.Unmarshal(event, &e); err != nil {
			return fmt.Errorf("failed to unmarshal chi: %w", err)
		}
		s.currentEventType = message.TypeChi
		s.onChi(&e)
	case message.TypePon:
		var e message.Pon
		if err := json.Unmarshal(event, &e); err != nil {
			return fmt.Errorf("failed to unmarshal pon: %w", err)
		}
		s.currentEventType = message.TypePon
		s.onPon(&e)
	case message.TypeDaiminkan:
		var e message.Daiminkan
		if err := json.Unmarshal(event, &e); err != nil {
			return fmt.Errorf("failed to unmarshal daiminkan: %w", err)
		}
		s.currentEventType = message.TypeDaiminkan
		s.onDaiminkan(&e)
	case message.TypeAnkan:
		var e message.Ankan
		if err := json.Unmarshal(event, &e); err != nil {
			return fmt.Errorf("failed to unmarshal ankan: %w", err)
		}
		s.currentEventType = message.TypeAnkan
		s.onAnkan(&e)
	case message.TypeKakan:
		var e message.Kakan
		if err := json.Unmarshal(event, &e); err != nil {
			return fmt.Errorf("failed to unmarshal kakan: %w", err)
		}
		s.currentEventType = message.TypeKakan
		s.onKakan(&e)
	case message.TypeDora:
		var e message.Dora
		if err := json.Unmarshal(event, &e); err != nil {
			return fmt.Errorf("failed to unmarshal dora: %w", err)
		}
		s.currentEventType = message.TypeDora
		s.onDora(&e)
	case message.TypeReach:
		var e message.Reach
		if err := json.Unmarshal(event, &e); err != nil {
			return fmt.Errorf("failed to unmarshal reach: %w", err)
		}
		s.currentEventType = message.TypeReach
		s.onReach(&e)
	case message.TypeReachAccepted:
		var e message.ReachAccepted
		if err := json.Unmarshal(event, &e); err != nil {
			return fmt.Errorf("failed to unmarshal reach_accepted: %w", err)
		}
		s.currentEventType = message.TypeReachAccepted
		s.onReachAccepted(&e)
	case message.TypeHora:
		var e message.Hora
		if err := json.Unmarshal(event, &e); err != nil {
			return fmt.Errorf("failed to unmarshal hora: %w", err)
		}
		s.currentEventType = message.TypeHora
		s.onHora(&e)
	case message.TypeRyukyoku:
		var e message.Ryukyoku
		if err := json.Unmarshal(event, &e); err != nil {
			return fmt.Errorf("failed to unmarshal ryukyoku: %w", err)
		}
		s.currentEventType = message.TypeRyukyoku
		s.onRyukyoku(&e)
	default:
		return fmt.Errorf("unknown event type: %v", event)
	}

	return nil
}

func (s *StateImpl) onStartKyoku(event *message.StartKyoku) error {
	if event == nil {
		return fmt.Errorf("start_kyoku message is nil")
	}

	bakaze, err := NewPaiWithName(event.Bakaze)
	if err != nil {
		return err
	}
	s.bakaze = *bakaze
	s.kyokuNum = event.Kyoku
	s.honba = event.Honba
	s.oya = &s.players[event.Oya]
	s.doraMarkers = make([]Pai, 0, MaxNumDoraMarkers)
	doraMarker, err := NewPaiWithName(event.DoraMarker)
	if err != nil {
		return err
	}
	s.doraMarkers = append(s.doraMarkers, *doraMarker)
	s.numPipais = NumInitPipais

	for i := range NumPlayers {
		tehais := make([]Pai, initTehaisSize)
		for j := range tehais {
			tehai, err := NewPaiWithName(event.Tehais[i][j])
			if err != nil {
				return err
			}
			tehais[j] = *tehai
		}

		var err error
		if event.Scores != nil {
			err = s.players[i].onStartKyoku(tehais, &event.Scores[i])
		} else {
			err = s.players[i].onStartKyoku(tehais, nil)
		}
		if err != nil {
			return err
		}
	}

	s.prevEventType = noEvent
	s.prevDahaiActor = noActor
	s.prevDahaiPai = nil

	s.lastActor = noActor
	s.lastActionType = noEvent

	s.kuikaePais = make([]Pai, 0, 3)
	s.missedRon = false
	s.isFuriten = false
	s.isRinshanTsumo = false

	return nil
}

func (s *StateImpl) onTsumo(event *message.Tsumo) error {
	if event == nil {
		return fmt.Errorf("tsumo message is nil")
	}

	if s.numPipais <= 0 {
		return fmt.Errorf("tsumo is not possible if numPipais is 0 or negative: %d", s.numPipais)
	}
	s.numPipais--

	pai, err := NewPaiWithName(event.Pai)
	if err != nil {
		return err
	}
	actor := event.Actor
	player := &s.players[actor]
	err = player.onTsumo(*pai)
	if err != nil {
		return err
	}

	s.lastActor = actor
	s.lastActionType = message.TypeTsumo

	return nil
}

func (s *StateImpl) onDahai(event *message.Dahai) error {
	if event == nil {
		return fmt.Errorf("dahai message is nil")
	}

	pai, err := NewPaiWithName(event.Pai)
	if err != nil {
		return err
	}
	actor := event.Actor
	player := &s.players[actor]
	err = player.onDahai(*pai)
	if err != nil {
		return err
	}

	s.prevDahaiActor = actor
	s.prevDahaiPai = pai

	s.lastActor = actor
	s.lastActionType = message.TypeDahai

	tehaiCounts, err := NewPaiSet(s.players[s.playerID].tehais)
	if err != nil {
		return err
	}

	if actor == s.playerID {
		s.kuikaePais = make([]Pai, 0, 3)
		if player.ReachState() != Accepted {
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

		if err := tehaiCounts.AddPai(pai, 1); err != nil {
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

func (s *StateImpl) onChi(event *message.Chi) error {
	if event == nil {
		return fmt.Errorf("chi message is nil")
	}

	if s.numPipais <= 0 {
		return fmt.Errorf("chi is not possible if numPipais is 0 or negative: %d", s.numPipais)
	}

	pai, err := NewPaiWithName(event.Pai)
	if err != nil {
		return err
	}
	var consumed [2]Pai
	for i, c := range event.Consumed {
		p, err := NewPaiWithName(c)
		if err != nil {
			return err
		}
		consumed[i] = *p
	}
	furo, err := NewChi(*pai, consumed, event.Target)
	if err != nil {
		return err
	}

	actor := event.Actor
	err = s.players[actor].onChi(furo)
	if err != nil {
		return err
	}

	target := event.Target
	err = s.players[target].onTargeted(furo)
	if err != nil {
		return err
	}

	s.lastActor = actor
	s.lastActionType = message.TypeChi

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

func (s *StateImpl) onPon(event *message.Pon) error {
	if event == nil {
		return fmt.Errorf("pon message is nil")
	}

	if s.numPipais <= 0 {
		return fmt.Errorf("pon is not possible if numPipais is 0 or negative: %d", s.numPipais)
	}

	pai, err := NewPaiWithName(event.Pai)
	if err != nil {
		return err
	}
	var consumed [2]Pai
	for i, c := range event.Consumed {
		p, err := NewPaiWithName(c)
		if err != nil {
			return err
		}
		consumed[i] = *p
	}
	furo, err := NewPon(*pai, consumed, event.Target)
	if err != nil {
		return err
	}

	actor := event.Actor
	err = s.players[actor].onPon(furo)
	if err != nil {
		return err
	}

	target := event.Target
	err = s.players[target].onTargeted(furo)
	if err != nil {
		return err
	}

	s.lastActor = actor
	s.lastActionType = message.TypePon

	if actor == s.playerID {
		s.kuikaePais = append(s.kuikaePais, *pai.RemoveRed())
		if !pai.IsTsupai() && pai.Number() == 5 {
			s.kuikaePais = append(s.kuikaePais, *pai.AddRed())
		}
	}

	return nil
}

func (s *StateImpl) onDaiminkan(event *message.Daiminkan) error {
	if event == nil {
		return fmt.Errorf("daiminkan message is nil")
	}

	if s.numPipais <= 0 {
		return fmt.Errorf("daiminkan is not possible if numPipais is 0 or negative: %d", s.numPipais)
	}

	pai, err := NewPaiWithName(event.Pai)
	if err != nil {
		return err
	}
	var consumed [3]Pai
	for i, c := range event.Consumed {
		p, err := NewPaiWithName(c)
		if err != nil {
			return err
		}
		consumed[i] = *p
	}
	furo, err := NewDaiminkan(*pai, consumed, event.Target)
	if err != nil {
		return err
	}

	actor := event.Actor
	err = s.players[actor].onDaiminkan(furo)
	if err != nil {
		return err
	}

	target := event.Target
	err = s.players[target].onTargeted(furo)
	if err != nil {
		return err
	}

	s.lastActor = actor
	s.lastActionType = message.TypeDaiminkan

	if actor == s.playerID {
		s.isRinshanTsumo = true
	}

	return nil
}

func (s *StateImpl) onAnkan(event *message.Ankan) error {
	if event == nil {
		return fmt.Errorf("ankan message is nil")
	}

	if s.numPipais <= 0 {
		return fmt.Errorf("ankan is not possible if numPipais is 0 or negative: %d", s.numPipais)
	}

	var consumed [4]Pai
	for i, c := range event.Consumed {
		p, err := NewPaiWithName(c)
		if err != nil {
			return err
		}
		consumed[i] = *p
	}
	furo, err := NewAnkan(consumed)
	if err != nil {
		return err
	}

	actor := event.Actor
	err = s.players[actor].onAnkan(furo)
	if err != nil {
		return err
	}

	s.lastActor = actor
	s.lastActionType = message.TypeAnkan

	if actor == s.playerID {
		s.isRinshanTsumo = true
	}

	return nil
}

func (s *StateImpl) onKakan(event *message.Kakan) error {
	if event == nil {
		return fmt.Errorf("kakan message is nil")
	}

	if s.numPipais <= 0 {
		return fmt.Errorf("kakan is not possible if numPipais is 0 or negative: %d", s.numPipais)
	}

	pai, err := NewPaiWithName(event.Pai)
	if err != nil {
		return err
	}
	var consumed [3]Pai
	for i, c := range event.Consumed {
		p, err := NewPaiWithName(c)
		if err != nil {
			return err
		}
		consumed[i] = *p
	}
	furo, err := NewKakan(*pai, consumed, nil)
	if err != nil {
		return err
	}

	actor := event.Actor
	err = s.players[actor].onKakan(furo)
	if err != nil {
		return err
	}

	// For chankan
	s.prevDahaiActor = actor
	s.prevDahaiPai = pai

	s.lastActor = actor
	s.lastActionType = message.TypeKakan

	if actor == s.playerID {
		s.isRinshanTsumo = true
	} else {
		tehaiCounts, err := NewPaiSet(s.players[s.playerID].tehais)
		if err != nil {
			return err
		}

		if s.missedRon {
			// The previous ron-able tile was missed
			s.isFuriten = true
		}

		if err := tehaiCounts.AddPai(pai, 1); err != nil {
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

func (s *StateImpl) onDora(event *message.Dora) error {
	if event == nil {
		return fmt.Errorf("dora message is nil")
	}

	if len(s.doraMarkers) >= MaxNumDoraMarkers {
		return fmt.Errorf("a 6th dora cannot be added")
	}

	pai, err := NewPaiWithName(event.DoraMarker)
	if err != nil {
		return err
	}
	s.doraMarkers = append(s.doraMarkers, *pai)

	return nil
}

func (s *StateImpl) onReach(event *message.Reach) error {
	if event == nil {
		return fmt.Errorf("reach message is nil")
	}

	if s.numPipais <= 0 {
		return fmt.Errorf("reach is not possible if numPipais is 0 or negative: %d", s.numPipais)
	}

	actor := event.Actor
	player := &s.players[actor]
	err := player.onReach()
	if err != nil {
		return err
	}

	s.lastActor = actor
	s.lastActionType = message.TypeReach

	return nil
}

func (s *StateImpl) onReachAccepted(event *message.ReachAccepted) error {
	if event == nil {
		return fmt.Errorf("reach_accepted message is nil")
	}

	actor := event.Actor
	player := &s.players[actor]
	var err error
	if event.Scores != nil {
		err = player.onReachAccepted(&event.Scores[actor])
	} else {
		err = player.onReachAccepted(nil)
	}
	if err != nil {
		return err
	}

	return nil
}

func (s *StateImpl) onHora(event *message.Hora) error {
	if event == nil {
		return fmt.Errorf("hora message is nil")
	}

	if event.Scores != nil {
		for i, score := range event.Scores {
			s.players[i].SetScore(score)
		}
	}

	// After hora, only end_kyoku comes, so reset the last action.
	s.lastActor = noActor
	s.lastActionType = noEvent

	return nil
}

func (s *StateImpl) onRyukyoku(event *message.Ryukyoku) error {
	if event == nil {
		return fmt.Errorf("ryukyoku message is nil")
	}

	if event.Scores != nil {
		for i, score := range event.Scores {
			s.players[i].SetScore(score)
		}
	}

	// After ryukyoku, only end_kyoku comes, so reset the last action.
	s.lastActor = noActor
	s.lastActionType = noEvent

	return nil
}
