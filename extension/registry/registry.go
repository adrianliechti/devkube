package registry

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/kubectl"
)

const (
	manifest = "https://raw.githubusercontent.com/adrianliechti/loop-registry/main/kubernetes/install.yaml"
)

func Install(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	return kubectl.Invoke(ctx, []string{"apply", "-f", manifest}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(namespace), kubectl.WithDefaultOutput())
}

func Uninstall(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	return kubectl.Invoke(ctx, []string{"delete", "-f", manifest}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(namespace), kubectl.WithDefaultOutput())
}
