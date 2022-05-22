package observability

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/devkube/pkg/kubectl"
)

const (
	tempo        = "tempo"
	tempoRepo    = "https://grafana.github.io/helm-charts"
	tempoChart   = "tempo"
	tempoVersion = "0.15.0"
)

func installTempo(ctx context.Context, kubeconfig, namespace string) error {
	values := map[string]any{
		"nameOverride": tempo,

		"persistence": map[string]any{
			"enabled": true,
			"size":    "10Gi",
		},

		"tempoQuery": map[string]any{
			"tag": "1.4.0",
		},
	}

	if err := helm.Install(ctx, kubeconfig, namespace, tempo, tempoRepo, tempoChart, tempoVersion, values); err != nil {
		return err
	}

	return nil
}

func uninstallTempo(ctx context.Context, kubeconfig, namespace string) error {
	if err := helm.Uninstall(ctx, kubeconfig, namespace, tempo); err != nil {
		return err
	}

	if err := kubectl.Invoke(ctx, kubeconfig, "delete", "pvc", "storage-"+tempo+"-0"); err != nil {
		return err
	}

	return nil
}
