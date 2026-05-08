package mjairuntime

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/ai"
)

type TCPConfig struct {
	Name  string
	URL   string
	Agent ai.Agent
	Log   io.Writer
}

type UsageError struct {
	err error
}

func (e *UsageError) Error() string {
	return e.err.Error()
}

func (e *UsageError) Unwrap() error {
	return e.err
}

func RunTCP(cfg TCPConfig) error {
	endpoint, err := parseMjsonpURL(cfg.URL)
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", endpoint.address)
	if err != nil {
		return err
	}
	defer conn.Close()

	return runTCPConn(cfg.Name, endpoint.room, cfg.Agent, conn, cfg.Log)
}

type mjsonpEndpoint struct {
	address string
	room    string
}

func parseMjsonpURL(rawURL string) (*mjsonpEndpoint, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, &UsageError{err: err}
	}
	if u.Scheme != "mjsonp" {
		return nil, &UsageError{err: fmt.Errorf("unsupported URL scheme %q", u.Scheme)}
	}
	if u.Host == "" {
		return nil, &UsageError{err: fmt.Errorf("mjsonp URL requires host:port")}
	}
	if u.Port() == "" {
		return nil, &UsageError{err: fmt.Errorf("mjsonp URL requires port")}
	}
	room := strings.TrimPrefix(u.EscapedPath(), "/")
	if room == "" || strings.Contains(room, "/") {
		return nil, &UsageError{err: fmt.Errorf("mjsonp URL requires room path")}
	}
	return &mjsonpEndpoint{
		address: u.Host,
		room:    room,
	}, nil
}

func runTCPConn(name string, room string, agent ai.Agent, conn net.Conn, log io.Writer) error {
	return runJSONLines(name, room, agent, conn, conn, log, jsonLinesPolicy{
		respondNoneOnNoReaction: true,
		stopOnEndGame:           true,
		errorOnEOFBeforeEndGame: true,
	})
}
