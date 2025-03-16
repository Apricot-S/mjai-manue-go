package message

import (
	"reflect"
	"testing"

	"github.com/go-json-experiment/json"
)

func TestAction_Marshal(t *testing.T) {
	actionMessage := Message{Type: "action"}

	tests := []struct {
		name    string
		args    *Action
		want    string
		wantErr bool
	}{
		{
			name:    "min_actor",
			args:    &Action{Message: actionMessage, Actor: 0},
			want:    `{"type":"action","actor":0}`,
			wantErr: false,
		},
		{
			name:    "max_actor",
			args:    &Action{Message: actionMessage, Actor: 3},
			want:    `{"type":"action","actor":3}`,
			wantErr: false,
		},
		{
			name:    "with_log",
			args:    &Action{Message: actionMessage, Actor: 1, Log: "hello"},
			want:    `{"type":"action","actor":1,"log":"hello"}`,
			wantErr: false,
		},
		{
			name:    "actor_out_of_range_lower",
			args:    &Action{Message: actionMessage, Actor: -1},
			want:    `{"type":"action","actor":-1}`,
			wantErr: true,
		},
		{
			name:    "actor_out_of_range_upper",
			args:    &Action{Message: actionMessage, Actor: 4},
			want:    `{"type":"action","actor":4}`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.args)
			if err != nil {
				t.Errorf("marshal error = %v, want %v", err, tt.wantErr)
			}
			if string(got) != tt.want {
				t.Errorf("Marshal() = %v, want %v", string(got), tt.want)
			}

			if err := messageValidator.Struct(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("validation error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestAction_Unmarshal(t *testing.T) {
	actionMessage := Message{Type: "action"}

	tests := []struct {
		name    string
		args    string
		want    Action
		wantErr bool
	}{
		{
			name:    "min_actor",
			args:    `{"type":"action","actor":0}`,
			want:    Action{Message: actionMessage, Actor: 0},
			wantErr: false,
		},
		{
			name:    "max_actor",
			args:    `{"type":"action","actor":3}`,
			want:    Action{Message: actionMessage, Actor: 3},
			wantErr: false,
		},
		{
			name:    "missing_actor_is_treated_as_0",
			args:    `{"type":"action"}`,
			want:    Action{Message: actionMessage},
			wantErr: false,
		},
		{
			name:    "null_is_treated_as_0",
			args:    `{"type":"action","actor":null}`,
			want:    Action{Message: actionMessage},
			wantErr: false,
		},
		{
			name:    "with_log",
			args:    `{"type":"action","actor":1,"log":"hello"}`,
			want:    Action{Message: actionMessage, Actor: 1, Log: "hello"},
			wantErr: false,
		},
		{
			name:    "actor_out_of_range_lower",
			args:    `{"type":"action","actor":-1}`,
			want:    Action{Message: actionMessage, Actor: -1},
			wantErr: true,
		},
		{
			name:    "actor_out_of_range_upper",
			args:    `{"type":"action","actor":4}`,
			want:    Action{Message: actionMessage, Actor: 4},
			wantErr: true,
		},
		{
			name:    "empty type",
			args:    `{"type":""}`,
			want:    Action{Message: Message{""}},
			wantErr: true,
		},
		{
			name:    "null type",
			args:    `{"type":null}`,
			want:    Action{Message: Message{""}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Action
			err := json.Unmarshal([]byte(tt.args), &got)
			if err != nil {
				t.Errorf("unmarshal error = %v, want %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unmarshal() = %v, want %v", got, tt.want)
			}

			if err := messageValidator.Struct(got); (err != nil) != tt.wantErr {
				t.Errorf("validation error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
