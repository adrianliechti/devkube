package cluster

import (
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

		Action: func(c *cli.Context) error {
			name := c.String("name")

			if _, _, err := docker.Tool(c.Context); err != nil {
				return err
			}

			if _, _, err := kind.Tool(c.Context); err != nil {
				return err
			}

			if _, _, err := helm.Tool(c.Context); err != nil {
				return err
			}

			if _, _, err := kubectl.Tool(c.Context); err != nil {
				return err
			}

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

			if err := kind.Create(c.Context, name, config); err != nil {
				return err
			}

			for _, image := range append(dashboard.Images, observability.Images...) {
				docker.Pull(c.Context, image)
				kind.LoadImage(c.Context, name, image)
			}

			namespace := "loop"
			kubeconfig := ""

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

			return nil
		},
	}
}
