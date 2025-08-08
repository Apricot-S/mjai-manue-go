package mjai

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/go-json-experiment/json"
)

func TestNewError(t *testing.T) {
	tests := []struct {
		name    string
		want    *Error
		wantErr bool
	}{
		{
			name: "test NewError()",
			want: &Error{
				Message: Message{Type: TypeError},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewError()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		args    *Error
		want    string
		wantErr bool
	}{
		{
			name: "valid",
			args: &Error{
				Message: Message{Type: TypeError},
			},
			want:    `{"type":"error"}`,
			wantErr: false,
		},
		{
			name: "empty type",
			args: &Error{
				Message: Message{Type: ""},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid type",
			args: &Error{
				Message: Message{Type: TypeHello},
			},
			want:    ``,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("marshal error = %v, want %v", err, tt.wantErr)
			}
			if string(got) != tt.want {
				t.Errorf("Marshal() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestError_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    Error
		wantErr bool
	}{
		{
			name: "valid",
			args: `{"type":"error"}`,
			want: Error{
				Message: Message{Type: TypeError},
			},
			wantErr: false,
		},
		{
			name: "with metadata",
			args: `{
				"type":"error",
				"metadata": {
					"foo": "bar"
				}
			}`,
			want: Error{
				Message: Message{Type: TypeError},
			},
			wantErr: false,
		},
		{
			name: "empty type",
			args: `{"type":""}`,
			want: Error{
				Message: Message{Type: ""},
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			args: `{"type":"hello"}`,
			want: Error{
				Message: Message{Type: TypeHello},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Error
			err := json.Unmarshal([]byte(tt.args), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshal error = %v, want %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unmarshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_ToEvent(t *testing.T) {
	type fields struct {
		Message Message
	}
	tests := []struct {
		name   string
		fields fields
		want   *inbound.Error
	}{
		{
			name: "valid",
			fields: fields{
				Message: Message{TypeError},
			},
			want: inbound.NewError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Error{
				Message: tt.fields.Message,
			}
			if got := m.ToEvent(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Error.ToEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
