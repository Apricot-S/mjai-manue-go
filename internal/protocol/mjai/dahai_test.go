package mjai

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/outbound"
	"github.com/go-json-experiment/json"
)

func TestNewDahai(t *testing.T) {
	type args struct {
		actor     int
		pai       string
		tsumogiri bool
		log       string
	}
	tests := []struct {
		name    string
		args    args
		want    *Dahai
		wantErr bool
	}{
		{
			name: "without log",
			args: args{
				actor:     0,
				pai:       "?",
				tsumogiri: false,
				log:       "",
			},
			want: &Dahai{
				Action: Action{
					Message: Message{TypeDahai},
					Actor:   0,
					Log:     "",
				},
				Pai:       "?",
				Tsumogiri: false,
			},
			wantErr: false,
		},
		{
			name: "with log",
			args: args{
				actor:     3,
				pai:       "5sr",
				tsumogiri: true,
				log:       "test",
			},
			want: &Dahai{
				Action: Action{
					Message: Message{TypeDahai},
					Actor:   3,
					Log:     "test",
				},
				Pai:       "5sr",
				Tsumogiri: true,
			},
			wantErr: false,
		},
		{
			name: "invalid actor min",
			args: args{
				actor:     -1,
				pai:       "?",
				tsumogiri: false,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: args{
				actor:     4,
				pai:       "?",
				tsumogiri: true,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty pai",
			args: args{
				actor:     0,
				pai:       "",
				tsumogiri: false,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid pai",
			args: args{
				actor:     1,
				pai:       "1z",
				tsumogiri: true,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDahai(tt.args.actor, tt.args.pai, tt.args.tsumogiri, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDahai() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDahai() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDahai_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		args    *Dahai
		want    string
		wantErr bool
	}{
		{
			name: "without log",
			args: &Dahai{
				Action: Action{
					Message: Message{TypeDahai},
					Actor:   0,
					Log:     "",
				},
				Pai:       "?",
				Tsumogiri: false,
			},
			want:    `{"type":"dahai","actor":0,"pai":"?","tsumogiri":false}`,
			wantErr: false,
		},
		{
			name: "with log",
			args: &Dahai{
				Action: Action{
					Message: Message{TypeDahai},
					Actor:   3,
					Log:     "test",
				},
				Pai:       "C",
				Tsumogiri: true,
			},
			want:    `{"type":"dahai","actor":3,"log":"test","pai":"C","tsumogiri":true}`,
			wantErr: false,
		},
		{
			name: "empty type",
			args: &Dahai{
				Action: Action{
					Message: Message{""},
					Actor:   0,
					Log:     "",
				},
				Pai: "9s",
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid type",
			args: &Dahai{
				Action: Action{
					Message: Message{TypeHello},
					Actor:   0,
					Log:     "",
				},
				Pai: "5p",
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid actor min",
			args: &Dahai{
				Action: Action{
					Message: Message{TypeDahai},
					Actor:   -1,
					Log:     "",
				},
				Pai: "1m",
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: &Dahai{
				Action: Action{
					Message: Message{TypeDahai},
					Actor:   4,
					Log:     "",
				},
				Pai: "E",
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "empty pai",
			args: &Dahai{
				Action: Action{
					Message: Message{TypeDahai},
					Actor:   1,
					Log:     "",
				},
				Pai: "",
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid pai",
			args: &Dahai{
				Action: Action{
					Message: Message{TypeDahai},
					Actor:   2,
					Log:     "",
				},
				Pai: "0m",
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

func TestDahai_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    Dahai
		wantErr bool
	}{
		{
			name: "without actor",
			args: `{"type":"dahai","pai":"?","tsumogiri":false}`,
			want: Dahai{
				Action: Action{
					Message: Message{TypeDahai},
					Actor:   0,
					Log:     "",
				},
				Pai:       "?",
				Tsumogiri: false,
			},
			wantErr: false,
		},
		{
			name: "without tsumogiri",
			args: `{"type":"dahai","actor":3,"pai":"?","log":"test"}`,
			want: Dahai{
				Action: Action{
					Message: Message{TypeDahai},
					Actor:   3,
					Log:     "test",
				},
				Pai:       "?",
				Tsumogiri: false,
			},
			wantErr: false,
		},
		{
			name: "without log",
			args: `{"type":"dahai","actor":0,"pai":"F","tsumogiri":true}`,
			want: Dahai{
				Action: Action{
					Message: Message{TypeDahai},
					Actor:   0,
					Log:     "",
				},
				Pai:       "F",
				Tsumogiri: true,
			},
			wantErr: false,
		},
		{
			name: "with log",
			args: `{"type":"dahai","actor":3,"pai":"?","tsumogiri":false,"log":"test"}`,
			want: Dahai{
				Action: Action{
					Message: Message{TypeDahai},
					Actor:   3,
					Log:     "test",
				},
				Pai:       "?",
				Tsumogiri: false,
			},
			wantErr: false,
		},
		{
			name: "empty type",
			args: `{"type":"","actor":0,"pai":"N","tsumogiri":false,"log":""}`,
			want: Dahai{
				Action: Action{
					Message: Message{""},
					Actor:   0,
					Log:     "",
				},
				Pai:       "N",
				Tsumogiri: false,
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			args: `{"type":"hello","actor":3,"pai":"N","tsumogiri":true}`,
			want: Dahai{
				Action: Action{
					Message: Message{TypeHello},
					Actor:   3,
					Log:     "",
				},
				Pai:       "N",
				Tsumogiri: true,
			},
			wantErr: true,
		},
		{
			name: "invalid actor min",
			args: `{"type":"dahai","actor":-1,"pai":"9m","tsumogiri":true}`,
			want: Dahai{
				Action: Action{
					Message: Message{TypeDahai},
					Actor:   -1,
					Log:     "",
				},
				Pai:       "9m",
				Tsumogiri: true,
			},
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: `{"type":"dahai","actor":4,"pai":"4s","tsumogiri":true}`,
			want: Dahai{
				Action: Action{
					Message: Message{TypeDahai},
					Actor:   4,
					Log:     "",
				},
				Pai:       "4s",
				Tsumogiri: true,
			},
			wantErr: true,
		},
		{
			name: "empty pai",
			args: `{"type":"dahai","actor":0,"pai":"","tsumogiri":false}`,
			want: Dahai{
				Action: Action{
					Message: Message{TypeDahai},
					Actor:   0,
					Log:     "",
				},
				Pai:       "",
				Tsumogiri: false,
			},
			wantErr: true,
		},
		{
			name: "invalid pai",
			args: `{"type":"dahai","actor":3,"pai":"4pr","tsumogiri":true}`,
			want: Dahai{
				Action: Action{
					Message: Message{TypeDahai},
					Actor:   3,
					Log:     "",
				},
				Pai:       "4pr",
				Tsumogiri: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Dahai
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

func TestDahai_ToEvent(t *testing.T) {
	tests := []struct {
		name    string
		args    *Dahai
		want    *inbound.Dahai
		wantErr bool
	}{
		{
			name: "valid",
			args: &Dahai{
				Action: Action{
					Message: Message{TypeDahai},
					Actor:   1,
					Log:     "",
				},
				Pai:       "P",
				Tsumogiri: true,
			},
			want: &inbound.Dahai{
				Actor:     1,
				Pai:       *mustPai("P"),
				Tsumogiri: true,
			},
			wantErr: false,
		},
		{
			name: "invalid",
			args: &Dahai{
				Action: Action{
					Message: Message{TypeDahai},
					Actor:   0,
					Log:     "",
				},
				Pai:       "?",
				Tsumogiri: false,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.ToEvent()
			if (err != nil) != tt.wantErr {
				t.Errorf("Dahai.ToEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dahai.ToEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDahaiFromEvent(t *testing.T) {
	valid, _ := outbound.NewDahai(1, *mustPai("E"), true, "test")
	invalid := *valid
	invalid.Actor = 4

	type args struct {
		ev *outbound.Dahai
	}
	tests := []struct {
		name    string
		args    args
		want    *Dahai
		wantErr bool
	}{
		{
			name: "valid",
			args: args{valid},
			want: &Dahai{
				Action: Action{
					Message: Message{Type: TypeDahai},
					Actor:   1,
					Log:     "test",
				},
				Pai:       "E",
				Tsumogiri: true,
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
			got, err := NewDahaiFromEvent(tt.args.ev)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDahaiFromEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDahaiFromEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
