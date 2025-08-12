package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/game"
)

func printAnalysisResults(writer io.Writer, paiSet *base.PaiSet, shantenNumber int, goals []game.Goal) {
	w := bufio.NewWriter(writer)
	defer w.Flush()

	paiStr := paiSet.ToString()
	fmt.Fprintf(w, "hand: %s\n", paiStr)
	fmt.Fprintf(w, "shanten number: %d\n", shantenNumber)
	fmt.Fprintf(w, "number of nearest winning hands: %d\n", len(goals))
	fmt.Fprintln(w, "nearest winning hands: [")
	for _, goal := range goals {
		fmt.Fprintln(w, "  [")
		fmt.Fprintf(w, "    shanten number: %d,\n", goal.Shanten)
		fmt.Fprintln(w, "    blocks: [")
		for _, mentsu := range goal.Mentsus {
			fmt.Fprintf(w, "      %v,\n", mentsu.ToString())
		}
		fmt.Fprintln(w, "    ],")
		fmt.Fprintf(w, "    winning hand tiles count: %v,\n", goal.CountVector)
		fmt.Fprintf(w, "    necesaary tiles count:    %v,\n", goal.RequiredVector)
		fmt.Fprintf(w, "    unnecesaary tiles count:  %v,\n", goal.ThrowableVector)
		fmt.Fprintln(w, "  ],")
	}
	fmt.Fprintln(w, "]")
}

func run(reader io.Reader, writer io.Writer) error {
	fmt.Fprint(writer, "Enter tiles (e.g., '1m 1m 1m 1m 2m 3m 4m 4m 4m 4m 1p 1p 1p 1p'): ")

	r := bufio.NewReader(reader)
	paiStr, err := r.ReadString('\n')
	if err != nil {
		return err
	}

	pais, err := base.StrToPais(paiStr)
	if err != nil {
		return err
	}

	paiSet, err := base.NewPaiSet(pais)
	if err != nil {
		return err
	}

	shantenNumber, goals, err := game.AnalyzeShanten(paiSet)
	if err != nil {
		return err
	}

	fmt.Fprintln(writer)
	printAnalysisResults(writer, paiSet, shantenNumber, goals)

	return nil
}

func main() {
	if err := run(os.Stdin, os.Stdout); err != nil {
		fmt.Println(err)
	}
}
