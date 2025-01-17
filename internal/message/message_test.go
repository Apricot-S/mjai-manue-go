package message

import (
	"reflect"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/go-playground/validator/v10"
)

func TestMessage_Serialize(t *testing.T) {
	type testCase struct {
		name     string
		input    string
		wantMsg  *Message
		wantJSON string
		wantErr  bool
	}
	tests := []testCase{
		{
			name:     "none",
			input:    `{"type":"none"}`,
			wantMsg:  &Message{Type: "none"},
			wantJSON: `{"type":"none"}`,
			wantErr:  false,
		},
		{
			name:     "start_game",
			input:    `{"type":"start_game","id":1,"names":["shanten","shanten","shanten","shanten"]}`,
			wantMsg:  &Message{Type: "start_game"},
			wantJSON: `{"type":"start_game"}`,
			wantErr:  false,
		},
		{
			name:     "empty",
			input:    `{}`,
			wantMsg:  &Message{Type: ""},
			wantJSON: `{"type":""}`,
			wantErr:  false,
		},
		{
			name:     "null",
			input:    `{"type":null}`,
			wantMsg:  &Message{Type: ""},
			wantJSON: `{"type":""}`,
			wantErr:  false,
		},
		{
			name:     "broken1",
			input:    `{"type":"none"`,
			wantMsg:  &Message{Type: "none"},
			wantJSON: `{"type":"none"}`,
			wantErr:  true,
		},
		{
			name:     "broken2",
			input:    `{"type":"none}`,
			wantMsg:  &Message{Type: ""},
			wantJSON: `{"type":""}`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var msg Message
			err := json.Unmarshal([]byte(tt.input), &msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshal error: %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(msg.Type, tt.wantMsg.Type) {
				t.Errorf("expected type '%s', got '%s'", tt.wantMsg.Type, msg.Type)
			}

			jsonData, err := json.Marshal(msg)
			if err != nil {
				t.Errorf("marshal error: %v", err)
			}

			if !reflect.DeepEqual(string(jsonData), tt.wantJSON) {
				t.Errorf("expected JSON '%s', got '%s'", tt.wantJSON, string(jsonData))
			}
		})
	}
}

func TestMessage_Validate(t *testing.T) {
	type testCase struct {
		name    string
		msg     *Message
		wantErr bool
	}
	tests := []testCase{
		{
			name:    "none",
			msg:     &Message{Type: "none"},
			wantErr: false,
		},
		{
			name:    "start_game",
			msg:     &Message{Type: "start_game"},
			wantErr: false,
		},
		{
			name:    "empty",
			msg:     &Message{Type: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validate := validator.New()

			err := validate.Struct(tt.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("validate.Struct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
