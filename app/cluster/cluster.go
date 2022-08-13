package cluster

import (
	"context"
	"errors"

	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/pkg/kind"
)

const (
	DefaultNamespace = "loop"
)

func SelectCluster(ctx context.Context) (string, error) {
	list, err := kind.List(ctx)

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

func MustCluster(ctx context.Context) string {
	cluster, err := SelectCluster(ctx)

	if err != nil {
		cli.Fatal(err)
	}

	return cluster
}
