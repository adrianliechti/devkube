package create

import (
	"context"
	"slices"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/app/setup"
	"github.com/adrianliechti/devkube/extension"
	"github.com/adrianliechti/go-cli"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "create cluster",

		Action: func(ctx context.Context, cmd *cli.Command) error {
			provider, cluster, err := app.Cluster(ctx, cmd)

			if err != nil {
				return err
			}

			clusters, err := provider.List(ctx)

			if err != nil {
				return err
			}

			// creating an existing cluster fails - keep create idempotent so it
			// can be re-run to (re-)install the extensions below
			if !slices.Contains(clusters, cluster) {
				cli.MustRun("Installing Kubernetes Cluster...", func() error {
					return provider.Create(ctx, cluster)
				})
			}

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
