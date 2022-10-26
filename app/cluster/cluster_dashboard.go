package cluster

import (
	"fmt"
	"time"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/pkg/kubectl"
)

func DashboardCommand() *cli.Command {
	return &cli.Command{
		Name:  "dashboard",
		Usage: "Open Dashboard",

		Category: app.ManageCategory,

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

			port := app.MustPortOrRandom(c, 9090)

			time.AfterFunc(3*time.Second, func() {
				url := fmt.Sprintf("http://127.0.0.1:%d", port)
				cli.OpenURL(url)
			})

			if err := kubectl.Invoke(c.Context, []string{"port-forward", "service/dashboard", fmt.Sprintf("%d:80", port)}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(app.DefaultNamespace), kubectl.WithDefaultOutput()); err != nil {
				return err
			}

			return nil
		},
	}
}
