package client

import (
	"io"

	"github.com/Apricot-S/mjai-manue-go/internal/agent"
)

type Batch struct {
	reader io.Reader
	writer io.Writer
	agent  *agent.Agent
}
