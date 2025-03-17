package message

import (
	"reflect"
	"testing"

	"github.com/go-json-experiment/json"
)

func TestNewJoin(t *testing.T) {
	type args struct {
		name string
		room string
	}
	tests := []struct {
		name    string
		args    args
		want    *Join
		wantErr bool
	}{
		{
			name: "without name and room",
			args: args{
				name: "",
				room: "",
			},
			want: &Join{
				Message: Message{Type: TypeJoin},
			},
			wantErr: false,
		},
		{
			name: "with name and room",
			args: args{
				name: "shanten",
				room: "default",
			},
			want: &Join{
				Message: Message{Type: TypeJoin},
				Name:    "shanten",
				Room:    "default",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewJoin(tt.args.name, tt.args.room)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewJoin() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJoin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJoin_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		args    *Join
		want    string
		wantErr bool
	}{
		{
			name: "without name and room",
			args: &Join{
				Message: Message{Type: TypeJoin},
			},
			want:    `{"type":"join"}`,
			wantErr: false,
		},
		{
			name: "with name and room",
			args: &Join{
				Message: Message{Type: TypeJoin},
				Name:    "shanten",
				Room:    "default",
			},
			want:    `{"type":"join","name":"shanten","room":"default"}`,
			wantErr: false,
		},
		{
			name: "empty type",
			args: &Join{
				Message: Message{Type: ""},
				Name:    "shanten",
				Room:    "default",
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid type",
			args: &Join{
				Message: Message{Type: TypeNone},
				Name:    "shanten",
				Room:    "default",
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

func TestJoin_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    Join
		wantErr bool
	}{
		{
			name: "without name and room",
			args: `{"type":"join"}`,
			want: Join{
				Message: Message{Type: TypeJoin},
				Name:    "",
				Room:    "",
			},
			wantErr: false,
		},
		{
			name: "with name and room",
			args: `{"type":"join","name":"shanten","room":"default"}`,
			want: Join{
				Message: Message{Type: TypeJoin},
				Name:    "shanten",
				Room:    "default",
			},
			wantErr: false,
		},
		{
			name: "empty type",
			args: `{"type":""}`,
			want: Join{
				Message: Message{Type: ""},
				Name:    "",
				Room:    "",
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			args: `{"type":"none"}`,
			want: Join{
				Message: Message{Type: TypeNone},
				Name:    "",
				Room:    "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Join
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
