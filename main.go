package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/app/connect"
	"github.com/adrianliechti/devkube/app/create"
	"github.com/adrianliechti/devkube/app/delete"
	"github.com/adrianliechti/devkube/app/grafana"
	"github.com/adrianliechti/devkube/app/load"
	"github.com/adrianliechti/devkube/app/logs"
	"github.com/adrianliechti/devkube/app/setup"
	"github.com/adrianliechti/devkube/pkg/cli"

	"github.com/lmittmann/tint"
)

var version string

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		Level:      slog.LevelInfo,
		TimeFormat: time.Kitchen,
	})))

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

		Flags: []cli.Flag{
			app.ProviderFlag,
			app.ClusterFlag,
		},

		Commands: []*cli.Command{
			create.Command(),
			delete.Command(),

			setup.Command(),
			connect.Command(),
			grafana.Command(),

			load.Command(),
			logs.Command(),
		},
	}
}
