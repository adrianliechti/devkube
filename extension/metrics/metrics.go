package metrics

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
)

var (
	metricsRepo = "https://kubernetes-sigs.github.io/metrics-server"

	metrics        = "metrics-server"
	metricsChart   = "metrics-server"
	metricsVersion = "3.8.2"
)

func Install(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	values := map[string]any{
		"args": []string{
			"--kubelet-insecure-tls",
			"--kubelet-preferred-address-types=InternalIP",
		},

		"metrics": map[string]any{
			"enabled": true,
		},

		"serviceMonitor": map[string]any{
			"enabled": true,
		},
	}

	if err := helm.Install(ctx, kubeconfig, namespace, metrics, metricsRepo, metricsChart, metricsVersion, values); err != nil {
		return err
	}

	return nil
}

func Uninstall(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	if err := helm.Uninstall(ctx, kubeconfig, namespace, metrics); err != nil {
		//return err
	}

	return nil
}
