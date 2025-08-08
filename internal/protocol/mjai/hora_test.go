package mjai

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/outbound"
	"github.com/go-json-experiment/json"
)

func TestNewHora(t *testing.T) {
	type args struct {
		actor      int
		target     int
		pai        string
		horaPoints int
		scores     []int
		log        string
	}
	tests := []struct {
		name    string
		args    args
		want    *Hora
		wantErr bool
	}{
		{
			name: "without hora_points",
			args: args{
				actor:      1,
				target:     0,
				pai:        "6s",
				horaPoints: 0,
				scores:     nil,
				log:        "",
			},
			want: &Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "",
				},
				Target:     0,
				Pai:        "6s",
				HoraPoints: 0,
				Scores:     nil,
			},
			wantErr: false,
		},
		{
			name: "without scores",
			args: args{
				actor:      1,
				target:     0,
				pai:        "6s",
				horaPoints: 2600,
				scores:     nil,
				log:        "",
			},
			want: &Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "",
				},
				Target:     0,
				Pai:        "6s",
				HoraPoints: 2600,
				Scores:     nil,
			},
			wantErr: false,
		},
		{
			name: "without log",
			args: args{
				actor:      1,
				target:     0,
				pai:        "6s",
				horaPoints: 2600,
				scores:     []int{27500, 22300, 24300, 25900},
				log:        "",
			},
			want: &Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "",
				},
				Target:     0,
				Pai:        "6s",
				HoraPoints: 2600,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			wantErr: false,
		},
		{
			name: "with log",
			args: args{
				actor:      3,
				target:     2,
				pai:        "5sr",
				horaPoints: 1000,
				scores:     []int{27500, 22300, 24300, 25900},
				log:        "test",
			},
			want: &Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   3,
					Log:     "test",
				},
				Target:     2,
				Pai:        "5sr",
				HoraPoints: 1000,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			wantErr: false,
		},
		{
			name: "invalid actor min",
			args: args{
				actor:      -1,
				target:     0,
				pai:        "6s",
				horaPoints: 2600,
				scores:     []int{27500, 22300, 24300, 25900},
				log:        "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: args{
				actor:      4,
				target:     3,
				pai:        "6s",
				horaPoints: 2600,
				scores:     []int{27500, 22300, 24300, 25900},
				log:        "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid target min",
			args: args{
				actor:      0,
				target:     -1,
				pai:        "6s",
				horaPoints: 2600,
				scores:     []int{27500, 22300, 24300, 25900},
				log:        "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid target max",
			args: args{
				actor:      0,
				target:     4,
				pai:        "6s",
				horaPoints: 2600,
				scores:     []int{27500, 22300, 24300, 25900},
				log:        "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty pai",
			args: args{
				actor:  0,
				target: 0,
				pai:    "",
				scores: []int{27500, 22300, 24300, 25900},
				log:    "",
			},
			want: &Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   0,
					Log:     "",
				},
				Target:     0,
				Pai:        "",
				HoraPoints: 0,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			wantErr: false,
		},
		{
			name: "invalid pai",
			args: args{
				actor:  0,
				target: 1,
				pai:    "6sr",
				scores: []int{27500, 22300, 24300, 25900},
				log:    "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid hora_points",
			args: args{
				actor:      0,
				target:     2,
				pai:        "6s",
				horaPoints: -1,
				scores:     []int{27500, 22300, 24300, 25900},
				log:        "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid scores empty",
			args: args{
				actor:  0,
				target: 3,
				pai:    "6s",
				scores: []int{},
				log:    "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid scores 3",
			args: args{
				actor:  0,
				target: 3,
				pai:    "6s",
				scores: []int{27500, 22300, 24300},
				log:    "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid scores 5",
			args: args{
				actor:  0,
				target: 3,
				pai:    "6s",
				scores: []int{27500, 22300, 24300, 25900, 25900},
				log:    "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewHora(tt.args.actor, tt.args.target, tt.args.pai, tt.args.horaPoints, tt.args.scores, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHora() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHora() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHora_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		args    *Hora
		want    string
		wantErr bool
	}{
		{
			name: "without hora_points",
			args: &Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "",
				},
				Target:     0,
				Pai:        "6s",
				HoraPoints: 0,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			want:    `{"type":"hora","actor":1,"target":0,"pai":"6s","scores":[27500,22300,24300,25900]}`,
			wantErr: false,
		},
		{
			name: "without scores",
			args: &Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "",
				},
				Target:     0,
				Pai:        "6s",
				HoraPoints: 0,
				Scores:     nil,
			},
			want:    `{"type":"hora","actor":1,"target":0,"pai":"6s"}`,
			wantErr: false,
		},
		{
			name: "without log",
			args: &Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "",
				},
				Target:     0,
				Pai:        "6s",
				HoraPoints: 2600,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			want:    `{"type":"hora","actor":1,"target":0,"pai":"6s","hora_points":2600,"scores":[27500,22300,24300,25900]}`,
			wantErr: false,
		},
		{
			name: "with log",
			args: &Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "test",
				},
				Target:     0,
				Pai:        "6s",
				HoraPoints: 2600,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			want:    `{"type":"hora","actor":1,"log":"test","target":0,"pai":"6s","hora_points":2600,"scores":[27500,22300,24300,25900]}`,
			wantErr: false,
		},
		{
			name: "empty type",
			args: &Hora{
				Action: Action{
					Message: Message{""},
					Actor:   1,
					Log:     "",
				},
				Target:     0,
				Pai:        "6s",
				HoraPoints: 2600,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid type",
			args: &Hora{
				Action: Action{
					Message: Message{TypeNone},
					Actor:   1,
					Log:     "",
				},
				Target:     0,
				Pai:        "6s",
				HoraPoints: 2600,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid actor min",
			args: &Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   -1,
					Log:     "",
				},
				Target:     0,
				Pai:        "6s",
				HoraPoints: 2600,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: &Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   4,
					Log:     "",
				},
				Target:     0,
				Pai:        "6s",
				HoraPoints: 2600,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid target min",
			args: &Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   0,
					Log:     "",
				},
				Target:     -1,
				Pai:        "6s",
				HoraPoints: 2600,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid target max",
			args: &Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   0,
					Log:     "",
				},
				Target:     4,
				Pai:        "6s",
				HoraPoints: 2600,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "empty pai",
			args: &Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "",
				},
				Target:     0,
				Pai:        "",
				HoraPoints: 2600,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			want:    `{"type":"hora","actor":1,"target":0,"hora_points":2600,"scores":[27500,22300,24300,25900]}`,
			wantErr: false,
		},
		{
			name: "invalid pai",
			args: &Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "",
				},
				Target:     0,
				Pai:        "6sr",
				HoraPoints: 2600,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid hora_points",
			args: &Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "",
				},
				Target:     0,
				Pai:        "6s",
				HoraPoints: -1,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid scores empty",
			args: &Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "",
				},
				Target:     0,
				Pai:        "6s",
				HoraPoints: 2600,
				Scores:     []int{},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid scores 3",
			args: &Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "",
				},
				Target:     0,
				Pai:        "6s",
				HoraPoints: 2600,
				Scores:     []int{27500, 22300, 24300},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid scores 5",
			args: &Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "",
				},
				Target:     0,
				Pai:        "6s",
				HoraPoints: 2600,
				Scores:     []int{27500, 22300, 24300, 22300, 24300},
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

