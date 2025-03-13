package client

import (
	"io"

	"github.com/Apricot-S/mjai-manue-go/internal/agent"
)

type Client struct {
	reader io.Reader
	writer io.Writer
	agent  *agent.Agent
}

func (c *Client) Run() error {
	// Dummy
	return nil
}
