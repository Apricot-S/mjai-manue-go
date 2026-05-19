package configs

import (
	"math"
	"testing"
)

func TestLoadDangerTree(t *testing.T) {
	t.Run("danger tree test", func(t *testing.T) {
		epsilon := 1e-15

		got, err := LoadDangerTree()
		if err != nil {
			t.Errorf("LoadDangerTree() error = %v", err)
			return
		}

		// root node
		if math.Abs(got.AverageProb-0.0977659128413358) > epsilon {
			t.Errorf("LoadDangerTree().AverageProb = %v, want %v", got.AverageProb, 0.0977659128413358)
		}
		if math.Abs(got.ConfInterval[0]-0.09699083321177084) > epsilon {
			t.Errorf("LoadDangerTree().ConfInterval[0] = %v, want %v", got.ConfInterval[0], 0.09699083321177084)
		}
		if math.Abs(got.ConfInterval[1]-0.09864226119626654) > epsilon {
			t.Errorf("LoadDangerTree().ConfInterval[1] = %v, want %v", got.ConfInterval[1], 0.09864226119626654)
		}
		if got.NumSamples != 20632 {
			t.Errorf("LoadDangerTree().NumSamples = %v, want %v", got.NumSamples, 20632)
		}
		if *got.FeatureName != "fonpai" {
			t.Errorf("LoadDangerTree().FeatureName = %v, want %v", got.FeatureName, "fonpai")
		}
		if got.Negative == nil {
			t.Errorf("LoadDangerTree().Negative = %v", got.Negative)
		}
		if got.Positive == nil {
			t.Errorf("LoadDangerTree().Positive = %v", got.Positive)
		}

		// child node
		got = got.Positive
		if math.Abs(got.AverageProb-0.02336508484195712) > epsilon {
			t.Errorf("LoadDangerTree().Positive.AverageProb = %v, want %v", got.AverageProb, 0.02336508484195712)
		}
		if math.Abs(got.ConfInterval[0]-0.02135843124124106) > epsilon {
			t.Errorf("LoadDangerTree().Positive.ConfInterval[0] = %v, want %v", got.ConfInterval[0], 0.02135843124124106)
		}
		if math.Abs(got.ConfInterval[1]-0.025510050399443273) > epsilon {
			t.Errorf("LoadDangerTree().Positive.ConfInterval[1] = %v, want %v", got.ConfInterval[1], 0.025510050399443273)
		}
		if got.NumSamples != 15726 {
			t.Errorf("LoadDangerTree().Positive.NumSamples = %v, want %v", got.NumSamples, 15726)
		}
		if got.FeatureName != nil {
			t.Errorf("LoadDangerTree().Positive.FeatureName = %v, want %v", got.FeatureName, nil)
		}
		if got.Negative != nil {
			t.Errorf("LoadDangerTree().Positive.Negative = %v, want %v", got.Negative, nil)
		}
		if got.Positive != nil {
			t.Errorf("LoadDangerTree().Positive.Positive = %v, want %v", got.Positive, nil)
		}
	})
}
