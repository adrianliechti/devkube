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

		"tempoQuery": map[string]any{
			"enabled": false,
		},
	}

	if err := helm.Install(ctx, tempo, grafanaRepo, tempoChart, tempoVersion, values, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace), helm.WithDefaultOutput()); err != nil {
		return err
	}

	return nil
}

func uninstallTempo(ctx context.Context, kubeconfig, namespace string) error {
	if err := helm.Uninstall(ctx, tempo, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace)); err != nil {
		//return err
	}

	if err := kubectl.Invoke(ctx, []string{"delete", "pvc", "-l", "app.kubernetes.io/instance=" + tempo}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(namespace)); err != nil {
		return err
	}

	return nil
}
