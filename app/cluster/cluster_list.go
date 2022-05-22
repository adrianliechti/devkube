package cluster

import (
	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/pkg/kind"
)

func ListCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List clusters",

		Action: func(c *cli.Context) error {
			ctx := c.Context

			list, err := kind.List(ctx)

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
