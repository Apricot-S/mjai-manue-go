package message

import (
	"reflect"
	"testing"

	"github.com/go-json-experiment/json"
)

func TestNewHello(t *testing.T) {
	type args struct {
		protocol        string
		protocolVersion int
	}
	tests := []struct {
		name    string
		args    args
		want    *Hello
		wantErr bool
	}{
		{
			name: "without protocol and protocol version",
			args: args{
				protocol:        "",
				protocolVersion: 0,
			},
			want: &Hello{
				Message: Message{Type: TypeHello},
			},
			wantErr: false,
		},
		{
			name: "with protocol and protocol version",
			args: args{
				protocol:        "mjsonp",
				protocolVersion: 1,
			},
			want: &Hello{
				Message:         Message{Type: TypeHello},
				Protocol:        "mjsonp",
				ProtocolVersion: 1,
			},
			wantErr: false,
		},
		{
			name: "invalid protocol version",
			args: args{
				protocol:        "mjsonp",
				protocolVersion: -1,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewHello(tt.args.protocol, tt.args.protocolVersion)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHello() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHello() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHello_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		args    *Hello
		want    string
		wantErr bool
	}{
		{
			name: "without protocol and protocol version",
			args: &Hello{
				Message: Message{Type: TypeHello},
			},
			want:    `{"type":"hello"}`,
			wantErr: false,
		},
		{
			name: "with protocol and protocol version",
			args: &Hello{
				Message:         Message{Type: TypeHello},
				Protocol:        "mjsonp",
				ProtocolVersion: 1,
			},
			want:    `{"type":"hello","protocol":"mjsonp","protocol_version":1}`,
			wantErr: false,
		},
		{
			name: "empty type",
			args: &Hello{
				Message:         Message{Type: ""},
				Protocol:        "mjsonp",
				ProtocolVersion: 0,
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid type",
			args: &Hello{
				Message:         Message{Type: TypeNone},
				Protocol:        "mjsonp",
				ProtocolVersion: 0,
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid protocol version",
			args: &Hello{
				Message:         Message{Type: TypeHello},
				Protocol:        "mjsonp",
				ProtocolVersion: -1,
			},
			want:    ``,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("marshal error = %v, want %v", err, tt.wantErr)
			}
			if string(b) != tt.want {
				t.Errorf("Marshal() = %v, want %v", string(b), tt.want)
			}
		})
	}
}

func TestHello_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    Hello
		wantErr bool
	}{
		{
			name: "without protocol and protocol version",
			args: `{"type":"hello"}`,
			want: Hello{
				Message:         Message{Type: TypeHello},
				Protocol:        "",
				ProtocolVersion: 0,
			},
			wantErr: false,
		},
		{
			name: "with protocol and protocol version",
			args: `{"type":"hello","protocol":"mjsonp","protocol_version":1}`,
			want: Hello{
				Message:         Message{Type: TypeHello},
				Protocol:        "mjsonp",
				ProtocolVersion: 1,
			},
			wantErr: false,
		},
		{
			name: "empty type",
			args: `{"type":""}`,
			want: Hello{
				Message:         Message{Type: ""},
				Protocol:        "",
				ProtocolVersion: 0,
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			args: `{"type":"none"}`,
			want: Hello{
				Message:         Message{Type: TypeNone},
				Protocol:        "",
				ProtocolVersion: 0,
			},
			wantErr: true,
		},
		{
			name: "invalid protocol version",
			args: `{"type":"hello","protocol":"mjsonp","protocol_version":-1}`,
			want: Hello{
				Message:         Message{Type: TypeHello},
				Protocol:        "mjsonp",
				ProtocolVersion: -1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Hello
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
