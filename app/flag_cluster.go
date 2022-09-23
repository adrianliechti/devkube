package app

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrianliechti/devkube/pkg/cli"

	"github.com/adrianliechti/devkube/provider"
	"github.com/adrianliechti/devkube/provider/aws"
	"github.com/adrianliechti/devkube/provider/azure"
	"github.com/adrianliechti/devkube/provider/digitalocean"
	"github.com/adrianliechti/devkube/provider/kind"
	"github.com/adrianliechti/devkube/provider/linode"
	"github.com/adrianliechti/devkube/provider/vultr"
)

var ClusterFlag = &cli.StringFlag{
	Name:  "cluster",
	Usage: "Cluster name",
}

func ListClusters(c *cli.Context) ([]string, error) {
	providers := map[string]provider.Provider{}

	if name := c.String(ProviderFlag.Name); name != "" {
		p, err := ProviderFromName(c.Context, name)

		if err != nil {
			return nil, err
		}

		providers[strings.ToLower(name)] = p
	} else {
		providers["local"] = kind.New()

		if p, err := aws.NewFromEnvironment(); err == nil {
			providers["aws"] = p
		}

		if p, err := azure.NewFromEnvironment(); err == nil {
			providers["azure"] = p
		}

		if p, err := digitalocean.NewFromEnvironment(); err == nil {
			providers["digitalocean"] = p
		}

		if p, err := linode.NewFromEnvironment(); err == nil {
			providers["linode"] = p
		}

		if p, err := vultr.NewFromEnvironment(); err == nil {
			providers["vultr"] = p
		}
	}

	var result []string

	for name, p := range providers {
		list, err := p.List(c.Context)

		if err != nil {
			continue
		}

		for _, c := range list {
			result = append(result, name+"/"+c)
		}
	}

	return result, nil
}

func SelectCluster(c *cli.Context) (provider.Provider, string, error) {
	list, err := ListClusters(c)

	if err != nil {
		return nil, "", err
	}

	var items []string

	filterProvider := c.String(ProviderFlag.Name)
	filterCluster := c.String(ClusterFlag.Name)

	for _, c := range list {
		pair := strings.Split(c, "/")

		if filterProvider != "" && !strings.EqualFold(pair[0], filterProvider) {
			continue
		}

		if filterCluster != "" && !strings.EqualFold(pair[1], filterCluster) {
			continue
		}

		items = append(items, c)
	}

	if len(items) == 0 {
		return nil, "", errors.New("no cluster found")
	}

	skipPrompt := len(items) == 1 && filterProvider != "" && filterCluster != ""

	if !skipPrompt {
		i, _, err := cli.Select("Select cluster", items)

		if err != nil {
			return nil, "", err
		}

		items = []string{items[i]}
	}

	pair := strings.Split(items[0], "/")

	provider, err := ProviderFromName(c.Context, pair[0])

	if err != nil {
		return nil, "", err
	}

	cluster := pair[1]

	return provider, cluster, nil
}

func MustCluster(c *cli.Context) (provider.Provider, string) {
	provider, cluster, err := SelectCluster(c)

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

	path := filepath.Join(dir, "kubeconfig")

	if err := provider.Export(c.Context, name, path); err != nil {
		closer()
		cli.Fatal(err)
	}

	return path, closer
}
