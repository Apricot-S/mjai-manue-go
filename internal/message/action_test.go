package message

import (
	"reflect"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/go-playground/validator/v10"
)

func TestAction_Serialize(t *testing.T) {
	type testCase struct {
		name     string
		input    string
		wantMsg  *Action
		wantJSON string
		wantErr  bool
	}
	tests := []testCase{
		{
			name:     "min_actor",
			input:    `{"type":"action","actor":0}`,
			wantMsg:  &Action{Type: "action", Actor: 0},
			wantJSON: `{"type":"action","actor":0}`,
			wantErr:  false,
		},
		{
			name:     "max_actor",
			input:    `{"type":"action","actor":3}`,
			wantMsg:  &Action{Type: "action", Actor: 3},
			wantJSON: `{"type":"action","actor":3}`,
			wantErr:  false,
		},
		{
			name:     "actor_out_of_range_lower",
			input:    `{"type":"action","actor":-1}`,
			wantMsg:  &Action{Type: "action", Actor: -1},
			wantJSON: `{"type":"action","actor":-1}`,
			wantErr:  false,
		},
		{
			name:     "actor_out_of_range_upper",
			input:    `{"type":"action","actor":4}`,
			wantMsg:  &Action{Type: "action", Actor: 4},
			wantJSON: `{"type":"action","actor":4}`,
			wantErr:  false,
		},
		{
			name:     "missing_actor_is_treated_as_0",
			input:    `{"type":"action"}`,
			wantMsg:  &Action{Type: "action"},
			wantJSON: `{"type":"action","actor":0}`,
			wantErr:  false,
		},
		{
			name:     "null_is_treated_as_0",
			input:    `{"type":"action","actor":null}`,
			wantMsg:  &Action{Type: "action"},
			wantJSON: `{"type":"action","actor":0}`,
			wantErr:  false,
		},
		{
			name:     "undefined",
			input:    `{"type":"action","actor":undefined}`,
			wantMsg:  &Action{Type: "action"},
			wantJSON: `{"type":"action","actor":0}`,
			wantErr:  true,
		},
		{
			name:     "invalid_json",
			input:    `{"type":"action","actor":}`,
			wantMsg:  &Action{Type: "action"},
			wantJSON: `{"type":"action","actor":0}`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var msg Action
			err := json.Unmarshal([]byte(tt.input), &msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshal error: %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(msg, *tt.wantMsg) {
				t.Errorf("expected message '%v', got '%v'", tt.wantMsg, msg)
			}

			jsonData, err := json.Marshal(msg)
			if err != nil {
				t.Errorf("marshal error: %v", err)
				return
			}

			if !reflect.DeepEqual(string(jsonData), tt.wantJSON) {
				t.Errorf("expected JSON '%v', got '%v'", tt.wantJSON, string(jsonData))
			}
		})
	}
}

func TestAction_Validate(t *testing.T) {
	type testCase struct {
		name    string
		msg     *Action
		wantErr bool
	}
	tests := []testCase{
		{
			name:    "min_actor",
			msg:     &Action{Type: "action", Actor: 0},
			wantErr: false,
		},
		{
			name:    "max_actor",
			msg:     &Action{Type: "action", Actor: 3},
			wantErr: false,
		},
		{
			name:    "actor_out_of_range_lower",
			msg:     &Action{Type: "action", Actor: -1},
			wantErr: true,
		},
		{
			name:    "actor_out_of_range_upper",
			msg:     &Action{Type: "action", Actor: 4},
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
