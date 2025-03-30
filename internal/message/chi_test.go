package message

import (
	"reflect"
	"testing"

	"github.com/go-json-experiment/json"
)

func TestNewChi(t *testing.T) {
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
		want    *Chi
		wantErr bool
	}{
		{
			name: "without log",
			args: args{
				actor:    1,
				target:   0,
				pai:      "6s",
				consumed: [2]string{"5sr", "7s"},
				log:      "",
			},
			want: &Chi{
				Action: Action{
					Message: Message{TypeChi},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "6s",
				Consumed: [2]string{"5sr", "7s"},
			},
			wantErr: false,
		},
		{
			name: "with log",
			args: args{
				actor:    3,
				target:   2,
				pai:      "5sr",
				consumed: [2]string{"4s", "6s"},
				log:      "test",
			},
			want: &Chi{
				Action: Action{
					Message: Message{TypeChi},
					Actor:   3,
					Log:     "test",
				},
				Target:   2,
				Pai:      "5sr",
				Consumed: [2]string{"4s", "6s"},
			},
			wantErr: false,
		},
		{
			name: "invalid actor min",
			args: args{
				actor:    -1,
				target:   0,
				pai:      "6s",
				consumed: [2]string{"5sr", "7s"},
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
				pai:      "6s",
				consumed: [2]string{"5sr", "7s"},
				log:      "",
			},
			wantErr: true,
		},
		{
			name: "invalid target min",
			args: args{
				actor:    0,
				target:   -1,
				pai:      "6s",
				consumed: [2]string{"5sr", "7s"},
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
				pai:      "6s",
				consumed: [2]string{"5sr", "7s"},
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
				consumed: [2]string{"5sr", "7s"},
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
				consumed: [2]string{"5sr", "7s"},
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
				pai:      "6s",
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
				pai:      "6s",
				consumed: [2]string{"5sr", "7sr"},
				log:      "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewChi(tt.args.actor, tt.args.target, tt.args.pai, tt.args.consumed, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewChi() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewChi() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChi_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		args    *Chi
		want    string
		wantErr bool
	}{
		{
			name: "without log",
			args: &Chi{
				Action: Action{
					Message: Message{TypeChi},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "6s",
				Consumed: [2]string{"5sr", "7s"},
			},
			want:    `{"type":"chi","actor":1,"target":0,"pai":"6s","consumed":["5sr","7s"]}`,
			wantErr: false,
		},
		{
			name: "with log",
			args: &Chi{
				Action: Action{
					Message: Message{TypeChi},
					Actor:   1,
					Log:     "test",
				},
				Target:   0,
				Pai:      "6s",
				Consumed: [2]string{"5sr", "7s"},
			},
			want:    `{"type":"chi","actor":1,"log":"test","target":0,"pai":"6s","consumed":["5sr","7s"]}`,
			wantErr: false,
		},
		{
			name: "empty type",
			args: &Chi{
				Action: Action{
					Message: Message{""},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "6s",
				Consumed: [2]string{"5sr", "7s"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid type",
			args: &Chi{
				Action: Action{
					Message: Message{TypeNone},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "6s",
				Consumed: [2]string{"5sr", "7s"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid actor min",
			args: &Chi{
				Action: Action{
					Message: Message{TypeChi},
					Actor:   -1,
					Log:     "",
				},
				Target:   0,
				Pai:      "6s",
				Consumed: [2]string{"5sr", "7s"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: &Chi{
				Action: Action{
					Message: Message{TypeChi},
					Actor:   4,
					Log:     "",
				},
				Target:   0,
				Pai:      "6s",
				Consumed: [2]string{"5sr", "7s"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid target min",
			args: &Chi{
				Action: Action{
					Message: Message{TypeChi},
					Actor:   0,
					Log:     "",
				},
				Target:   -1,
				Pai:      "6s",
				Consumed: [2]string{"5sr", "7s"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid target max",
			args: &Chi{
				Action: Action{
					Message: Message{TypeChi},
					Actor:   0,
					Log:     "",
				},
				Target:   4,
				Pai:      "6s",
				Consumed: [2]string{"5sr", "7s"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "empty pai",
			args: &Chi{
				Action: Action{
					Message: Message{TypeChi},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "",
				Consumed: [2]string{"5sr", "7s"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid pai",
			args: &Chi{
				Action: Action{
					Message: Message{TypeChi},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "6sr",
				Consumed: [2]string{"5sr", "7s"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "empty consumed",
			args: &Chi{
				Action: Action{
					Message: Message{TypeChi},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "6s",
				Consumed: [2]string{"", ""},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid consumed",
			args: &Chi{
				Action: Action{
					Message: Message{TypeChi},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "6s",
				Consumed: [2]string{"5sr", "7sr"},
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

func TestChi_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    Chi
		wantErr bool
	}{
		{
			name: "without actor",
			args: `{"type":"chi","target":3,"pai":"6s","consumed":["5sr","7s"]}`,
			want: Chi{
				Action: Action{
					Message: Message{TypeChi},
					Actor:   0,
					Log:     "",
				},
				Target:   3,
				Pai:      "6s",
				Consumed: [2]string{"5sr", "7s"},
			},
			wantErr: false,
		},
		{
			name: "without target",
			args: `{"type":"chi","actor":1,"pai":"6s","consumed":["5sr","7s"]}`,
			want: Chi{
				Action: Action{
					Message: Message{TypeChi},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "6s",
				Consumed: [2]string{"5sr", "7s"},
			},
			wantErr: false,
		},
		{
			name: "without log",
			args: `{"type":"chi","actor":1,"target":0,"pai":"6s","consumed":["5sr","7s"]}`,
			want: Chi{
				Action: Action{
					Message: Message{TypeChi},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "6s",
				Consumed: [2]string{"5sr", "7s"},
			},
			wantErr: false,
		},
		{
			name: "with log",
			args: `{"type":"chi","actor":1,"target":0,"pai":"6s","consumed":["5sr","7s"],"log":"test"}`,
			want: Chi{
				Action: Action{
					Message: Message{TypeChi},
					Actor:   1,
					Log:     "test",
				},
				Target:   0,
				Pai:      "6s",
				Consumed: [2]string{"5sr", "7s"},
			},
			wantErr: false,
		},
		{
			name: "empty type",
			args: `{"type":"","actor":1,"target":0,"pai":"6s","consumed":["5sr","7s"]}`,
			want: Chi{
				Action: Action{
					Message: Message{""},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "6s",
				Consumed: [2]string{"5sr", "7s"},
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			args: `{"type":"none","actor":1,"target":0,"pai":"6s","consumed":["5sr","7s"]}`,
			want: Chi{
				Action: Action{
					Message: Message{TypeNone},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "6s",
				Consumed: [2]string{"5sr", "7s"},
			},
			wantErr: true,
		},
		{
			name: "invalid actor min",
			args: `{"type":"chi","actor":-1,"target":0,"pai":"6s","consumed":["5sr","7s"]}`,
			want: Chi{
				Action: Action{
					Message: Message{TypeChi},
					Actor:   -1,
					Log:     "",
				},
				Target:   0,
				Pai:      "6s",
				Consumed: [2]string{"5sr", "7s"},
			},
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: `{"type":"chi","actor":4,"target":0,"pai":"6s","consumed":["5sr","7s"]}`,
			want: Chi{
				Action: Action{
					Message: Message{TypeChi},
					Actor:   4,
					Log:     "",
				},
				Target:   0,
				Pai:      "6s",
				Consumed: [2]string{"5sr", "7s"},
			},
			wantErr: true,
		},
		{
			name: "invalid target min",
			args: `{"type":"chi","actor":0,"target":-1,"pai":"6s","consumed":["5sr","7s"]}`,
			want: Chi{
				Action: Action{
					Message: Message{TypeChi},
					Actor:   0,
					Log:     "",
				},
				Target:   -1,
				Pai:      "6s",
				Consumed: [2]string{"5sr", "7s"},
			},
			wantErr: true,
		},
		{
			name: "invalid target max",
			args: `{"type":"chi","actor":0,"target":4,"pai":"6s","consumed":["5sr","7s"]}`,
			want: Chi{
				Action: Action{
					Message: Message{TypeChi},
					Actor:   0,
					Log:     "",
				},
				Target:   4,
				Pai:      "6s",
				Consumed: [2]string{"5sr", "7s"},
			},
			wantErr: true,
		},
		{
			name: "empty pai",
			args: `{"type":"chi","actor":1,"target":0,"pai":"","consumed":["5sr","7s"]}`,
			want: Chi{
				Action: Action{
					Message: Message{TypeChi},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "",
				Consumed: [2]string{"5sr", "7s"},
			},
			wantErr: true,
		},
		{
			name: "invalid pai",
			args: `{"type":"chi","actor":1,"target":0,"pai":"6sr","consumed":["5sr","7s"]}`,
			want: Chi{
				Action: Action{
					Message: Message{TypeChi},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "6sr",
				Consumed: [2]string{"5sr", "7s"},
			},
			wantErr: true,
		},
		{
			name: "empty consumed pai",
			args: `{"type":"chi","actor":1,"target":0,"pai":"6s","consumed":["",""]}`,
			want: Chi{
				Action: Action{
					Message: Message{TypeChi},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "6s",
				Consumed: [2]string{"", ""},
			},
			wantErr: true,
		},
		{
			name: "invalid consumed pai",
			args: `{"type":"chi","actor":1,"target":0,"pai":"6s","consumed":["5sr","7sr"]}`,
			want: Chi{
				Action: Action{
					Message: Message{TypeChi},
					Actor:   1,
					Log:     "",
				},
				Target:   0,
				Pai:      "6s",
				Consumed: [2]string{"5sr", "7sr"},
			},
			wantErr: true,
		},
		{
			name:    "invalid consumed len 1",
			args:    `{"type":"chi","actor":1,"target":0,"pai":"6s","consumed":["5sr"]}`,
			want:    Chi{},
			wantErr: true,
		},
		{
			name:    "invalid consumed len 3",
			args:    `{"type":"chi","actor":1,"target":0,"pai":"6s","consumed":["5sr","6s","7s"]}`,
			want:    Chi{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Chi
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
