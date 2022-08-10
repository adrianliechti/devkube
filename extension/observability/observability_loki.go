package observability

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/devkube/pkg/kubectl"
)

const (
	loki        = "loki"
	lokiRepo    = "https://grafana.github.io/helm-charts"
	lokiChart   = "loki"
	lokiVersion = "2.13.3"
)

func installLoki(ctx context.Context, kubeconfig, namespace string) error {
	values := map[string]any{
		"nameOverride": loki,

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

	if err := helm.Install(ctx, kubeconfig, namespace, loki, lokiRepo, lokiChart, lokiVersion, values); err != nil {
		return err
	}

	return nil
}

func uninstallLoki(ctx context.Context, kubeconfig, namespace string) error {
	if err := helm.Uninstall(ctx, kubeconfig, namespace, loki); err != nil {
		//return err
	}

	if err := kubectl.Invoke(ctx, kubeconfig, "delete", "pvc", "-n", namespace, "storage-"+loki+"-0"); err != nil {
		//return err
	}

	return nil
}
