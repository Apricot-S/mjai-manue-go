package main

import (
	"fmt"
	"io"
	"maps"
	"os"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/configs"
)

type SupaiCriteriaEntry struct {
	Base Criterion
	Test []Criterion
}

var TsupaiCriteria = []Criterion{
	{"tsupai": true},
	{"tsupai": true, "sangenpai": true},
	{"tsupai": true, "sangenpai": false},
	{"tsupai": true, "fanpai": true},
	{"tsupai": true, "fanpai": false},
	{"tsupai": true, "visible>=3": true},
	{"tsupai": true, "visible>=3": false, "visible>=2": true},
	{"tsupai": true, "visible>=2": false, "visible>=1": true},
	{"tsupai": true, "visible>=1": false},
}

var SupaiCriteria = []SupaiCriteriaEntry{
	{
		Base: Criterion{"tsupai": false, "suji": false},
		Test: []Criterion{
			{"tsupai": false, "suji": true},
		},
	},

	{
		Base: Criterion{"tsupai": false, "suji": false},
		Test: []Criterion{
			{"tsupai": false, "suji": false, "weak_suji": true},
			{"tsupai": false, "suji": false, "weak_suji": false},
		},
	},

	{
		Base: Criterion{"tsupai": false, "suji": true},
		Test: []Criterion{
			{"tsupai": false, "suji": true, "4<=n<=6": false}, // Omotesuji (表筋)
			{"tsupai": false, "suji": true, "4<=n<=6": true},  // Nakasuji (中筋)
			{"tsupai": false, "suji": true, "reach_suji": true},
			{"tsupai": false, "suji": true, "prereach_suji": true},
			{"tsupai": false, "suji": true, "prereach_suji": false},
		},
	},

	{
		Base: Criterion{"tsupai": false, "suji": false},
		Test: []Criterion{
			{"tsupai": false, "suji": false, "outer_early_sutehai": true},
			{"tsupai": false, "suji": false, "outer_prereach_sutehai": true},

			{"tsupai": false, "suji": false, "urasuji": true},
			{"tsupai": false, "suji": false, "early_urasuji": true},
			{"tsupai": false, "suji": false, "reach_urasuji": true},
			{"tsupai": false, "suji": false, "urasuji_of_5": true},
			{"tsupai": false, "suji": false, "aida4ken": true},
			{"tsupai": false, "suji": false, "matagisuji": true},
			{"tsupai": false, "suji": false, "early_matagisuji": true},
			{"tsupai": false, "suji": false, "late_matagisuji": true},
			{"tsupai": false, "suji": false, "reach_matagisuji": true},
			{"tsupai": false, "suji": false, "senkisuji": true},
			{"tsupai": false, "suji": false, "early_senkisuji": true},

			{"tsupai": false, "suji": false, "chances<=0": true},
			{"tsupai": false, "suji": false, "chances<=0": false, "chances<=1": true},
			{"tsupai": false, "suji": false, "chances<=1": false, "chances<=2": true},
			{"tsupai": false, "suji": false, "chances<=2": false, "chances<=3": true},
			{"tsupai": false, "suji": false, "chances<=3": false},

			{"tsupai": false, "suji": false, "visible>=3": true},
			{"tsupai": false, "suji": false, "visible>=3": false, "visible>=2": true},
			{"tsupai": false, "suji": false, "visible>=2": false, "visible>=1": true},
			{"tsupai": false, "suji": false, "visible>=1": false},

			{"tsupai": false, "suji": false, "suji_visible<=0": true},
			{"tsupai": false, "suji": false, "suji_visible<=0": false, "suji_visible<=1": true},
			{"tsupai": false, "suji": false, "suji_visible<=1": false, "suji_visible<=2": true},
			{"tsupai": false, "suji": false, "suji_visible<=2": false, "suji_visible<=3": true},
			{"tsupai": false, "suji": false, "suji_visible<=3": false},

			{"tsupai": false, "suji": false, "dora": true},
			{"tsupai": false, "suji": false, "dora_suji": true},
			{"tsupai": false, "suji": false, "dora_matagi": true},

			{"tsupai": false, "suji": false, "in_tehais>=4": true},
			{"tsupai": false, "suji": false, "in_tehais>=4": false, "in_tehais>=3": true},
			{"tsupai": false, "suji": false, "in_tehais>=3": false, "in_tehais>=2": true},
			{"tsupai": false, "suji": false, "in_tehais>=2": false},

			{"tsupai": false, "suji": false, "suji_in_tehais>=4": true},
			{"tsupai": false, "suji": false, "suji_in_tehais>=4": false, "suji_in_tehais>=3": true},
			{"tsupai": false, "suji": false, "suji_in_tehais>=3": false, "suji_in_tehais>=2": true},
			{"tsupai": false, "suji": false, "suji_in_tehais>=2": false, "suji_in_tehais>=1": true},
			{"tsupai": false, "suji": false, "suji_in_tehais>=1": false},

			// same_type_in_prereach>=5 is too rare.
			{"tsupai": false, "suji": false, "same_type_in_prereach>=4": true},
			{"tsupai": false, "suji": false, "same_type_in_prereach>=4": false, "same_type_in_prereach>=3": true},
			{"tsupai": false, "suji": false, "same_type_in_prereach>=3": false, "same_type_in_prereach>=2": true},
			{"tsupai": false, "suji": false, "same_type_in_prereach>=2": false, "same_type_in_prereach>=1": true},
			{"tsupai": false, "suji": false, "same_type_in_prereach>=1": false},
		},
	},
}

