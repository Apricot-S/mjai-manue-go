package message

import (
	"reflect"
	"testing"

	"github.com/go-json-experiment/json"
)

func TestNewDora(t *testing.T) {
	type args struct {
		doraMarker string
	}
	tests := []struct {
		name    string
		args    args
		want    *Dora
		wantErr bool
	}{
		{
			name: "valid 1m",
			args: args{
				doraMarker: "1m",
			},
			want: &Dora{
				Message:    Message{TypeDora},
				DoraMarker: "1m",
			},
			wantErr: false,
		},
		{
			name: "valid ?",
			args: args{
				doraMarker: "?",
			},
			want: &Dora{
				Message:    Message{TypeDora},
				DoraMarker: "?",
			},
			wantErr: false,
		},
		{
			name: "empty pai",
			args: args{
				doraMarker: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid pai",
			args: args{
				doraMarker: "1z",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDora(tt.args.doraMarker)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDora() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDora() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDora_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		args    *Dora
		want    string
		wantErr bool
	}{
		{
			name: "valid ?",
			args: &Dora{
				Message:    Message{TypeDora},
				DoraMarker: "?",
			},
			want:    `{"type":"dora","dora_marker":"?"}`,
			wantErr: false,
		},
		{
			name: "valid C",
			args: &Dora{
				Message:    Message{TypeDora},
				DoraMarker: "C",
			},
			want:    `{"type":"dora","dora_marker":"C"}`,
			wantErr: false,
		},
		{
			name: "empty type",
			args: &Dora{
				Message:    Message{""},
				DoraMarker: "9s",
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid type",
			args: &Dora{
				Message:    Message{TypeHello},
				DoraMarker: "5p",
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "empty pai",
			args: &Dora{
				Message:    Message{TypeDora},
				DoraMarker: "",
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid pai",
			args: &Dora{
				Message:    Message{TypeDora},
				DoraMarker: "0m",
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

func TestDora_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    Dora
		wantErr bool
	}{
		{
			name: "valid ?",
			args: `{"type":"dora","dora_marker":"?"}`,
			want: Dora{
				Message:    Message{TypeDora},
				DoraMarker: "?",
			},
			wantErr: false,
		},
		{
			name: "valid F",
			args: `{"type":"dora","dora_marker":"F"}`,
			want: Dora{
				Message:    Message{TypeDora},
				DoraMarker: "F",
			},
			wantErr: false,
		},
		{
			name: "empty type",
			args: `{"type":"","dora_marker":"N"}`,
			want: Dora{
				Message:    Message{""},
				DoraMarker: "N",
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			args: `{"type":"hello","dora_marker":"N"}`,
			want: Dora{
				Message:    Message{TypeHello},
				DoraMarker: "N",
			},
			wantErr: true,
		},
		{
			name: "empty pai",
			args: `{"type":"dora","dora_marker":""}`,
			want: Dora{
				Message:    Message{TypeDora},
				DoraMarker: "",
			},
			wantErr: true,
		},
		{
			name: "invalid pai",
			args: `{"type":"dora","dora_marker":"4pr"}`,
			want: Dora{
				Message:    Message{TypeDora},
				DoraMarker: "4pr",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Dora
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
