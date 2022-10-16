package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/adrianliechti/devkube/app/cluster"
	"github.com/adrianliechti/devkube/app/feature"
	"github.com/adrianliechti/devkube/pkg/cli"
)

var version string

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	app := initApp()

	if err := app.RunContext(ctx, os.Args); err != nil {
		cli.Fatal(err)
	}
}

func initApp() cli.App {
	return cli.App{
		Usage: "DevKube",

		Suggest: true,
		Version: version,

		HideHelpCommand: true,

		Commands: []*cli.Command{
			cluster.ListCommand(),

			cluster.CreateCommand(),
			cluster.DeleteCommand(),

			cluster.SetupCommand(),
			cluster.TrustCommand(),

			cluster.RegistryCommand(),
			cluster.IngressCommand(),

			cluster.GrafanaCommand(),
			cluster.DashboardCommand(),

			feature.EnableCommand(),
			feature.DisableCommand(),
		},
	}
}
