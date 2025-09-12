package core

import (
	"reflect"
	"slices"
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
		hm1 := NewHashMap[[]float64]()
		hm1.Set([]float64{1, 2, 0, 0}, 0.3)
		hm1.Set([]float64{2, 3, 0, 0}, 0.4)
		hm1.Set([]float64{3, 4, 0, 0}, 0.3)

		mapper := func(v []float64) float64 {
			return v[0] + v[1]
		}

		wantDist := NewHashMap[float64]()
		wantDist.Set(3.0, 0.3)
		wantDist.Set(5.0, 0.4)
		wantDist.Set(7.0, 0.3)
		want := &ScalarProbDist{dist: wantDist}

		tests = append(tests, testCase{
			name:   "map to sum of first two elements",
			arg:    hm1,
			mapper: mapper,
			want:   want,
		})
	}

	{
		hm1 := NewHashMap[[]float64]()
		hm1.Set([]float64{1, 1, 0, 0}, 0.2)
		hm1.Set([]float64{1, 2, 0, 0}, 0.3)
		hm1.Set([]float64{2, 1, 0, 0}, 0.3)
		hm1.Set([]float64{2, 2, 0, 0}, 0.2)

		mapper := func(v []float64) float64 {
			return slices.Max(v)
		}

		wantDist := NewHashMap[float64]()
		wantDist.Set(1.0, 0.2)
		wantDist.Set(2.0, 0.8)
		want := &ScalarProbDist{dist: wantDist}

		tests = append(tests, testCase{
			name:   "map to max of first two elements with overlaps",
			arg:    hm1,
			mapper: mapper,
			want:   want,
		})
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
		hm1 := NewHashMap[[]float64]()
		hm1.Set([]float64{1, 2, 0, 0}, 0.5)
		hm1.Set([]float64{3, 4, 0, 0}, 0.5)

		mapper := func(v []float64) []float64 {
			vv := make([]float64, len(v))
			for i, e := range v {
				vv[i] = e * 2
			}
			return vv
		}

		wantDist := NewHashMap[[]float64]()
		wantDist.Set([]float64{2, 4, 0, 0}, 0.5)
		wantDist.Set([]float64{6, 8, 0, 0}, 0.5)
		want := &VectorProbDist{dist: wantDist}

		tests = append(tests, testCase{
			name:   "map to double each element",
			arg:    hm1,
			mapper: mapper,
			want:   want,
		})
	}

	{
		hm1 := NewHashMap[[]float64]()
		hm1.Set([]float64{1, 2, 0, 0}, 0.2)
		hm1.Set([]float64{2, 1, 0, 0}, 0.3)
		hm1.Set([]float64{3, 4, 0, 0}, 0.5)

		mapper := func(v []float64) []float64 {
			vv := slices.Clone(v)
			slices.Sort(vv)
			return vv
		}

		wantDist := NewHashMap[[]float64]()
		wantDist.Set([]float64{0, 0, 1, 2}, 0.5)
		wantDist.Set([]float64{0, 0, 3, 4}, 0.5)
		want := &VectorProbDist{dist: wantDist}

		tests = append(tests, testCase{
			name:   "map to sorted vector with overlaps",
			arg:    hm1,
			mapper: mapper,
			want:   want,
		})
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
