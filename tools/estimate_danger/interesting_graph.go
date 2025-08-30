package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"os/exec"

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

func generateGnuplotSpec(id int, baseTitle, testTitle string) string {
	return fmt.Sprintf(`
 set encoding utf8
 set terminal png size 640,480 font "IPAGothic"
 set output "exp/graphs/%d.graph.png"
 set xrange [0:6]
 set yrange [0:25]
 set xlabel "牌の数字"
 set ylabel "放銃率 [%%]"
 set xtics ("1,9" 1, "2,8" 2, "3,7" 3, "4,6" 4, "5" 5)
 plot  "exp/graphs/%d.base.points" using 1:2:3:4 with yerrorbars title "%s", \
   "exp/graphs/%d.test.points" using 1:2:3:4 with yerrorbars title "%s"
`,
		id,
		id,
		baseTitle,
		id,
		testTitle,
	)
}

func executeGnuplot(id int, spec string, outputDir string) error {
	plotFile := fmt.Sprintf("%s/%d.plot", outputDir, id)
	f, err := os.Create(plotFile)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(spec); err != nil {
		return err
	}

	cmd := exec.Command("gnuplot", plotFile)
	return cmd.Run()
}

func createGraph(probs map[string]*configs.DecisionNode, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

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

			baseTitle := fmt.Sprintf("%v", entry.Base)
			testTitle := fmt.Sprintf("%v", testCriterion)
			spec := generateGnuplotSpec(id, baseTitle, testTitle)
			if err := executeGnuplot(id, spec, outputDir); err != nil {
				return err
			}
			id++
		}
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
