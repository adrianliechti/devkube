package cluster

import (
	"fmt"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/pkg/kubernetes"
)

func GrafanaCommand() *cli.Command {
	return &cli.Command{
		Name:  "grafana",
		Usage: "Open Grafana",

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

			port := app.MustPortOrRandom(c, 3000)

			ports := map[int]int{
				port: 3000, // 80
			}

			ready := make(chan struct{})

			go func() {
				<-ready

				url := fmt.Sprintf("http://127.0.0.1:%d", port)
				cli.OpenURL(url)
			}()

			if err := client.ServicePortForward(c.Context, DefaultNamespace, "grafana", "", ports, ready); err != nil {
				return err
			}

			return nil
		},
	}
}
