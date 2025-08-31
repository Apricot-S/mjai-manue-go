package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"slices"

	"github.com/go-json-experiment/json"
)

type Options struct {
	Verbose bool
	Start   string
	Num     int
	Output  string
	MinGap  float64
	Filter  string
}

func parseOptions(action string, args []string) (*Options, []string, error) {
	opts := Options{}

	name := fmt.Sprintf("estimate_danger %s", action)
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	switch action {
	case "extract":
		fs.StringVar(&opts.Output, "o", "", "output filepath")
		fs.BoolVar(&opts.Verbose, "v", false, "enable verbose mode")
		fs.StringVar(&opts.Start, "start", "", "start filepath")
		fs.IntVar(&opts.Num, "n", 0, "limit number of files")
		fs.StringVar(&opts.Filter, "filter", "", "filter expression")
	case "single":
		// no options
	case "interesting":
		fs.StringVar(&opts.Output, "o", "", "output filepath")
	case "interesting_graph":
		// no options
	case "benchmark":
		// no options
	case "tree":
		fs.StringVar(&opts.Output, "o", "", "output filepath")
		fs.Float64Var(&opts.MinGap, "min_gap", 0.0, "minimum gap percentage")
	case "dump_tree":
		// no options
	case "dump_tree_json":
		fs.StringVar(&opts.Output, "o", "", "output filepath")
	default:
		return nil, nil, fmt.Errorf("unknown action: %s", action)
	}

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

func runExtract(paths []string, opts *Options, w io.Writer) error {
	paths = filterInputPaths(paths, opts)
	if len(paths) == 0 {
		return fmt.Errorf("there are no files to process")
	}

	if opts.Output == "" {
		return fmt.Errorf("-o is missing")
	}

	var listener Listener = nil
	if opts.Filter != "" {
		listener = NewDumpListener(opts.Filter)
	}

	if err := ExtractFeaturesFromFiles(paths, opts.Output, listener, opts.Verbose, w); err != nil {
		return err
	}

	return nil
}

func runInteresting(path string, opts *Options, w io.Writer) error {
	probs, err := CalculateInterestingProbabilities(path, w)
	if err != nil {
		return err
	}
	if opts.Output == "" {
		return nil
	}
	if err := DumpProbabilities(probs, opts.Output); err != nil {
		return err
	}
	return nil
}

func runTree(path string, opts *Options, w io.Writer) error {
	root, err := GenerateDecisionTree(path, w, opts.MinGap)
	if err != nil {
		log.Fatal(err)
	}
	RenderDecisionTree(w, root, "all", 0)
	if opts.Output == "" {
		return nil
	}
	if err := DumpDecisionTree(root, opts.Output); err != nil {
		return err
	}
	return nil
}

func runDumpTree(path string, w io.Writer) error {
	root, err := LoadDecisionTree(path)
	if err != nil {
		return err
	}
	RenderDecisionTree(w, root, "all", 0)
	return nil
}

func runDumpTreeJSON(path string, opts *Options) error {
	if opts.Output == "" {
		return fmt.Errorf("-o is missing")
	}

	root, err := LoadDecisionTree(path)
	if err != nil {
		return err
	}

	f, err := os.Create(opts.Output)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer f.Close()

	if err := json.MarshalWrite(f, root, json.Deterministic(true)); err != nil {
		return fmt.Errorf("failed to encode tree to JSON: %w", err)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "missing action argument")
		os.Exit(2)
	}
	action := os.Args[1]
	args := os.Args[2:]

	opts, paths, err := parseOptions(action, args)
	if err != nil {
		log.Fatal(err)
	}
	if len(paths) == 0 {
		log.Fatal("no file specified for processing")
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	switch action {
	case "extract":
		if err := runExtract(paths, opts, w); err != nil {
			log.Fatal(err)
		}
	case "single":
		if err := CalculateSingleProbabilities(paths[0], w); err != nil {
			log.Fatal(err)
		}
	case "interesting":
		if err := runInteresting(paths[0], opts, w); err != nil {
			log.Fatal(err)
		}
	case "interesting_graph":
		if err := RunInterestingGraph(paths[0]); err != nil {
			log.Fatal(err)
		}
	case "benchmark":
		if err := RunBenchmark(paths[0]); err != nil {
			log.Fatal(err)
		}
	case "tree":
		if err := runTree(paths[0], opts, w); err != nil {
			log.Fatal(err)
		}
	case "dump_tree":
		if err := runDumpTree(paths[0], w); err != nil {
			log.Fatal(err)
		}
	case "dump_tree_json":
		if err := runDumpTreeJSON(paths[0], opts); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("unknown action: %s\n", action)
	}
}
