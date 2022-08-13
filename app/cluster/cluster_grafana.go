package cluster

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
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

			if err := kind.Kubeconfig(c.Context, name, kubeconfig); err != nil {
				return err
			}

			time.AfterFunc(3*time.Second, func() {
				url := fmt.Sprintf("http://127.0.0.1:%d", port)
				cli.OpenURL(url)
			})

			namespace := DefaultNamespace

			if err := kubectl.Invoke(c.Context, kubeconfig, "port-forward", "-n", namespace, "service/grafana", fmt.Sprintf("%d:80", port)); err != nil {
				return err
			}

			return nil
		},
	}
}
