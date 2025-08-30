package main

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/Apricot-S/mjai-manue-go/configs"
	"github.com/go-json-experiment/json"
)

const rootDir = "exp/graphs"

func createPointsFile(path string, nodes []*configs.DecisionNode, gap float64) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	for i, node := range nodes {
		if node == nil {
			continue
		}

		line := fmt.Sprintf(
			"%f\t%f\t%f\t%f\n",
			float64(i+1)+gap,
			node.AverageProb*100.0,
			node.ConfInterval[0]*100.0,
			node.ConfInterval[1]*100.0,
		)

		if _, err := f.WriteString(line); err != nil {
			return err
		}
	}

	return nil
}

func createGraph(probs map[string]*configs.DecisionNode, outputDir string) error {
	id := 0
	for _, entry := range SupaiCriteria {
		for _, testCriterion := range entry.Test {
			baseNCriteria := GetNumberCriteria(entry.Base)
			testNCriteria := GetNumberCriteria(testCriterion)
			baseNodes := make([]*configs.DecisionNode, 5)
			testNodes := make([]*configs.DecisionNode, 5)
			for i := range 5 {
				baseKey, err := json.Marshal(baseNCriteria[i], json.Deterministic(true))
				if err != nil {
					return fmt.Errorf("failed to encode criterion: %w", err)
				}
				baseNodes[i] = probs[string(baseKey)]

				testKey, err := json.Marshal(testNCriteria[i], json.Deterministic(true))
				if err != nil {
					return fmt.Errorf("failed to encode criterion: %w", err)
				}
				testNodes[i] = probs[string(testKey)]
			}

			baseFileName := fmt.Sprintf("%s/%d.base.points", outputDir, id)
			testFileName := fmt.Sprintf("%s/%d.test.points", outputDir, id)
			if err := createPointsFile(baseFileName, baseNodes, 0.0); err != nil {
				return err
			}
			if err := createPointsFile(testFileName, testNodes, 0.05); err != nil {
				return err
			}
		}
		id++
	}

	f, err := os.Create(fmt.Sprintf("%s/graphs.html", outputDir))
	if err != nil {
		return err
	}
	defer f.Close()

	for i := range id {
		fmt.Fprintf(f, "<div><img src='%d.graph.png'></div>\n", i)
	}

	return nil
}

func RunInterestingGraph(probsPath string) error {
	r, err := os.Open(probsPath)
	if err != nil {
		return fmt.Errorf("failed to open probabilities file: %w", err)
	}
	defer r.Close()

	decoder := gob.NewDecoder(r)

	var probs map[string]*configs.DecisionNode
	if err := decoder.Decode(&probs); err != nil {
		return fmt.Errorf("failed to load probabilities %w", err)
	}

	return createGraph(probs, rootDir)
}
