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
