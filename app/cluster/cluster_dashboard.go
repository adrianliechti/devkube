package cluster

import (
	"fmt"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/pkg/kubernetes"
)

func DashboardCommand() *cli.Command {
	return &cli.Command{
		Name:  "dashboard",
		Usage: "Open Dashboard",

		Flags: []cli.Flag{
			app.ProviderFlag,
			app.ClusterFlag,
			app.PortFlag,
		},

		Action: func(c *cli.Context) error {
			provider, cluster := app.MustCluster(c)

			kubeconfig, closer := app.MustClusterKubeconfig(c, provider, cluster)
			defer closer()

			client, err := kubernetes.NewFromConfig(kubeconfig)

			if err != nil {
				return err
			}

			port := app.MustPortOrRandom(c, 9090)

			ports := map[int]int{
				port: 9090, // 80
			}

			ready := make(chan struct{})

			go func() {
				<-ready

				url := fmt.Sprintf("http://127.0.0.1:%d", port)
				cli.OpenURL(url)
			}()

			if err := client.ServicePortForward(c.Context, DefaultNamespace, "dashboard", "", ports, ready); err != nil {
				return err
			}

			return nil
		},
	}
}
