package protocol

import (
	"context"
)

type (
	// Client is a client
	Client struct {
		addr string
	}
)

// NewClient returns a new client
func NewClient(addr string) *Client {
	return &Client{addr: addr}
}

// ExecScenario asks wisdom from remote server
func (c *Client) ExecScenario(ctx context.Context, scenario Scenario) (string, error) {
	conn, err := Dial(c.addr)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = conn.Close()
	}()
	quote, err := scenario.Execute(ctx, conn)
	if err != nil {
		return "", err
	}
	return quote.Text, nil
}
