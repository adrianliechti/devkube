package cluster

import (
	"context"
	"errors"
	"os"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/provider"
	"github.com/adrianliechti/loop/pkg/kubernetes"
)

func SetupCommand() *cli.Command {
	return &cli.Command{
		Name:  "setup",
		Usage: "Setup cluster",

		Action: func(c *cli.Context) error {
			provider, cluster := app.MustCluster(c)

			return ExportConfig(c.Context, provider, cluster, "")
		},
	}
}

func ExportConfig(ctx context.Context, provider provider.Provider, cluster string, path string) error {
	if path == "" {
		path = kubernetes.ConfigPath()
	}

	if path == "" {
		return errors.New("invalid path")
	}

	data, err := provider.Config(ctx, cluster)

	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}
