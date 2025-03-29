package ai

import (
	"fmt"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/game"
)

func TestNewBitVector(t *testing.T) {
	type args struct {
		countVector *game.PaiSet
		threshold   int
	}
	type testCase struct {
		name string
		args args
		want BitVector
	}
	tests := []testCase{}

	{
		ps := game.PaiSet{}
		tests = append(tests, testCase{
			name: "empty",
			args: args{countVector: &ps, threshold: 0},
			want: BitVector(0x3_FFFF_FFFF),
		})
		for i := 1; i <= 4; i++ {
			tests = append(tests, testCase{
				name: fmt.Sprintf("empty %d", i),
				args: args{countVector: &ps, threshold: i},
				want: BitVector(0),
			})
		}
	}

	{
		ps := game.PaiSet{0, 1, 2, 3, 4}
		tests = append(tests, testCase{
			name: "threshold 0",
			args: args{countVector: &ps, threshold: 0},
			want: BitVector(0x3_FFFF_FFFF),
		})
		tests = append(tests, testCase{
			name: "threshold 1",
			args: args{countVector: &ps, threshold: 1},
			want: BitVector(0x1E), // 11110 in binary
		})
		tests = append(tests, testCase{
			name: "threshold 2",
			args: args{countVector: &ps, threshold: 2},
			want: BitVector(0x1C), // 11100 in binary
		})
		tests = append(tests, testCase{
			name: "threshold 3",
			args: args{countVector: &ps, threshold: 3},
			want: BitVector(0x18), // 11000 in binary
		})
		tests = append(tests, testCase{
			name: "threshold 4",
			args: args{countVector: &ps, threshold: 4},
			want: BitVector(0x10), // 10000 in binary
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBitVector(tt.args.countVector, tt.args.threshold); got != tt.want {
				t.Errorf("NewBitVector() = %v, want %v", got, tt.want)
			}
		})
	}
}
