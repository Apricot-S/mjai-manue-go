package client

import (
	"io"

	"github.com/Apricot-S/mjai-manue-go/internal/agent"
)

type Client struct {
	reader io.Reader
	writer io.Writer
	agent  agent.Agent
}

func NewClient(reader io.Reader, writer io.Writer, agent agent.Agent) *Client {
	return &Client{reader, writer, agent}
}

func (c *Client) Run() error {
	return nil
}
