package cluster

import (
	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
)

func ListCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List clusters",

		Flags: []cli.Flag{
			app.ProviderFlag,
		},

		Action: func(c *cli.Context) error {
			list, err := app.ListClusters(c)

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
