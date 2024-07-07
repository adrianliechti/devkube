package kind

import (
	"context"
	"os"
	"time"

	"github.com/adrianliechti/devkube/provider"

	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/log"
)

type Provider struct {
}

func New() provider.Provider {
	return new(Provider)
}

func (p *Provider) List(ctx context.Context) ([]string, error) {
	logger := log.NoopLogger{}

	provider := cluster.NewProvider(
		cluster.ProviderWithLogger(logger),
		//runtime.GetDefault(logger),
	)

	return provider.List()
}

func (p *Provider) Create(ctx context.Context, name string) error {
	dir, err := os.MkdirTemp("", "kubeconfig-")

	if err != nil {
		return err
	}

	kubeconfig := dir + "/.config"

	defer os.RemoveAll(dir)

	logger := log.NoopLogger{}

	provider := cluster.NewProvider(
		cluster.ProviderWithLogger(logger),
	)

	opts := []cluster.CreateOption{
		cluster.CreateWithWaitForReady(time.Duration(0)),
		cluster.CreateWithKubeconfigPath(kubeconfig),
	}

	return provider.Create(name, opts...)
}

func (p *Provider) Delete(ctx context.Context, name string) error {
	logger := log.NoopLogger{}

	provider := cluster.NewProvider(
		cluster.ProviderWithLogger(logger),
		//runtime.GetDefault(logger),
	)

	return provider.Delete(name, "")
}

func (p *Provider) Config(ctx context.Context, name string) ([]byte, error) {
	logger := log.NoopLogger{}

	provider := cluster.NewProvider(
		cluster.ProviderWithLogger(logger),
		//runtime.GetDefault(logger),
	)

	data, err := provider.KubeConfig(name, false)

	if err != nil {
		return nil, err
	}

	return []byte(data), nil
}
