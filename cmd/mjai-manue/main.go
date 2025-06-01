package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/Apricot-S/mjai-manue-go/internal/agent"
	"github.com/Apricot-S/mjai-manue-go/internal/ai"
	"github.com/Apricot-S/mjai-manue-go/internal/client"
)

const defaultName = "Manue020"
const defaultPort = "11600"

func getRoom(rawURL string) (string, error) {
	if rawURL == "" {
		return "", nil
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %v", err)
	}

	room := path.Base(u.Path)
	if room == "." || room == "/" {
		room = ""
	}
	return room, nil
}

func getHost(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %v", err)
	}

	host := u.Host
	if !strings.Contains(host, ":") {
		host += fmt.Sprintf(":%s", defaultPort)
	}

	return host, nil
}

func runPipeMode(agent agent.Agent) error {
	c := client.NewClient(os.Stdin, os.Stdout, agent)
	return c.Run()
}

func runTCPClientMode(rawURL string, agent agent.Agent) error {
	host, err := getHost(rawURL)
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", host)
	if err != nil {
		return fmt.Errorf("error accepting connection: %v", err)
	}
	defer conn.Close()

	fmt.Fprintf(os.Stderr, "connected: %s\n", host)

	c := client.NewClient(conn, conn, agent)
	if err := c.Run(); err != nil {
		return fmt.Errorf("client error: %v", err)
	}

	return nil
}

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

	room, err := getRoom(rawURL)
	if err != nil {
		log.Fatal(err)
	}

	ai, err := ai.NewManueAI()
	if err != nil {
		log.Fatalf("failed to create AI: %v", err)
	}
	agent := agent.NewAIAgent(name, room, ai)

	if usePipe {
		err = runPipeMode(agent)
	} else {
		err = runTCPClientMode(rawURL, agent)
	}

	if err != nil {
		log.Fatalf("error running client: %v", err)
	}
}
