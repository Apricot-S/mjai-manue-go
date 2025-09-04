package mjai

import (
	"encoding/json/v2"
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
)

func TestNewTsumo(t *testing.T) {
	type args struct {
		actor int
		pai   string
		log   string
	}
	tests := []struct {
		name    string
		args    args
		want    *Tsumo
		wantErr bool
	}{
		{
			name: "without log",
			args: args{
				actor: 0,
				pai:   "?",
				log:   "",
			},
			want: &Tsumo{
				Action: Action{
					Message: Message{TypeTsumo},
					Actor:   0,
					Log:     "",
				},
				Pai: "?",
			},
			wantErr: false,
		},
		{
			name: "with log",
			args: args{
				actor: 3,
				pai:   "5sr",
				log:   "test",
			},
			want: &Tsumo{
				Action: Action{
					Message: Message{TypeTsumo},
					Actor:   3,
					Log:     "test",
				},
				Pai: "5sr",
			},
			wantErr: false,
		},
		{
			name: "invalid actor min",
			args: args{
				actor: -1,
				pai:   "?",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: args{
				actor: 4,
				pai:   "?",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty pai",
			args: args{
				actor: 0,
				pai:   "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid pai",
			args: args{
				actor: 1,
				pai:   "1z",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTsumo(tt.args.actor, tt.args.pai, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTsumo() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTsumo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTsumo_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		args    *Tsumo
		want    string
		wantErr bool
	}{
		{
			name: "without log",
			args: &Tsumo{
				Action: Action{
					Message: Message{TypeTsumo},
					Actor:   0,
					Log:     "",
				},
				Pai: "?",
			},
			want:    `{"type":"tsumo","actor":0,"pai":"?"}`,
			wantErr: false,
		},
		{
			name: "with log",
			args: &Tsumo{
				Action: Action{
					Message: Message{TypeTsumo},
					Actor:   3,
					Log:     "test",
				},
				Pai: "C",
			},
			want:    `{"type":"tsumo","actor":3,"log":"test","pai":"C"}`,
			wantErr: false,
		},
		{
			name: "empty type",
			args: &Tsumo{
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
			args: &Tsumo{
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
			args: &Tsumo{
				Action: Action{
					Message: Message{TypeTsumo},
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
			args: &Tsumo{
				Action: Action{
					Message: Message{TypeTsumo},
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
			args: &Tsumo{
				Action: Action{
					Message: Message{TypeTsumo},
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
			args: &Tsumo{
				Action: Action{
					Message: Message{TypeTsumo},
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

func TestTsumo_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    Tsumo
		wantErr bool
	}{
		{
			name: "without actor",
			args: `{"type":"tsumo","pai":"?"}`,
			want: Tsumo{
				Action: Action{
					Message: Message{TypeTsumo},
					Actor:   0,
					Log:     "",
				},
				Pai: "?",
			},
			wantErr: false,
		},
		{
			name: "without log",
			args: `{"type":"tsumo","actor":0,"pai":"F"}`,
			want: Tsumo{
				Action: Action{
					Message: Message{TypeTsumo},
					Actor:   0,
					Log:     "",
				},
				Pai: "F",
			},
			wantErr: false,
		},
		{
			name: "with log",
			args: `{"type":"tsumo","actor":3,"pai":"?","log":"test"}`,
			want: Tsumo{
				Action: Action{
					Message: Message{TypeTsumo},
					Actor:   3,
					Log:     "test",
				},
				Pai: "?",
			},
			wantErr: false,
		},
		{
			name: "empty type",
			args: `{"type":"","actor":0,"pai":"N"}`,
			want: Tsumo{
				Action: Action{
					Message: Message{""},
					Actor:   0,
					Log:     "",
				},
				Pai: "N",
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			args: `{"type":"hello","actor":3,"pai":"N"}`,
			want: Tsumo{
				Action: Action{
					Message: Message{TypeHello},
					Actor:   3,
					Log:     "",
				},
				Pai: "N",
			},
			wantErr: true,
		},
		{
			name: "invalid actor min",
			args: `{"type":"tsumo","actor":-1,"pai":"9m"}`,
			want: Tsumo{
				Action: Action{
					Message: Message{TypeTsumo},
					Actor:   -1,
					Log:     "",
				},
				Pai: "9m",
			},
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: `{"type":"tsumo","actor":4,"pai":"4s"}`,
			want: Tsumo{
				Action: Action{
					Message: Message{TypeTsumo},
					Actor:   4,
					Log:     "",
				},
				Pai: "4s",
			},
			wantErr: true,
		},
		{
			name: "empty pai",
			args: `{"type":"tsumo","actor":0,"pai":""}`,
			want: Tsumo{
				Action: Action{
					Message: Message{TypeTsumo},
					Actor:   0,
					Log:     "",
				},
				Pai: "",
			},
			wantErr: true,
		},
		{
			name: "invalid pai",
			args: `{"type":"tsumo","actor":3,"pai":"4pr"}`,
			want: Tsumo{
				Action: Action{
					Message: Message{TypeTsumo},
					Actor:   3,
					Log:     "",
				},
				Pai: "4pr",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Tsumo
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

func TestTsumo_ToEvent(t *testing.T) {
	tests := []struct {
		name    string
		args    *Tsumo
		want    *inbound.Tsumo
		wantErr bool
	}{
		{
			name: "valid",
			args: &Tsumo{
				Action: Action{
					Message: Message{TypeTsumo},
					Actor:   1,
					Log:     "",
				},
				Pai: "?",
			},
			want: &inbound.Tsumo{
				Actor: 1,
				Pai:   *mustPai("?"),
			},
			wantErr: false,
		},
		{
			name: "invalid",
			args: &Tsumo{
				Action: Action{
					Message: Message{TypeTsumo},
					Actor:   -1,
					Log:     "",
				},
				Pai: "4s",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.ToEvent()
			if (err != nil) != tt.wantErr {
				t.Errorf("Tsumo.ToEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tsumo.ToEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
