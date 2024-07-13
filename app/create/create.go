package create

import (
	"context"
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

		Action: func(ctx context.Context, cmd *cli.Command) error {
			provider, cluster, _ := app.Cluster(ctx, cmd)

			if provider == nil {
				return errors.New("invalid provider specified")
			}

			if cluster == "" {
				return errors.New("invalid cluster specified")
			}

			cli.MustRun("Installing Kubernetes Cluster...", func() error {
				provider.Create(ctx, cluster)
				return nil
			})

			client := app.MustClient(ctx, cmd)

			for _, e := range extension.Default {
				cli.MustRun("Installing "+e.Title+"...", func() error {
					return e.Ensure(ctx, client)
				})
			}

			return setup.Export(ctx, provider, cluster, "")
		},
	}
}
