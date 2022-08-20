package cluster

import (
	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
)

func DeleteCommand() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete cluster",

		Flags: []cli.Flag{
			app.ProviderFlag,
			app.ClusterFlag,
		},

		Action: func(c *cli.Context) error {
			provider, cluster := app.MustCluster(c)

			if ok, _ := cli.Confirm("Are you sure you want to delete cluster \""+cluster+"\"", false); ok {
				return provider.Delete(c.Context, cluster)
			}

			return nil
		},
	}
}
