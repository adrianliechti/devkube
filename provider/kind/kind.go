package kind

import (
	"context"
	"os"

	"github.com/adrianliechti/devkube/pkg/kind"
	"github.com/adrianliechti/devkube/provider"
)

type Provider struct {
}

func New() provider.Provider {
	return new(Provider)
}

func (p *Provider) List(ctx context.Context) ([]string, error) {
	return kind.List(ctx)
}

func (p *Provider) Create(ctx context.Context, name string) error {
	dir, err := os.MkdirTemp("", "kubeconfig-")

	if err != nil {
		return err
	}

	path := dir + "/.config"

	defer os.RemoveAll(dir)

	if err := kind.Create(ctx, name, nil, path, kind.WithDefaultOutput()); err != nil {
		return err
	}

	return nil
}

func (p *Provider) Delete(ctx context.Context, name string) error {
	return kind.Delete(ctx, name)
}

func (p *Provider) Config(ctx context.Context, name string) ([]byte, error) {
	return kind.Config(ctx, name)
}
