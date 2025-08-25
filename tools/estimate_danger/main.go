package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
)

type Options struct {
	Verbose bool
	Start   string
	Num     int
	Output  string
	MinGap  float64
	Filter  string
}

func parseOptions(args []string) (*Options, []string, error) {
	opts := Options{}

	fs := flag.NewFlagSet("estimate_danger", flag.ExitOnError)

	fs.BoolVar(&opts.Verbose, "v", false, "enable verbose mode")
	fs.StringVar(&opts.Start, "start", "", "start filepath")
	fs.IntVar(&opts.Num, "n", 0, "limit number of files")
	fs.StringVar(&opts.Output, "o", "", "output filepath")
	fs.Float64Var(&opts.MinGap, "min_gap", 0.0, "minimum gap percentage")
	fs.StringVar(&opts.Filter, "filter", "", "filter expression")

	if err := fs.Parse(args); err != nil {
		return nil, nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	// Convert min_gap from percent to decimal
	opts.MinGap = opts.MinGap / 100.0

	paths := fs.Args()

	return &opts, paths, nil
}

func filterInputPaths(paths []string, opts *Options) []string {
	if opts.Start != "" {
		startIndex := slices.Index(paths, opts.Start)
		if startIndex >= 0 {
			paths = paths[startIndex:]
		}
	}

	if opts.Num > 0 && len(paths) > opts.Num {
		paths = paths[:opts.Num]
	}

	return paths
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "missing action argument")
		os.Exit(2)
	}
	action := os.Args[1]
	args := os.Args[2:]

	opts, paths, err := parseOptions(args)
	if err != nil {
		log.Fatal(err)
	}

	paths = filterInputPaths(paths, opts)

	switch action {
	case "extract":
		if opts.Output == "" {
			log.Fatal("-o is missing")
		}

		var listener Listener = nil
		if opts.Filter != "" {
			listener = NewDumpListener(opts.Filter)
		}

		w := bufio.NewWriter(os.Stdout)
		defer w.Flush()

		if err := ExtractFeaturesFromFiles(paths, opts.Output, listener, opts.Verbose, w); err != nil {
			log.Fatal(err)
		}
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
