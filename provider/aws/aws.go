package aws

import (
	"context"
	"errors"
	"os"

	"github.com/adrianliechti/devkube/pkg/eksctl"
	"github.com/adrianliechti/devkube/provider"
)

type Provider struct {
}

func New() provider.Provider {
	return &Provider{}
}

func NewFromEnvironment() (provider.Provider, error) {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	accessSecret := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if accessKey == "" {
		return nil, errors.New("AWS_ACCESS_KEY_ID is not set")
	}

	if accessSecret == "" {
		return nil, errors.New("AWS_SECRET_ACCESS_KEY is not set")
	}

	return New(), nil
}

func (p *Provider) List(ctx context.Context) ([]string, error) {
	return eksctl.List(ctx)
}

func (p *Provider) Create(ctx context.Context, name string, kubeconfig string) error {
	return eksctl.Create(ctx, name, kubeconfig, eksctl.WithDefaultOutput())
}

func (p *Provider) Delete(ctx context.Context, name string) error {
	return eksctl.Delete(ctx, name, eksctl.WithDefaultOutput())
}

func (p *Provider) Export(ctx context.Context, name, path string) error {
	return eksctl.Export(ctx, name, path)
}
