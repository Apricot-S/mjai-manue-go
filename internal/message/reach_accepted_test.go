package message

import (
	"reflect"
	"testing"

	"github.com/go-json-experiment/json"
)

func TestNewReachAccepted(t *testing.T) {
	type args struct {
		actor int
		log   string
	}
	tests := []struct {
		name    string
		args    args
		want    *ReachAccepted
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				actor: 0,
				log:   "",
			},
			want: &ReachAccepted{
				Message: Message{TypeReachAccepted},
				Actor:   0,
			},
			wantErr: false,
		},
		{
			name: "invalid actor min",
			args: args{
				actor: -1,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: args{
				actor: 4,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewReachAccepted(tt.args.actor)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewReachAccepted() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReachAccepted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReachAccepted_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		args    *ReachAccepted
		want    string
		wantErr bool
	}{
		{
			name: "valid",
			args: &ReachAccepted{
				Message: Message{TypeReachAccepted},
				Actor:   0,
			},
			want:    `{"type":"reach_accepted","actor":0}`,
			wantErr: false,
		},
		{
			name: "empty type",
			args: &ReachAccepted{
				Message: Message{""},
				Actor:   0,
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid type",
			args: &ReachAccepted{
				Message: Message{TypeHello},
				Actor:   0,
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid actor min",
			args: &ReachAccepted{
				Message: Message{TypeReachAccepted},
				Actor:   -1,
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: &ReachAccepted{
				Message: Message{TypeReachAccepted},
				Actor:   4,
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

func TestReachAccepted_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    ReachAccepted
		wantErr bool
	}{
		{
			name: "without actor",
			args: `{"type":"reach_accepted"}`,
			want: ReachAccepted{
				Message: Message{TypeReachAccepted},
				Actor:   0,
			},
			wantErr: false,
		},
		{
			name: "valid",
			args: `{"type":"reach_accepted","actor":0}`,
			want: ReachAccepted{
				Message: Message{TypeReachAccepted},
				Actor:   0,
			},
			wantErr: false,
		},
		{
			name: "empty type",
			args: `{"type":"","actor":0}`,
			want: ReachAccepted{
				Message: Message{""},
				Actor:   0,
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			args: `{"type":"hello","actor":3}`,
			want: ReachAccepted{
				Message: Message{TypeHello},
				Actor:   3,
			},
			wantErr: true,
		},
		{
			name: "invalid actor min",
			args: `{"type":"reach_accepted","actor":-1}`,
			want: ReachAccepted{
				Message: Message{TypeReachAccepted},
				Actor:   -1,
			},
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: `{"type":"reach_accepted","actor":4}`,
			want: ReachAccepted{
				Message: Message{TypeReachAccepted},
				Actor:   4,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got ReachAccepted
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
