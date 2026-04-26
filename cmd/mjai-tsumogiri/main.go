package main

import (
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
	if flags.NArg() > 0 {
		fmt.Fprintln(errOut, "URL mode is not implemented yet")
		return 2
	}

	if err := runtime.RunStdio(runtime.StdioConfig{
		Name:  *name,
		Room:  "default",
		Agent: ai.NewTsumogiriAgent(),
		In:    in,
		Out:   out,
	}); err != nil {
		fmt.Fprintln(errOut, err)
		return 1
	}
	return 0
}
