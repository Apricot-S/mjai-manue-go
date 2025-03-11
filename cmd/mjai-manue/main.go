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

	name := "mjai-manue"
	fmt.Println("Hello World!", name)
}
