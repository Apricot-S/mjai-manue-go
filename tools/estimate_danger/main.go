package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
)

type Options struct {
	Verbose        bool
	Start          string
	Num            int
	Output         string
	MinGap         float64
	Filter         string
	FilterSet      bool
	ExcludePlayers stringListFlag
}

type stringListFlag []string

func (f *stringListFlag) String() string {
	return fmt.Sprint([]string(*f))
}

func (f *stringListFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
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
		fs.Var(&opts.ExcludePlayers, "exclude_player", "player name to exclude; may be specified multiple times")
	case "single":
		// no options
	case "interesting":
		// not implemented yet
	case "interesting_graph":
		// not implemented yet
	case "benchmark":
		// not implemented yet
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
	fs.Visit(func(f *flag.Flag) {
		if f.Name == "filter" {
			opts.FilterSet = true
		}
	})

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

	if 0 < opts.Num && opts.Num < len(paths) {
		paths = paths[:opts.Num]
	}

	return paths
}

func runExtract(paths []string, opts *Options, w io.Writer) error {
	if opts.Output == "" {
		return fmt.Errorf("-o is missing")
	}

	paths = filterInputPaths(paths, opts)
	if len(paths) == 0 {
		return fmt.Errorf("there are no files to process")
	}

	var listener Listener = nil
	if opts.FilterSet {
		listener = NewDumpListener(opts.Filter)
	}

	return ExtractFeaturesFromFiles(paths, opts.Output, listener, opts.Verbose, w, opts.ExcludePlayers)
}

func runTree(path string, opts *Options, w io.Writer) error {
	root, err := GenerateDecisionTree(path, w, opts.MinGap)
	if err != nil {
		return err
	}
	RenderDecisionTree(w, root, "all", 0)
	if opts.Output == "" {
		return nil
	}
	return DumpDecisionTree(root, opts.Output)
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

	return DumpDecisionTreeJSON(root, opts.Output)
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

	var runErr error
	switch action {
	case "extract":
		runErr = runExtract(paths, opts, w)
	case "tree":
		runErr = runTree(paths[0], opts, w)
	case "dump_tree":
		runErr = runDumpTree(paths[0], w)
	case "dump_tree_json":
		runErr = runDumpTreeJSON(paths[0], opts)
	case "single", "interesting", "interesting_graph", "benchmark":
		runErr = fmt.Errorf("%s is not implemented yet", action)
	default:
		runErr = fmt.Errorf("unknown action: %s", action)
	}
	if runErr != nil {
		log.Fatal(runErr)
	}
}
