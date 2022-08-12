package observability

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
)

const (
	falco        = "falco"
	falcoChart   = "falco"
	falcoVersion = "2.0.16"

	falcoExporter        = "falco-exporter"
	falcoExporterChart   = "falco-exporter"
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

		"driver": map[string]any{
			"enabled": true,
			"kind":    "ebpf",

			"loader": map[string]any{
				"enabled": false,
			},
		},
	}

	if err := helm.Install(ctx, kubeconfig, namespace, falco, falcoRepo, falcoChart, falcoVersion, values); err != nil {
		return err
	}

	if err := helm.Install(ctx, kubeconfig, namespace, falcoExporter, falcoRepo, falcoExporterChart, falcoExporterVersion, nil); err != nil {
		return err
	}

	return nil
}

func uninstallFalco(ctx context.Context, kubeconfig, namespace string) error {
	if err := helm.Uninstall(ctx, kubeconfig, namespace, falco); err != nil {
		//return err
	}

	if err := helm.Uninstall(ctx, kubeconfig, namespace, falcoExporter); err != nil {
		//return err
	}

	return nil
}
