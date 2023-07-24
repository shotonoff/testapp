package main

import (
	"context"

	"github.com/shotonoff/worldOfWisdom/internal/cli"
)

func main() {
	ctx := context.Background()
	cmd := cli.RootCommand()
	cmd.AddCommand(cli.ServerCommand(), cli.ClientCommand())
	err := cmd.ExecuteContext(ctx)
	if err != nil {
		return
	}
}
