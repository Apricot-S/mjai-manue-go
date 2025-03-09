package game

import "testing"

func Test_getDistance(t *testing.T) {
	type args struct {
		p1 *Player
		p2 *Player
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "p1 and p2 are the same, p1 is 0",
			args: args{p1: &Player{id: 0}, p2: &Player{id: 0}},
			want: 0,
		},
		{
			name: "p1 and p2 are the same, p1 is 1",
			args: args{p1: &Player{id: 1}, p2: &Player{id: 1}},
			want: 0,
		},
		{
			name: "p1 and p2 are the same, p1 is 2",
			args: args{p1: &Player{id: 2}, p2: &Player{id: 2}},
			want: 0,
		},
		{
			name: "p1 and p2 are the same, p1 is 3",
			args: args{p1: &Player{id: 3}, p2: &Player{id: 3}},
			want: 0,
		},
		{
			name: "p2 is shimocha of p1, p1 is 0",
			args: args{p1: &Player{id: 0}, p2: &Player{id: 1}},
			want: 3,
		},
		{
			name: "p2 is shimocha of p1, p1 is 1",
			args: args{p1: &Player{id: 1}, p2: &Player{id: 2}},
			want: 3,
		},
		{
			name: "p2 is shimocha of p1, p1 is 2",
			args: args{p1: &Player{id: 2}, p2: &Player{id: 3}},
			want: 3,
		},
		{
			name: "p2 is shimocha of p1, p1 is 3",
			args: args{p1: &Player{id: 3}, p2: &Player{id: 0}},
			want: 3,
		},
		{
			name: "p2 is toimen of p1, p1 is 0",
			args: args{p1: &Player{id: 0}, p2: &Player{id: 2}},
			want: 2,
		},
		{
			name: "p2 is toimen of p1, p1 is 1",
			args: args{p1: &Player{id: 1}, p2: &Player{id: 3}},
			want: 2,
		},
		{
			name: "p2 is toimen of p1, p1 is 2",
			args: args{p1: &Player{id: 2}, p2: &Player{id: 0}},
			want: 2,
		},
		{
			name: "p2 is toimen of p1, p1 is 3",
			args: args{p1: &Player{id: 3}, p2: &Player{id: 1}},
			want: 2,
		},
		{
			name: "p2 is kamicha of p1, p1 is 0",
			args: args{p1: &Player{id: 0}, p2: &Player{id: 3}},
			want: 1,
		},
		{
			name: "p2 is kamicha of p1, p1 is 1",
			args: args{p1: &Player{id: 1}, p2: &Player{id: 0}},
			want: 1,
		},
		{
			name: "p2 is kamicha of p1, p1 is 2",
			args: args{p1: &Player{id: 2}, p2: &Player{id: 1}},
			want: 1,
		},
		{
			name: "p2 is kamicha of p1, p1 is 3",
			args: args{p1: &Player{id: 3}, p2: &Player{id: 2}},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDistance(tt.args.p1, tt.args.p2); got != tt.want {
				t.Errorf("getDistance() = %v, want %v", got, tt.want)
			}
		})
	}
}
