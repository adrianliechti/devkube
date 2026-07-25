package connect

import (
	"context"
	"errors"
	"log/slog"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/go-cli"
	"github.com/adrianliechti/loop/pkg/catapult"
	"github.com/adrianliechti/loop/pkg/gateway"
	"github.com/adrianliechti/loop/pkg/kubernetes"
	"github.com/adrianliechti/loop/pkg/system"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "connect",
		Usage: "forward Kubernetes services",

		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:  "namespace",
				Usage: "filter namespace(s)",
			},

			&cli.StringFlag{
				Name:  "scope",
				Usage: "set namespace scope",
			},
		},

		Action: func(ctx context.Context, cmd *cli.Command) error {
			elevated, err := system.IsElevated()

			if err != nil {
				return err
			}

			if !elevated {
				cli.Fatal("This command must be run as root!")
			}

			client := app.MustClient(ctx, cmd)

			return Catapult(ctx, client, cmd.StringSlice("namespace"), cmd.String("scope"))
		},
	}
}

func Catapult(ctx context.Context, client kubernetes.Client, namespaces []string, scope string) error {
	if scope == "" && len(namespaces) > 0 {
		scope = namespaces[0]
	}

	if scope == "" {
		scope = client.Namespace()
	}

	addFunc := func(address string, hosts []string, ports []int) {
		slog.InfoContext(ctx, "adding tunnel", "address", address, "hosts", hosts, "ports", ports)
	}

	deleteFunc := func(address string, hosts []string, ports []int) {
		slog.InfoContext(ctx, "removing tunnel", "address", address, "hosts", hosts, "ports", ports)
	}

	catapult, err := catapult.New(client, catapult.CatapultOptions{
		Scope:      scope,
		Namespaces: namespaces,

		Logger: slog.Default(),

		AddFunc:    addFunc,
		DeleteFunc: deleteFunc,
	})

	if err != nil {
		return err
	}

	gateway, err := gateway.New(client, gateway.GatewayOptions{
		Namespaces: namespaces,

		Logger: slog.Default(),

		AddFunc:    addFunc,
		DeleteFunc: deleteFunc,
	})

	if err != nil {
		return err
	}

	// if one of them stops, tear down the other one as well
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errs := make(chan error, 2)

	for _, start := range []func(context.Context) error{catapult.Start, gateway.Start} {
		go func() {
			defer cancel()
			errs <- start(ctx)
		}()
	}

	return errors.Join(<-errs, <-errs)
}
