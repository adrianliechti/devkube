package cluster

import (
	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/pkg/kind"
)

func SetupCommand() *cli.Command {
	return &cli.Command{
		Name:  "setup",
		Usage: "Setup cluster",

		Flags: []cli.Flag{
			app.NameFlag,
		},

		Action: func(c *cli.Context) error {
			name := c.String("name")

			if name == "" {
				name = MustCluster(c.Context)
			}

			if err := kind.Kubeconfig(c.Context, name, ""); err != nil {
				return err
			}

			return nil
		},
	}
}
