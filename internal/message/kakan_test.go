package message

import (
	"reflect"
	"testing"

	"github.com/go-json-experiment/json"
)

func TestNewKakan(t *testing.T) {
	type args struct {
		actor    int
		pai      string
		consumed [3]string
		log      string
	}
	tests := []struct {
		name    string
		args    args
		want    *Kakan
		wantErr bool
	}{
		{
			name: "without log",
			args: args{
				actor:    1,
				pai:      "C",
				consumed: [3]string{"C", "C", "C"},
				log:      "",
			},
			want: &Kakan{
				Action: Action{
					Message: Message{TypeKakan},
					Actor:   1,
					Log:     "",
				},
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			wantErr: false,
		},
		{
			name: "with log",
			args: args{
				actor:    3,
				pai:      "5sr",
				consumed: [3]string{"5s", "5s", "5s"},
				log:      "test",
			},
			want: &Kakan{
				Action: Action{
					Message: Message{TypeKakan},
					Actor:   3,
					Log:     "test",
				},
				Pai:      "5sr",
				Consumed: [3]string{"5s", "5s", "5s"},
			},
			wantErr: false,
		},
		{
			name: "invalid actor min",
			args: args{
				actor:    -1,
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
			got, err := NewKakan(tt.args.actor, tt.args.pai, tt.args.consumed, tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewKakan() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKakan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKakan_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		args    *Kakan
		want    string
		wantErr bool
	}{
		{
			name: "without log",
			args: &Kakan{
				Action: Action{
					Message: Message{TypeKakan},
					Actor:   1,
					Log:     "",
				},
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			want:    `{"type":"kakan","actor":1,"pai":"C","consumed":["C","C","C"]}`,
			wantErr: false,
		},
		{
			name: "with log",
			args: &Kakan{
				Action: Action{
					Message: Message{TypeKakan},
					Actor:   1,
					Log:     "test",
				},
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			want:    `{"type":"kakan","actor":1,"log":"test","pai":"C","consumed":["C","C","C"]}`,
			wantErr: false,
		},
		{
			name: "empty type",
			args: &Kakan{
				Action: Action{
					Message: Message{""},
					Actor:   1,
					Log:     "",
				},
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid type",
			args: &Kakan{
				Action: Action{
					Message: Message{TypeNone},
					Actor:   1,
					Log:     "",
				},
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid actor min",
			args: &Kakan{
				Action: Action{
					Message: Message{TypeKakan},
					Actor:   -1,
					Log:     "",
				},
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: &Kakan{
				Action: Action{
					Message: Message{TypeKakan},
					Actor:   4,
					Log:     "",
				},
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "empty pai",
			args: &Kakan{
				Action: Action{
					Message: Message{TypeKakan},
					Actor:   1,
					Log:     "",
				},
				Pai:      "",
				Consumed: [3]string{"C", "C", "C"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid pai",
			args: &Kakan{
				Action: Action{
					Message: Message{TypeKakan},
					Actor:   1,
					Log:     "",
				},
				Pai:      "6sr",
				Consumed: [3]string{"6s", "6s", "6s"},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "empty consumed",
			args: &Kakan{
				Action: Action{
					Message: Message{TypeKakan},
					Actor:   1,
					Log:     "",
				},
				Pai:      "C",
				Consumed: [3]string{"", "", ""},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid consumed",
			args: &Kakan{
				Action: Action{
					Message: Message{TypeKakan},
					Actor:   1,
					Log:     "",
				},
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

func TestKakan_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    Kakan
		wantErr bool
	}{
		{
			name: "without actor",
			args: `{"type":"kakan","pai":"C","consumed":["C","C","C"]}`,
			want: Kakan{
				Action: Action{
					Message: Message{TypeKakan},
					Actor:   0,
					Log:     "",
				},
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			wantErr: false,
		},
		{
			name: "without log",
			args: `{"type":"kakan","actor":1,"pai":"C","consumed":["C","C","C"]}`,
			want: Kakan{
				Action: Action{
					Message: Message{TypeKakan},
					Actor:   1,
					Log:     "",
				},
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			wantErr: false,
		},
		{
			name: "with log",
			args: `{"type":"kakan","actor":1,"pai":"C","consumed":["C","C","C"],"log":"test"}`,
			want: Kakan{
				Action: Action{
					Message: Message{TypeKakan},
					Actor:   1,
					Log:     "test",
				},
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			wantErr: false,
		},
		{
			name: "empty type",
			args: `{"type":"","actor":1,"pai":"C","consumed":["C","C","C"]}`,
			want: Kakan{
				Action: Action{
					Message: Message{""},
					Actor:   1,
					Log:     "",
				},
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			args: `{"type":"none","actor":1,"pai":"C","consumed":["C","C","C"]}`,
			want: Kakan{
				Action: Action{
					Message: Message{TypeNone},
					Actor:   1,
					Log:     "",
				},
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			wantErr: true,
		},
		{
			name: "invalid actor min",
			args: `{"type":"kakan","actor":-1,"pai":"C","consumed":["C","C","C"]}`,
			want: Kakan{
				Action: Action{
					Message: Message{TypeKakan},
					Actor:   -1,
					Log:     "",
				},
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: `{"type":"kakan","actor":4,"pai":"C","consumed":["C","C","C"]}`,
			want: Kakan{
				Action: Action{
					Message: Message{TypeKakan},
					Actor:   4,
					Log:     "",
				},
				Pai:      "C",
				Consumed: [3]string{"C", "C", "C"},
			},
			wantErr: true,
		},
		{
			name: "empty pai",
			args: `{"type":"kakan","actor":1,"pai":"","consumed":["C","C","C"]}`,
			want: Kakan{
				Action: Action{
					Message: Message{TypeKakan},
					Actor:   1,
					Log:     "",
				},
				Pai:      "",
				Consumed: [3]string{"C", "C", "C"},
			},
			wantErr: true,
		},
		{
			name: "invalid pai",
			args: `{"type":"kakan","actor":1,"pai":"6sr","consumed":["6s","6s","6s"]}`,
			want: Kakan{
				Action: Action{
					Message: Message{TypeKakan},
					Actor:   1,
					Log:     "",
				},
				Pai:      "6sr",
				Consumed: [3]string{"6s", "6s", "6s"},
			},
			wantErr: true,
		},
		{
			name: "empty consumed pai",
			args: `{"type":"kakan","actor":1,"pai":"C","consumed":["","",""]}`,
			want: Kakan{
				Action: Action{
					Message: Message{TypeKakan},
					Actor:   1,
					Log:     "",
				},
				Pai:      "C",
				Consumed: [3]string{"", "", ""},
			},
			wantErr: true,
		},
		{
			name: "invalid consumed pai",
			args: `{"type":"kakan","actor":1,"pai":"6s","consumed":["6sr","6s","6s"]}`,
			want: Kakan{
				Action: Action{
					Message: Message{TypeKakan},
					Actor:   1,
					Log:     "",
				},
				Pai:      "6s",
				Consumed: [3]string{"6sr", "6s", "6s"},
			},
			wantErr: true,
		},
		{
			name:    "invalid consumed len 2",
			args:    `{"type":"kakan","actor":1,"pai":"5s","consumed":["5sr","5s"]}`,
			want:    Kakan{},
			wantErr: true,
		},
		{
			name:    "invalid consumed len 4",
			args:    `{"type":"kakan","actor":1,"pai":"6s","consumed":["6s","6s","6s","6s"]}`,
			want:    Kakan{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Kakan
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
