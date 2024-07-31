package stop

import (
	"context"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "stop",
		Usage: "stop cluster",

		Action: func(ctx context.Context, cmd *cli.Command) error {
			provider, cluster := app.MustCluster(ctx, cmd)

			cli.MustRun("Stopping Kubernetes Cluster...", func() error {
				return provider.Stop(ctx, cluster)
			})

			return nil
		},
	}
}
