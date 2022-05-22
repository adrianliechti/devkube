package observability

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/devkube/pkg/kubectl"
)

const (
	prometheus        = "prometheus"
	prometheusRepo    = "https://prometheus-community.github.io/helm-charts"
	prometheusChart   = "prometheus"
	prometheusVersion = "15.8.7"
)

func installPrometheus(ctx context.Context, kubeconfig, namespace string) error {
	values := map[string]any{
		"nameOverride": prometheus,

		"server": map[string]any{
			"statefulSet": map[string]any{
				"enabled": true,
			},

			"persistentVolume": map[string]any{
				"enabled": true,
				"size":    "10Gi",
			},
		},

		"alertmanager": map[string]any{
			"enabled": false,
		},

		"pushgateway": map[string]any{
			"enabled": false,
		},

		"configmapReload": map[string]any{
			"prometheus": map[string]any{
				"enabled": false,
			},
			"alertmanager": map[string]any{
				"enabled": false,
			},
		},
	}

	if err := helm.Install(ctx, kubeconfig, namespace, prometheus, prometheusRepo, prometheusChart, prometheusVersion, values); err != nil {
		return err
	}

	return nil
}

func uninstallPrometheus(ctx context.Context, kubeconfig, namespace string) error {
	if err := helm.Uninstall(ctx, kubeconfig, namespace, prometheus); err != nil {
		return err
	}

	if err := kubectl.Invoke(ctx, kubeconfig, "delete", "pvc", "storage-volume-"+prometheus+"-server-0"); err != nil {
		return err
	}

	return nil
}
