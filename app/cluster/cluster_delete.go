package cluster

import (
	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/pkg/kind"
)

func DeleteCommand() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete cluster",

		Flags: []cli.Flag{
			app.NameFlag,
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
