package observability

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
)

const (
	falcoExporter        = "falco"
	falcoRepo    = "https://falcosecurity.github.io/charts"
	falcoExporterChart   = "falco"
	falcoVersion = "2.0.16"

	        = "falco-exporter"
	falcoExporterVersion = "0.8.2"
)

func installFalco(ctx context.Context, kubeconfig, namespace string) error {
	values := map[string]any{
		"falco": map[string]any{
			"grpc": map[string]any{
				"enabled": true,
			},

			"grpc_output": map[string]any{
				"enabled": true,
			},
		},
	}

	if err := helm.Install(ctx, kubeconfig, namespace, "grafana", grafanaRepo, grafanaChart, grafanaVersion, values); err != nil {
		return err
	}

	return nil
}

func uninstallGrafana(ctx context.Context, kubeconfig, namespace string) error {
	if err := helm.Uninstall(ctx, kubeconfig, namespace, grafana); err != nil {
		//return err
	}

	return nil
}
