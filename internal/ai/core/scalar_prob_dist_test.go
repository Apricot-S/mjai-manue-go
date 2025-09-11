package core

import (
	"testing"
)

func TestScalarProbDist_Expected(t *testing.T) {
	hm1 := NewHashMap[float64]()
	hm1.Set(0, 0.5)
	hm1.Set(8000, 0.5)

	hm2 := NewHashMap[float64]()
	hm2.Set(0, 0.25)
	hm2.Set(5000, 0.25)
	hm2.Set(-2000, 0.25)
	hm2.Set(1000, 0.25)

	tests := []struct {
		name string
		arg  HashMap[float64]
		want float64
	}{
		{
			name: "Test 1",
			arg:  hm1,
			want: 4000,
		},
		{
			name: "Test 2",
			arg:  hm2,
			want: 1000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewScalarProbDist(tt.arg)
			got := p.Expected()
			if got != tt.want {
				t.Errorf("Expected() = %v, want %v", got, tt.want)
			}
		})
	}
}
