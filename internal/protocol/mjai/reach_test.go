package mjai

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/game/event/outbound"
	"github.com/go-json-experiment/json"
)

func TestNewReach(t *testing.T) {
	type args struct {
		actor int
		log   string
	}
	tests := []struct {
		name    string
		args    args
		want    *Reach
		wantErr bool
	}{
		{
			name: "without log",
			args: args{
				actor: 0,
				log:   "",
			},
			want: &Reach{
				Action: Action{
					Message: Message{TypeReach},
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
			want: &Reach{
				Action: Action{
					Message: Message{TypeReach},
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
			got, err := NewReach(tt.args.actor, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewReach() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReach() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReach_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		args    *Reach
		want    string
		wantErr bool
	}{
		{
			name: "without log",
			args: &Reach{
				Action: Action{
					Message: Message{TypeReach},
					Actor:   0,
					Log:     "",
				},
			},
			want:    `{"type":"reach","actor":0}`,
			wantErr: false,
		},
		{
			name: "with log",
			args: &Reach{
				Action: Action{
					Message: Message{TypeReach},
					Actor:   3,
					Log:     "test",
				},
			},
			want:    `{"type":"reach","actor":3,"log":"test"}`,
			wantErr: false,
		},
		{
			name: "empty type",
			args: &Reach{
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
			args: &Reach{
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
			args: &Reach{
				Action: Action{
					Message: Message{TypeReach},
					Actor:   -1,
					Log:     "",
				},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: &Reach{
				Action: Action{
					Message: Message{TypeReach},
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

func TestReach_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    Reach
		wantErr bool
	}{
		{
			name: "without actor",
			args: `{"type":"reach"}`,
			want: Reach{
				Action: Action{
					Message: Message{TypeReach},
					Actor:   0,
					Log:     "",
				},
			},
			wantErr: false,
		},
		{
			name: "without log",
			args: `{"type":"reach","actor":0,"pai":"F"}`,
			want: Reach{
				Action: Action{
					Message: Message{TypeReach},
					Actor:   0,
					Log:     "",
				},
			},
			wantErr: false,
		},
		{
			name: "with log",
			args: `{"type":"reach","actor":3,"log":"test"}`,
			want: Reach{
				Action: Action{
					Message: Message{TypeReach},
					Actor:   3,
					Log:     "test",
				},
			},
			wantErr: false,
		},
		{
			name: "empty type",
			args: `{"type":"","actor":0,"pai":"N"}`,
			want: Reach{
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
			args: `{"type":"hello","actor":3,"pai":"N"}`,
			want: Reach{
				Action: Action{
					Message: Message{TypeHello},
					Actor:   3,
					Log:     "",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid actor min",
			args: `{"type":"reach","actor":-1,"pai":"9m"}`,
			want: Reach{
				Action: Action{
					Message: Message{TypeReach},
					Actor:   -1,
					Log:     "",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: `{"type":"reach","actor":4,"pai":"4s"}`,
			want: Reach{
				Action: Action{
					Message: Message{TypeReach},
					Actor:   4,
					Log:     "",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Reach
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

func TestNewReachFromEvent(t *testing.T) {
	valid, _ := outbound.NewReach(1, "test")
	invalid := *valid
	invalid.Actor = 4

	type args struct {
		ev *outbound.Reach
	}
	tests := []struct {
		name    string
		args    args
		want    *Reach
		wantErr bool
	}{
		{
			name: "valid",
			args: args{valid},
			want: &Reach{
				Action: Action{
					Message: Message{Type: TypeReach},
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
			got, err := NewReachFromEvent(tt.args.ev)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewReachFromEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReachFromEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
