package mjai

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/outbound"
	"github.com/go-json-experiment/json"
)

type MjaiAdapter struct{}

func (a *MjaiAdapter) MessageToEvent(rawMsg []byte) (inbound.Event, error) {
	var msg Message
	if err := json.Unmarshal(rawMsg, &msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	switch msg.Type {
	case TypeHello:
		var hello Hello
		if err := json.Unmarshal(rawMsg, &hello); err != nil {
			return nil, fmt.Errorf("failed to unmarshal hello message: %w", err)
		}
		return inbound.NewHello(), nil
	case TypeStartGame:
		var startGame StartGame
		if err := json.Unmarshal(rawMsg, &startGame); err != nil {
			return nil, fmt.Errorf("failed to unmarshal start_game message: %w", err)
		}

		names := [4]string{"", "", "", ""}
		if startGame.Names != nil {
			names = [4]string(startGame.Names)
		}
		return inbound.NewStartGame(startGame.ID, names)
	case TypeStartKyoku:
		var startKyoku StartKyoku
		if err := json.Unmarshal(rawMsg, &startKyoku); err != nil {
			return nil, fmt.Errorf("failed to unmarshal start_kyoku message: %w", err)
		}

		bakaze, err := base.NewPaiWithName(startKyoku.Bakaze)
		if err != nil {
			return nil, err
		}
		doraMarker, err := base.NewPaiWithName(startKyoku.DoraMarker)
		if err != nil {
			return nil, err
		}

		var scores *[4]int = nil
		if startKyoku.Scores != nil {
			scores = (*[4]int)(startKyoku.Scores)
		}

		tehais := [4][13]base.Pai{}
		for i, tehai := range startKyoku.Tehais {
			for n, ts := range tehai {
				tp, err := base.NewPaiWithName(ts)
				if err != nil {
					return nil, err
				}
				tehais[i][n] = *tp
			}
		}

		return inbound.NewStartKyoku(
			*bakaze,
			startKyoku.Kyoku,
			startKyoku.Honba,
			startKyoku.Kyotaku,
			startKyoku.Oya,
			*doraMarker,
			scores,
			tehais,
		)
	case TypeTsumo:
		var tsumo Tsumo
		if err := json.Unmarshal(rawMsg, &tsumo); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tsumo message: %w", err)
		}

		pai, err := base.NewPaiWithName(tsumo.Pai)
		if err != nil {
			return nil, err
		}

		return inbound.NewTsumo(tsumo.Actor, *pai)
	case TypeDahai:
		var dahai Dahai
		if err := json.Unmarshal(rawMsg, &dahai); err != nil {
			return nil, fmt.Errorf("failed to unmarshal dahai message: %w", err)
		}

		pai, err := base.NewPaiWithName(dahai.Pai)
		if err != nil {
			return nil, err
		}

		return inbound.NewDahai(dahai.Actor, *pai, dahai.Tsumogiri)
	case TypeChi:
		var chi Chi
		if err := json.Unmarshal(rawMsg, &chi); err != nil {
			return nil, fmt.Errorf("failed to unmarshal chi message: %w", err)
		}

		taken, err := base.NewPaiWithName(chi.Pai)
		if err != nil {
			return nil, err
		}

		consumed := [2]base.Pai{}
		for i, c := range chi.Consumed {
			p, err := base.NewPaiWithName(c)
			if err != nil {
				return nil, err
			}
			consumed[i] = *p
		}

		return inbound.NewChi(chi.Actor, chi.Target, *taken, consumed)
	case TypePon:
		var pon Pon
		if err := json.Unmarshal(rawMsg, &pon); err != nil {
			return nil, fmt.Errorf("failed to unmarshal pon message: %w", err)
		}

		taken, err := base.NewPaiWithName(pon.Pai)
		if err != nil {
			return nil, err
		}

		consumed := [2]base.Pai{}
		for i, c := range pon.Consumed {
			p, err := base.NewPaiWithName(c)
			if err != nil {
				return nil, err
			}
			consumed[i] = *p
		}

		return inbound.NewPon(pon.Actor, pon.Target, *taken, consumed)
	case TypeDaiminkan:
		var daiminkan Daiminkan
		if err := json.Unmarshal(rawMsg, &daiminkan); err != nil {
			return nil, fmt.Errorf("failed to unmarshal daiminkan message: %w", err)
		}

		taken, err := base.NewPaiWithName(daiminkan.Pai)
		if err != nil {
			return nil, err
		}

		consumed := [3]base.Pai{}
		for i, c := range daiminkan.Consumed {
			p, err := base.NewPaiWithName(c)
			if err != nil {
				return nil, err
			}
			consumed[i] = *p
		}

		return inbound.NewDaiminkan(daiminkan.Actor, daiminkan.Target, *taken, consumed)
	case TypeAnkan:
		var ankan Ankan
		if err := json.Unmarshal(rawMsg, &ankan); err != nil {
			return nil, fmt.Errorf("failed to unmarshal ankan message: %w", err)
		}
		return ankan.ToEvent()
	case TypeKakan:
		var kakan Kakan
		if err := json.Unmarshal(rawMsg, &kakan); err != nil {
			return nil, fmt.Errorf("failed to unmarshal kakan message: %w", err)
		}

		added, err := base.NewPaiWithName(kakan.Pai)
		if err != nil {
			return nil, err
		}

		consumed := [2]base.Pai{}
		for i, c := range kakan.Consumed[:2] {
			p, err := base.NewPaiWithName(c)
			if err != nil {
				return nil, err
			}
			consumed[i] = *p
		}

		taken, err := base.NewPaiWithName(kakan.Consumed[2])
		if err != nil {
			return nil, err
		}

		return inbound.NewKakan(kakan.Actor, kakan.Actor, *taken, consumed, *added)
	default:
		return nil, fmt.Errorf("unknown message type: %s", msg.Type)
	}
}

