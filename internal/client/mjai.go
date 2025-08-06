package client

import (
	"fmt"
	"io"
	"os"

	"github.com/Apricot-S/mjai-manue-go/internal/agent"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/protocol/mjai"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type MjaiClient struct {
	reader io.Reader
	writer io.Writer
	agent  agent.Agent
}

func NewMjaiClient(reader io.Reader, writer io.Writer, agent agent.Agent) *MjaiClient {
	return &MjaiClient{reader, writer, agent}
}

func (c *MjaiClient) Run() error {
	adapter := mjai.MjaiAdapter{}
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
		fmt.Fprintf(os.Stderr, "<-\t%s\n", raw)

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

		events := make([]inbound.Event, len(msgs))
		for i, msg := range msgs {
			ev, err := adapter.MessageToEvent(msg)
			if err != nil {
				return err
			}
			events[i] = ev
		}

		resEv, err := c.agent.Respond(events)
		if err != nil {
			return fmt.Errorf("failed to respond from agent: %w", err)
		}
		res, err := adapter.EventToMessage(resEv)
		if err != nil {
			return fmt.Errorf("failed to convert event to message: %w", err)
		}
		fmt.Fprintf(os.Stderr, "->\t%s\n", res)

		res = append(res, '\n')
		if _, err := c.writer.Write(res); err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
	}
}
