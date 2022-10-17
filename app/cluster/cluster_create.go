package cluster

import (
	"os"
	"path/filepath"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/devkube/pkg/kubectl"

	"github.com/adrianliechti/devkube/extension/certmanager"
	"github.com/adrianliechti/devkube/extension/dashboard"
	"github.com/adrianliechti/devkube/extension/ingress"
	"github.com/adrianliechti/devkube/extension/metrics"
	"github.com/adrianliechti/devkube/extension/observability"
	"github.com/adrianliechti/devkube/extension/registry"
)

func CreateCommand() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create cluster",

		Category: app.ClusterCategory,

		Flags: []cli.Flag{
			app.ProviderFlag,
			app.ClusterFlag,
		},

		Before: func(c *cli.Context) error {
			if _, _, err := helm.Info(c.Context); err != nil {
				return err
			}

			if _, _, err := kubectl.Info(c.Context); err != nil {
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

			kubeconfig := filepath.Join(dir, "kubeconfig")

			if err := provider.Create(c.Context, cluster, kubeconfig); err != nil {
				return err
			}

			kubectl.Invoke(c.Context, []string{"create", "namespace", app.DefaultNamespace}, kubectl.WithKubeconfig(kubeconfig))

			if err := observability.InstallCRD(c.Context, kubeconfig, app.DefaultNamespace); err != nil {
				return err
			}

			if err := certmanager.Install(c.Context, kubeconfig, app.DefaultNamespace); err != nil {
				return err
			}

			if err := metrics.Install(c.Context, kubeconfig, app.DefaultNamespace); err != nil {
				return err
			}

			if err := dashboard.Install(c.Context, kubeconfig, app.DefaultNamespace); err != nil {
				return err
			}

			if err := registry.Install(c.Context, kubeconfig, app.DefaultNamespace); err != nil {
				return err
			}

			if err := ingress.Install(c.Context, kubeconfig, app.DefaultNamespace); err != nil {
				return err
			}

			if err := observability.Install(c.Context, kubeconfig, app.DefaultNamespace); err != nil {
				return err
			}

			return provider.Export(c.Context, cluster, "")
		},
	}
}
