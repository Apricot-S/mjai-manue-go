package cli

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/Apricot-S/mjai-manue-go/internal/agent"
	"github.com/Apricot-S/mjai-manue-go/internal/client"
)

const defaultPort = "11600"

func GetRoom(rawURL string) (string, error) {
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

func RunPipeMode(agent agent.Agent) error {
	c := client.NewMjaiClient(os.Stdin, os.Stdout, agent)
	return c.Run()
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

func RunTCPClientMode(rawURL string, agent agent.Agent) error {
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

	c := client.NewMjaiClient(conn, conn, agent)
	if err := c.Run(); err != nil {
		return fmt.Errorf("client error: %v", err)
	}

	return nil
}
