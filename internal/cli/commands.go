package cli

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/shotonoff/worldOfWisdom/internal/log"
	"github.com/shotonoff/worldOfWisdom/internal/protocol"
	"github.com/shotonoff/worldOfWisdom/internal/qoute"
)

// RootCommand is the root CLI command
func RootCommand() *cobra.Command {
	cmd := &cobra.Command{
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return cmd
}

// ServerCommand returns a CLI command that starts a server
func ServerCommand() *cobra.Command {
	var addr string
	cmd := &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := log.New()
			logger.Info("Starting server", "addr", addr)
			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
			srv, err := protocol.NewServer(addr)
			if err != nil {
				return err
			}
			go func() {
				<-sigs
				logger.Info("Received interrupt signal. Stopping server...")
				srv.Stop()
			}()
			logger.Info("Server started")
			err = srv.Serve(
				cmd.Context(),
				protocol.ScenarioHandler(logger, protocol.ServerScenario(qoute.New())),
			)
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&addr, "addr", "a", "0.0.0.0:5001", "listen address")
	return cmd
}

// ClientCommand returns a CLI command that runs a client
func ClientCommand() *cobra.Command {
	var addr string
	cmd := &cobra.Command{
		Use: "ask-quote",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := protocol.NewClient(addr)
			quote, err := client.ExecScenario(cmd.Context(), protocol.ClientScenario())
			if err != nil {
				return err
			}
			cmd.Println(quote)
			return nil
		},
	}
	cmd.Flags().StringVarP(&addr, "addr", "a", "localhost:5001", "server address")
	return cmd
}
