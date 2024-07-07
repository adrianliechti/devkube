package registry

import (
	"context"
	_ "embed"
	"strings"

	"github.com/adrianliechti/devkube/pkg/apply"
	"github.com/adrianliechti/loop/pkg/kubernetes"
)

const (
	namespace = "default"
)

var (
	//go:embed manifest.yaml
	manifest string
)

func Ensure(ctx context.Context, client kubernetes.Client) error {
	if err := apply.Apply(ctx, client, namespace, strings.NewReader(manifest)); err != nil {
		return err
	}

	return nil
}
