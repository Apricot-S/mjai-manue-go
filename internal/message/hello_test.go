package message

import (
	"reflect"
	"testing"
)

func TestNewHello(t *testing.T) {
	type args struct {
		protocol        string
		protocolVersion int
	}
	tests := []struct {
		name     string
		args     args
		want     *Hello
		wantJSON string
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
			wantJSON: `{"type":"hello"}`,
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
			wantJSON: `{"type":"hello","protocol":"mjsonp","protocol_version":1}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHello(tt.args.protocol, tt.args.protocolVersion); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHello() = %v, want %v", got, tt.want)
			}
		})
	}
}
