package main

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/configs"
)

func main() {
	configs.InitializeStats()
	fmt.Printf("%d\n", configs.Stats.NumHoras)
	fmt.Printf("%.16f\n", configs.Stats.WinProbsMap["S4,3,2"]["-100"])

	name := "mjai-manue"
	fmt.Println("Hello World!", name)
}
