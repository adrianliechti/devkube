package cluster

import (
	"io/ioutil"
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

			if name == "" {
				name = "devkube"
			}

			dir, err := ioutil.TempDir("", "kind")

			if err != nil {
				return err
			}

			defer os.RemoveAll(dir)

			config := map[string]any{
				"kind":       "Cluster",
				"apiVersion": "kind.x-k8s.io/v1alpha4",

				"kubeadmConfigPatches": []string{
					`kind: ClusterConfiguration
controllerManager:
  extraArgs:
    bind-address: 0.0.0.0
scheduler:
  extraArgs:
    bind-address: 0.0.0.0
etcd:
  local:
    extraArgs:
      listen-metrics-urls: http://0.0.0.0:2381
`,
					`kind: KubeProxyConfiguration
metricsBindAddress: 0.0.0.0
`,
				},
			}

			kubeconfig := path.Join(dir, "kubeconfig")
			namespace := DefaultNamespace

			if err := kind.Create(c.Context, name, config, kubeconfig); err != nil {
				return err
			}

			if err := observability.InstallCRD(c.Context, kubeconfig, namespace); err != nil {
				return err
			}

			if err := metrics.Install(c.Context, kubeconfig, namespace); err != nil {
				return err
			}

			if err := dashboard.Install(c.Context, kubeconfig, namespace); err != nil {
				return err
			}

			if err := observability.Install(c.Context, kubeconfig, namespace); err != nil {
				return err
			}

			if err := kind.Kubeconfig(c.Context, name, ""); err != nil {
				return err
			}

			return nil
		},
	}
}
