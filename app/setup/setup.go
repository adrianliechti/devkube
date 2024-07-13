package setup

import (
	"context"
	"errors"
	"os"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/pkg/kubeconfig"
	"github.com/adrianliechti/devkube/provider"
	"github.com/adrianliechti/loop/pkg/kubernetes"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "setup",
		Usage: "setup kubectl config",

		Action: func(ctx context.Context, cmd *cli.Command) error {
			provider, cluster := app.MustCluster(ctx, cmd)

			return Export(ctx, provider, cluster, "")
		},
	}
}

func Export(ctx context.Context, provider provider.Provider, cluster string, path string) error {
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

	var configs [][]byte

	if source, _ := os.ReadFile(path); source != nil {
		configs = append(configs, source)
	}

	configs = append(configs, data)

	result, err := kubeconfig.Merge(configs...)

	if err != nil {
		return os.WriteFile(path, data, 0600)
	}

	return os.WriteFile(path, result, 0600)
}
