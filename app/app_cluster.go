package app

import (
	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/provider"
	"github.com/adrianliechti/loop/pkg/kubernetes"
)

var ClusterFlag = &cli.StringFlag{
	Name:  "cluster",
	Usage: "cluster instance",
}

func MustCluster(c *cli.Context) (provider.Provider, string) {
	provider, cluster, err := Cluster(c)

	if err != nil {
		cli.Fatal(err)
	}

	return provider, cluster
}

func Cluster(c *cli.Context) (provider.Provider, string, error) {
	provider, err := Provider(c)

	if err != nil {
		return nil, "", err
	}

	cluster := c.String(ClusterFlag.Name)

	if cluster == "" {
		cluster = "devkube"
	}

	return provider, cluster, nil
}

func MustClient(c *cli.Context) kubernetes.Client {
	client, err := Client(c)

	if err != nil {
		cli.Fatal(err)
	}

	return client
}

func Client(c *cli.Context) (kubernetes.Client, error) {
	provider, cluster, err := Cluster(c)

	if err != nil {
		return nil, err
	}

	data, err := provider.Config(c.Context, cluster)

	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewFromBytes(data)

	if err != nil {
		return nil, err
	}

	return client, nil
}
