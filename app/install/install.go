package install

import (
	"errors"
	"strings"

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

			var e extension.Extension

			if c.Args().Len() > 0 {
				if c.Args().Len() > 1 {
					return errors.New("too many arguments")
				}

				for _, i := range extension.Optional {
					if strings.EqualFold(c.Args().Get(0), i.Name) {
						e = i
					}
				}
			} else {
				var labels []string

				for _, i := range extension.Optional {
					labels = append(labels, i.Title)
				}

				i, _ := cli.MustSelect("Extension", labels)
				e = extension.Optional[i]
			}

			if e.Ensure == nil {
				return errors.New("unknown extension")
			}

			cli.MustRun("Installing "+e.Title+"...", func() error {
				return e.Ensure(c.Context, client)
			})

			return nil
		},
	}
}
