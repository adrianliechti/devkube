package cluster

import (
	"os"
	"path"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/pkg/docker"
	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/devkube/pkg/kind"
	"github.com/adrianliechti/devkube/pkg/kubectl"

	"github.com/adrianliechti/devkube/extension/dashboard"
	"github.com/adrianliechti/devkube/extension/metrics"
	"github.com/adrianliechti/devkube/extension/observability"
)

func CreateCommand() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create cluster",

		Flags: []cli.Flag{
			app.NameFlag,
		},

		Before: func(c *cli.Context) error {
			if _, _, err := docker.Info(c.Context); err != nil {
				return err
			}

			if _, _, err := kind.Info(c.Context); err != nil {
				return err
			}

			if _, _, err := helm.Info(c.Context); err != nil {
				return err
			}

			if _, _, err := kubectl.Info(c.Context); err != nil {
				return err
			}

			return nil
		},

		Action: func(c *cli.Context) error {
			name := c.String("name")

			provider := MustProvider(c.Context)

			if name == "" {
				name = "devkube"
			}

			dir, err := os.MkdirTemp("", "devkube")

			if err != nil {
				return err
			}

			defer os.RemoveAll(dir)

			kubeconfig := path.Join(dir, "kubeconfig")

			if err := provider.Create(c.Context, name, kubeconfig); err != nil {
				return err
			}

			if err := observability.InstallCRD(c.Context, kubeconfig, DefaultNamespace); err != nil {
				return err
			}

			if err := metrics.Install(c.Context, kubeconfig, DefaultNamespace); err != nil {
				return err
			}

			if err := dashboard.Install(c.Context, kubeconfig, DefaultNamespace); err != nil {
				return err
			}

			if err := observability.Install(c.Context, kubeconfig, DefaultNamespace); err != nil {
				return err
			}

			return provider.ExportConfig(c.Context, name, "")
		},
	}
}
