package connect

import (
	"context"
	"errors"
	"log/slog"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
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

			var tunnelScope string = "default"
			var tunnelNamespaces []string = nil

			if val := cmd.StringSlice("namespace"); len(val) > 0 {
				tunnelScope = val[0]
				tunnelNamespaces = val
			}

			if val := cmd.String("scope"); val != "" {
				tunnelScope = val
			}

			return Catapult(ctx, client, tunnelNamespaces, tunnelScope)
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

	catapult, err := catapult.New(client, catapult.CatapultOptions{
		Scope:      scope,
		Namespaces: namespaces,

		Logger: slog.Default(),

		AddFunc: func(address string, hosts []string, ports []int) {
			slog.InfoContext(ctx, "adding tunnel", "address", address, "hosts", hosts, "ports", ports)
		},

		DeleteFunc: func(address string, hosts []string, ports []int) {
			slog.InfoContext(ctx, "removing tunnel", "address", address, "hosts", hosts, "ports", ports)
		},
	})

	if err != nil {
		return err
	}

	gateway, err := gateway.New(client, gateway.GatewayOptions{
		Namespaces: namespaces,

		Logger: slog.Default(),

		AddFunc: func(address string, hosts []string, ports []int) {
			slog.InfoContext(ctx, "adding tunnel", "address", address, "hosts", hosts, "ports", ports)
		},

		DeleteFunc: func(address string, hosts []string, ports []int) {
			slog.InfoContext(ctx, "removing tunnel", "address", address, "hosts", hosts, "ports", ports)
		},
	})

	if err != nil {
		return err
	}

	err1 := make(chan error)
	err2 := make(chan error)

	go func() {
		err1 <- catapult.Start(ctx)
	}()

	go func() {
		err2 <- gateway.Start(ctx)
	}()

	return errors.Join(<-err1, <-err2)
}
