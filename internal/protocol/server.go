package protocol

import (
	"context"
	"github.com/shotonoff/worldOfWisdom/internal/log"
	"net"
	"sync"
)

const (
	Type       = "tcp"
	Difficulty = 3
)

type (
	// HandlerFunc is a function that handles a connection
	HandlerFunc func(context.Context, Connection) error
	// MiddlewareFunc is a function that wraps a handler
	MiddlewareFunc func(HandlerFunc) HandlerFunc
)

type (
	// Server is a tcp server
	Server struct {
		logger   log.Logger
		listener net.Listener
		quit     chan struct{}
		wg       sync.WaitGroup
	}
	// OptionFunc is a function that sets an option
	OptionFunc func(*Server)
)

// NewServer creates a new server
func NewServer(addr string, opts ...OptionFunc) (*Server, error) {
	listener, err := net.Listen(Type, addr)
	if err != nil {
		return nil, err
	}
	srv := &Server{
		logger:   log.NewNop(),
		listener: listener,
		quit:     make(chan struct{}),
	}
	for _, opt := range opts {
		opt(srv)
	}
	srv.wg.Add(1)
	return srv, nil
}

// Serve serves the server
func (s *Server) Serve(ctx context.Context, hd HandlerFunc) error {
	defer s.wg.Done()
	for {
		err := s.accept(ctx, hd)
		if err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		// stop the server if received a quit signal
		case <-s.quit:
			return nil
		default:
		}
	}
}

// Stop stops the server
func (s *Server) Stop() {
	close(s.quit)
	_ = s.listener.Close()
	s.wg.Wait()
}

func (s *Server) accept(ctx context.Context, hd HandlerFunc) error {
	conn, err := s.listener.Accept()
	if err != nil {
		return err
	}
	s.logger.Debug("Accepted connection")
	go func() {
		s.wg.Add(1)
		defer s.wg.Done()
		_ = hd(ctx, &Conn{conn})
	}()
	return nil
}

// ScenarioHandler is a handler that executes a scenario
func ScenarioHandler(logger log.Logger, scenario Scenario) HandlerFunc {
	return func(ctx context.Context, conn Connection) error {
		quote, err := scenario.Execute(ctx, conn)
		if err != nil {
			return err
		}
		logger.Info("Selected quote", "quote", quote.Text)
		return nil
	}
}
