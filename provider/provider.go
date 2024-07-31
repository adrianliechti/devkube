package provider

import (
	"context"
)

type Provider interface {
	List(ctx context.Context) ([]string, error)

	Create(ctx context.Context, name string) error
	Delete(ctx context.Context, name string) error

	Start(ctx context.Context, name string) error
	Stop(ctx context.Context, name string) error

	Config(ctx context.Context, name string) ([]byte, error)
}
