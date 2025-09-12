package core

import (
	"math"
	"reflect"
	"testing"
)

func TestAddVectorVector(t *testing.T) {
	hm3 := NewHashMap[[]float64]()
	hm3.Set([]float64{0, 1, 0, 0}, 0.5)
	hm3.Set([]float64{2, 3, 0, 0}, 0.5)
	pb3 := NewVectorProbDist(hm3)

	wantHm := NewHashMap[[]float64]()
	wantHm.Set([]float64{0, 2, 0, 0}, 0.25)
	wantHm.Set([]float64{2, 4, 0, 0}, 0.5)
	wantHm.Set([]float64{4, 6, 0, 0}, 0.25)
	wantPd := NewVectorProbDist(wantHm)

	tests := []struct {
		name string
		lhs  *VectorProbDist
		rhs  *VectorProbDist
		want *VectorProbDist
	}{
		{
			name: "Add two vector distributions",
			lhs:  pb3,
			rhs:  pb3,
			want: wantPd,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddVectorVector(tt.lhs, tt.rhs)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddVectorVector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMultScalarVector(t *testing.T) {
	hm4 := NewHashMap[float64]()
	hm4.Set(1, 0.5)
	hm4.Set(-1, 0.5)
	pb4 := NewScalarProbDist(hm4)

	hm5 := NewHashMap[[]float64]()
	hm5.Set([]float64{1, 2, 0, 0}, 0.5)
	hm5.Set([]float64{4, 8, 0, 0}, 0.5)
	pb5 := NewVectorProbDist(hm5)

	negZero := math.Copysign(0.0, -1.0)
	wantHm := NewHashMap[[]float64]()
	wantHm.Set([]float64{1, 2, 0, 0}, 0.25)
	wantHm.Set([]float64{-1, -2, negZero, negZero}, 0.25)
	wantHm.Set([]float64{4, 8, 0, 0}, 0.25)
	wantHm.Set([]float64{-4, -8, negZero, negZero}, 0.25)
	wantPd := NewVectorProbDist(wantHm)

	tests := []struct {
		name string
		lhs  *ScalarProbDist
		rhs  *VectorProbDist
		want *VectorProbDist
	}{
		{
			name: "Multiply scalar and vector distributions",
			lhs:  pb4,
			rhs:  pb5,
			want: wantPd,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MultScalarVector(tt.lhs, tt.rhs)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MultScalarVector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeVector(t *testing.T) {
	hm1 := NewHashMap[[]float64]()
	hm1.Set([]float64{0, 0, 0, 0}, 0.5)
	hm1.Set([]float64{8000, 0, 0, 0}, 0.5)
	pb1 := NewVectorProbDist(hm1)

	hm2 := NewHashMap[[]float64]()
	hm2.Set([]float64{0, 0, 0, 0}, 0.5)
	hm2.Set([]float64{-2000, 0, 0, 0}, 0.5)
	pb2 := NewVectorProbDist(hm2)

	wantHm := NewHashMap[[]float64]()
	wantHm.Set([]float64{0, 0, 0, 0}, 0.5)
	wantHm.Set([]float64{8000, 0, 0, 0}, 0.25)
	wantHm.Set([]float64{-2000, 0, 0, 0}, 0.25)
	wantPd := NewVectorProbDist(wantHm)

	tests := []struct {
		name  string
		items []WeightedVectorProbDist
		want  *VectorProbDist
	}{
		{
			name: "Merge two vector distributions",
			items: []WeightedVectorProbDist{
				{pb1, 0.5},
				{pb2, 0.5},
			},
			want: wantPd,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MergeVector(tt.items)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeVector() = %v, want %v", got, tt.want)
			}
		})
	}
}
