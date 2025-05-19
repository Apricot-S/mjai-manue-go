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
	decoder := jsontext.NewDecoder(c.reader)
	var raw jsontext.Value
	var msgs []jsontext.Value

	for {
		if err := json.UnmarshalDecode(decoder, &raw); err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("failed to read message: %w", err)
		}

		switch raw.Kind() {
		case '{':
			// single object
			msgs = []jsontext.Value{raw}
		case '[':
			// array
			if err := json.Unmarshal(raw, &msgs); err != nil {
				return fmt.Errorf("failed to unmarshal messages: %w", err)
			}
		default:
			return fmt.Errorf("invalid message: %v", raw)
		}

		res, err := c.agent.Respond(msgs)
		if err != nil {
			return fmt.Errorf("failed to respond from agent: %w", err)
		}

		res = append(res, '\n')

		if _, err := c.writer.Write(res); err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	}
}
