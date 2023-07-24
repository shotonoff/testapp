package protocol

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"math/rand"
	"net"

	"github.com/shotonoff/worldOfWisdom/internal/pow"
	"github.com/shotonoff/worldOfWisdom/internal/qoute"
)

func init() {
	gob.Register(Challenge{})
	gob.Register(Solution{})
	gob.Register(Quote{})
}

type (
	// Connection is an interface for a transport layer
	Connection interface {
		Writer
		Reader
		Closer
	}
	// Scenario is a list of tasks
	Scenario []TaskFunc
	// ExecutionState is a task execution context
	ExecutionState struct {
		Challenge Challenge
		Solution  Solution
		Quote     Quote
	}
	// TaskFunc is a task function that handles a connection
	TaskFunc func(context.Context, *ExecutionState, Connection) error
	// Challenge is a message sent from the server to the client
	Challenge struct {
		Difficulty int
		Data       []byte
	}
	// Solution is a message sent from the client to the server
	Solution struct {
		Nonce int
		Hash  []byte
	}
	// Quote is a message sent from the server to the client
	Quote struct {
		Text string
	}
	Echo struct {
		Data []byte
	}
	// Reader is an interface for a reader
	Reader interface {
		Read(ctx context.Context) ([]byte, error)
	}
	// Writer is an interface for a writer
	Writer interface {
		Write(ctx context.Context, data []byte) error
	}
	Closer interface {
		Close() error
	}
)

var (
	ErrInvalidProof = errors.New("invalid proof")
)

// Execute executes a scenario
func (s Scenario) Execute(ctx context.Context, conn Connection) (Quote, error) {
	var execState ExecutionState
	for _, task := range s {
		err := task(ctx, &execState, conn)
		if err != nil {
			return Quote{}, err
		}
	}
	return execState.Quote, nil
}

// ClientScenario is a scenario execution for a client
func ClientScenario() Scenario {
	return Scenario{
		receiveChallenge,
		solveChallenge,
		receiveQuote,
	}
}

// ServerScenario is a scenario execution for a server
func ServerScenario(quote *qoute.Store) Scenario {
	return Scenario{
		generateChallenge,
		receiveSolution,
		sendQuote(quote),
	}
}

// Marshal serializes a message into a byte slice
func Marshal(msg any) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(msg)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Unmarshal deserializes a message from a byte slice
func Unmarshal(data []byte, msg any) error {
	buf := bytes.NewBuffer(data)
	return gob.NewDecoder(buf).Decode(msg)
}

type Conn struct {
	conn net.Conn
}

// Read reads a message from a connection
func (c *Conn) Read(ctx context.Context) ([]byte, error) {
	header := make([]byte, 4)
	_, err := c.conn.Read(header)
	if err != nil {
		return nil, err
	}
	msgSize := int(binary.BigEndian.Uint32(header))
	buf := bytes.NewBuffer(make([]byte, 0, msgSize))
	tmp := make([]byte, 256)
	totalRead := 0
	for totalRead < msgSize {
		n, err := c.conn.Read(tmp)
		if err != nil {
			return nil, err
		}
		totalRead += n
		buf.Write(tmp[:n])
	}
	return buf.Bytes(), nil
}

// Write writes a message to a connection
func (c *Conn) Write(ctx context.Context, data []byte) error {
	msgSize := len(data)
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(msgSize))
	_, err := c.conn.Write(header)
	if err != nil {
		return err
	}
	totalWrite := 0
	for totalWrite < msgSize {
		n, err := c.conn.Write(data)
		if err != nil {
			return err
		}
		totalWrite += n
		data = data[n:]
	}
	return nil
}

// Close closes the connection
func (c *Conn) Close() error {
	return c.conn.Close()
}

// Dial dials a connection
func Dial(addr string) (Connection, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	return &Conn{conn: conn}, nil
}

func receiveChallenge(ctx context.Context, execState *ExecutionState, conn Connection) error {
	data, err := conn.Read(ctx)
	if err != nil {
		return err
	}
	return Unmarshal(data, &execState.Challenge)
}

func solveChallenge(ctx context.Context, execState *ExecutionState, conn Connection) error {
	nonce, hash, err := pow.Compute(execState.Challenge.Data, execState.Challenge.Difficulty)
	if err != nil {
		return err
	}
	execState.Solution.Nonce = nonce
	execState.Solution.Hash = hash
	data, err := Marshal(execState.Solution)
	if err != nil {
		return err
	}
	return conn.Write(ctx, data)
}

func receiveQuote(ctx context.Context, execState *ExecutionState, conn Connection) error {
	data, err := conn.Read(ctx)
	if err != nil {
		return err
	}
	return Unmarshal(data, &execState.Quote)
}

func generateChallenge(ctx context.Context, execState *ExecutionState, conn Connection) error {
	challenge := make([]byte, 8)
	binary.BigEndian.PutUint64(challenge, uint64(rand.Int63()))
	execState.Challenge = Challenge{
		Difficulty: Difficulty,
		Data:       challenge,
	}
	data, err := Marshal(execState.Challenge)
	if err != nil {
		return err
	}
	return conn.Write(ctx, data)
}

func receiveSolution(ctx context.Context, execState *ExecutionState, conn Connection) error {
	data, err := conn.Read(ctx)
	if err != nil {
		return err
	}
	err = Unmarshal(data, &execState.Solution)
	if err != nil {
		return err
	}
	if !pow.Verify(
		execState.Challenge.Data,
		execState.Solution.Hash,
		execState.Challenge.Difficulty,
		execState.Solution.Nonce,
	) {
		return ErrInvalidProof
	}
	return nil
}

func sendQuote(quote *qoute.Store) TaskFunc {
	return func(ctx context.Context, execState *ExecutionState, conn Connection) error {
		execState.Quote.Text = quote.Random()
		data, err := Marshal(execState.Quote)
		if err != nil {
			return err
		}
		return conn.Write(ctx, data)
	}
}
