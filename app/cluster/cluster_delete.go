package cluster

import (
	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"

	"github.com/adrianliechti/devkube/extension/dashboard"
	"github.com/adrianliechti/devkube/extension/metrics"
	"github.com/adrianliechti/devkube/extension/observability"
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

			if ok, _ := cli.Confirm("Are you sure you want to delete cluster \""+cluster+"\"", false); !ok {
				return nil
			}

			kubeconfig, closer := app.MustClusterKubeconfig(c, provider, cluster)
			defer closer()

			if err := dashboard.Uninstall(c.Context, kubeconfig, DefaultNamespace); err != nil {
				//return err
			}

			if err := metrics.Uninstall(c.Context, kubeconfig, DefaultNamespace); err != nil {
				//return err
			}

			if err := observability.Uninstall(c.Context, kubeconfig, DefaultNamespace); err != nil {
				//return err
			}

			return provider.Delete(c.Context, cluster)

		},
	}
}
