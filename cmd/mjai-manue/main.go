package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/Apricot-S/mjai-manue-go/configs"
	mjairuntime "github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/runtime"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/ai"
)

const (
	defaultName = "Manue030"
	defaultSeed = uint64(0)

	exitOK           = 0
	exitRuntimeError = 1
	exitUsageError   = 2
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, in io.Reader, out io.Writer, errOut io.Writer) int {
	flags := flag.NewFlagSet("mjai-manue", flag.ContinueOnError)
	flags.SetOutput(errOut)
	name := flags.String("name", defaultName, "player name")
	seed := flags.Uint64("seed", defaultSeed, "random seed")
	if err := flags.Parse(args); err != nil {
		return exitUsageError
	}
	if flags.NArg() > 1 {
		fmt.Fprintln(errOut, "too many arguments")
		return exitUsageError
	}

	stats, err := configs.LoadGameStats()
	if err != nil {
		fmt.Fprintf(errOut, "failed to load game stats: %v\n", err)
		return exitRuntimeError
	}
	dangerTree, err := configs.LoadDangerTree()
	if err != nil {
		fmt.Fprintf(errOut, "failed to load danger tree: %v\n", err)
		return exitRuntimeError
	}
	agent, err := ai.NewManueAgent(*seed, ai.ManueAgentDeps{
		Stats:  stats,
		Danger: ai.NewDangerEstimator(dangerTree),
	})
	if err != nil {
		fmt.Fprintln(errOut, err)
		return exitRuntimeError
	}

	if flags.NArg() == 1 {
		err = mjairuntime.RunTCP(mjairuntime.TCPConfig{
			Name:  *name,
			URL:   flags.Arg(0),
			Agent: agent,
			Log:   errOut,
		})
	} else {
		err = mjairuntime.RunStdio(mjairuntime.StdioConfig{
			Name:  *name,
			Room:  "default",
			Agent: agent,
			In:    in,
			Out:   out,
			Log:   errOut,
		})
	}
	if err != nil {
		fmt.Fprintln(errOut, err)
		if _, ok := errors.AsType[*mjairuntime.UsageError](err); ok {
			return exitUsageError
		}
		return exitRuntimeError
	}
	return exitOK
}
