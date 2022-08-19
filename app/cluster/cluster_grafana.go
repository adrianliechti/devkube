package cluster

import (
	"fmt"
	"time"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/pkg/docker"
	"github.com/adrianliechti/devkube/pkg/kind"
	"github.com/adrianliechti/devkube/pkg/kubectl"
)

func GrafanaCommand() *cli.Command {
	return &cli.Command{
		Name:  "grafana",
		Usage: "Open Grafana",

		Flags: []cli.Flag{
			app.NameFlag,
			app.PortFlag,
		},

		Before: func(c *cli.Context) error {
			if _, _, err := docker.Info(c.Context); err != nil {
				return err
			}

			if _, _, err := kind.Info(c.Context); err != nil {
				return err
			}

			if _, _, err := kubectl.Info(c.Context); err != nil {
				return err
			}

			return nil
		},

		Action: func(c *cli.Context) error {
			port := app.MustPortOrRandom(c, 3000)
			name := c.String("name")

			provider := MustProvider(c.Context)

			kubeconfig, closer := MustTempKubeconfig(c.Context, provider, name)
			defer closer()

			time.AfterFunc(3*time.Second, func() {
				url := fmt.Sprintf("http://127.0.0.1:%d", port)
				cli.OpenURL(url)
			})

			if err := kubectl.Invoke(c.Context, []string{"port-forward", "service/grafana", fmt.Sprintf("%d:80", port)}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(DefaultNamespace), kubectl.WithDefaultOutput()); err != nil {
				return err
			}

			return nil
		},
	}
}
