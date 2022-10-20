package catalog

import (
	"context"
	_ "embed"

	"github.com/adrianliechti/devkube/pkg/kubectl"
)

const (
	catalog          = "catalog-controller-manager"
	catalogNamespace = "loop"

	manifestURL = "https://raw.githubusercontent.com/adrianliechti/loop-catalog/main/kubernetes/install.yaml"
)

func Install(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	if err := kubectl.Invoke(ctx, []string{"apply", "-f", manifestURL}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(catalogNamespace), kubectl.WithDefaultOutput()); err != nil {
		return err
	}

	if err := kubectl.Invoke(ctx, []string{"delete", "clusterrolebinding", catalog}, kubectl.WithKubeconfig(kubeconfig)); err != nil {
		// ignore error
	}

	if err := kubectl.Invoke(ctx, []string{"create", "clusterrolebinding", catalog, "--clusterrole=cluster-admin", "--serviceaccount=" + namespace + ":" + catalog}, kubectl.WithKubeconfig(kubeconfig)); err != nil {
		return err
	}

	return nil
}

func Uninstall(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	if err := kubectl.Invoke(ctx, []string{"delete", "clusterrolebinding", catalog}, kubectl.WithKubeconfig(kubeconfig)); err != nil {
		// return err
	}

	if err := kubectl.Invoke(ctx, []string{"delete", "-f", manifestURL}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(catalogNamespace), kubectl.WithDefaultOutput()); err != nil {
		// return err
	}

	return nil
}
