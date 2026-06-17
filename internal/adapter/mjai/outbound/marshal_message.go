package outbound

import (
	"encoding/json/v2"
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
)

func MarshalMessage(msg Message) ([]byte, error) {
	return json.Marshal(msg)
}

func ToMessage(a action.Action, log string) (Message, error) {
	switch a := a.(type) {
	case *action.Pass:
		return NewPass(a, log), nil
	case *action.Discard:
		return NewDahai(a, log), nil
	case *action.Chii:
		return NewChi(a, log), nil
	case *action.Pon:
		return NewPon(a, log), nil
	case *action.CalledKan:
		return NewDaiminkan(a, log), nil
	case *action.ConcealedKan:
		return NewAnkan(a, log), nil
	case *action.PromotedKan:
		return NewKakan(a, log), nil
	case *action.Riichi:
		return NewReach(a, log), nil
	case *action.Win:
		return NewHora(a, log), nil
	case *action.Kyushukyuhai:
		return NewKyushukyuhai(a, log), nil
	default:
		return nil, fmt.Errorf("unsupported action type: %T", a)
	}
}
