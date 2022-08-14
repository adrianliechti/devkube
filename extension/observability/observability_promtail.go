package observability

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
)

const (
	promtail        = "promtail"
	promtailChart   = "promtail"
	promtailVersion = "6.2.2"
)

func installPromtail(ctx context.Context, kubeconfig, namespace string) error {
	values := map[string]any{
		"config": map[string]any{
			"clients": []map[string]any{
				{
					"url": "http://" + loki + ":3100/loki/api/v1/push",
				},
			},
		},
	}

	if err := helm.Install(ctx, promtail, grafanaRepo, promtailChart, promtailVersion, values, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace), helm.WithDefaultOutput()); err != nil {
		return err
	}

	return nil
}

func uninstallPromtail(ctx context.Context, kubeconfig, namespace string) error {
	if err := helm.Uninstall(ctx, promtail, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace)); err != nil {
		//return err
	}

	return nil
}
