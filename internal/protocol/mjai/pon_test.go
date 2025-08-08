package mjai

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/outbound"
	"github.com/go-json-experiment/json"
)

func TestNewPon(t *testing.T) {
	type args struct {
		actor    int
		target   int
		pai      string
		consumed [2]string
		log      string
	}
	tests := []struct {
		name    string
		args    args
		want    *Pon
		wantErr bool
	}{
		{
			name: "without log",
			args: args{
				actor:    1,
				target:   0,
				pai:      "C",
				consumed: [2]string{"C", "C"},
				log:      "",
			},
			want: &Pon{
				Action: Action{
					Message: Message{TypePon},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [2]string{"C", "C"},
			},
			wantErr: false,
		},
		{
			name: "with log",
			args: args{
				actor:    3,
				target:   2,
				pai:      "5sr",
				consumed: [2]string{"5s", "5s"},
				log:      "test",
			},
			want: &Pon{
				Action: Action{
					Message: Message{TypePon},
					Actor:   3,
					Log:     "test",
				},
				Target:   2,
				Pai:      "5sr",
				Consumed: [2]string{"5s", "5s"},
			},
			wantErr: false,
		},
		{
			name: "invalid actor min",
			args: args{
				actor:    -1,
				target:   0,
				pai:      "C",
				consumed: [2]string{"C", "C"},
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: args{
				actor:    4,
				target:   3,
				pai:      "C",
				consumed: [2]string{"C", "C"},
				log:      "",
			},
			wantErr: true,
		},
		{
			name: "invalid target min",
			args: args{
				actor:    0,
				target:   -1,
				pai:      "C",
				consumed: [2]string{"C", "C"},
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid target max",
			args: args{
				actor:    0,
				target:   4,
				pai:      "C",
				consumed: [2]string{"C", "C"},
				log:      "",
			},
			wantErr: true,
		},
		{
			name: "empty pai",
			args: args{
				actor:    0,
				target:   3,
				pai:      "",
				consumed: [2]string{"C", "C"},
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid pai",
			args: args{
				actor:    0,
				target:   3,
				pai:      "6sr",
				consumed: [2]string{"6s", "6s"},
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty consumed",
			args: args{
				actor:    0,
				target:   3,
				pai:      "C",
				consumed: [2]string{"", ""},
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid consumed",
			args: args{
				actor:    0,
				target:   3,
				pai:      "7s",
				consumed: [2]string{"7s", "7sr"},
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPon(tt.args.actor, tt.args.target, tt.args.pai, tt.args.consumed, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPon() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPon() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPon_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		args    *Pon
		want    string
		wantErr bool
	}{
		{
			name: "without log",
			args: &Pon{
				Action: Action{
					Message: Message{TypePon},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [2]string{"C", "C"},
			},
			want:    `{"type":"pon","actor":1,"target":0,"pai":"C","consumed":["C","C"]}`,
			wantErr: false,
		},
		{
			name: "with log",
			args: &Pon{
				Action: Action{
					Message: Message{TypePon},
					Actor:   1,
					Log:     "test",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [2]string{"C", "C"},
			},
			want:    `{"type":"pon","actor":1,"log":"test","target":0,"pai":"C","consumed":["C","C"]}`,
			wantErr: false,
		},
		{
			name: "empty type",
			args: &Pon{
				Action: Action{
					Message: Message{""},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [2]string{"C", "C"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid type",
			args: &Pon{
				Action: Action{
					Message: Message{TypeNone},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [2]string{"C", "C"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid actor min",
			args: &Pon{
				Action: Action{
					Message: Message{TypePon},
					Actor:   -1,
					Log:     "",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [2]string{"C", "C"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: &Pon{
				Action: Action{
					Message: Message{TypePon},
					Actor:   4,
					Log:     "",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [2]string{"C", "C"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid target min",
			args: &Pon{
				Action: Action{
					Message: Message{TypePon},
					Actor:   0,
					Log:     "",
				},
				Target:   -1,
				Pai:      "C",
				Consumed: [2]string{"C", "C"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid target max",
			args: &Pon{
				Action: Action{
					Message: Message{TypePon},
					Actor:   0,
					Log:     "",
				},
				Target:   4,
				Pai:      "C",
				Consumed: [2]string{"C", "C"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "empty pai",
			args: &Pon{
				Action: Action{
					Message: Message{TypePon},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "",
				Consumed: [2]string{"C", "C"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid pai",
			args: &Pon{
				Action: Action{
					Message: Message{TypePon},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "6sr",
				Consumed: [2]string{"6s", "6s"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "empty consumed",
			args: &Pon{
				Action: Action{
					Message: Message{TypePon},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [2]string{"", ""},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid consumed",
			args: &Pon{
				Action: Action{
					Message: Message{TypePon},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "7s",
				Consumed: [2]string{"7s", "7sr"},
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

func TestPon_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    Pon
		wantErr bool
	}{
		{
			name: "without actor",
			args: `{"type":"pon","target":3,"pai":"C","consumed":["C","C"]}`,
			want: Pon{
				Action: Action{
					Message: Message{TypePon},
					Actor:   0,
					Log:     "",
				},
				Target:   3,
				Pai:      "C",
				Consumed: [2]string{"C", "C"},
			},
			wantErr: false,
		},
		{
			name: "without target",
			args: `{"type":"pon","actor":1,"pai":"C","consumed":["C","C"]}`,
			want: Pon{
				Action: Action{
					Message: Message{TypePon},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [2]string{"C", "C"},
			},
			wantErr: false,
		},
		{
			name: "without log",
			args: `{"type":"pon","actor":1,"target":0,"pai":"C","consumed":["C","C"]}`,
			want: Pon{
				Action: Action{
					Message: Message{TypePon},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [2]string{"C", "C"},
			},
			wantErr: false,
		},
		{
			name: "with log",
			args: `{"type":"pon","actor":1,"target":0,"pai":"C","consumed":["C","C"],"log":"test"}`,
			want: Pon{
				Action: Action{
					Message: Message{TypePon},
					Actor:   1,
					Log:     "test",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [2]string{"C", "C"},
			},
			wantErr: false,
		},
		{
			name: "empty type",
			args: `{"type":"","actor":1,"target":0,"pai":"C","consumed":["C","C"]}`,
			want: Pon{
				Action: Action{
					Message: Message{""},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [2]string{"C", "C"},
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			args: `{"type":"none","actor":1,"target":0,"pai":"C","consumed":["C","C"]}`,
			want: Pon{
				Action: Action{
					Message: Message{TypeNone},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [2]string{"C", "C"},
			},
			wantErr: true,
		},
		{
			name: "invalid actor min",
			args: `{"type":"pon","actor":-1,"target":0,"pai":"C","consumed":["C","C"]}`,
			want: Pon{
				Action: Action{
					Message: Message{TypePon},
					Actor:   -1,
					Log:     "",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [2]string{"C", "C"},
			},
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: `{"type":"pon","actor":4,"target":0,"pai":"C","consumed":["C","C"]}`,
			want: Pon{
				Action: Action{
					Message: Message{TypePon},
					Actor:   4,
					Log:     "",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [2]string{"C", "C"},
			},
			wantErr: true,
		},
		{
			name: "invalid target min",
			args: `{"type":"pon","actor":0,"target":-1,"pai":"C","consumed":["C","C"]}`,
			want: Pon{
				Action: Action{
					Message: Message{TypePon},
					Actor:   0,
					Log:     "",
				},
				Target:   -1,
				Pai:      "C",
				Consumed: [2]string{"C", "C"},
			},
			wantErr: true,
		},
		{
			name: "invalid target max",
			args: `{"type":"pon","actor":0,"target":4,"pai":"C","consumed":["C","C"]}`,
			want: Pon{
				Action: Action{
					Message: Message{TypePon},
					Actor:   0,
					Log:     "",
				},
				Target:   4,
				Pai:      "C",
				Consumed: [2]string{"C", "C"},
			},
			wantErr: true,
		},
		{
			name: "empty pai",
			args: `{"type":"pon","actor":1,"target":0,"pai":"","consumed":["C","C"]}`,
			want: Pon{
				Action: Action{
					Message: Message{TypePon},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "",
				Consumed: [2]string{"C", "C"},
			},
			wantErr: true,
		},
		{
			name: "invalid pai",
			args: `{"type":"pon","actor":1,"target":0,"pai":"6sr","consumed":["6s","6s"]}`,
			want: Pon{
				Action: Action{
					Message: Message{TypePon},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "6sr",
				Consumed: [2]string{"6s", "6s"},
			},
			wantErr: true,
		},
		{
			name: "empty consumed pai",
			args: `{"type":"pon","actor":1,"target":0,"pai":"C","consumed":["",""]}`,
			want: Pon{
				Action: Action{
					Message: Message{TypePon},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [2]string{"", ""},
			},
			wantErr: true,
		},
		{
			name: "invalid consumed pai",
			args: `{"type":"pon","actor":1,"target":0,"pai":"6s","consumed":["6sr","6s"]}`,
			want: Pon{
				Action: Action{
					Message: Message{TypePon},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "6s",
				Consumed: [2]string{"6sr", "6s"},
			},
			wantErr: true,
		},
		{
			name:    "invalid consumed len 1",
			args:    `{"type":"pon","actor":1,"target":0,"pai":"5s","consumed":["5sr"]}`,
			want:    Pon{},
			wantErr: true,
		},
		{
			name:    "invalid consumed len 3",
			args:    `{"type":"pon","actor":1,"target":0,"pai":"6s","consumed":["6s","6s","6s"]}`,
			want:    Pon{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Pon
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

func TestPon_ToEvent(t *testing.T) {
	tests := []struct {
		name    string
		args    *Pon
		want    *inbound.Pon
		wantErr bool
	}{
		{
			name: "valid",
			args: &Pon{
				Action: Action{
					Message: Message{TypePon},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "6s",
				Consumed: [2]string{"6s", "6s"},
			},
			want: &inbound.Pon{
				Actor:    1,
				Target:   0,
				Taken:    *mustPai("6s"),
				Consumed: [2]base.Pai(mustPais("6s", "6s")),
			},
			wantErr: false,
		},
		{
			name: "invalid",
			args: &Pon{
				Action: Action{
					Message: Message{TypePon},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "4s",
				Consumed: [2]string{"6s", "6s"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.ToEvent()
			if (err != nil) != tt.wantErr {
				t.Errorf("Pon.ToEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pon.ToEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPonFromEvent(t *testing.T) {
	valid, _ := outbound.NewPon(1, 0, *mustPai("1m"), [2]base.Pai(mustPais("1m", "1m")), "test")
	invalid := *valid
	invalid.Actor = -1

	type args struct {
		ev *outbound.Pon
	}
	tests := []struct {
		name    string
		args    args
		want    *Pon
		wantErr bool
	}{
		{
			name: "valid",
			args: args{valid},
			want: &Pon{
				Action: Action{
					Message: Message{TypePon},
					Actor:   1,
					Log:     "test",
				},
				Target:   0,
				Pai:      "1m",
				Consumed: [2]string{"1m", "1m"},
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
			got, err := NewPonFromEvent(tt.args.ev)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPonFromEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPonFromEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
