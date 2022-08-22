package app

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/provider"
)

var ClusterFlag = &cli.StringFlag{
	Name:  "cluster",
	Usage: "Cluster name",
}

func SelectCluster(c *cli.Context, provider provider.Provider) (string, error) {
	cluster := c.String(ClusterFlag.Name)

	list, err := provider.List(c.Context)

	var items []string

	if err != nil {
		return "", err
	}

	for _, c := range list {
		if cluster != "" && !strings.EqualFold(c, cluster) {
			continue
		}

		items = append(items, c)
	}

	if len(items) == 0 {
		if cluster != "" {
			return "", fmt.Errorf("cluster %q not found", cluster)
		}

		return "", errors.New("no cluster found")
	}

	if len(items) == 1 {
		return items[0], nil
	}

	i, _, err := cli.Select("Select cluster", items)

	if err != nil {
		return "", err
	}

	return list[i], nil
}

func MustCluster(c *cli.Context) (provider.Provider, string) {
	provider := MustProvider(c)

	cluster, err := SelectCluster(c, provider)

	if err != nil {
		cli.Fatal(err)
	}

	return provider, cluster
}

func MustClusterKubeconfig(c *cli.Context, provider provider.Provider, name string) (string, func()) {
	dir, err := os.MkdirTemp("", "devkube")

	if err != nil {
		cli.Fatal(err)
	}

	closer := func() {
		os.RemoveAll(dir)
	}

	path := path.Join(dir, "kubeconfig")

	if err := provider.Export(c.Context, name, path); err != nil {
		closer()
		cli.Fatal(err)
	}

	return path, closer
}
