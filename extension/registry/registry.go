package registry

import (
	"context"
	_ "embed"
	"strings"

	"github.com/adrianliechti/devkube/pkg/kubectl"
)

//go:embed registry.yaml
var manifest string

func Install(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	return kubectl.Invoke(ctx, []string{"apply", "-f", "-"}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(namespace), kubectl.WithInput(strings.NewReader(manifest)), kubectl.WithDefaultOutput())
}

func Uninstall(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	return kubectl.Invoke(ctx, []string{"delete", "-f", "-"}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(namespace), kubectl.WithInput(strings.NewReader(manifest)), kubectl.WithDefaultOutput())
}
