package delete

import (
	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "delete cluster",

		Action: func(c *cli.Context) error {
			provider, cluster := app.MustCluster(c)

			if ok := cli.MustConfirm("Are you sure you want to delete the cluster?", false); !ok {
				return nil
			}

			cli.MustRun("Deleting Kubernetes Cluster...", func() error {
				return provider.Delete(c.Context, cluster)
			})

			return nil
		},
	}
}
