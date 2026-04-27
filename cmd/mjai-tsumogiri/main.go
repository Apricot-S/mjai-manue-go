package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/runtime"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/ai"
)

const defaultName = "tsumogiri"

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, in io.Reader, out io.Writer, errOut io.Writer) int {
	flags := flag.NewFlagSet("mjai-tsumogiri", flag.ContinueOnError)
	flags.SetOutput(errOut)
	name := flags.String("name", defaultName, "player name")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if flags.NArg() > 1 {
		fmt.Fprintln(errOut, "too many arguments")
		return 2
	}

	agent := ai.NewTsumogiriAgent()
	var err error
	if flags.NArg() == 1 {
		err = runtime.RunTCP(runtime.TCPConfig{
			Name:  *name,
			URL:   flags.Arg(0),
			Agent: agent,
		})
	} else {
		err = runtime.RunStdio(runtime.StdioConfig{
			Name:  *name,
			Room:  "default",
			Agent: agent,
			In:    in,
			Out:   out,
		})
	}
	if err != nil {
		fmt.Fprintln(errOut, err)
		if _, ok := errors.AsType[*runtime.UsageError](err); ok {
			return 2
		}
		return 1
	}
	return 0
}
