package observability

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
)

const (
	promtail        = "promtail"
	promtailRepo    = "https://grafana.github.io/helm-charts"
	promtailChart   = "promtail"
	promtailVersion = "6.2.2"
)

func installPromtail(ctx context.Context, kubeconfig, namespace string) error {
	values := map[string]any{
		"nameOverride": promtail,

		"config": map[string]any{
			"lokiAddress": "http://" + loki + ":3100/loki/api/v1/push",
		},
	}

	if err := helm.Install(ctx, kubeconfig, namespace, promtail, promtailRepo, promtailChart, promtailVersion, values); err != nil {
		return err
	}

	return nil
}

func uninstallPromtail(ctx context.Context, kubeconfig, namespace string) error {
	if err := helm.Uninstall(ctx, kubeconfig, namespace, promtail); err != nil {
		//return err
	}

	return nil
}
