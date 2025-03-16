package message

import (
	"reflect"
	"testing"

	"github.com/go-json-experiment/json"
)

func TestMessage_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		args    *Message
		want    string
		wantErr bool
	}{
		{
			name:    "valid",
			args:    &Message{Type: TypeNone},
			want:    `{"type":"none"}`,
			wantErr: false,
		},
		{
			name:    "empty type",
			args:    &Message{Type: ""},
			want:    `{"type":""}`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.args)
			if err != nil {
				t.Errorf("marshal error = %v, want %v", err, tt.wantErr)
			}
			if string(got) != tt.want {
				t.Errorf("Marshal() = %v, want %v", string(got), tt.want)
			}

			if err := messageValidator.Struct(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("marshal error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestMessage_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    Message
		wantErr bool
	}{
		{
			name:    "valid",
			args:    `{"type":"none"}`,
			want:    Message{Type: TypeNone},
			wantErr: false,
		},
		{
			name: "with metadata",
			args: `{
				"type":"none",
				"metadata": {
					"foo": "bar"
				}
			}`,
			want:    Message{Type: TypeNone},
			wantErr: false,
		},
		{
			name:    "empty type",
			args:    `{"type":""}`,
			want:    Message{Type: ""},
			wantErr: true,
		},
		{
			name:    "null type",
			args:    `{"type":null}`,
			want:    Message{Type: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Message
			err := json.Unmarshal([]byte(tt.args), &got)
			if err != nil {
				t.Errorf("unmarshal error = %v, want %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unmarshal() = %v, want %v", got, tt.want)
			}

			if err := messageValidator.Struct(got); (err != nil) != tt.wantErr {
				t.Errorf("validation error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
