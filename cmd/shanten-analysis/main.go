package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Apricot-S/mjai-manue-go/internal/game"
)

func printShantenAnalysis(paiStr string) {
	pais, err := game.StrToPais(paiStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	paiSet, err := game.NewPaiSetWithPais(&pais)
	if err != nil {
		fmt.Println(err)
		return
	}

	shantenNumber, goals, err := game.AnalyzeShanten(paiSet)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("hand: %s\n", paiStr)
	fmt.Printf("shanten number: %d\n", shantenNumber)
	fmt.Printf("number of goals: %d\n", len(goals))
	fmt.Println("goals: [")
	for _, goal := range goals {
		fmt.Println("  [")
		fmt.Printf("    shanten number: %d,\n", goal.Shanten)
		fmt.Println("    blocks: [")
		for _, mentsu := range goal.Mentsus {
			fmt.Printf("      %v,\n", mentsu.ToString())
		}
		fmt.Println("    ],")
		fmt.Printf("    winning form tiles count: %v,\n", goal.CountVector)
		fmt.Printf("    necesaary tiles count:    %v,\n", goal.RequiredVector)
		fmt.Printf("    unnecesaary tiles count:  %v,\n", goal.ThrowableVector)
		fmt.Println("  ],")
	}
	fmt.Println("]")
}

func main() {
	fmt.Print("Enter tiles (e.g., '1m 1m 1m 1m 2m 3m 4m 4m 4m 4m 1p 1p 1p 1p'): ")

	r := bufio.NewReader(os.Stdin)
	paiStr, err := r.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println()
	printShantenAnalysis(paiStr)
}
