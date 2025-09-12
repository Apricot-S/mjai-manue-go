package core

import (
	"reflect"
	"testing"
)

func TestHashMap(t *testing.T) {
	t.Run("single element keys", func(t *testing.T) {
		h := NewHashMap[float64]()

		// Set & Get
		h.Set(1.0, 100.0)
		if got := h.Get(1.0, 0.0); got != 100.0 {
			t.Errorf("Get(1.0) = %v, want 100.0", got)
		}

		// Default value
		if got := h.Get(2.0, -1.0); got != -1.0 {
			t.Errorf("Get(2.0) with default = %v, want -1.0", got)
		}

		// HasKey
		if !h.HasKey(1.0) {
			t.Error("HasKey(1.0) = false, want true")
		}
		if h.HasKey(2.0) {
			t.Error("HasKey(2.0) = true, want false")
		}
	})

	t.Run("multiple element keys", func(t *testing.T) {
		h := NewHashMap[[4]float64]()

		key1 := [4]float64{1.0, 2.0}
		key2 := [4]float64{3.0, 4.0}

		// Set & Get
		h.Set(key1, 100.0)
		if got := h.Get(key1, 0.0); got != 100.0 {
			t.Errorf("Get(%v) = %v, want 100.0", key1, got)
		}

		// Default value
		if got := h.Get(key2, -1.0); got != -1.0 {
			t.Errorf("Get(%v) with default = %v, want -1.0", key2, got)
		}

		// HasKey
		if !h.HasKey(key1) {
			t.Errorf("HasKey(%v) = false, want true", key1)
		}
		if h.HasKey(key2) {
			t.Errorf("HasKey(%v) = true, want false", key2)
		}
	})

	t.Run("ForEach with single element keys", func(t *testing.T) {
		h := NewHashMap[float64]()
		h.Set(1.0, 100.0)
		h.Set(2.0, 200.0)

		sum := 0.0
		h.ForEach(func(key float64, value float64) {
			sum += value
		})

		if sum != 300.0 {
			t.Errorf("ForEach sum = %v, want 300.0", sum)
		}
	})

	t.Run("ForEach with multiple element keys", func(t *testing.T) {
		h := NewHashMap[[4]float64]()
		key1 := [4]float64{1.0, 2.0}
		key2 := [4]float64{3.0, 4.0}
		h.Set(key1, 100.0)
		h.Set(key2, 200.0)

		sum := 0.0
		keySum := []float64{0.0, 0.0}
		h.ForEach(func(key [4]float64, value float64) {
			sum += value
			keySum[0] += key[0]
			keySum[1] += key[1]
		})

		if sum != 300.0 {
			t.Errorf("ForEach sum = %v, want 300.0", sum)
		}
		wantKeySum := []float64{4.0, 6.0} // (1.0 + 3.0, 2.0 + 4.0)
		if !reflect.DeepEqual(keySum, wantKeySum) {
			t.Errorf("ForEach keySum = %v, want %v", keySum, wantKeySum)
		}
	})
}
