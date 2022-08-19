package cluster

import (
	"context"
	"errors"
	"os"
	"path"

	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/provider"
	"github.com/adrianliechti/devkube/provider/kind"
)

const (
	DefaultNamespace = "loop"
)

func MustProvider(ctx context.Context) provider.Provider {
	p := kind.New()
	return p
}

func SelectCluster(ctx context.Context, provider provider.Provider) (string, error) {
	list, err := provider.List(ctx)

	var items []string

	if err != nil {
		return "", err
	}

	for _, c := range list {
		items = append(items, c)
	}

	if len(items) == 0 {
		return "", errors.New("no instances found")
	}

	if len(items) == 1 {
		return items[0], nil
	}

	i, _, err := cli.Select("Select instance", items)

	if err != nil {
		return "", err
	}

	return list[i], nil
}

func MustCluster(ctx context.Context, provider provider.Provider) string {
	cluster, err := SelectCluster(ctx, provider)

	if err != nil {
		cli.Fatal(err)
	}

	return cluster
}

func MustTempKubeconfig(ctx context.Context, provider provider.Provider, name string) (string, func()) {
	if name == "" {
		name = MustCluster(ctx, provider)
	}

	dir, err := os.MkdirTemp("", "devkube")

	if err != nil {
		cli.Fatal(err)
	}

	closer := func() {
		os.RemoveAll(dir)
	}

	path := path.Join(dir, "kubeconfig")

	if err := provider.ExportConfig(ctx, name, path); err != nil {
		cli.Fatal(err)
	}

	return path, closer
}
