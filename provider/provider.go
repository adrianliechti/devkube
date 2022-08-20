package provider

import "context"

type Provider interface {
	List(ctx context.Context) ([]string, error)

	Create(ctx context.Context, name string, kubeconfig string) error
	Delete(ctx context.Context, name string) error

	ExportConfig(ctx context.Context, name, path string) error
}
