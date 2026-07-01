package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
)

const (
	defaultPlayerName = "Manue014"
	defaultSeed       = uint64(0)

	exitOK       = 0
	exitMismatch = 1
	exitRunError = 2
)

type config struct {
	playerName string
	seed       uint64
	limit      int
	showMatch  bool
	patterns   []string
}

func parseConfig(args []string, errOut io.Writer) (config, error) {
	flags := flag.NewFlagSet("compare", flag.ContinueOnError)
	flags.SetOutput(errOut)
	playerName := flags.String("player-name", defaultPlayerName, "original player name")
	seed := flags.Uint64("seed", defaultSeed, "Go port random seed")
	limit := flags.Int("limit", 0, "maximum number of mismatches to report; 0 means unlimited")
	showMatch := flags.Bool("show-matches", false, "print matched decisions")
	if err := flags.Parse(args); err != nil {
		return config{}, err
	}
	if flags.NArg() == 0 {
		return config{}, errors.New("usage: compare [OPTIONS] <LOG_GLOB_PATTERNS>...")
	}
	if *limit < 0 {
		return config{}, errors.New("--limit must be >= 0")
	}
	return config{
		playerName: *playerName,
		seed:       *seed,
		limit:      *limit,
		showMatch:  *showMatch,
		patterns:   flags.Args(),
	}, nil
}

func globAll(patterns []string) ([]string, error) {
	var paths []string
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid glob pattern %q: %w", pattern, err)
		}
		if len(matches) == 0 {
			if _, err := os.Stat(pattern); err == nil {
				matches = []string{pattern}
			}
		}
		if len(matches) == 0 {
			return nil, fmt.Errorf("no files match %q", pattern)
		}
		paths = append(paths, matches...)
	}
	slices.Sort(paths)
	return slices.Compact(paths), nil
}
