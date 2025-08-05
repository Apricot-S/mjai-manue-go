package mjai

import (
	"reflect"
	"testing"

	"github.com/go-json-experiment/json"
)

func TestNewNone(t *testing.T) {
	tests := []struct {
		name    string
		want    *None
		wantErr bool
	}{
		{
			name: "test NewNone()",
			want: &None{
				Message: Message{Type: TypeNone},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewNone()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNone_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		args    *None
		want    string
		wantErr bool
	}{
		{
			name: "valid",
			args: &None{
				Message: Message{Type: TypeNone},
			},
			want:    `{"type":"none"}`,
			wantErr: false,
		},
		{
			name: "empty type",
			args: &None{
				Message: Message{Type: ""},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid type",
			args: &None{
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

func TestNone_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    None
		wantErr bool
	}{
		{
			name: "valid",
			args: `{"type":"none"}`,
			want: None{
				Message: Message{Type: TypeNone},
			},
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
			want: None{
				Message: Message{Type: TypeNone},
			},
			wantErr: false,
		},
		{
			name: "empty type",
			args: `{"type":""}`,
			want: None{
				Message: Message{Type: ""},
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			args: `{"type":"hello"}`,
			want: None{
				Message: Message{Type: TypeHello},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got None
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
