package kind

import (
	"context"
	_ "embed"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/adrianliechti/devkube/provider"

	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type kind struct {
	provider *cluster.Provider
}

var (
	//go:embed config.yaml
	config []byte
)

func New() (provider.Provider, error) {
	logger := log.NoopLogger{}

	opts := []cluster.ProviderOption{
		cluster.ProviderWithLogger(logger),
	}

	if o, err := cluster.DetectNodeProvider(); err == nil {
		opts = append(opts, o)
	}

	return &kind{
		provider: cluster.NewProvider(opts...),
	}, nil
}

func (k *kind) List(ctx context.Context) ([]string, error) {
	return k.provider.List()
}

func (k *kind) Create(ctx context.Context, name string) error {
	dir, err := os.MkdirTemp("", "kubeconfig-")

	if err != nil {
		return err
	}

	kubeconfig := dir + "/.config"

	defer os.RemoveAll(dir)

	opts := []cluster.CreateOption{
		cluster.CreateWithRawConfig(config),
		cluster.CreateWithKubeconfigPath(kubeconfig),
		cluster.CreateWithWaitForReady(time.Duration(0)),
	}

	return k.provider.Create(name, opts...)
}

func (k *kind) Delete(ctx context.Context, name string) error {
	return k.provider.Delete(name, "")
}

func (k *kind) Start(ctx context.Context, name string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		return err
	}

	containerID, err := k.clusterContainer(ctx, cli, name)

	if err != nil {
		return err
	}

	return cli.ContainerStart(ctx, containerID, container.StartOptions{})
}

func (k *kind) Stop(ctx context.Context, name string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		return err
	}

	containerID, err := k.clusterContainer(ctx, cli, name)

	if err != nil {
		return err
	}

	return cli.ContainerStop(ctx, containerID, container.StopOptions{})
}

func (k *kind) Config(ctx context.Context, name string) ([]byte, error) {
	data, err := k.provider.KubeConfig(name, false)

	if err != nil {
		return nil, err
	}

	return []byte(data), nil
}

func (k *kind) clusterContainer(ctx context.Context, client *client.Client, name string) (string, error) {
	filter := filters.NewArgs(
		filters.KeyValuePair{
			Key:   "label",
			Value: "io.x-k8s.kind.cluster=" + strings.ToLower(name),
		},
		filters.KeyValuePair{
			Key:   "label",
			Value: "io.x-k8s.kind.role=control-plane",
		},
	)

	containers, err := client.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: filter,
	})

	if err != nil {
		return "", err
	}

	if len(containers) != 1 {
		return "", errors.New("container not found")
	}

	return containers[0].ID, nil
}
