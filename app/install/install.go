package install

import (
	"context"
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

		Action: func(ctx context.Context, cmd *cli.Command) error {
			provider, cluster, _ := app.Cluster(ctx, cmd)

			if provider == nil {
				return errors.New("invalid provider specified")
			}

			if cluster == "" {
				return errors.New("invalid cluster specified")
			}

			client := app.MustClient(ctx, cmd)

			var e extension.Extension

			if cmd.Args().Len() > 0 {
				if cmd.Args().Len() > 1 {
					return errors.New("too many arguments")
				}

				for _, i := range extension.Optional {
					if strings.EqualFold(cmd.Args().Get(0), i.Name) {
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
				return e.Ensure(ctx, client)
			})

			return nil
		},
	}
}
