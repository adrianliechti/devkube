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

			return nil
		},

		Action: func(c *cli.Context) error {
			name := c.String("name")

			provider := MustProvider(c.Context)

			if name == "" {
				name = MustCluster(c.Context, provider)
			}

			if ok, _ := cli.Confirm("Are you sure you want to delete cluster "+name, false); ok {
				return provider.Delete(c.Context, name)
			}

			return nil
		},
	}
}
