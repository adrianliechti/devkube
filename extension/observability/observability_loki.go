package observability

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/devkube/pkg/kubectl"
)

const (
	loki        = "loki"
	lokiChart   = "loki"
	lokiVersion = "3.1.0"
)

func installLoki(ctx context.Context, kubeconfig, namespace string) error {
	values := map[string]any{
		"loki": map[string]any{
			"commonConfig": map[string]any{
				"replication_factor": 1,
			},

			"storage": map[string]any{
				"type": "filesystem",
			},
		},

		"singleBinary": map[string]any{
			"persistence": map[string]any{
				"size": "10Gi",
			},
		},
	}

	if err := helm.Install(ctx, loki, grafanaRepo, lokiChart, lokiVersion, values, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace), helm.WithDefaultOutput()); err != nil {
		return err
	}

	return nil
}

func uninstallLoki(ctx context.Context, kubeconfig, namespace string) error {
	if err := helm.Uninstall(ctx, loki, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace)); err != nil {
		//return err
	}

	if err := kubectl.Invoke(ctx, []string{"delete", "pvc", "-l", "release=" + loki}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(namespace)); err != nil {
		return err
	}

	return nil
}
