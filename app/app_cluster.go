package app

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/provider"
	"github.com/adrianliechti/loop/pkg/kubernetes"
)

var ClusterFlag = &cli.StringFlag{
	Name:  "cluster",
	Usage: "cluster instance",
}

func MustCluster(ctx context.Context, cmd *cli.Command) (provider.Provider, string) {
	provider, cluster, err := Cluster(ctx, cmd)

	if err != nil {
		cli.Fatal(err)
	}

	return provider, cluster
}

func Cluster(ctx context.Context, cmd *cli.Command) (provider.Provider, string, error) {
	provider, err := Provider(ctx, cmd)

	if err != nil {
		return nil, "", err
	}

	cluster := cmd.String(ClusterFlag.Name)

	if cluster == "" {
		cluster = "devkube"
	}

	return provider, cluster, nil
}

func MustClient(ctx context.Context, cmd *cli.Command) kubernetes.Client {
	client, err := Client(ctx, cmd)

	if err != nil {
		cli.Fatal(err)
	}

	return client
}

func Client(ctx context.Context, cmd *cli.Command) (kubernetes.Client, error) {
	provider, cluster, err := Cluster(ctx, cmd)

	if err != nil {
		return nil, err
	}

	data, err := provider.Config(ctx, cluster)

	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewFromBytes(data)

	if err != nil {
		return nil, err
	}

	return client, nil
}
