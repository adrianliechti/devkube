package none

import (
	"context"
	"os"
	"path/filepath"

	"github.com/adrianliechti/devkube/provider"
)

type Provider struct {
	kubeconfig string
}

func New(kubeconfig string) provider.Provider {
	return &Provider{
		kubeconfig: kubeconfig,
	}
}

func NewFromEnvironment() (provider.Provider, error) {
	path := os.Getenv("KUBECONFIG")

	if path == "" {
		dir, err := os.UserHomeDir()

		if err != nil {
			return nil, err
		}

		path = filepath.Join(dir, ".kube", "config")
	}

	if _, err := os.Stat(path); err != nil {
		return nil, err
	}

	return New(path), nil
}

func (p *Provider) List(ctx context.Context) ([]string, error) {
	if _, err := os.Stat(p.kubeconfig); err != nil {
		return nil, err
	}

	return []string{"current"}, nil
}

func (p *Provider) Create(ctx context.Context, name string) error {
	return nil
}

func (p *Provider) Delete(ctx context.Context, name string) error {
	return nil
}

func (p *Provider) Config(ctx context.Context, name string) ([]byte, error) {
	data, err := os.ReadFile(p.kubeconfig)

	if err != nil {
		return nil, err
	}

	return data, nil
}
