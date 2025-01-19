package message

import (
	"reflect"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/go-playground/validator/v10"
)

func TestNone_Serialize(t *testing.T) {
	type testCase struct {
		name     string
		input    string
		wantMsg  *None
		wantJSON string
		wantErr  bool
	}
	tests := []testCase{
		{
			name:     "none",
			input:    `{"type":"none"}`,
			wantMsg:  &None{Type: "none"},
			wantJSON: `{"type":"none"}`,
			wantErr:  false,
		},
		{
			name:     "start_game",
			input:    `{"type":"start_game","id":1,"names":["shanten","shanten","shanten","shanten"]}`,
			wantMsg:  &None{Type: "start_game"},
			wantJSON: `{"type":"start_game"}`,
			wantErr:  false,
		},
		{
			name:     "empty",
			input:    `{}`,
			wantMsg:  &None{Type: ""},
			wantJSON: `{"type":""}`,
			wantErr:  false,
		},
		{
			name:     "null",
			input:    `{"type":null}`,
			wantMsg:  &None{Type: ""},
			wantJSON: `{"type":""}`,
			wantErr:  false,
		},
		{
			name:     "broken1",
			input:    `{"type":"none"`,
			wantMsg:  &None{Type: "none"},
			wantJSON: `{"type":"none"}`,
			wantErr:  true,
		},
		{
			name:     "broken2",
			input:    `{"type":"none}`,
			wantMsg:  &None{Type: ""},
			wantJSON: `{"type":""}`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var msg None
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

func TestNone_Validate(t *testing.T) {
	type testCase struct {
		name    string
		msg     *None
		wantErr bool
	}
	tests := []testCase{
		{
			name:    "none",
			msg:     &None{Type: "none"},
			wantErr: false,
		},
		{
			name:    "start_game",
			msg:     &None{Type: "start_game"},
			wantErr: true,
		},
		{
			name:    "empty",
			msg:     &None{Type: ""},
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
