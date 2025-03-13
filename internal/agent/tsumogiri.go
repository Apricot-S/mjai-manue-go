package agent

import (
	"github.com/go-json-experiment/json/jsontext"
)

type TsumogiriAgent struct {
	name     string
	playerID int
}

func (a *TsumogiriAgent) Respond(msg *jsontext.Value) (jsontext.Value, error) {
	// Dummy implementation
	return []byte{}, nil
}