func TestHora_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    Hora
		wantErr bool
	}{
		{
			name: "without actor",
			args: `{"type":"hora","target":3,"pai":"6s","scores":[27500,22300,24300,25900]}`,
			want: Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   0,
					Log:     "",
				},
				Target:     3,
				Pai:        "6s",
				HoraPoints: 0,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			wantErr: false,
		},
		{
			name: "without target",
			args: `{"type":"hora","actor":1,"pai":"6s","scores":[27500,22300,24300,25900]}`,
			want: Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "",
				},
				Target:     0,
				Pai:        "6s",
				HoraPoints: 0,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			wantErr: false,
		},
		{
			name: "without hora_points",
			args: `{"type":"hora","actor":1,"target":3,"pai":"6s","scores":[27500,22300,24300,25900]}`,
			want: Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "",
				},
				Target:     3,
				Pai:        "6s",
				HoraPoints: 0,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			wantErr: false,
		},
		{
			name: "without scores",
			args: `{"type":"hora","actor":1,"target":3,"pai":"6s"}`,
			want: Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "",
				},
				Target:     3,
				Pai:        "6s",
				HoraPoints: 0,
				Scores:     nil,
			},
			wantErr: false,
		},
		{
			name: "without log",
			args: `{"type":"hora","actor":1,"target":2,"pai":"6s","hora_points":2600,"scores":[27500,22300,24300,25900]}`,
			want: Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "",
				},
				Target:     2,
				Pai:        "6s",
				HoraPoints: 2600,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			wantErr: false,
		},
		{
			name: "with log",
			args: `{"type":"hora","actor":1,"target":0,"pai":"6s","scores":[27500,22300,24300,25900],"log":"test"}`,
			want: Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "test",
				},
				Target:     0,
				Pai:        "6s",
				HoraPoints: 0,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			wantErr: false,
		},
		{
			name: "empty type",
			args: `{"type":"","actor":1,"target":0,"pai":"6s","scores":[27500,22300,24300,25900]}`,
			want: Hora{
				Action: Action{
					Message: Message{""},
					Actor:   1,
					Log:     "",
				},
				Target:     0,
				Pai:        "6s",
				HoraPoints: 0,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			args: `{"type":"none","actor":1,"target":0,"pai":"6s","scores":[27500,22300,24300,25900]}`,
			want: Hora{
				Action: Action{
					Message: Message{TypeNone},
					Actor:   1,
					Log:     "",
				},
				Target:     0,
				Pai:        "6s",
				HoraPoints: 0,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			wantErr: true,
		},
		{
			name: "invalid actor min",
			args: `{"type":"hora","actor":-1,"target":0,"pai":"6s","scores":[27500,22300,24300,25900]}`,
			want: Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   -1,
					Log:     "",
				},
				Target:     0,
				Pai:        "6s",
				HoraPoints: 0,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: `{"type":"hora","actor":4,"target":0,"pai":"6s","scores":[27500,22300,24300,25900]}`,
			want: Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   4,
					Log:     "",
				},
				Target:     0,
				Pai:        "6s",
				HoraPoints: 0,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			wantErr: true,
		},
		{
			name: "invalid target min",
			args: `{"type":"hora","actor":0,"target":-1,"pai":"6s","scores":[27500,22300,24300,25900]}`,
			want: Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   0,
					Log:     "",
				},
				Target:     -1,
				Pai:        "6s",
				HoraPoints: 0,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			wantErr: true,
		},
		{
			name: "invalid target max",
			args: `{"type":"hora","actor":0,"target":4,"pai":"6s","scores":[27500,22300,24300,25900]}`,
			want: Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   0,
					Log:     "",
				},
				Target:     4,
				Pai:        "6s",
				HoraPoints: 0,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			wantErr: true,
		},
		{
			name: "empty pai",
			args: `{"type":"hora","actor":1,"target":0,"pai":"","scores":[27500,22300,24300,25900]}`,
			want: Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "",
				},
				Target:     0,
				Pai:        "",
				HoraPoints: 0,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			wantErr: false,
		},
		{
			name: "invalid pai",
			args: `{"type":"hora","actor":1,"target":0,"pai":"6sr","scores":[27500,22300,24300,25900]}`,
			want: Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "",
				},
				Target:     0,
				Pai:        "6sr",
				HoraPoints: 0,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			wantErr: true,
		},
		{
			name: "invalid hora_points",
			args: `{"type":"hora","actor":1,"target":2,"pai":"6s","hora_points":-1,"scores":[27500,22300,24300,25900]}`,
			want: Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "",
				},
				Target:     2,
				Pai:        "6s",
				HoraPoints: -1,
				Scores:     []int{27500, 22300, 24300, 25900},
			},
			wantErr: true,
		},
		{
			name: "invalid scores empty",
			args: `{"type":"hora","actor":1,"target":0,"pai":"6s","scores":[]}`,
			want: Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "",
				},
				Target:     0,
				Pai:        "6s",
				HoraPoints: 0,
				Scores:     []int{},
			},
			wantErr: true,
		},
		{
			name: "invalid scores 3",
			args: `{"type":"hora","actor":1,"target":0,"pai":"6s","scores":[27500,22300,24300]}`,
			want: Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "",
				},
				Target:     0,
				Pai:        "6s",
				HoraPoints: 0,
				Scores:     []int{27500, 22300, 24300},
			},
			wantErr: true,
		},
		{
			name: "invalid scores 5",
			args: `{"type":"hora","actor":1,"target":0,"pai":"6s","scores":[27500,22300,24300,25900,25900]}`,
			want: Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "",
				},
				Target:     0,
				Pai:        "6s",
				HoraPoints: 0,
				Scores:     []int{27500, 22300, 24300, 25900, 25900},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Hora
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

