package registry

import (
	"bytes"
	"context"

	"github.com/adrianliechti/devkube/pkg/kubectl"

	_ "embed"
)

//go:embed registry.yaml
var manifest string

func Install(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	data := bytes.NewReader([]byte(manifest))
	return kubectl.Invoke(ctx, []string{"apply", "-f", "-"}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(namespace), kubectl.WithInput(data), kubectl.WithDefaultOutput())
}

func Uninstall(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	data := bytes.NewReader([]byte(manifest))
	return kubectl.Invoke(ctx, []string{"delete", "-f", "-"}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(namespace), kubectl.WithInput(data), kubectl.WithDefaultOutput())
}
