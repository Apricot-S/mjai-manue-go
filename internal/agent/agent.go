package agent

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/message"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Agent interface {
	Respond(msgs []jsontext.Value) (jsontext.Value, error)
}

func makeNoneResponse() (jsontext.Value, error) {
	none := message.NewNone()
	res, err := json.Marshal(&none)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal none message: %w", err)
	}
	return res, nil
}

func makeJoinResponse(name string, room string) (jsontext.Value, error) {
	join := message.NewJoin(name, room)
	res, err := json.Marshal(&join)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal join message: %w", err)
	}
	return res, nil
}
