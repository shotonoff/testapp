package protocol

import (
	"context"
)

type (
	// Client is a client
	Client struct {
		addr string
		conn Connection
	}
)

// NewClient returns a new client
func NewClient(addr string) *Client {
	return &Client{addr: addr}
}

// Connect connects to remote server
func (c *Client) Connect() error {
	var err error
	c.conn, err = Dial(c.addr)
	return err
}

// Disconnect disconnects from remote server
func (c *Client) Disconnect() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Send sends a message to remote server
func (c *Client) Send(ctx context.Context, msg any) error {
	data, err := Marshal(msg)
	if err != nil {
		return err
	}
	return c.conn.Write(ctx, data)
}

// Read reads a message from remote server
func (c *Client) Read(ctx context.Context, msg any) error {
	data, err := c.conn.Read(ctx)
	if err != nil {
		return err
	}
	return Unmarshal(data, msg)
}

// ExecScenario asks wisdom from remote server
func (c *Client) ExecScenario(ctx context.Context, scenario Scenario) (string, error) {
	err := c.Connect()
	if err != nil {
		return "", err
	}
	defer func() {
		_ = c.Disconnect()
	}()
	quote, err := scenario.Execute(ctx, c.conn)
	if err != nil {
		return "", err
	}
	return quote.Text, nil
}
