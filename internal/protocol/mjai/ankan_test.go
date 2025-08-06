package mjai

import (
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/go-json-experiment/json"
)

func TestNewAnkan(t *testing.T) {
	type args struct {
		actor    int
		consumed [4]string
		log      string
	}
	tests := []struct {
		name    string
		args    args
		want    *Ankan
		wantErr bool
	}{
		{
			name: "without log",
			args: args{
				actor:    1,
				consumed: [4]string{"C", "C", "C", "C"},
				log:      "",
			},
			want: &Ankan{
				Action: Action{
					Message: Message{TypeAnkan},
					Actor:   1,
					Log:     "",
				},
				Consumed: [4]string{"C", "C", "C", "C"},
			},
			wantErr: false,
		},
		{
			name: "with log",
			args: args{
				actor:    3,
				consumed: [4]string{"5sr", "5s", "5s", "5s"},
				log:      "test",
			},
			want: &Ankan{
				Action: Action{
					Message: Message{TypeAnkan},
					Actor:   3,
					Log:     "test",
				},
				Consumed: [4]string{"5sr", "5s", "5s", "5s"},
			},
			wantErr: false,
		},
		{
			name: "invalid actor min",
			args: args{
				actor:    -1,
				consumed: [4]string{"C", "C", "C", "C"},
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: args{
				actor:    4,
				consumed: [4]string{"C", "C", "C", "C"},
				log:      "",
			},
			wantErr: true,
		},
		{
			name: "empty consumed",
			args: args{
				actor:    0,
				consumed: [4]string{"", "", "", ""},
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid consumed",
			args: args{
				actor:    0,
				consumed: [4]string{"7s", "7sr", "7s", "7s"},
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAnkan(tt.args.actor, tt.args.consumed, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAnkan() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAnkan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnkan_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		args    *Ankan
		want    string
		wantErr bool
	}{
		{
			name: "without log",
			args: &Ankan{
				Action: Action{
					Message: Message{TypeAnkan},
					Actor:   1,
					Log:     "",
				},
				Consumed: [4]string{"C", "C", "C", "C"},
			},
			want:    `{"type":"ankan","actor":1,"consumed":["C","C","C","C"]}`,
			wantErr: false,
		},
		{
			name: "with log",
			args: &Ankan{
				Action: Action{
					Message: Message{TypeAnkan},
					Actor:   1,
					Log:     "test",
				},
				Consumed: [4]string{"C", "C", "C", "C"},
			},
			want:    `{"type":"ankan","actor":1,"log":"test","consumed":["C","C","C","C"]}`,
			wantErr: false,
		},
		{
			name: "empty type",
			args: &Ankan{
				Action: Action{
					Message: Message{""},
					Actor:   1,
					Log:     "",
				},
				Consumed: [4]string{"C", "C", "C", "C"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid type",
			args: &Ankan{
				Action: Action{
					Message: Message{TypeNone},
					Actor:   1,
					Log:     "",
				},
				Consumed: [4]string{"C", "C", "C", "C"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid actor min",
			args: &Ankan{
				Action: Action{
					Message: Message{TypeAnkan},
					Actor:   -1,
					Log:     "",
				},
				Consumed: [4]string{"C", "C", "C", "C"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: &Ankan{
				Action: Action{
					Message: Message{TypeAnkan},
					Actor:   4,
					Log:     "",
				},
				Consumed: [4]string{"C", "C", "C", "C"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "empty consumed",
			args: &Ankan{
				Action: Action{
					Message: Message{TypeAnkan},
					Actor:   1,
					Log:     "",
				},
				Consumed: [4]string{"", "", "", ""},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid consumed",
			args: &Ankan{
				Action: Action{
					Message: Message{TypeAnkan},
					Actor:   1,
					Log:     "",
				},
				Consumed: [4]string{"7s", "7sr", "7s", "7s"},
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

func TestAnkan_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    Ankan
		wantErr bool
	}{
		{
			name: "without actor",
			args: `{"type":"ankan","consumed":["C","C","C","C"]}`,
			want: Ankan{
				Action: Action{
					Message: Message{TypeAnkan},
					Actor:   0,
					Log:     "",
				},
				Consumed: [4]string{"C", "C", "C", "C"},
			},
			wantErr: false,
		},
		{
			name: "without log",
			args: `{"type":"ankan","actor":1,"consumed":["C","C","C","C"]}`,
			want: Ankan{
				Action: Action{
					Message: Message{TypeAnkan},
					Actor:   1,
					Log:     "",
				},
				Consumed: [4]string{"C", "C", "C", "C"},
			},
			wantErr: false,
		},
		{
			name: "with log",
			args: `{"type":"ankan","actor":1,"consumed":["C","C","C","C"],"log":"test"}`,
			want: Ankan{
				Action: Action{
					Message: Message{TypeAnkan},
					Actor:   1,
					Log:     "test",
				},
				Consumed: [4]string{"C", "C", "C", "C"},
			},
			wantErr: false,
		},
		{
			name: "empty type",
			args: `{"type":"","actor":1,"consumed":["C","C","C","C"]}`,
			want: Ankan{
				Action: Action{
					Message: Message{""},
					Actor:   1,
					Log:     "",
				},
				Consumed: [4]string{"C", "C", "C", "C"},
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			args: `{"type":"none","actor":1,"consumed":["C","C","C","C"]}`,
			want: Ankan{
				Action: Action{
					Message: Message{TypeNone},
					Actor:   1,
					Log:     "",
				},
				Consumed: [4]string{"C", "C", "C", "C"},
			},
			wantErr: true,
		},
		{
			name: "invalid actor min",
			args: `{"type":"ankan","actor":-1,"consumed":["C","C","C","C"]}`,
			want: Ankan{
				Action: Action{
					Message: Message{TypeAnkan},
					Actor:   -1,
					Log:     "",
				},
				Consumed: [4]string{"C", "C", "C", "C"},
			},
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: `{"type":"ankan","actor":4,"consumed":["C","C","C","C"]}`,
			want: Ankan{
				Action: Action{
					Message: Message{TypeAnkan},
					Actor:   4,
					Log:     "",
				},
				Consumed: [4]string{"C", "C", "C", "C"},
			},
			wantErr: true,
		},
		{
			name: "empty consumed pai",
			args: `{"type":"ankan","actor":1,"consumed":["","","",""]}`,
			want: Ankan{
				Action: Action{
					Message: Message{TypeAnkan},
					Actor:   1,
					Log:     "",
				},
				Consumed: [4]string{"", "", "", ""},
			},
			wantErr: true,
		},
		{
			name: "invalid consumed pai",
			args: `{"type":"ankan","actor":1,"consumed":["6sr","6s","6s","6s"]}`,
			want: Ankan{
				Action: Action{
					Message: Message{TypeAnkan},
					Actor:   1,
					Log:     "",
				},
				Consumed: [4]string{"6sr", "6s", "6s", "6s"},
			},
			wantErr: true,
		},
		{
			name:    "invalid consumed len 3",
			args:    `{"type":"ankan","actor":1,"consumed":["5sr","5s","5s"]}`,
			want:    Ankan{},
			wantErr: true,
		},
		{
			name:    "invalid consumed len 5",
			args:    `{"type":"ankan","actor":1,"consumed":["6s","6s","6s","6s","6s"]}`,
			want:    Ankan{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Ankan
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

func TestAnkan_ToEvent(t *testing.T) {
	tests := []struct {
		name    string
		args    *Ankan
		want    *inbound.Ankan
		wantErr bool
	}{
		{
			name: "valid",
			args: &Ankan{
				Action: Action{
					Message: Message{TypeAnkan},
					Actor:   1,
					Log:     "",
				},
				Consumed: [4]string{"C", "C", "C", "C"},
			},
			want: &inbound.Ankan{
				Actor:    1,
				Consumed: [4]base.Pai(mustPais([]string{"C", "C", "C", "C"})),
			},
			wantErr: false,
		},
		{
			name: "invalid",
			args: &Ankan{
				Action: Action{
					Message: Message{TypeAnkan},
					Actor:   1,
					Log:     "",
				},
				Consumed: [4]string{"C", "C", "C", "P"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.ToEvent()
			if (err != nil) != tt.wantErr {
				t.Errorf("Ankan.ToEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ankan.ToEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
