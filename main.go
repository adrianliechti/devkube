package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/adrianliechti/devkube/app/cluster"
	"github.com/adrianliechti/devkube/pkg/cli"
)

var version string

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGPIPE)
	defer stop()

	app := initApp()

	if err := app.RunContext(ctx, os.Args); err != nil {
		cli.Fatal(err)
	}
}

func initApp() cli.App {
	return cli.App{
		Version: version,
		// Usage:   "DevOps Loop",

		HideHelpCommand: true,

		Commands: []*cli.Command{
			cluster.ListCommand(),

			cluster.CreateCommand(),
			cluster.DeleteCommand(),
			cluster.SetupCommand(),

			cluster.GrafanaCommand(),
			cluster.DashboardCommand(),
		},
	}
}