func TestHora_ToEvent(t *testing.T) {
	horaPoints := new(int)
	*horaPoints = 2600

	tests := []struct {
		name    string
		args    *Hora
		want    *inbound.Hora
		wantErr bool
	}{
		{
			name: "without hora_points and scores",
			args: &Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "",
				},
				Target:     1,
				Pai:        "2p",
				HoraPoints: 0,
				Scores:     nil,
			},
			want: &inbound.Hora{
				Actor:      1,
				Target:     1,
				Pai:        mustPai("2p"),
				HoraPoints: nil,
				Scores:     nil,
			},
			wantErr: false,
		},
		{
			name: "with hora_points",
			args: &Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "",
				},
				Target:     1,
				Pai:        "2p",
				HoraPoints: 2600,
				Scores:     nil,
			},
			want: &inbound.Hora{
				Actor:      1,
				Target:     1,
				Pai:        mustPai("2p"),
				HoraPoints: horaPoints,
				Scores:     nil,
			},
			wantErr: false,
		},
		{
			name: "with scores",
			args: &Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "",
				},
				Target:     1,
				Pai:        "2p",
				HoraPoints: 0,
				Scores:     []int{26000, 24000, 23000, 24000},
			},
			want: &inbound.Hora{
				Actor:      1,
				Target:     1,
				Pai:        mustPai("2p"),
				HoraPoints: nil,
				Scores:     &[4]int{26000, 24000, 23000, 24000},
			},
			wantErr: false,
		},
		{
			name: "invalid",
			args: &Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "",
				},
				Target:     1,
				Pai:        "2p",
				HoraPoints: -1,
				Scores:     nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.ToEvent()
			if (err != nil) != tt.wantErr {
				t.Errorf("Hora.ToEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Hora.ToEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewHoraFromEvent(t *testing.T) {
	valid, _ := outbound.NewHora(1, 0, *mustPai("1m"), "test")
	invalid := *valid
	invalid.Actor = 4

	type args struct {
		ev *outbound.Hora
	}
	tests := []struct {
		name    string
		args    args
		want    *Hora
		wantErr bool
	}{
		{
			name: "valid",
			args: args{valid},
			want: &Hora{
				Action: Action{
					Message: Message{TypeHora},
					Actor:   1,
					Log:     "test",
				},
				Target: 0,
				Pai:    "1m",
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
			got, err := NewHoraFromEvent(tt.args.ev)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHoraFromEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHoraFromEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
