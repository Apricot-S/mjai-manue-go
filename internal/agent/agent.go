package agent

import (
	"github.com/go-json-experiment/json/jsontext"
)

type Agent interface {
	Respond(msgs []jsontext.Value) (jsontext.Value, error)
}
