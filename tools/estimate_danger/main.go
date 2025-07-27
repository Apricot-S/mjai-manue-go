package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "missing action argument")
		os.Exit(2)
	}
	action := os.Args[1]
	args := os.Args[2:]

	// サブコマンド共通のオプション群
	fs := flag.NewFlagSet(action, flag.ExitOnError)
	var (
		verbose bool
		start   string
		n       int
		o       string
		minGap  float64
		filter  string
	)
	fs.BoolVar(&verbose, "v", false, "verbose mode")
	fs.StringVar(&start, "start", "", "start filename")
	fs.IntVar(&n, "n", 0, "limit number of files")
	fs.StringVar(&o, "o", "", "output path")
	fs.Float64Var(&minGap, "min_gap", 0.0, "minimum gap percentage")
	fs.StringVar(&filter, "filter", "", "filter expression")

	if err := fs.Parse(args); err != nil {
		log.Fatalf("failed to parse flags: %v\n", err)
	}

	switch action {
	case "extract":
		panic("extract not implemented")
	case "single":
		panic("single not implemented")
	case "interesting":
		panic("interesting not implemented")
	case "interesting_graph":
		panic("interesting_graph not implemented")
	case "benchmark":
		panic("benchmark not implemented")
	case "tree":
		panic("tree not implemented")
	case "dump_tree":
		panic("dump_tree not implemented")
	case "dump_tree_json":
		panic("dump_tree_json not implemented")
	default:
		log.Fatalf("unknown action: %s\n", action)
	}
}
