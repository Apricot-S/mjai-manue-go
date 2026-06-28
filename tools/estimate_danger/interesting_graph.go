package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/Apricot-S/mjai-manue-go/configs"
)

const interestingGraphOutputDir = "exp/graphs"

var (
	criterionTitleTypeRegexp       = regexp.MustCompile(`^main\.Criterion\{`)
	criterionTitleQuoteRegexp      = regexp.MustCompile(`"`)
	criterionTitleUnderscoreRegexp = regexp.MustCompile(`_`)
	criterionTitleCloseRegexp      = regexp.MustCompile(`}`)
)

type gnuplotRunner func(id int, spec string, outputDir string) error

func createPointsFile(path string, nodes []*configs.DecisionNode, gap float64) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create points file: %w", err)
	}
	defer f.Close()

	for i, node := range nodes {
		if node == nil {
			continue
		}
		if _, err := fmt.Fprintf(
			f,
			"%f\t%f\t%f\t%f\n",
			float64(i+1)+gap,
			node.AverageProb*100.0,
			node.ConfInterval[0]*100.0,
			node.ConfInterval[1]*100.0,
		); err != nil {
			return fmt.Errorf("failed to write points file: %w", err)
		}
	}
	return nil
}

func formatCriterionTitle(c Criterion) string {
	title := fmt.Sprintf("%#v", c)
	title = criterionTitleTypeRegexp.ReplaceAllString(title, `\\{`)
	title = criterionTitleQuoteRegexp.ReplaceAllString(title, `\"`)
	title = criterionTitleUnderscoreRegexp.ReplaceAllString(title, `\\_`)
	title = criterionTitleCloseRegexp.ReplaceAllString(title, `\\}`)
	return title
}

func generateGnuplotSpec(id int, baseTitle, testTitle string, outputDir string) string {
	return fmt.Sprintf(`
 set encoding utf8
 set terminal pngcairo size 640,480 font "IPAGothic"
 set output "%s/%d.graph.png"
 set xrange [0:6]
 set yrange [0:25]
 set xlabel "牌の数字"
 set ylabel "放銃率 [%%]"
 set xtics ("1,9" 1, "2,8" 2, "3,7" 3, "4,6" 4, "5" 5)
 plot  "%s/%d.base.points" using 1:2:3:4 with yerrorbars title "%s", \
   "%s/%d.test.points" using 1:2:3:4 with yerrorbars title "%s"
`,
		outputDir,
		id,
		outputDir,
		id,
		baseTitle,
		outputDir,
		id,
		testTitle,
	)
}

func executeGnuplot(id int, spec string, outputDir string) error {
	plotFile := fmt.Sprintf("%s/%d.plot", outputDir, id)
	if err := os.WriteFile(plotFile, []byte(spec), 0644); err != nil {
		return fmt.Errorf("failed to write plot file: %w", err)
	}

	if err := exec.Command("gnuplot", plotFile).Run(); err != nil {
		return fmt.Errorf("failed to execute gnuplot: %w", err)
	}
	return nil
}

func createGraphsHTML(numGraphs int, outputDir string) error {
	f, err := os.Create(fmt.Sprintf("%s/graphs.html", outputDir))
	if err != nil {
		return fmt.Errorf("failed to create graphs html: %w", err)
	}
	defer f.Close()

	for i := range numGraphs {
		if _, err := fmt.Fprintf(f, "<div><img src='%d.graph.png'></div>\n", i); err != nil {
			return fmt.Errorf("failed to write graphs html: %w", err)
		}
	}
	return nil
}

func probabilityNodesForNumberCriteria(
	probs map[string]*configs.DecisionNode,
	criteria []Criterion,
) ([]*configs.DecisionNode, error) {
	nodes := make([]*configs.DecisionNode, len(criteria))
	for i, criterion := range criteria {
		key, err := encodeCriterion(criterion)
		if err != nil {
			return nil, err
		}
		nodes[i] = probs[key]
	}
	return nodes, nil
}

func createInterestingGraphs(probs map[string]*configs.DecisionNode, outputDir string, runGnuplot gnuplotRunner) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create graph output dir: %w", err)
	}

	id := 0
	for _, entry := range supaiCriteria {
		for _, testCriterion := range entry.Test {
			baseNodes, err := probabilityNodesForNumberCriteria(probs, getNumberCriteria(entry.Base))
			if err != nil {
				return err
			}
			testNodes, err := probabilityNodesForNumberCriteria(probs, getNumberCriteria(testCriterion))
			if err != nil {
				return err
			}

			if err := createPointsFile(fmt.Sprintf("%s/%d.base.points", outputDir, id), baseNodes, 0.0); err != nil {
				return err
			}
			if err := createPointsFile(fmt.Sprintf("%s/%d.test.points", outputDir, id), testNodes, 0.05); err != nil {
				return err
			}

			spec := generateGnuplotSpec(
				id,
				formatCriterionTitle(entry.Base),
				formatCriterionTitle(testCriterion),
				outputDir,
			)
			if err := runGnuplot(id, spec, outputDir); err != nil {
				return err
			}
			id++
		}
	}

	return createGraphsHTML(id, outputDir)
}

func RunInterestingGraph(probsPath string) error {
	f, err := os.Open(probsPath)
	if err != nil {
		return fmt.Errorf("failed to open probabilities file: %w", err)
	}
	defer f.Close()

	var probs map[string]*configs.DecisionNode
	if err := gob.NewDecoder(f).Decode(&probs); err != nil {
		return fmt.Errorf("failed to load probabilities: %w", err)
	}

	return createInterestingGraphs(probs, interestingGraphOutputDir, executeGnuplot)
}
