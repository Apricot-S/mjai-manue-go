package core

import (
	"reflect"
	"testing"
)

func TestVectorProbDist_Expected(t *testing.T) {
	hm3 := NewHashMap[[]float64]()
	hm3.Set([]float64{0, 1, 0, 0}, 0.5)
	hm3.Set([]float64{2, 3, 0, 0}, 0.5)

	tests := []struct {
		name string
		arg  HashMap[[]float64]
		want []float64
	}{
		{
			name: "simple case",
			arg:  hm3,
			want: []float64{1, 2, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewVectorProbDist(tt.arg)
			got := p.Expected()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Expected() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVectorProbDist_Replace(t *testing.T) {
	type testCase struct {
		name     string
		arg      HashMap[[]float64]
		oldValue []float64
		newPb    *VectorProbDist
		want     *VectorProbDist
	}
	tests := []testCase{}

	{
		hm1 := NewHashMap[[]float64]()
		hm1.Set([]float64{1, 2, 0, 0}, 0.4)
		hm1.Set([]float64{3, 4, 0, 0}, 0.6)

		hm2 := NewHashMap[[]float64]()
		hm2.Set([]float64{5, 6, 0, 0}, 1.0)
		newPb := NewVectorProbDist(hm2)

		wantDist := NewHashMap[[]float64]()
		wantDist.Set([]float64{3, 4, 0, 0}, 0.6)
		wantDist.Set([]float64{5, 6, 0, 0}, 0.4)
		want := NewVectorProbDist(wantDist)

		tests = append(tests, testCase{
			name:     "replace existing value",
			arg:      hm1,
			oldValue: []float64{1, 2, 0, 0},
			newPb:    newPb,
			want:     want,
		})
	}

	{
		hm1 := NewHashMap[[]float64]()
		hm1.Set([]float64{1, 2, 0, 0}, 0.4)
		hm1.Set([]float64{1, 2, 0, 0}, 0.6)

		hm2 := NewHashMap[[]float64]()
		hm2.Set([]float64{5, 6, 0, 0}, 0.5)
		newPb := NewVectorProbDist(hm2)

		wantDist := NewHashMap[[]float64]()
		wantDist.Set([]float64{5, 6, 0, 0}, 0.2)
		wantDist.Set([]float64{5, 6, 0, 0}, 0.3)
		want := NewVectorProbDist(wantDist)

		tests = append(tests, testCase{
			name:     "replace existing value with overlap",
			arg:      hm1,
			oldValue: []float64{1, 2, 0, 0},
			newPb:    newPb,
			want:     want,
		})
	}

	{
		hm1 := NewHashMap[[]float64]()
		hm1.Set([]float64{1, 2, 0, 0}, 0.4)
		hm1.Set([]float64{3, 4, 0, 0}, 0.6)

		hm2 := NewHashMap[[]float64]()
		hm2.Set([]float64{7, 8, 0, 0}, 1.0)
		newPb := NewVectorProbDist(hm2)

		wantDist := NewHashMap[[]float64]()
		want := NewVectorProbDist(wantDist)
		want.dist.Set([]float64{1, 2, 0, 0}, 0.4)
		want.dist.Set([]float64{3, 4, 0, 0}, 0.6)
		want.dist.Set([]float64{7, 8, 0, 0}, 0.0)

		tests = append(tests, testCase{
			name:     "replace non-existing value",
			arg:      hm1,
			oldValue: []float64{5, 6, 0, 0},
			newPb:    newPb,
			want:     want,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewVectorProbDist(tt.arg)
			got := p.Replace(tt.oldValue, tt.newPb)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Replace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVectorProbDist_MapValueScalar(t *testing.T) {
	type testCase struct {
		name   string
		arg    HashMap[[]float64]
		mapper func([]float64) float64
		want   *ScalarProbDist
	}
	tests := []testCase{}

	{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewVectorProbDist(tt.arg)
			got := p.MapValueScalar(tt.mapper)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapValueScalar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVectorProbDist_MapValueVector(t *testing.T) {
	type testCase struct {
		name   string
		arg    HashMap[[]float64]
		mapper func([]float64) []float64
		want   *VectorProbDist
	}
	tests := []testCase{}

	{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewVectorProbDist(tt.arg)
			got := p.MapValueVector(tt.mapper)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapValueVector() = %v, want %v", got, tt.want)
			}
		})
	}
}
