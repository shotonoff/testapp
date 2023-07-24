package protocol

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/shotonoff/worldOfWisdom/internal/log"
	"github.com/shotonoff/worldOfWisdom/internal/qoute"
)

type ServerTestSuite struct {
	suite.Suite

	srv    *Server
	wg     sync.WaitGroup
	logger log.Logger
	client *Client

	serverEchoScenario  Scenario
	serverEmptyScenario Scenario
	clientEchoScenario  Scenario
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (suite *ServerTestSuite) SetupSuite() {
	suite.logger = log.NewNop()
	suite.client = NewClient("localhost:5001")
	suite.serverEmptyScenario = Scenario{}
	suite.serverEchoScenario = Scenario{
		func(ctx context.Context, state *ExecutionState, conn Connection) error {
			return conn.SendMsg(ctx, Echo{Data: []byte("echo")})
		},
	}
	suite.clientEchoScenario = Scenario{
		func(ctx context.Context, state *ExecutionState, conn Connection) error {
			echo := Echo{}
			return conn.ReceiveMsg(ctx, &echo)
		},
	}
}

func (suite *ServerTestSuite) TestEcho() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	suite.serverUp()
	go func() {
		defer suite.wg.Done()
		err := suite.srv.Serve(ctx, ScenarioHandler(suite.logger, suite.serverEchoScenario))
		suite.Require().NoError(err)
	}()
	_, err := suite.client.ExecScenario(ctx, suite.clientEchoScenario)
	suite.Require().NoError(err)
	suite.serverDown()
}

func (suite *ServerTestSuite) TestContextCancel() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	suite.serverUp()
	go func() {
		defer suite.wg.Done()
		err := suite.srv.Serve(ctx, ScenarioHandler(suite.logger, suite.serverEmptyScenario))
		suite.Require().Error(err)
	}()
	_, err := suite.client.ExecScenario(ctx, suite.clientEchoScenario)
	suite.Require().ErrorIs(err, context.Canceled)
	suite.serverDown()
}

func (suite *ServerTestSuite) TestContextTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	suite.serverUp(WithTimeout(10 * time.Millisecond))
	go func() {
		defer suite.wg.Done()
		err := suite.srv.Serve(ctx, ScenarioHandler(suite.logger, suite.serverEmptyScenario))
		suite.Require().Error(err)
	}()
	_, err := suite.client.ExecScenario(ctx, suite.clientEchoScenario)
	suite.Require().ErrorIs(err, context.DeadlineExceeded)
	suite.serverDown()
}

func (suite *ServerTestSuite) TestBasicScenario() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	suite.serverUp()
	quotes := []string{"custom quote"}
	quoteStore := qoute.New(qoute.WithQuotes(quotes))
	logger := log.NewNop()
	go func() {
		defer suite.wg.Done()
		err := suite.srv.Serve(ctx, ScenarioHandler(logger, ServerScenario(quoteStore)))
		suite.Require().NoError(err)
	}()

	result, err := suite.client.ExecScenario(ctx, ClientScenario())
	suite.Require().NoError(err)
	suite.Require().Equal(quotes[0], result)
	suite.serverDown()
}

func (suite *ServerTestSuite) serverUp(opts ...ServerOption) {
	var err error
	suite.srv, err = NewServer("0.0.0.0:5001", opts...)
	suite.Require().NoError(err)
	suite.wg.Add(1)
}

func (suite *ServerTestSuite) serverDown() {
	suite.srv.Stop()
	suite.wg.Wait()
}
