package cluster

import (
	"fmt"
	"runtime"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/pkg/kubectl"
)

func RegistryCommand() *cli.Command {
	return &cli.Command{
		Name:  "registry",
		Usage: "Connect Grafana",

		Flags: []cli.Flag{
			app.ProviderFlag,
			app.ClusterFlag,
			app.PortFlag,
		},

		Before: func(c *cli.Context) error {
			if _, _, err := kubectl.Info(c.Context); err != nil {
				return err
			}

			return nil
		},

		Action: func(c *cli.Context) error {
			provider, cluster := app.MustCluster(c)

			kubeconfig, closer := app.MustClusterKubeconfig(c, provider, cluster)
			defer closer()

			port := 5000

			if runtime.GOOS == "darwin" {
				port = 5001
			}

			port = app.MustPortOrRandom(c, port)

			if err := kubectl.Invoke(c.Context, []string{"port-forward", "service/registry", fmt.Sprintf("%d:80", port)}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(DefaultNamespace), kubectl.WithDefaultOutput()); err != nil {
				return err
			}

			return nil
		},
	}
}
