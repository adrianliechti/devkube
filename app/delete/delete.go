package delete

import (
	"context"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/go-cli"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "delete cluster",

		Action: func(ctx context.Context, cmd *cli.Command) error {
			provider, cluster := app.MustCluster(ctx, cmd)

			if ok := cli.MustConfirm("Are you sure you want to delete the cluster?", false); !ok {
				return nil
			}

			cli.MustRun("Deleting Kubernetes Cluster...", func() error {
				return provider.Delete(ctx, cluster)
			})

			return nil
		},
	}
}
