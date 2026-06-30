package main

import (
	"fmt"
	"io"
	"os"

	"github.com/Apricot-S/mjai-manue-go/configs"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/ai"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, out io.Writer, errOut io.Writer) int {
	cfg, err := parseConfig(args, errOut)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return exitRunError
	}

	paths, err := globAll(cfg.patterns)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return exitRunError
	}

	stats, err := configs.LoadGameStats()
	if err != nil {
		fmt.Fprintf(errOut, "failed to load game stats: %v\n", err)
		return exitRunError
	}
	dangerTree, err := configs.LoadDangerTree()
	if err != nil {
		fmt.Fprintf(errOut, "failed to load danger tree: %v\n", err)
		return exitRunError
	}

	c := comparer{
		cfg: cfg,
		deps: ai.ManueAgentDeps{
			Stats:  stats,
			Danger: ai.NewDangerEstimator(dangerTree),
		},
		out: out,
		log: errOut,
	}

	s := summary{files: len(paths)}
	for _, path := range paths {
		fileSummary, err := c.compareFile(path)
		s.add(fileSummary)
		if err != nil {
			s.errors++
			fmt.Fprintf(out, "error: %s: %v\n", path, err)
		}
	}
	fmt.Fprintf(out, "summary: files=%d decisions=%d matches=%d mismatches=%d errors=%d\n",
		s.files, s.decisions, s.matches, s.mismatches, s.errors)
	if s.errors > 0 {
		return exitRunError
	}
	if s.mismatches > 0 {
		return exitMismatch
	}
	return exitOK
}
