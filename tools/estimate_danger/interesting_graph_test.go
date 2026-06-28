package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/configs"
)

func TestCreatePointsFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "points")
	nodes := []*configs.DecisionNode{
		{
			AverageProb:  0.25,
			ConfInterval: [2]float64{0.1, 0.4},
			NumSamples:   8,
		},
		nil,
		{
			AverageProb:  0.5,
			ConfInterval: [2]float64{0.2, 0.8},
			NumSamples:   4,
		},
	}

	if err := createPointsFile(path, nodes, 0.05); err != nil {
		t.Fatalf("createPointsFile() error = %v", err)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	want := "1.050000\t25.000000\t10.000000\t40.000000\n" +
		"3.050000\t50.000000\t20.000000\t80.000000\n"
	if got := string(b); got != want {
		t.Errorf("points file = %q, want %q", got, want)
	}
}

func TestCreateInterestingGraphs(t *testing.T) {
	probs := make(map[string]*configs.DecisionNode)
	entry := supaiCriteria[0]
	baseCriteria := getNumberCriteria(entry.Base)
	testCriteria := getNumberCriteria(entry.Test[0])
	for i, criterion := range baseCriteria {
		if criterion != nil {
			key, err := encodeCriterion(criterion)
			if err != nil {
				t.Fatalf("encodeCriterion(base[%d]) error = %v", i, err)
			}
			probs[key] = &configs.DecisionNode{AverageProb: float64(i+1) / 100.0}
		}
	}
	for i, criterion := range testCriteria {
		if criterion != nil {
			key, err := encodeCriterion(criterion)
			if err != nil {
				t.Fatalf("encodeCriterion(test[%d]) error = %v", i, err)
			}
			probs[key] = &configs.DecisionNode{AverageProb: float64(i+11) / 100.0}
		}
	}

	var gnuplotCalls int
	outputDir := t.TempDir()
	err := createInterestingGraphs(probs, outputDir, func(id int, spec string, outputDir string) error {
		gnuplotCalls++
		if !strings.Contains(spec, fmt.Sprintf("%d.base.points", id)) {
			t.Errorf("gnuplot spec for %d does not reference base points: %q", id, spec)
		}
		return os.WriteFile(filepath.Join(outputDir, fmt.Sprintf("%d.plot", id)), []byte(spec), 0644)
	})
	if err != nil {
		t.Fatalf("createInterestingGraphs() error = %v", err)
	}

	if gnuplotCalls == 0 {
		t.Fatal("gnuplot runner was not called")
	}
	for _, name := range []string{"0.base.points", "0.test.points", "0.plot", "graphs.html"} {
		if _, err := os.Stat(filepath.Join(outputDir, name)); err != nil {
			t.Errorf("expected %s to be created: %v", name, err)
		}
	}
}
