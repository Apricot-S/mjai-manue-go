package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Apricot-S/mjai-manue-go/internal/agent"
	"github.com/Apricot-S/mjai-manue-go/internal/ai"
	"github.com/Apricot-S/mjai-manue-go/internal/cli"
)

const defaultName = "Manue020"

func main() {
	var name string
	flag.StringVar(&name, "name", defaultName, "Player's name")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [--name NAME] [url]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	var rawURL string
	var usePipe bool
	args := flag.Args()
	switch len(args) {
	case 0:
		usePipe = true
	case 1:
		rawURL = args[0]
		usePipe = false
	default:
		flag.Usage()
		os.Exit(2)
	}

	room, err := cli.GetRoom(rawURL)
	if err != nil {
		log.Fatal(err)
	}

	ai, err := ai.NewManueAIDefault()
	if err != nil {
		log.Fatalf("failed to create AI: %v", err)
	}
	agent := agent.NewAIAgentDefault(name, room, ai)

	if usePipe {
		err = cli.RunPipeMode(agent)
	} else {
		err = cli.RunTCPClientMode(rawURL, agent)
	}

	if err != nil {
		log.Fatalf("error running client: %v", err)
	}
}
