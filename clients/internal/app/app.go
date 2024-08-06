package app

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/greenblat17/yet-another-messenger/clients/internal/cli"
)

func Run(commands *cli.CLI) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	commands.SetRunning(true)
	defer commands.SetRunning(false)

	go func() {
		commands.Run(ctx)
	}()

	<-ctx.Done()
	fmt.Println("Received shutdown signal")
	commands.Close()
	fmt.Println("stopped program")
}
