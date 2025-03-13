package client

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/agent"
	"github.com/go-json-experiment/json/jsontext"
)

// For test.
type EchoAgent struct{}

func (a *EchoAgent) Respond(msg *jsontext.Value) (jsontext.Value, error) {
	return *msg, nil
}

func TestClient_Run(t *testing.T) {
	type fields struct {
		reader io.Reader
		writer io.Writer
		agent  agent.Agent
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "empty reader should return error",
			fields: fields{
				reader: strings.NewReader(""),
				writer: &bytes.Buffer{},
				agent:  &EchoAgent{},
			},
			wantErr: true,
		},
		{
			name: "invalid JSON input should return error",
			fields: fields{
				reader: strings.NewReader("invalid json"),
				writer: &bytes.Buffer{},
				agent:  &EchoAgent{},
			},
			wantErr: true,
		},
		{
			name: "single message",
			fields: fields{
				reader: strings.NewReader(`{"type":"none"}`),
				writer: &bytes.Buffer{},
				agent:  &EchoAgent{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				reader: tt.fields.reader,
				writer: tt.fields.writer,
				agent:  tt.fields.agent,
			}
			if err := c.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Client.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
