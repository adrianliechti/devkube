package install

import (
	"errors"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/extension"
	"github.com/adrianliechti/devkube/pkg/cli"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "install",
		Usage: "install extension",

		Hidden: true,

		Action: func(c *cli.Context) error {
			provider, cluster, _ := app.Cluster(c)

			if provider == nil {
				return errors.New("invalid provider specified")
			}

			if cluster == "" {
				return errors.New("invalid cluster specified")
			}

			client := app.MustClient(c)

			var labels []string

			for _, i := range extension.Optional {
				labels = append(labels, i.Title)
			}

			i, _ := cli.MustSelect("Extension", labels)
			e := extension.Optional[i]

			cli.Info("★ installing " + e.Title + "...")

			if err := e.Ensure(c.Context, client); err != nil {
				return err
			}

			return nil
		},
	}
}