func (a *MjaiAdapter) EventToMessage(ev outbound.Event) ([]byte, error) {
	switch e := ev.(type) {
	case *outbound.None:
		msg := NewNoneFromEvent(e)
		return json.Marshal(msg)
	case *outbound.Join:
		msg := NewJoinFromEvent(e)
		return json.Marshal(msg)
	case *outbound.Dahai:
		msg, err := NewDahaiFromEvent(e)
		if err != nil {
			return nil, fmt.Errorf("failed to create dahai message: %w", err)
		}
		return json.Marshal(msg)
	case *outbound.Skip:
		msg, err := NewSkipFromEvent(e)
		if err != nil {
			return nil, fmt.Errorf("failed to create skip message: %w", err)
		}
		return json.Marshal(msg)
	case *outbound.Chi:
		msg, err := NewChiFromEvent(e)
		if err != nil {
			return nil, fmt.Errorf("failed to create chi message: %w", err)
		}
		return json.Marshal(msg)
	case *outbound.Pon:
		msg, err := NewPonFromEvent(e)
		if err != nil {
			return nil, fmt.Errorf("failed to create pon message: %w", err)
		}
		return json.Marshal(msg)
	case *outbound.Daiminkan:
		consumed := [3]string{}
		for i, p := range e.Consumed {
			consumed[i] = p.ToString()
		}

		msg, err := NewDaiminkan(e.Actor, e.Target, e.Taken.ToString(), consumed, e.Log)
		if err != nil {
			return nil, fmt.Errorf("failed to create daiminkan message: %w", err)
		}
		return json.Marshal(msg)
	case *outbound.Ankan:
		consumed := [4]string{}
		for i, p := range e.Consumed {
			consumed[i] = p.ToString()
		}

		msg, err := NewAnkan(e.Actor, consumed, e.Log)
		if err != nil {
			return nil, fmt.Errorf("failed to create ankan message: %w", err)
		}
		return json.Marshal(msg)
	case *outbound.Kakan:
		consumed := [3]string{}
		consumed[0] = e.Consumed[0].ToString()
		consumed[1] = e.Consumed[1].ToString()
		consumed[2] = e.Taken.ToString()
		msg, err := NewKakan(e.Actor, e.Added.ToString(), consumed, e.Log)
		if err != nil {
			return nil, fmt.Errorf("failed to create kakan message: %w", err)
		}
		return json.Marshal(msg)
	case *outbound.Reach:
		msg, err := NewReach(e.Actor, e.Log)
		if err != nil {
			return nil, fmt.Errorf("failed to create reach message: %w", err)
		}
		return json.Marshal(msg)
	case *outbound.Hora:
		msg, err := NewHora(e.Actor, e.Target, e.Pai.ToString(), 0, nil, e.Log)
		if err != nil {
			return nil, fmt.Errorf("failed to create hora message: %w", err)
		}
		return json.Marshal(msg)
	default:
		return nil, fmt.Errorf("unknown event type: %T", ev)
	}
}
