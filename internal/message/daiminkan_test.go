package message

import (
	"reflect"
	"testing"

	"github.com/go-json-experiment/json"
)

func TestNewDaiminkan(t *testing.T) {
	type args struct {
		actor    int
		target   int
		pai      string
		consumed [3]string
		log      string
	}
	tests := []struct {
		name    string
		args    args
		want    *Daiminkan
		wantErr bool
	}{
		{
			name: "without log",
			args: args{
				actor:    1,
				target:   0,
				pai:      "C",
				consumed: [3]string{"C", "C", "C"},
				log:      "",
			},
			want: &Daiminkan{
				Action: Action{
					Message: Message{TypeDaiminkan},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			wantErr: false,
		},
		{
			name: "with log",
			args: args{
				actor:    3,
				target:   2,
				pai:      "5sr",
				consumed: [3]string{"5s", "5s", "5s"},
				log:      "test",
			},
			want: &Daiminkan{
				Action: Action{
					Message: Message{TypeDaiminkan},
					Actor:   3,
					Log:     "test",
				},
				Target:   2,
				Pai:      "5sr",
				Consumed: [3]string{"5s", "5s", "5s"},
			},
			wantErr: false,
		},
		{
			name: "invalid actor min",
			args: args{
				actor:    -1,
				target:   0,
				pai:      "C",
				consumed: [3]string{"C", "C", "C"},
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
				consumed: [3]string{"C", "C", "C"},
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
				consumed: [3]string{"C", "C", "C"},
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
				consumed: [3]string{"C", "C", "C"},
				log:      "",
			},
			wantErr: true,
		},
		{
			name: "empty pai",
			args: args{
				actor:    0,
				target:   4,
				pai:      "",
				consumed: [3]string{"C", "C", "C"},
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid pai",
			args: args{
				actor:    0,
				target:   4,
				pai:      "6sr",
				consumed: [3]string{"6s", "6s", "6s"},
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty consumed",
			args: args{
				actor:    0,
				target:   4,
				pai:      "C",
				consumed: [3]string{"", "", ""},
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid consumed",
			args: args{
				actor:    0,
				target:   4,
				pai:      "7s",
				consumed: [3]string{"7s", "7sr", "7s"},
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDaiminkan(tt.args.actor, tt.args.target, tt.args.pai, tt.args.consumed, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDaiminkan() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDaiminkan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDaiminkan_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		args    *Daiminkan
		want    string
		wantErr bool
	}{
		{
			name: "without log",
			args: &Daiminkan{
				Action: Action{
					Message: Message{TypeDaiminkan},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			want:    `{"type":"daiminkan","actor":1,"target":0,"pai":"C","consumed":["C","C","C"]}`,
			wantErr: false,
		},
		{
			name: "with log",
			args: &Daiminkan{
				Action: Action{
					Message: Message{TypeDaiminkan},
					Actor:   1,
					Log:     "test",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			want:    `{"type":"daiminkan","actor":1,"log":"test","target":0,"pai":"C","consumed":["C","C","C"]}`,
			wantErr: false,
		},
		{
			name: "empty type",
			args: &Daiminkan{
				Action: Action{
					Message: Message{""},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid type",
			args: &Daiminkan{
				Action: Action{
					Message: Message{TypeNone},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid actor min",
			args: &Daiminkan{
				Action: Action{
					Message: Message{TypeDaiminkan},
					Actor:   -1,
					Log:     "",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: &Daiminkan{
				Action: Action{
					Message: Message{TypeDaiminkan},
					Actor:   4,
					Log:     "",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid target min",
			args: &Daiminkan{
				Action: Action{
					Message: Message{TypeDaiminkan},
					Actor:   0,
					Log:     "",
				},
				Target:   -1,
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid target max",
			args: &Daiminkan{
				Action: Action{
					Message: Message{TypeDaiminkan},
					Actor:   0,
					Log:     "",
				},
				Target:   4,
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "empty pai",
			args: &Daiminkan{
				Action: Action{
					Message: Message{TypeDaiminkan},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "",
				Consumed: [3]string{"C", "C", "C"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid pai",
			args: &Daiminkan{
				Action: Action{
					Message: Message{TypeDaiminkan},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "6sr",
				Consumed: [3]string{"6s", "6s", "6s"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "empty consumed",
			args: &Daiminkan{
				Action: Action{
					Message: Message{TypeDaiminkan},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [3]string{"", "", ""},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid consumed",
			args: &Daiminkan{
				Action: Action{
					Message: Message{TypeDaiminkan},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "7s",
				Consumed: [3]string{"7s", "7sr", "7s"},
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

func TestDaiminkan_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    Daiminkan
		wantErr bool
	}{
		{
			name: "without actor",
			args: `{"type":"daiminkan","target":3,"pai":"C","consumed":["C","C","C"]}`,
			want: Daiminkan{
				Action: Action{
					Message: Message{TypeDaiminkan},
					Actor:   0,
					Log:     "",
				},
				Target:   3,
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			wantErr: false,
		},
		{
			name: "without target",
			args: `{"type":"daiminkan","actor":1,"pai":"C","consumed":["C","C","C"]}`,
			want: Daiminkan{
				Action: Action{
					Message: Message{TypeDaiminkan},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			wantErr: false,
		},
		{
			name: "without log",
			args: `{"type":"daiminkan","actor":1,"target":0,"pai":"C","consumed":["C","C","C"]}`,
			want: Daiminkan{
				Action: Action{
					Message: Message{TypeDaiminkan},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			wantErr: false,
		},
		{
			name: "with log",
			args: `{"type":"daiminkan","actor":1,"target":0,"pai":"C","consumed":["C","C","C"],"log":"test"}`,
			want: Daiminkan{
				Action: Action{
					Message: Message{TypeDaiminkan},
					Actor:   1,
					Log:     "test",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			wantErr: false,
		},
		{
			name: "empty type",
			args: `{"type":"","actor":1,"target":0,"pai":"C","consumed":["C","C","C"]}`,
			want: Daiminkan{
				Action: Action{
					Message: Message{""},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			args: `{"type":"none","actor":1,"target":0,"pai":"C","consumed":["C","C","C"]}`,
			want: Daiminkan{
				Action: Action{
					Message: Message{TypeNone},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			wantErr: true,
		},
		{
			name: "invalid actor min",
			args: `{"type":"daiminkan","actor":-1,"target":0,"pai":"C","consumed":["C","C","C"]}`,
			want: Daiminkan{
				Action: Action{
					Message: Message{TypeDaiminkan},
					Actor:   -1,
					Log:     "",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: `{"type":"daiminkan","actor":4,"target":0,"pai":"C","consumed":["C","C","C"]}`,
			want: Daiminkan{
				Action: Action{
					Message: Message{TypeDaiminkan},
					Actor:   4,
					Log:     "",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			wantErr: true,
		},
		{
			name: "invalid target min",
			args: `{"type":"daiminkan","actor":0,"target":-1,"pai":"C","consumed":["C","C","C"]}`,
			want: Daiminkan{
				Action: Action{
					Message: Message{TypeDaiminkan},
					Actor:   0,
					Log:     "",
				},
				Target:   -1,
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			wantErr: true,
		},
		{
			name: "invalid target max",
			args: `{"type":"daiminkan","actor":0,"target":4,"pai":"C","consumed":["C","C","C"]}`,
			want: Daiminkan{
				Action: Action{
					Message: Message{TypeDaiminkan},
					Actor:   0,
					Log:     "",
				},
				Target:   4,
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			wantErr: true,
		},
		{
			name: "empty pai",
			args: `{"type":"daiminkan","actor":1,"target":0,"pai":"","consumed":["C","C","C"]}`,
			want: Daiminkan{
				Action: Action{
					Message: Message{TypeDaiminkan},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "",
				Consumed: [3]string{"C", "C", "C"},
			},
			wantErr: true,
		},
		{
			name: "invalid pai",
			args: `{"type":"daiminkan","actor":1,"target":0,"pai":"6sr","consumed":["6s","6s","6s"]}`,
			want: Daiminkan{
				Action: Action{
					Message: Message{TypeDaiminkan},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "6sr",
				Consumed: [3]string{"6s", "6s", "6s"},
			},
			wantErr: true,
		},
		{
			name: "empty consumed pai",
			args: `{"type":"daiminkan","actor":1,"target":0,"pai":"C","consumed":["","",""]}`,
			want: Daiminkan{
				Action: Action{
					Message: Message{TypeDaiminkan},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "C",
				Consumed: [3]string{"", "", ""},
			},
			wantErr: true,
		},
		{
			name: "invalid consumed pai",
			args: `{"type":"daiminkan","actor":1,"target":0,"pai":"6s","consumed":["6sr","6s","6s"]}`,
			want: Daiminkan{
				Action: Action{
					Message: Message{TypeDaiminkan},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "6s",
				Consumed: [3]string{"6sr", "6s", "6s"},
			},
			wantErr: true,
		},
		{
			name:    "invalid consumed len 2",
			args:    `{"type":"daiminkan","actor":1,"target":0,"pai":"5s","consumed":["5sr","5s"]}`,
			want:    Daiminkan{},
			wantErr: true,
		},
		{
			name:    "invalid consumed len 4",
			args:    `{"type":"daiminkan","actor":1,"target":0,"pai":"6s","consumed":["6s","6s","6s","6s"]}`,
			want:    Daiminkan{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Daiminkan
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
