package start

import (
	"context"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/go-cli"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "start",
		Usage: "start cluster",

		Action: func(ctx context.Context, cmd *cli.Command) error {
			provider, cluster := app.MustCluster(ctx, cmd)

			cli.MustRun("Starting Kubernetes Cluster...", func() error {
				return provider.Start(ctx, cluster)
			})

			return nil
		},
	}
}
