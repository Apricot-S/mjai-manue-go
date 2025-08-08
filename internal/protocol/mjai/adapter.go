package mjai

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/outbound"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type MjaiAdapter struct{}

func (a *MjaiAdapter) messageToEvent(rawMsg []byte) (inbound.Event, error) {
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
		return hello.ToEvent(), nil
	case TypeStartGame:
		var startGame StartGame
		if err := json.Unmarshal(rawMsg, &startGame); err != nil {
			return nil, fmt.Errorf("failed to unmarshal start_game message: %w", err)
		}
		return startGame.ToEvent()
	case TypeStartKyoku:
		var startKyoku StartKyoku
		if err := json.Unmarshal(rawMsg, &startKyoku); err != nil {
			return nil, fmt.Errorf("failed to unmarshal start_kyoku message: %w", err)
		}
		return startKyoku.ToEvent()
	case TypeTsumo:
		var tsumo Tsumo
		if err := json.Unmarshal(rawMsg, &tsumo); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tsumo message: %w", err)
		}
		return tsumo.ToEvent()
	case TypeDahai:
		var dahai Dahai
		if err := json.Unmarshal(rawMsg, &dahai); err != nil {
			return nil, fmt.Errorf("failed to unmarshal dahai message: %w", err)
		}
		return dahai.ToEvent()
	case TypeChi:
		var chi Chi
		if err := json.Unmarshal(rawMsg, &chi); err != nil {
			return nil, fmt.Errorf("failed to unmarshal chi message: %w", err)
		}
		return chi.ToEvent()
	case TypePon:
		var pon Pon
		if err := json.Unmarshal(rawMsg, &pon); err != nil {
			return nil, fmt.Errorf("failed to unmarshal pon message: %w", err)
		}
		return pon.ToEvent()
	case TypeDaiminkan:
		var daiminkan Daiminkan
		if err := json.Unmarshal(rawMsg, &daiminkan); err != nil {
			return nil, fmt.Errorf("failed to unmarshal daiminkan message: %w", err)
		}
		return daiminkan.ToEvent()
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
		return kakan.ToEvent()
	case TypeDora:
		var dora Dora
		if err := json.Unmarshal(rawMsg, &dora); err != nil {
			return nil, fmt.Errorf("failed to unmarshal dora message: %w", err)
		}
		return dora.ToEvent()
	case TypeReach:
		var reach Reach
		if err := json.Unmarshal(rawMsg, &reach); err != nil {
			return nil, fmt.Errorf("failed to unmarshal reach message: %w", err)
		}
		return reach.ToEvent()
	case TypeReachAccepted:
		var reachAccepted ReachAccepted
		if err := json.Unmarshal(rawMsg, &reachAccepted); err != nil {
			return nil, fmt.Errorf("failed to unmarshal reach accepted message: %w", err)
		}
		return reachAccepted.ToEvent()
	case TypeHora:
		panic("not implemented")
	case TypeRyukyoku:
		var ryukyoku Ryukyoku
		if err := json.Unmarshal(rawMsg, &ryukyoku); err != nil {
			return nil, fmt.Errorf("failed to unmarshal ryukyoku message: %w", err)
		}
		return ryukyoku.ToEvent(), nil
	case TypeEndKyoku:
		var endKyoku EndKyoku
		if err := json.Unmarshal(rawMsg, &endKyoku); err != nil {
			return nil, fmt.Errorf("failed to unmarshal end_kyoku message: %w", err)
		}
		return endKyoku.ToEvent(), nil
	case TypeEndGame:
		var endGame EndGame
		if err := json.Unmarshal(rawMsg, &endGame); err != nil {
			return nil, fmt.Errorf("failed to unmarshal end_game message: %w", err)
		}
		return endGame.ToEvent(), nil
	case TypeError:
		var errorMsg Error
		if err := json.Unmarshal(rawMsg, &errorMsg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal error message: %w", err)
		}
		return errorMsg.ToEvent(), nil
	default:
		return nil, fmt.Errorf("unknown message type: %s", msg.Type)
	}
}

func (a *MjaiAdapter) DecodeMessages(msg []byte) ([]inbound.Event, error) {
	var msgs []jsontext.Value

	switch jsontext.Value(msg).Kind() {
	case '{':
		// single object
		msgs = []jsontext.Value{msg}
	case '[':
		// array
		if err := json.Unmarshal(msg, &msgs); err != nil {
			return nil, fmt.Errorf("failed to unmarshal messages: %w", err)
		}
	default:
		return nil, fmt.Errorf("invalid message: %v", msg)
	}

	events := make([]inbound.Event, len(msgs))
	for i, msg := range msgs {
		ev, err := a.messageToEvent(msg)
		if err != nil {
			return nil, err
		}
		events[i] = ev
	}

	return events, nil
}

func (a *MjaiAdapter) EncodeResponse(ev outbound.Event) ([]byte, error) {
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
		msg, err := NewDaiminkanFromEvent(e)
		if err != nil {
			return nil, fmt.Errorf("failed to create daiminkan message: %w", err)
		}
		return json.Marshal(msg)
	case *outbound.Ankan:
		msg, err := NewAnkanFromEvent(e)
		if err != nil {
			return nil, fmt.Errorf("failed to create ankan message: %w", err)
		}
		return json.Marshal(msg)
	case *outbound.Kakan:
		msg, err := NewKakanFromEvent(e)
		if err != nil {
			return nil, fmt.Errorf("failed to create kakan message: %w", err)
		}
		return json.Marshal(msg)
	case *outbound.Reach:
		msg, err := NewReachFromEvent(e)
		if err != nil {
			return nil, fmt.Errorf("failed to create reach message: %w", err)
		}
		return json.Marshal(msg)
	case *outbound.Hora:
		msg, err := NewHoraFromEvent(e)
		if err != nil {
			return nil, fmt.Errorf("failed to create hora message: %w", err)
		}
		return json.Marshal(msg)
	default:
		return nil, fmt.Errorf("unknown event type: %T", ev)
	}
}
