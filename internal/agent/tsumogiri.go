package agent

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/message"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type TsumogiriAgent struct {
	name     string
	room     string
	playerID int
}

func NewTsumogiriAgent(name string, room string) *TsumogiriAgent {
	return &TsumogiriAgent{
		name: name,
		room: room,
	}
}

func (a *TsumogiriAgent) Respond(msgs []jsontext.Value) (jsontext.Value, error) {
	lastMsg := msgs[len(msgs)-1]
	var msg message.Message
	if err := json.Unmarshal(lastMsg, &msg); err != nil {
		return nil, err
	}

	switch msg.Type {
	case message.TypeTsumo:
		var tsumo message.Tsumo
		if err := json.Unmarshal(lastMsg, &tsumo); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tsumo message: %w", err)
		}

		if tsumo.Actor != a.playerID {
			// Not self tsumo
			return makeNoneResponse()
		}

		// Dummy implementation
		res := []byte{}
		return res, nil
	case message.TypeHello:
		return makeJoinResponse(a.name, a.room)
	default:
		return makeNoneResponse()
	}
}
