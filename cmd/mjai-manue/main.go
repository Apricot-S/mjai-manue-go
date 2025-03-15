package main

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/configs"
)

func main() {
	stats, err := configs.GetStats()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d\n", stats.NumHoras)
	fmt.Printf("%d\n", stats.RyukyokuTenpaiStat.TenpaiTurnDistribution["null"])
	fmt.Printf("%.16f\n", stats.WinProbsMap["S4,3,2"]["-100"])
	fmt.Println()

	root, err := configs.GetDangerTree()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%.16f\n", root.AverageProb)
	fmt.Printf("%.16f\n", root.ConfInterval[0])
	fmt.Printf("%s\n", *root.FeatureName)
	fmt.Println()

	name := "Manue020"
	fmt.Println("Hello World!", name)
}
