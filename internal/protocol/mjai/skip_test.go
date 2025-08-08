package mjai

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/game/event/outbound"
	"github.com/go-json-experiment/json"
)

func TestNewSkip(t *testing.T) {
	type args struct {
		actor int
		log   string
	}
	tests := []struct {
		name    string
		args    args
		want    *Skip
		wantErr bool
	}{
		{
			name: "without log",
			args: args{
				actor: 0,
				log:   "",
			},
			want: &Skip{
				Action: Action{
					Message: Message{TypeNone},
					Actor:   0,
					Log:     "",
				},
			},
			wantErr: false,
		},
		{
			name: "with log",
			args: args{
				actor: 3,
				log:   "test",
			},
			want: &Skip{
				Action: Action{
					Message: Message{TypeNone},
					Actor:   3,
					Log:     "test",
				},
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
			got, err := NewSkip(tt.args.actor, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSkip() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSkip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSkip_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		args    *Skip
		want    string
		wantErr bool
	}{
		{
			name: "without log",
			args: &Skip{
				Action: Action{
					Message: Message{TypeNone},
					Actor:   0,
					Log:     "",
				},
			},
			want:    `{"type":"none","actor":0}`,
			wantErr: false,
		},
		{
			name: "with log",
			args: &Skip{
				Action: Action{
					Message: Message{TypeNone},
					Actor:   3,
					Log:     "test",
				},
			},
			want:    `{"type":"none","actor":3,"log":"test"}`,
			wantErr: false,
		},
		{
			name: "empty type",
			args: &Skip{
				Action: Action{
					Message: Message{""},
					Actor:   0,
					Log:     "",
				},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid type",
			args: &Skip{
				Action: Action{
					Message: Message{TypeHello},
					Actor:   0,
					Log:     "",
				},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid actor min",
			args: &Skip{
				Action: Action{
					Message: Message{TypeNone},
					Actor:   -1,
					Log:     "",
				},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: &Skip{
				Action: Action{
					Message: Message{TypeNone},
					Actor:   4,
					Log:     "",
				},
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

func TestSkip_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    Skip
		wantErr bool
	}{
		{
			name: "without actor",
			args: `{"type":"none"}`,
			want: Skip{
				Action: Action{
					Message: Message{TypeNone},
					Actor:   0,
					Log:     "",
				},
			},
			wantErr: false,
		},
		{
			name: "without log",
			args: `{"type":"none","actor":0}`,
			want: Skip{
				Action: Action{
					Message: Message{TypeNone},
					Actor:   0,
					Log:     "",
				},
			},
			wantErr: false,
		},
		{
			name: "with log",
			args: `{"type":"none","actor":3,"log":"test"}`,
			want: Skip{
				Action: Action{
					Message: Message{TypeNone},
					Actor:   3,
					Log:     "test",
				},
			},
			wantErr: false,
		},
		{
			name: "empty type",
			args: `{"type":""}`,
			want: Skip{
				Action: Action{
					Message: Message{""},
					Actor:   0,
					Log:     "",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			args: `{"type":"hello"}`,
			want: Skip{
				Action: Action{
					Message: Message{TypeHello},
					Actor:   0,
					Log:     "",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid actor min",
			args: `{"type":"none","actor":-1}`,
			want: Skip{
				Action: Action{
					Message: Message{TypeNone},
					Actor:   -1,
					Log:     "",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: `{"type":"none","actor":4}`,
			want: Skip{
				Action: Action{
					Message: Message{TypeNone},
					Actor:   4,
					Log:     "",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Skip
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

func TestNewSkipFromEvent(t *testing.T) {
	valid, _ := outbound.NewSkip(1, "test")
	invalid := *valid
	invalid.Actor = 4

	type args struct {
		ev *outbound.Skip
	}
	tests := []struct {
		name    string
		args    args
		want    *Skip
		wantErr bool
	}{
		{
			name: "valid",
			args: args{valid},
			want: &Skip{
				Action: Action{
					Message: Message{Type: TypeNone},
					Actor:   1,
					Log:     "test",
				},
			},
			wantErr: false,
		},
		{
			name:    "invalid",
			args:    args{&invalid},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSkipFromEvent(tt.args.ev)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSkipFromEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSkipFromEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
