package observability

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/devkube/pkg/kubectl"
)

const (
	tempo        = "tempo"
	tempoChart   = "tempo"
	tempoVersion = "0.15.8"
)

func installTempo(ctx context.Context, kubeconfig, namespace string) error {
	values := map[string]any{
		"persistence": map[string]any{
			"enabled": true,
			"size":    "10Gi",
		},
	}

	if err := helm.Install(ctx, kubeconfig, namespace, tempo, grafanaRepo, tempoChart, tempoVersion, values); err != nil {
		return err
	}

	return nil
}

func uninstallTempo(ctx context.Context, kubeconfig, namespace string) error {
	if err := helm.Uninstall(ctx, kubeconfig, namespace, tempo); err != nil {
		//return err
	}

	if err := kubectl.Invoke(ctx, kubeconfig, "delete", "pvc", "-n", namespace, "storage-"+tempo+"-0"); err != nil {
		//return err
	}

	return nil
}
