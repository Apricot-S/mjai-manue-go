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
		if err == io.EOF {
			return nil
		}
		return fmt.Errorf("failed to read message: %w", err)
	}

	switch raw.Kind() {
	case '{':
		// single object
		res, err := c.agent.Respond(&raw)
		if err != nil {
			return err
		}

		rawRes, err := json.Marshal(&res)
		if err != nil {
			return fmt.Errorf("failed to marshal response: %w", err)
		}

		if _, err := c.writer.Write(rawRes); err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	case '[':
		// array
		var array []jsontext.Value
		if err := json.Unmarshal(raw, &array); err != nil {
			return fmt.Errorf("failed to unmarshal array: %w", err)
		}

		var lastRawRes jsontext.Value
		for _, item := range array {
			res, err := c.agent.Respond(&item)
			if err != nil {
				return err
			}

			rawRes, err := json.Marshal(&res)
			if err != nil {
				return fmt.Errorf("failed to marshal response: %w", err)
			}
			lastRawRes = rawRes
		}

		if _, err := c.writer.Write(lastRawRes); err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	default:
		return fmt.Errorf("invalid message: %v", raw)
	}

	return nil
}
