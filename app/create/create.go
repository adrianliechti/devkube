package create

import (
	"errors"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/app/setup"
	"github.com/adrianliechti/devkube/extension"
	"github.com/adrianliechti/devkube/pkg/cli"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "create cluster",

		Action: func(c *cli.Context) error {
			provider, cluster, _ := app.Cluster(c)

			if provider == nil {
				return errors.New("invalid provider specified")
			}

			if cluster == "" {
				return errors.New("invalid cluster specified")
			}

			cli.MustRun("Installing Kubernetes Cluster...", func() error {
				provider.Create(c.Context, cluster)
				return nil
			})

			client := app.MustClient(c)

			for _, e := range extension.Default {
				cli.MustRun("Installing "+e.Title+"...", func() error {
					return e.Ensure(c.Context, client)
				})
			}

			return setup.Export(c.Context, provider, cluster, "")
		},
	}
}
