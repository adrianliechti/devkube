package cluster

import (
	"os"
	"path"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/pkg/helm"

	"github.com/adrianliechti/devkube/extension/dashboard"
	"github.com/adrianliechti/devkube/extension/metrics"
	"github.com/adrianliechti/devkube/extension/observability"
)

func CreateCommand() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create cluster",

		Flags: []cli.Flag{
			app.ProviderFlag,
			app.ClusterFlag,
		},

		Before: func(c *cli.Context) error {
			if _, _, err := helm.Info(c.Context); err != nil {
				return err
			}

			return nil
		},

		Action: func(c *cli.Context) error {
			provider := app.MustProvider(c)

			cluster := c.String(app.ClusterFlag.Name)

			if cluster == "" {
				cluster = "devkube"
			}

			dir, err := os.MkdirTemp("", "devkube")

			if err != nil {
				return err
			}

			defer os.RemoveAll(dir)

			kubeconfig := path.Join(dir, "kubeconfig")

			if err := provider.Create(c.Context, cluster, kubeconfig); err != nil {
				return err
			}

			if err := observability.Install(c.Context, kubeconfig, DefaultNamespace); err != nil {
				return err
			}

			if err := metrics.Install(c.Context, kubeconfig, DefaultNamespace); err != nil {
				return err
			}

			if err := dashboard.Install(c.Context, kubeconfig, DefaultNamespace); err != nil {
				return err
			}

			return provider.Export(c.Context, cluster, "")
		},
	}
}
