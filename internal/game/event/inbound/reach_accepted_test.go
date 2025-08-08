package inbound

import (
	"reflect"
	"testing"
)

func TestNewReachAccepted(t *testing.T) {
	type args struct {
		actor  int
		scores *[4]int
	}
	tests := []struct {
		name    string
		args    args
		want    *ReachAccepted
		wantErr bool
	}{
		{
			name: "valid actor and no scores",
			args: args{
				actor:  0,
				scores: nil,
			},
			want: &ReachAccepted{
				Actor:  0,
				Scores: nil,
			},
			wantErr: false,
		},
		{
			name: "valid actor and scores",
			args: args{
				actor:  3,
				scores: &[4]int{25000, 24000, 27000, 23000},
			},
			want: &ReachAccepted{
				Actor:  3,
				Scores: &[4]int{25000, 24000, 27000, 23000},
			},
			wantErr: false,
		},
		{
			name: "invalid actor min",
			args: args{
				actor:  -1,
				scores: &[4]int{25000, 24000, 27000, 23000},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid actor max",
			args: args{
				actor:  4,
				scores: &[4]int{25000, 24000, 27000, 23000},
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
