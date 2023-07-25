package protocol

import (
	"context"
)

type (
	// Client is a client implementation of the protocol
	Client struct {
		addr string
	}
)

// NewClient returns a new client
func NewClient(addr string) *Client {
	return &Client{addr: addr}
}

// ExecScenario executes the scenario of client part of the protocol and returns a quote
// otherwise returns an empty string and an error
func (c *Client) ExecScenario(ctx context.Context, scenario Scenario) (string, error) {
	// Open a connection to the server
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
