package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/Apricot-S/mjai-manue-go/internal/agent"
	"github.com/Apricot-S/mjai-manue-go/internal/client"
)

func parseOptions() (name string, room string, rawURL string, usePipe bool) {
	flag.StringVar(&name, "name", "Tsumogiri", "Player's name")
	flag.StringVar(&room, "room", "", "Room name (overrides the last path segment of URL if specified)")
	flag.StringVar(&rawURL, "url", "", "Server URL (e.g., http://localhost:11600/default)")
	flag.BoolVar(&usePipe, "pipe", false, "Use pipe instead of HTTP (ignore --url if specified)")
	flag.Parse()
	return
}

func getRoomName(rawURL, room string) (string, error) {
	if room != "" {
		return room, nil
	}

	if rawURL == "" {
		return room, nil
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %v", err)
	}

	room = path.Base(u.Path)
	if room == "." || room == "/" {
		room = ""
	}
	return room, nil
}

func getHostAndPort(rawURL string) (string, string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", "", fmt.Errorf("invalid URL: %v", err)
	}

	port := "11600"
	host := u.Host

	if h, p, err := net.SplitHostPort(u.Host); err == nil {
		host = h
		port = p
	}

	return host, port, nil
}

func runPipeMode(agent agent.Agent) error {
	c := client.NewClient(os.Stdin, os.Stdout, true, agent)
	return c.Run()
}

func runServerMode(rawURL string, agent agent.Agent) error {
	_, port, err := getHostAndPort(rawURL)
	if err != nil {
		return err
	}

	http.HandleFunc(rawURL, func(w http.ResponseWriter, r *http.Request) {
		c := client.NewClient(r.Body, w, false, agent)
		if err := c.Run(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	addr := ":" + port
	fmt.Fprintf(os.Stderr, "Server is running at %s (port: %s)\n", rawURL, port)
	return http.ListenAndServe(addr, nil)
}

func main() {
	name, room, rawURL, usePipe := parseOptions()

	if !usePipe && rawURL == "" {
		log.Fatal("specify --url or --pipe")
	}

	room, err := getRoomName(rawURL, room)
	if err != nil {
		log.Fatal(err)
	}

	agent := agent.NewTsumogiriAgent(name, room)

	if usePipe {
		if err := runPipeMode(agent); err != nil {
			log.Fatalf("error running client: %v", err)
		}
	} else {
		if err := runServerMode(rawURL, agent); err != nil {
			log.Fatalf("error starting server: %v", err)
		}
	}
}
