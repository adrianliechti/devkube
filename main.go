package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/app/bridge"
	"github.com/adrianliechti/devkube/app/build"
	"github.com/adrianliechti/devkube/app/connect"
	"github.com/adrianliechti/devkube/app/create"
	"github.com/adrianliechti/devkube/app/delete"
	"github.com/adrianliechti/devkube/app/grafana"
	"github.com/adrianliechti/devkube/app/install"
	"github.com/adrianliechti/devkube/app/load"
	"github.com/adrianliechti/devkube/app/otel"
	"github.com/adrianliechti/devkube/app/setup"
	"github.com/adrianliechti/devkube/app/start"
	"github.com/adrianliechti/devkube/app/stop"
	"github.com/adrianliechti/go-cli"

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

	if err := app.Run(ctx, os.Args); err != nil {
		cli.Fatal(err)
	}
}

func initApp() cli.Command {
	return cli.Command{
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

			start.Command(),
			stop.Command(),

			setup.Command(),
			connect.Command(),

			install.Command(),

			bridge.Command(),
			grafana.Command(),
			otel.Command(),

			load.Command(),
			build.Command(),
		},
	}
}
