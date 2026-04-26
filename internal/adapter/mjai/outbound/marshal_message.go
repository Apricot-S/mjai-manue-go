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
		return NewNone(log), nil
	case *action.Discard:
		return NewDahai(a, log), nil
	default:
		return nil, fmt.Errorf("unsupported action type: %T", a)
	}
}
