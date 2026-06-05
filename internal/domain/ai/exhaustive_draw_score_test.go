package ai

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
)

func TestTenpaiProb_ReturnsOneWithRiichi(t *testing.T) {
	got := tenpaiProb(stubManueStats{}, true, 10, 0)

	if got != 1.0 {
		t.Errorf("tenpaiProb() = %v, want 1", got)
	}
}

func TestTenpaiProb_ReturnsYamitenRatio(t *testing.T) {
	got := tenpaiProb(stubManueStats{
		yamitenCounts: map[string]yamitenCount{
			"10,2": {total: 20, tenpai: 5},
		},
	}, false, 10, 2)

	want := 0.25
	if got != want {
		t.Errorf("tenpaiProb() = %v, want %v", got, want)
	}
}

func TestTenpaiProb_ReturnsOneWithoutStats(t *testing.T) {
	got := tenpaiProb(stubManueStats{}, false, 10, 2)

	if got != 1.0 {
		t.Errorf("tenpaiProb() = %v, want 1", got)
	}
}

func TestNotenExhaustiveDrawTenpaiProb(t *testing.T) {
	got, err := notenExhaustiveDrawTenpaiProb(stubManueStats{
		exhaustiveDrawNotenCount: 100,
		exhaustiveDrawTenpaiTurnFreqs: map[string]int{
			"16.25": 30,
			"16.5":  20,
			"16.75": 10,
			"17":    40,
			"17.25": 0,
			"17.5":  0,
		},
	}, 16)
	if err != nil {
		t.Fatalf("notenExhaustiveDrawTenpaiProb() failed: %v", err)
	}

	want := 0.5
	if got != want {
		t.Errorf("notenExhaustiveDrawTenpaiProb() = %v, want %v", got, want)
	}
}

func TestNotenExhaustiveDrawTenpaiProb_UsesFutureTurnsOnly(t *testing.T) {
	got, err := notenExhaustiveDrawTenpaiProb(stubManueStats{
		exhaustiveDrawNotenCount: 100,
		exhaustiveDrawTenpaiTurnFreqs: map[string]int{
			"15.75": 1000,
			"16":    1000,
			"16.25": 100,
			"16.5":  0,
			"16.75": 0,
			"17":    0,
			"17.25": 0,
			"17.5":  0,
		},
	}, 16)
	if err != nil {
		t.Fatalf("notenExhaustiveDrawTenpaiProb() failed: %v", err)
	}

	want := 0.5
	if got != want {
		t.Errorf("notenExhaustiveDrawTenpaiProb() = %v, want %v", got, want)
	}
}

func TestNotenExhaustiveDrawTenpaiProb_AllowsExistingZeroFreq(t *testing.T) {
	got, err := notenExhaustiveDrawTenpaiProb(stubManueStats{
		exhaustiveDrawNotenCount: 100,
		exhaustiveDrawTenpaiTurnFreqs: map[string]int{
			"16.25": 0,
			"16.5":  100,
			"16.75": 0,
			"17":    0,
			"17.25": 0,
			"17.5":  0,
		},
	}, 16)
	if err != nil {
		t.Fatalf("notenExhaustiveDrawTenpaiProb() failed: %v", err)
	}

	want := 0.5
	if got != want {
		t.Errorf("notenExhaustiveDrawTenpaiProb() = %v, want %v", got, want)
	}
}

func TestNotenExhaustiveDrawTenpaiProb_ReturnsErrorWithoutFreqs(t *testing.T) {
	_, err := notenExhaustiveDrawTenpaiProb(stubManueStats{}, 16)
	if err == nil {
		t.Fatal("notenExhaustiveDrawTenpaiProb() succeeded unexpectedly")
	}
}

func TestExhaustiveDrawTenpaiProbs(t *testing.T) {
	got := exhaustiveDrawTenpaiProbs([common.NumPlayers]float64{0, 0.25, 0.5, 1}, 0.4)
	want := [common.NumPlayers]float64{0.4, 0.55, 0.7, 1}
	if got != want {
		t.Errorf("exhaustiveDrawTenpaiProbs() = %v, want %v", got, want)
	}
}

func TestRyukyokuScoreDelta(t *testing.T) {
	got := ryukyokuScoreDelta([common.NumPlayers]bool{true, false, true, false})
	want := scoreDelta{1500, -1500, 1500, -1500}
	if got != want {
		t.Errorf("ryukyokuScoreDelta() = %v, want %v", got, want)
	}
}

func TestExhaustiveDrawScoreDeltaDist(t *testing.T) {
	got := exhaustiveDrawScoreDeltaDist([common.NumPlayers]float64{1, 0, 0.5, 0})
	want := scoreDeltaProbDist{
		{3000, -1000, -1000, -1000}: 0.5,
		{1500, -1500, 1500, -1500}:  0.5,
	}
	assertScoreDeltaProbDist(t, got, want)
}

func TestExhaustiveDrawAvgPts(t *testing.T) {
	got := exhaustiveDrawAvgPts(0, [common.NumPlayers]float64{1, 0, 0.5, 0})
	want := 2250.0
	if got != want {
		t.Errorf("exhaustiveDrawAvgPts() = %v, want %v", got, want)
	}
}
