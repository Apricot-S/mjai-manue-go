package mjai

import (
	"encoding/json/v2"
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
)

func TestNewReachAccepted(t *testing.T) {
	type args struct {
		actor  int
		scores []int
	}
	tests := []struct {
		name    string
		args    args
		want    *ReachAccepted
		wantErr bool
	}{
		{
			name: "without scores",
			args: args{
				actor:  0,
				scores: nil,
			},
			want: &ReachAccepted{
				Message: Message{TypeReachAccepted},
				Actor:   0,
				Scores:  nil,
			},
			wantErr: false,
		},
		{
			name: "with scores",
			args: args{
				actor:  0,
				scores: []int{28000, 23000, 24000, 24000},
			},
			want: &ReachAccepted{
				Message: Message{TypeReachAccepted},
				Actor:   0,
				Scores:  []int{28000, 23000, 24000, 24000},
			},
			wantErr: false,
		},
		{
			name: "invalid actor min",
			args: args{
				actor: -1,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: args{
				actor: 4,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid scores empty",
			args: args{
				actor:  0,
				scores: []int{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid scores 3",
			args: args{
				actor:  0,
				scores: []int{28000, 23000, 24000},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid scores 5",
			args: args{
				actor:  0,
				scores: []int{28000, 23000, 24000, 24000, 24000},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewReachAccepted(tt.args.actor, tt.args.scores)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewReachAccepted() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReachAccepted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReachAccepted_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		args    *ReachAccepted
		want    string
		wantErr bool
	}{
		{
			name: "without scores",
			args: &ReachAccepted{
				Message: Message{TypeReachAccepted},
				Actor:   0,
				Scores:  nil,
			},
			want:    `{"type":"reach_accepted","actor":0}`,
			wantErr: false,
		},
		{
			name: "with scores",
			args: &ReachAccepted{
				Message: Message{TypeReachAccepted},
				Actor:   0,
				Scores:  []int{28000, 23000, 24000, 24000},
			},
			want:    `{"type":"reach_accepted","actor":0,"scores":[28000,23000,24000,24000]}`,
			wantErr: false,
		},
		{
			name: "empty type",
			args: &ReachAccepted{
				Message: Message{""},
				Actor:   0,
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid type",
			args: &ReachAccepted{
				Message: Message{TypeHello},
				Actor:   0,
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid actor min",
			args: &ReachAccepted{
				Message: Message{TypeReachAccepted},
				Actor:   -1,
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: &ReachAccepted{
				Message: Message{TypeReachAccepted},
				Actor:   4,
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid scores empty",
			args: &ReachAccepted{
				Message: Message{TypeReachAccepted},
				Actor:   0,
				Scores:  []int{},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid scores 3",
			args: &ReachAccepted{
				Message: Message{TypeReachAccepted},
				Actor:   0,
				Scores:  []int{28000, 23000, 24000},
			},
			want:    ``,
			wantErr: true,
		},
		{
			name: "invalid scores 5",
			args: &ReachAccepted{
				Message: Message{TypeReachAccepted},
				Actor:   0,
				Scores:  []int{28000, 23000, 24000, 24000, 24000},
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

func TestReachAccepted_Unmarshal(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    ReachAccepted
		wantErr bool
	}{
		{
			name: "without actor",
			args: `{"type":"reach_accepted"}`,
			want: ReachAccepted{
				Message: Message{TypeReachAccepted},
				Actor:   0,
				Scores:  nil,
			},
			wantErr: false,
		},
		{
			name: "without scores",
			args: `{"type":"reach_accepted","actor":0}`,
			want: ReachAccepted{
				Message: Message{TypeReachAccepted},
				Actor:   0,
				Scores:  nil,
			},
			wantErr: false,
		},
		{
			name: "with scores",
			args: `{"type":"reach_accepted","actor":1,"deltas":[0,-1000,0,0],"scores":[28000,23000,24000,24000]}`,
			want: ReachAccepted{
				Message: Message{TypeReachAccepted},
				Actor:   1,
				Scores:  []int{28000, 23000, 24000, 24000},
			},
			wantErr: false,
		},
		{
			name: "empty type",
			args: `{"type":"","actor":0}`,
			want: ReachAccepted{
				Message: Message{""},
				Actor:   0,
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			args: `{"type":"hello","actor":3}`,
			want: ReachAccepted{
				Message: Message{TypeHello},
				Actor:   3,
			},
			wantErr: true,
		},
		{
			name: "invalid actor min",
			args: `{"type":"reach_accepted","actor":-1}`,
			want: ReachAccepted{
				Message: Message{TypeReachAccepted},
				Actor:   -1,
			},
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: `{"type":"reach_accepted","actor":4}`,
			want: ReachAccepted{
				Message: Message{TypeReachAccepted},
				Actor:   4,
			},
			wantErr: true,
		},
		{
			name: "invalid scores empty",
			args: `{"type":"reach_accepted","actor":1,"deltas":[0,-1000,0,0],"scores":[]}`,
			want: ReachAccepted{
				Message: Message{TypeReachAccepted},
				Actor:   1,
				Scores:  []int{},
			},
			wantErr: true,
		},
		{
			name: "invalid scores 3",
			args: `{"type":"reach_accepted","actor":1,"deltas":[0,-1000,0,0],"scores":[28000,23000,24000]}`,
			want: ReachAccepted{
				Message: Message{TypeReachAccepted},
				Actor:   1,
				Scores:  []int{28000, 23000, 24000},
			},
			wantErr: true,
		},
		{
			name: "invalid scores 5",
			args: `{"type":"reach_accepted","actor":1,"deltas":[0,-1000,0,0],"scores":[28000,23000,24000,24000,24000]}`,
			want: ReachAccepted{
				Message: Message{TypeReachAccepted},
				Actor:   1,
				Scores:  []int{28000, 23000, 24000, 24000, 24000},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got ReachAccepted
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

func TestReachAccepted_ToEvent(t *testing.T) {
	tests := []struct {
		name    string
		args    *ReachAccepted
		want    *inbound.ReachAccepted
		wantErr bool
	}{
		{
			name: "without scores",
			args: &ReachAccepted{
				Message: Message{Type: TypeReachAccepted},
				Actor:   1,
				Scores:  nil,
			},
			want: &inbound.ReachAccepted{
				Actor:  1,
				Scores: nil,
			},
			wantErr: false,
		},
		{
			name: "with scores",
			args: &ReachAccepted{
				Message: Message{Type: TypeReachAccepted},
				Actor:   2,
				Scores:  []int{26000, 24000, 23000, 24000},
			},
			want: &inbound.ReachAccepted{
				Actor:  2,
				Scores: &[4]int{26000, 24000, 23000, 24000},
			},
			wantErr: false,
		},
		{
			name: "invalid",
			args: &ReachAccepted{
				Message: Message{Type: TypeReachAccepted},
				Actor:   4,
				Scores:  []int{26000, 24000, 23000, 24000},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.ToEvent()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReachAccepted.ToEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReachAccepted.ToEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
