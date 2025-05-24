package core

import (
	"fmt"
	"reflect"
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

	{
		ps := game.PaiSet{33: 4}
		tests = append(tests, testCase{
			name: "threshold 0",
			args: args{countVector: &ps, threshold: 0},
			want: BitVector(0x3_FFFF_FFFF),
		})
		tests = append(tests, testCase{
			name: "threshold 1",
			args: args{countVector: &ps, threshold: 1},
			want: BitVector(0x2_0000_0000),
		})
		tests = append(tests, testCase{
			name: "threshold 2",
			args: args{countVector: &ps, threshold: 2},
			want: BitVector(0x2_0000_0000),
		})
		tests = append(tests, testCase{
			name: "threshold 3",
			args: args{countVector: &ps, threshold: 3},
			want: BitVector(0x2_0000_0000),
		})
		tests = append(tests, testCase{
			name: "threshold 4",
			args: args{countVector: &ps, threshold: 4},
			want: BitVector(0x2_0000_0000),
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

func TestBitVector_IsSubsetOf(t *testing.T) {
	type args struct {
		other BitVector
	}
	type testCase struct {
		name string
		bv   BitVector
		args args
		want bool
	}
	tests := []testCase{
		{
			name: "empty subset of empty",
			bv:   BitVector(0),
			args: args{other: BitVector(0)},
			want: true,
		},
		{
			name: "empty subset of non-empty",
			bv:   BitVector(0),
			args: args{other: BitVector(0xFF)},
			want: true,
		},
		{
			name: "same sets",
			bv:   BitVector(0xFF),
			args: args{other: BitVector(0xFF)},
			want: true,
		},
		{
			name: "proper subset",
			bv:   BitVector(0x0F),
			args: args{other: BitVector(0xFF)},
			want: true,
		},
		{
			name: "not a subset - different bits",
			bv:   BitVector(0xF0),
			args: args{other: BitVector(0x0F)},
			want: false,
		},
		{
			name: "not a subset - extra bits",
			bv:   BitVector(0xFF),
			args: args{other: BitVector(0x0F)},
			want: false,
		},
		{
			name: "large numbers",
			bv:   BitVector(0x2_0000_0000),
			args: args{other: BitVector(0x3_FFFF_FFFF)},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bv.IsSubsetOf(tt.args.other); got != tt.want {
				t.Errorf("BitVector.IsSubsetOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitVector_HasIntersectionWith(t *testing.T) {
	type args struct {
		other BitVector
	}
	type testCase struct {
		name string
		bv   BitVector
		args args
		want bool
	}
	tests := []testCase{
		{
			name: "empty vectors",
			bv:   BitVector(0),
			args: args{other: BitVector(0)},
			want: false,
		},
		{
			name: "one empty, one non-empty",
			bv:   BitVector(0),
			args: args{other: BitVector(0xFF)},
			want: false,
		},
		{
			name: "overlapping vectors",
			bv:   BitVector(0x0F),
			args: args{other: BitVector(0xFF)},
			want: true,
		},
		{
			name: "disjoint vectors",
			bv:   BitVector(0xF0),
			args: args{other: BitVector(0x0F)},
			want: false,
		},
		{
			name: "single bit overlap",
			bv:   BitVector(0x01),
			args: args{other: BitVector(0x01)},
			want: true,
		},
		{
			name: "large numbers with overlap",
			bv:   BitVector(0x2_0000_0000),
			args: args{other: BitVector(0x3_0000_0000)},
			want: true,
		},
		{
			name: "large numbers without overlap",
			bv:   BitVector(0x2_0000_0000),
			args: args{other: BitVector(0x1_0000_0000)},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bv.HasIntersectionWith(tt.args.other); got != tt.want {
				t.Errorf("BitVector.HasIntersectionWith() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCountVectorToBitVectors(t *testing.T) {
	type args struct {
		countVector *game.PaiSet
	}
	type testCase struct {
		name string
		args args
		want [4]BitVector
	}
	tests := []testCase{}

	{
		ps := game.PaiSet{0, 1, 2, 3, 4}
		want := [4]BitVector{
			BitVector(0x1E), // 11110 in binary
			BitVector(0x1C), // 11100 in binary
			BitVector(0x18), // 11000 in binary
			BitVector(0x10), // 10000 in binary
		}
		tests = append(tests, testCase{
			name: "12345m",
			args: args{countVector: &ps},
			want: want,
		})
	}

	{
		ps := game.PaiSet{33: 4}
		want := [4]BitVector{
			BitVector(0x2_0000_0000),
			BitVector(0x2_0000_0000),
			BitVector(0x2_0000_0000),
			BitVector(0x2_0000_0000),
		}
		tests = append(tests, testCase{
			name: "all pais",
			args: args{countVector: &ps},
			want: want,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CountVectorToBitVectors(tt.args.countVector); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CountVectorToBitVectors() = %v, want %v", got, tt.want)
			}
		})
	}
}
