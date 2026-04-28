package mjairuntime

import (
	"bufio"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/adapter/mjai/outbound"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/ai"
)

type TCPConfig struct {
	Name  string
	URL   string
	Agent ai.Agent
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

	return runTCPConn(cfg.Name, endpoint.room, cfg.Agent, conn)
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

func runTCPConn(name string, room string, agent ai.Agent, conn net.Conn) error {
	r := bufio.NewScanner(conn)
	w := bufio.NewWriter(conn)
	defer w.Flush()

	driver := NewDriver(name, room, agent)
	for r.Scan() {
		line := r.Bytes()
		if len(line) == 0 {
			return fmt.Errorf("empty input line")
		}

		msg, err := inbound.ParseMessage(line)
		if err != nil {
			return err
		}
		outMsg, err := driver.Handle(msg)
		if err != nil {
			return err
		}
		if driver.Ended() {
			if err := driver.FinalizeEndGame(); err != nil {
				return err
			}
			return nil
		}
		if outMsg == nil {
			outMsg = outbound.NewNone("")
		}
		if err := writeMessage(w, outMsg); err != nil {
			return err
		}
	}
	if err := r.Err(); err != nil {
		return err
	}
	return fmt.Errorf("connection closed before end_game")
}
