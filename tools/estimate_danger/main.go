package main

import (
	"bufio"
	"encoding/gob"
	"flag"
	"fmt"
	"io"
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
	case "tree":
		fs.StringVar(&opts.Output, "o", "", "output filepath")
		fs.Float64Var(&opts.MinGap, "min_gap", 0.0, "minimum gap percentage")
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

func runInteresting(featuresPath string, opts *Options, w io.Writer) error {
	r, err := os.Open(featuresPath)
	if err != nil {
		return fmt.Errorf("failed to open features file: %w", err)
	}
	defer r.Close()

	stat, err := r.Stat()
	if err != nil {
		return err
	}

	fn := FeatureNames()
	criteria := BuildAllCriteria()
	criteria = slices.DeleteFunc(criteria, func(c Criterion) bool {
		return c == nil
	})
	result, err := CalculateProbabilities(r, w, stat.Size(), fn, criteria)
	if err != nil {
		return err
	}

	if opts.Output == "" {
		return nil
	}

	f, err := os.Create(opts.Output)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer f.Close()

	encoder := gob.NewEncoder(f)
	return encoder.Encode(result)
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

	switch action {
	case "extract":
		w := bufio.NewWriter(os.Stdout)
		defer w.Flush()

		if err := runExtract(paths, opts, w); err != nil {
			log.Fatal(err)
		}
	case "single":
		w := bufio.NewWriter(os.Stdout)
		defer w.Flush()

		if err := CalculateSingleProbabilities(paths[0], w); err != nil {
			log.Fatal(err)
		}
	case "interesting":
		w := bufio.NewWriter(os.Stdout)
		defer w.Flush()

		if err := runInteresting(paths[0], opts, w); err != nil {
			log.Fatal(err)
		}
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
