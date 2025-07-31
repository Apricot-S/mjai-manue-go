package game

import (
	"log"
	"reflect"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

func mustPlayer(id int) *base.Player {
	if id < base.MinPlayerID || base.MaxPlayerID < id {
		log.Panicf("player ID is invalid: %d", id)
	}
	p, _ := base.NewPlayer(id, "", InitScore)
	return p
}

func Test_GetPlayerDistance(t *testing.T) {
	type args struct {
		p1 *base.Player
		p2 *base.Player
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "p1 and p2 are the same, p1 is 0",
			args: args{p1: mustPlayer(0), p2: mustPlayer(0)},
			want: 0,
		},
		{
			name: "p1 and p2 are the same, p1 is 1",
			args: args{p1: mustPlayer(1), p2: mustPlayer(1)},
			want: 0,
		},
		{
			name: "p1 and p2 are the same, p1 is 2",
			args: args{p1: mustPlayer(2), p2: mustPlayer(2)},
			want: 0,
		},
		{
			name: "p1 and p2 are the same, p1 is 3",
			args: args{p1: mustPlayer(3), p2: mustPlayer(3)},
			want: 0,
		},
		{
			name: "p2 is shimocha of p1, p1 is 0",
			args: args{p1: mustPlayer(0), p2: mustPlayer(1)},
			want: 3,
		},
		{
			name: "p2 is shimocha of p1, p1 is 1",
			args: args{p1: mustPlayer(1), p2: mustPlayer(2)},
			want: 3,
		},
		{
			name: "p2 is shimocha of p1, p1 is 2",
			args: args{p1: mustPlayer(2), p2: mustPlayer(3)},
			want: 3,
		},
		{
			name: "p2 is shimocha of p1, p1 is 3",
			args: args{p1: mustPlayer(3), p2: mustPlayer(0)},
			want: 3,
		},
		{
			name: "p2 is toimen of p1, p1 is 0",
			args: args{p1: mustPlayer(0), p2: mustPlayer(2)},
			want: 2,
		},
		{
			name: "p2 is toimen of p1, p1 is 1",
			args: args{p1: mustPlayer(1), p2: mustPlayer(3)},
			want: 2,
		},
		{
			name: "p2 is toimen of p1, p1 is 2",
			args: args{p1: mustPlayer(2), p2: mustPlayer(0)},
			want: 2,
		},
		{
			name: "p2 is toimen of p1, p1 is 3",
			args: args{p1: mustPlayer(3), p2: mustPlayer(1)},
			want: 2,
		},
		{
			name: "p2 is kamicha of p1, p1 is 0",
			args: args{p1: mustPlayer(0), p2: mustPlayer(3)},
			want: 1,
		},
		{
			name: "p2 is kamicha of p1, p1 is 1",
			args: args{p1: mustPlayer(1), p2: mustPlayer(0)},
			want: 1,
		},
		{
			name: "p2 is kamicha of p1, p1 is 2",
			args: args{p1: mustPlayer(2), p2: mustPlayer(1)},
			want: 1,
		},
		{
			name: "p2 is kamicha of p1, p1 is 3",
			args: args{p1: mustPlayer(3), p2: mustPlayer(2)},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPlayerDistance(tt.args.p1, tt.args.p2); got != tt.want {
				t.Errorf("GetPlayerDistance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNextKyoku(t *testing.T) {
	east, _ := base.NewPaiWithName("E")
	south, _ := base.NewPaiWithName("S")
	west, _ := base.NewPaiWithName("W")
	north, _ := base.NewPaiWithName("N")

	type args struct {
		bakaze   *base.Pai
		kyokuNum int
	}
	tests := []struct {
		name  string
		args  args
		want  *base.Pai
		want1 int
	}{
		{
			name:  "E1 -> E2",
			args:  args{bakaze: east, kyokuNum: 1},
			want:  east,
			want1: 2,
		},
		{
			name:  "E4 -> S1",
			args:  args{bakaze: east, kyokuNum: 4},
			want:  south,
			want1: 1,
		},
		{
			name:  "S2 -> S3",
			args:  args{bakaze: south, kyokuNum: 2},
			want:  south,
			want1: 3,
		},
		{
			name:  "S4 -> W1",
			args:  args{bakaze: south, kyokuNum: 4},
			want:  west,
			want1: 1,
		},
		{
			name:  "W3 -> W4",
			args:  args{bakaze: west, kyokuNum: 3},
			want:  west,
			want1: 4,
		},
		{
			name:  "W4 -> N1",
			args:  args{bakaze: west, kyokuNum: 4},
			want:  north,
			want1: 1,
		},
		{
			name:  "N4 -> E1",
			args:  args{bakaze: north, kyokuNum: 4},
			want:  east,
			want1: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getNextKyoku(tt.args.bakaze, tt.args.kyokuNum)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getNextKyoku() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getNextKyoku() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
