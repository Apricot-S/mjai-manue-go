package client

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/agent"
	"github.com/go-json-experiment/json/jsontext"
)

// For test.
type EchoAgent struct{}

func (a *EchoAgent) Respond(msg []jsontext.Value) (jsontext.Value, error) {
	return msg[len(msg)-1], nil
}

func TestClient_Run_OutputValidation(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantOutput string
		agent      agent.Agent
		wantErr    bool
	}{
		{
			name:       "single message",
			input:      `{"type":"none"}`,
			wantOutput: `{"type":"none"}`,
			agent:      &EchoAgent{},
			wantErr:    false,
		},
		{
			name:       "array messages should write last response",
			input:      `[{"type":"none"},{"type":"end_game"}]`,
			wantOutput: `{"type":"end_game"}`,
			agent:      &EchoAgent{},
			wantErr:    false,
		},
		{
			name:       "empty reader",
			input:      ``,
			wantOutput: ``,
			agent:      &EchoAgent{},
			wantErr:    false,
		},
		{
			name:       "invalid JSON input should return error",
			input:      `invalid json`,
			wantOutput: ``,
			agent:      &EchoAgent{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			c := &Client{
				reader: strings.NewReader(tt.input),
				writer: buf,
				agent:  tt.agent,
			}

			err := c.Run()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got := buf.String(); got != tt.wantOutput {
				t.Errorf("Client.Run() output = %q, want %q", got, tt.wantOutput)
			}
		})
	}
}

type ErrorAgent struct {
	err error
}

func (a *ErrorAgent) Respond(msg []jsontext.Value) (jsontext.Value, error) {
	return nil, a.err
}

func TestClient_Run_AgentError(t *testing.T) {
	expectedErr := fmt.Errorf("agent error")
	c := &Client{
		reader: strings.NewReader(`{"type":"none"}`),
		writer: &bytes.Buffer{},
		agent:  &ErrorAgent{err: expectedErr},
	}

	err := c.Run()
	if err == nil {
		t.Error("Client.Run() expected error but got nil")
	}
	if !strings.Contains(err.Error(), expectedErr.Error()) {
		t.Errorf("Client.Run() error = %v, should contain %v", err, expectedErr)
	}
}
