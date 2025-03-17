package agent

import (
	"github.com/Apricot-S/mjai-manue-go/internal/message"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type TsumogiriAgent struct {
	name     string
	room     string
	playerID int
}

func (a *TsumogiriAgent) Respond(msgs []jsontext.Value) (jsontext.Value, error) {
	lastMsg := msgs[len(msgs)-1]
	var tsumo message.Tsumo
	err := json.Unmarshal(lastMsg, &tsumo)

	if err != nil {
		// Not tsumo
		return makeNone()
	}

	if tsumo.Actor != a.playerID {
		// Not self tsumo
		return makeNone()
	}

	// Dummy implementation
	res := []byte{}
	return res, nil
}
