package cluster

import (
	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/pkg/docker"
	"github.com/adrianliechti/devkube/pkg/kind"
)

func DeleteCommand() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete cluster",

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

			// if _, _, err := helm.Info(c.Context); err != nil {
			// 	return err
			// }

			// if _, _, err := kubectl.Info(c.Context); err != nil {
			// 	return err
			// }

			return nil
		},

		Action: func(c *cli.Context) error {
			name := c.String("name")

			if name == "" {
				name = MustCluster(c.Context)
			}

			return kind.Delete(c.Context, name)
		},
	}
}
