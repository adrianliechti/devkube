package observability

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/devkube/pkg/kubectl"
)

const (
	loki        = "loki"
	lokiChart   = "loki"
	lokiVersion = "2.16.0"
)

func installLoki(ctx context.Context, kubeconfig, namespace string) error {
	values := map[string]any{
		"rbac": map[string]any{
			"pspEnabled": false,
		},

		"persistence": map[string]any{
			"enabled": true,
			"size":    "10Gi",
		},

		"ruler": map[string]any{
			"storage": map[string]any{
				"type": "local",
				"local": map[string]any{
					"directory": "/rules",
				},
				"rule_path":        "/tmp/scratch",
				"alertmanager_url": "http://" + prometheus + "-alertmanager:9093",
				"ring": map[string]any{
					"kvstore": map[string]any{
						"store": "inmemory",
					},
				},
				"enable_api": true,
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
