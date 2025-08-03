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
		name    string
		args    args
		want    *Ryukyoku
		wantErr bool
	}{
		{
			name: "no scores",
			args: args{
				scores: nil,
			},
			want: &Ryukyoku{
				Scores: nil,
			},
			wantErr: false,
		},
		{
			name: "with scores",
			args: args{
				scores: &[4]int{25000, 24000, 27000, 23000},
			},
			want: &Ryukyoku{
				Scores: &[4]int{25000, 24000, 27000, 23000},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRyukyoku(tt.args.scores)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRyukyoku() error = %v, want %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRyukyoku() = %v, want %v", got, tt.want)
			}
		})
	}
}
