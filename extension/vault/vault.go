package vault

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
)

const (
	// csiRepo = "https://kubernetes-sigs.github.io/secrets-store-csi-driver/charts"

	// csi        = "secrets-store-csi-driver"
	// csiChart   = "secrets-store-csi-driver"
	// csiVersion = "1.2.4"

	vaultRepo = "https://helm.releases.hashicorp.com"

	vault        = "vault"
	vaultChart   = "vault"
	vaultVersion = "0.22.0"
)

func Install(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	// csiValues := map[string]any{}

	// if err := helm.Install(ctx, csi, csiRepo, csiChart, csiVersion, csiValues, helm.WithKubeconfig(kubeconfig), helm.WithNamespace("kube-system"), helm.WithWait(true), helm.WithDefaultOutput()); err != nil {
	// 	return err
	// }

	vaultValues := map[string]any{
		"server": map[string]any{
			"enabled": true,

			"dev": map[string]any{
				"enabled": true,
			},
		},

		"csi": map[string]any{
			"enabled": true,
		},

		"injector": map[string]any{
			"enabled": false,
		},
	}

	if err := helm.Install(ctx, vault, vaultRepo, vaultChart, vaultVersion, vaultValues, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace), helm.WithDefaultOutput()); err != nil {
		return err
	}

	return nil
}

func Uninstall(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	if err := helm.Uninstall(ctx, vault, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace)); err != nil {
		//return err
	}

	// if err := helm.Uninstall(ctx, csi, helm.WithKubeconfig(kubeconfig), helm.WithNamespace("kube-system")); err != nil {
	// 	//return err
	// }

	return nil
}
