package none

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/adrianliechti/devkube/provider"
)

type Provider struct {
}

func New() provider.Provider {
	return new(Provider)
}

func (p *Provider) List(ctx context.Context) ([]string, error) {
	_, err := kubeconfig()

	if err != nil {
		return nil, err
	}

	return []string{"Current"}, nil
}

func (p *Provider) Create(ctx context.Context, name string, kubeconfig string) error {
	return p.Export(ctx, name, kubeconfig)
}

func (p *Provider) Delete(ctx context.Context, name string) error {
	return nil
}

func (p *Provider) Export(ctx context.Context, name, path string) error {
	kubeconfig, err := kubeconfig()

	if err != nil {
		return err
	}

	if path == "" {
		path = kubeconfig
	}

	data, err := ioutil.ReadFile(kubeconfig)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0600)
}

func kubeconfig() (string, error) {
	path := os.Getenv("KUBECONFIG")

	if path == "" {
		dir, err := os.UserHomeDir()

		if err != nil {
			return "", err
		}

		path = filepath.Join(dir, ".kube", "config")
	}

	_, err := os.Stat(path)
	return path, err
}
