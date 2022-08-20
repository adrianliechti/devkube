package cluster

import (
	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
)

func ListCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List clusters",

		Action: func(c *cli.Context) error {
			provider := app.MustProvider(c)

			list, err := provider.List(c.Context)

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
