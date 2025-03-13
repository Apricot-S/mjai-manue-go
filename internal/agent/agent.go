package agent

import (
	"github.com/go-json-experiment/json/jsontext"
)

type Agent interface {
	Respond(message *jsontext.Value) (jsontext.Value, error)
}
