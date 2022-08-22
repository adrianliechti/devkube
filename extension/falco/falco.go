package falco

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
)

const (
	falco        = "falco"
	falcoChart   = "falco"
	falcoVersion = "2.0.16"

	exporter        = "falco-exporter"
	exporterChart   = "falco-exporter"
	exporterVersion = "0.8.2"

	falcoRepo = "https://falcosecurity.github.io/charts"
)

func Install(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	falcoValues := map[string]any{
		"falco": map[string]any{
			"grpc": map[string]any{
				"enabled": true,
			},

			"grpc_output": map[string]any{
				"enabled": true,
			},
		},
	}

	if err := helm.Install(ctx, falco, falcoRepo, falcoChart, falcoVersion, falcoValues, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace), helm.WithDefaultOutput()); err != nil {
		return err
	}

	exporterValues := map[string]any{
		"serviceMonitor": map[string]any{
			"enabled": true,
		},

		"grafanaDashboard": map[string]any{
			"enabled":   true,
			"namespace": nil,
		},

		"prometheusRules": map[string]any{
			"enabled": true,
		},
	}

	if err := helm.Install(ctx, exporter, falcoRepo, exporterChart, exporterVersion, exporterValues, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace), helm.WithDefaultOutput()); err != nil {
		return err
	}

	return nil
}

func Uninstall(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	if err := helm.Uninstall(ctx, falco, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace)); err != nil {
		//return err
	}

	if err := helm.Uninstall(ctx, exporter, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace)); err != nil {
		//return err
	}

	return nil
}
