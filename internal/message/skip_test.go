package message

import (
	"reflect"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/go-playground/validator/v10"
)

func TestNewSkip(t *testing.T) {
	type args struct {
		actor int
		log   string
	}
	tests := []struct {
		name     string
		args     args
		want     *Skip
		wantJSON string
		wantErr  bool
	}{
		{
			name: "without log",
			args: args{
				actor: 0,
			},
			want: &Skip{
				Action: Action{
					Message: Message{TypeNone},
					Actor:   0,
					Log:     "",
				},
			},
			wantJSON: `{"type":"none","actor":0}`,
			wantErr:  false,
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
			wantJSON: `{"type":"none","actor":3,"log":"test"}`,
			wantErr:  false,
		},
		{
			name: "invalid actor",
			args: args{
				actor: -1,
			},
			want: &Skip{
				Action: Action{
					Message: Message{TypeNone},
					Actor:   -1,
					Log:     "",
				},
			},
			wantJSON: `{"type":"none","actor":-1}`,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSkip(tt.args.actor, tt.args.log)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSkip() = %v, want %v", got, tt.want)
			}
			if got.Type != TypeNone {
				t.Errorf("Type = %v, want %v", got.Type, TypeNone)
			}

			jsonData, err := json.Marshal(got)
			if err != nil {
				t.Errorf("marshal error: %v", err)
				return
			}

			if !reflect.DeepEqual(string(jsonData), tt.wantJSON) {
				t.Errorf("expected JSON '%v', got '%v'", tt.wantJSON, string(jsonData))
			}

			validator := validator.New()
			err = validator.Struct(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSkip() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
