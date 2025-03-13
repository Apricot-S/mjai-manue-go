package agent

import (
	"github.com/go-json-experiment/json/jsontext"
)

type Agent interface {
	Respond(msg *jsontext.Value) (jsontext.Value, error)
}
