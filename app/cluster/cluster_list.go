package cluster

import (
	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/pkg/docker"
	"github.com/adrianliechti/devkube/pkg/kind"
)

func ListCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List clusters",

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
			ctx := c.Context

			provider := MustProvider(c.Context)

			list, err := provider.List(ctx)

			if err != nil {
				return err
			}

			for _, c := range list {
				cli.Info(c)
			}

			return nil
		},
	}
}
