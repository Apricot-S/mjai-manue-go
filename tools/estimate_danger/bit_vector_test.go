package main

import (
	"math/big"
	"reflect"
	"testing"
)

func TestBoolArrayToBitVector(t *testing.T) {
	tests := []struct {
		name string
		args []bool
		want *BitVector
	}{
		{
			name: "000",
			args: []bool{false, false, false},
			want: big.NewInt(0b0),
		},
		{
			name: "10010",
			args: []bool{false, true, false, false, true},
			want: big.NewInt(0b10010),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BoolArrayToBitVector(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BoolArrayToBitVector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatches(t *testing.T) {
	type args struct {
		featureVector *BitVector
		positiveMask  *BitVector
		negativeMask  *BitVector
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "feature: 0, pos: 0, neg: 0",
			args: args{
				featureVector: big.NewInt(0),
				positiveMask:  big.NewInt(0),
				negativeMask:  big.NewInt(0),
			},
			want: true,
		},
		{
			name: "feature: 1, pos: 1, neg: 1",
			args: args{
				featureVector: big.NewInt(1),
				positiveMask:  big.NewInt(1),
				negativeMask:  big.NewInt(1),
			},
			want: true,
		},
		{
			name: "feature: 1, pos: 1, neg: 0",
			args: args{
				featureVector: big.NewInt(1),
				positiveMask:  big.NewInt(1),
				negativeMask:  big.NewInt(0),
			},
			want: false,
		},
		{
			name: "feature: 1, pos: 0, neg: 1",
			args: args{
				featureVector: big.NewInt(1),
				positiveMask:  big.NewInt(0),
				negativeMask:  big.NewInt(1),
			},
			want: true,
		},
		{
			name: "feature: 0b1010, pos: 0b1010, neg: 0b1000",
			args: args{
				featureVector: big.NewInt(0b1010),
				positiveMask:  big.NewInt(0b1010),
				negativeMask:  big.NewInt(0b1000),
			},
			want: false,
		},
		{
			name: "feature: 0b1010, pos: 0b0010, neg: 0b1000",
			args: args{
				featureVector: big.NewInt(0b1010),
				positiveMask:  big.NewInt(0b0010),
				negativeMask:  big.NewInt(0b1000),
			},
			want: false,
		},
		{
			name: "feature: 0b1010, pos: 0b0010, neg: 0b0100",
			args: args{
				featureVector: big.NewInt(0b1010),
				positiveMask:  big.NewInt(0b0010),
				negativeMask:  big.NewInt(0b0100),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Matches(tt.args.featureVector, tt.args.positiveMask, tt.args.negativeMask); got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}
