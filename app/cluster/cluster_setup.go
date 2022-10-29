package cluster

import (
	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
)

func SetupCommand() *cli.Command {
	return &cli.Command{
		Name:  "setup",
		Usage: "Setup cluster",

		Category: app.ClusterCategory,

		Flags: []cli.Flag{
			app.ProviderFlag,
			app.ClusterFlag,
		},

		Action: func(c *cli.Context) error {
			provider, cluster := app.MustCluster(c)

			return provider.Export(c.Context, cluster, "")
		},
	}
}
