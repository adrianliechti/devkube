package app

import (
	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/provider"
)

func MustCluster(c *cli.Context) (provider.Provider, string) {
	provider, cluster, err := Cluster(c)

	if err != nil {
		cli.Fatal(err)
	}

	return provider, cluster
}

func Cluster(c *cli.Context) (provider.Provider, string, error) {
	provider := MustProvider(c)

	cluster := ""

	if cluster == "" {
		cluster = "devkube"
	}

	return provider, cluster, nil
}
