package message

import (
	"reflect"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/go-playground/validator/v10"
)

func TestNone(t *testing.T) {
	type testCase struct {
		name     string
		input    string
		wantMsg  *None
		wantJSON string
		wantErr  bool
	}
	tests := []testCase{
		{
			name:     "only type",
			input:    `{"type":"none"}`,
			wantMsg:  &None{Message{Type: "none"}},
			wantJSON: `{"type":"none"}`,
			wantErr:  false,
		},
		{
			name: "with metadata",
			input: `{
				"type":"none",
				"metadata": {
					"foo": "bar"
				}
			}`,
			wantMsg:  &None{Message{Type: "none"}},
			wantJSON: `{"type":"none"}`,
			wantErr:  false,
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

func TestNone_Validate(t *testing.T) {
	type testCase struct {
		name    string
		msg     *None
		wantErr bool
	}
	tests := []testCase{
		{
			name:    "valid",
			msg:     &None{Message{Type: TypeNone}},
			wantErr: false,
		},
		{
			name:    "invalid",
			msg:     &None{Message{Type: ""}},
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

func TestNewNone(t *testing.T) {
	msg := NewNone()
	if msg.Type != TypeNone {
		t.Errorf("expected type '%s', got '%s'", TypeNone, msg.Type)
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		t.Errorf("marshal error: %v", err)
	}

	wantJSON := `{"type":"none"}`
	if !reflect.DeepEqual(string(jsonData), wantJSON) {
		t.Errorf("expected JSON '%s', got '%s'", wantJSON, string(jsonData))
	}
}