func GetNumberCriteria(baseCriterion Criterion) []Criterion {
	result := make([]Criterion, 5)
	for i := 1; i <= 5; i++ {
		criterion := make(Criterion)
		maps.Copy(criterion, baseCriterion)

		if i > 1 {
			name := fmt.Sprintf("%d<=n<=%d", i, 10-i)
			if val, ok := criterion[name]; ok && !val {
				result[i-1] = nil
				continue
			}
			criterion[name] = true
		}

		if i < 5 {
			name := fmt.Sprintf("%d<=n<=%d", i+1, 10-(i+1))
			if val, ok := criterion[name]; ok && val {
				result[i-1] = nil
				continue
			}
			criterion[name] = false
		}

		result[i-1] = criterion
	}
	return result
}

func BuildInterestingCriteria() []Criterion {
	var criteria []Criterion
	criteria = slices.Clone(TsupaiCriteria)

	for _, entry := range SupaiCriteria {
		criteria = append(criteria, entry.Base)
		criteria = slices.Concat(criteria, entry.Test)
	}

	for _, entry := range SupaiCriteria {
		criteria = slices.Concat(criteria, GetNumberCriteria(entry.Base))
		for _, testCriterion := range entry.Test {
			criteria = slices.Concat(criteria, GetNumberCriteria(testCriterion))
		}
	}

	criteria = slices.DeleteFunc(criteria, func(c Criterion) bool {
		return c == nil
	})

	return criteria
}

func CalculateInterestingProbabilities(featuresPath string, w io.Writer) (map[string]*configs.DecisionNode, error) {
	r, err := os.Open(featuresPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open features file: %w", err)
	}
	defer r.Close()

	stat, err := r.Stat()
	if err != nil {
		return nil, err
	}

	fn := FeatureNames()
	storedKyokus, err := LoadStoredKyokus(r, stat.Size(), fn)
	if err != nil {
		return nil, err
	}

	criteria := BuildInterestingCriteria()
	result, err := CalculateProbabilities(w, storedKyokus, fn, criteria)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func RunBenchmark(featuresPath string) error {
	r, err := os.Open(featuresPath)
	if err != nil {
		return fmt.Errorf("failed to open features file: %w", err)
	}
	defer r.Close()

	stat, err := r.Stat()
	if err != nil {
		return err
	}

	fn := FeatureNames()
	storedKyokus, err := LoadStoredKyokus(r, stat.Size(), fn)
	if err != nil {
		return err
	}

	criteria := BuildInterestingCriteria()
	if _, err := CreateKyokuProbsMap(storedKyokus, fn, criteria); err != nil {
		return err
	}

	return nil
}
