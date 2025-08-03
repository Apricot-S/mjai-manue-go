package inbound

import (
	"reflect"
	"testing"
)

func TestNewRyukyoku(t *testing.T) {
	type args struct {
		scores *[4]int
	}
	tests := []struct {
		name string
		args args
		want *Ryukyoku
	}{
		{
			name: "no scores",
			args: args{
				scores: nil,
			},
			want: &Ryukyoku{
				Scores: nil,
			},
		},
		{
			name: "with scores",
			args: args{
				scores: &[4]int{25000, 24000, 27000, 23000},
			},
			want: &Ryukyoku{
				Scores: &[4]int{25000, 24000, 27000, 23000},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRyukyoku(tt.args.scores)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRyukyoku() = %v, want %v", got, tt.want)
			}
		})
	}
}
