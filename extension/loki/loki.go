package loki

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/loop/pkg/kubernetes"
)

const (
	name      = "loki"
	namespace = "monitoring"

	repoURL      = "https://grafana.github.io/helm-charts"
	chartName    = "loki"
	chartVersion = "5.48.0"
)

func Ensure(ctx context.Context, client kubernetes.Client) error {
	values := map[string]any{
		"loki": map[string]any{
			"commonConfig": map[string]any{
				"replication_factor": 1,
			},

			"auth_enabled": false,

			"storage": map[string]any{
				"type": "filesystem",
			},

			"querier": map[string]any{
				"max_concurrent": 8,
			},
		},

		"singleBinary": map[string]any{
			"replicas": 1,

			"persistence": map[string]any{
				"size": "10Gi",
			},
		},

		"tableManager": map[string]any{
			"retention_deletes_enabled": true,
			"retention_period":          "7d",
		},

		"gateway": map[string]any{
			"enabled": false,
		},

		"test": map[string]any{
			"enabled": false,
		},

		"monitoring": map[string]any{
			"dashboards": map[string]any{
				"enabled": false,
			},

			"rules": map[string]any{
				"enabled": false,
			},

			"serviceMonitor": map[string]any{
				"enabled": false,
			},

			"selfMonitoring": map[string]any{
				"enabled": false,

				"grafanaAgent": map[string]any{
					"installOperator": false,
				},
			},

			"lokiCanary": map[string]any{
				"enabled": false,
			},
		},
	}

	if err := helm.Ensure(ctx, client, namespace, name, repoURL, chartName, chartVersion, values); err != nil {
		return err
	}

	return nil
}
