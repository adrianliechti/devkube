package registry

import (
	"context"
	_ "embed"
	"strings"

	"github.com/adrianliechti/devkube/pkg/kube"
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
	if err := kube.EnsureNamespace(ctx, client, namespace); err != nil {
		return err
	}

	if err := client.Apply(ctx, namespace, strings.NewReader(manifest)); err != nil {
		return err
	}

	return nil
}
