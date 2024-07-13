package otel

import (
	"context"
	_ "embed"
	"strings"

	"github.com/adrianliechti/loop/pkg/kubernetes"
)

const (
	namespace = "platform"
)

var (
	//go:embed manifest.yaml
	manifest string
)

func Ensure(ctx context.Context, client kubernetes.Client) error {
	if err := client.Apply(ctx, namespace, strings.NewReader(manifest)); err != nil {
		return err
	}

	return nil
}
