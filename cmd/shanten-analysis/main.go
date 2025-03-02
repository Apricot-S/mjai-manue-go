package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Apricot-S/mjai-manue-go/internal/game"
)

func printAnalysisResults(paiSet *game.PaiSet, shantenNumber int, goals []game.Goal) {
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	paiStr := paiSet.ToString()
	fmt.Fprintf(writer, "hand: %s\n", paiStr)
	fmt.Fprintf(writer, "shanten number: %d\n", shantenNumber)
	fmt.Fprintf(writer, "number of nearest winning hands: %d\n", len(goals))
	fmt.Fprintln(writer, "nearest winning hands: [")
	for _, goal := range goals {
		fmt.Fprintln(writer, "  [")
		fmt.Fprintf(writer, "    shanten number: %d,\n", goal.Shanten)
		fmt.Fprintln(writer, "    blocks: [")
		for _, mentsu := range goal.Mentsus {
			fmt.Fprintf(writer, "      %v,\n", mentsu.ToString())
		}
		fmt.Fprintln(writer, "    ],")
		fmt.Fprintf(writer, "    winning hand tiles count: %v,\n", goal.CountVector)
		fmt.Fprintf(writer, "    necesaary tiles count:    %v,\n", goal.RequiredVector)
		fmt.Fprintf(writer, "    unnecesaary tiles count:  %v,\n", goal.ThrowableVector)
		fmt.Fprintln(writer, "  ],")
	}
	fmt.Fprintln(writer, "]")
}

func main() {
	fmt.Print("Enter tiles (e.g., '1m 1m 1m 1m 2m 3m 4m 4m 4m 4m 1p 1p 1p 1p'): ")

	r := bufio.NewReader(os.Stdin)
	paiStr, err := r.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

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

	fmt.Println()
	printAnalysisResults(paiSet, shantenNumber, goals)
}
