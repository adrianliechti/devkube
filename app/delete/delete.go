package delete

import (
	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete cluster",

		Action: func(c *cli.Context) error {
			provider, cluster := app.MustCluster(c)

			return provider.Delete(c.Context, cluster)
		},
	}
}
