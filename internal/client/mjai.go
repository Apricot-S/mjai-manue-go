package client

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"io"
	"os"

	"github.com/Apricot-S/mjai-manue-go/internal/agent"
	"github.com/Apricot-S/mjai-manue-go/internal/protocol/mjai"
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
	decoder := jsontext.NewDecoder(c.reader)
	var raw jsontext.Value

	for {
		if err := json.UnmarshalDecode(decoder, &raw); err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("failed to read message: %w", err)
		}
		fmt.Fprintf(os.Stderr, "<-\t%s\n", raw)

		events, err := mjai.Adapter.DecodeMessages(raw)
		if err != nil {
			return err
		}

		resEv, err := c.agent.Respond(events)
		if err != nil {
			return fmt.Errorf("failed to respond from agent: %w", err)
		}

		res, err := mjai.Adapter.EncodeResponse(resEv)
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
