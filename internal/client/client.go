package client

import (
	"fmt"
	"io"

	"github.com/Apricot-S/mjai-manue-go/internal/agent"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Client struct {
	reader io.Reader
	writer io.Writer
	agent  agent.Agent
}

func NewClient(reader io.Reader, writer io.Writer, agent agent.Agent) *Client {
	return &Client{reader, writer, agent}
}

func (c *Client) Run() error {
	var raw jsontext.Value
	if err := json.UnmarshalRead(c.reader, &raw); err != nil {
		return fmt.Errorf("failed to read message: %w", err)
	}

	switch raw.Kind() {
	case '{':
		// single object
	case '[':
		// array
	default:
		return fmt.Errorf("invalid message: %v", raw)
	}

	return nil
}
