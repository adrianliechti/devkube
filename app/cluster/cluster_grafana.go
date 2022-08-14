package cluster

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
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

			if name == "" {
				name = MustCluster(c.Context)
			}

			dir, err := ioutil.TempDir("", "kind")

			if err != nil {
				return err
			}

			defer os.RemoveAll(dir)
			kubeconfig := path.Join(dir, "kubeconfig")

			if err := kind.ExportConfig(c.Context, name, kubeconfig); err != nil {
				return err
			}

			time.AfterFunc(3*time.Second, func() {
				url := fmt.Sprintf("http://127.0.0.1:%d", port)
				cli.OpenURL(url)
			})

			namespace := DefaultNamespace

			if err := kubectl.Invoke(c.Context, []string{"port-forward", "service/grafana", fmt.Sprintf("%d:80", port)}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(namespace), kubectl.WithDefaultOutput()); err != nil {
				return err
			}

			return nil
		},
	}
}
