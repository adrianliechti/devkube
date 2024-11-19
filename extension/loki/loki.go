package loki

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/loop/pkg/kubernetes"
)

const (
	name      = "loki"
	namespace = "platform"

	// https://artifacthub.io/packages/helm/grafana/loki
	repoURL      = "https://grafana.github.io/helm-charts"
	chartName    = "loki"
	chartVersion = "6.20.0"
)

func Ensure(ctx context.Context, client kubernetes.Client) error {
	values := map[string]any{
		"deploymentMode": "SingleBinary",

		"loki": map[string]any{
			"auth_enabled": false,

			"commonConfig": map[string]any{
				"replication_factor": 1,
			},

			"storage": map[string]any{
				"type": "filesystem",
			},

			"schemaConfig": map[string]any{
				"configs": []map[string]any{
					{
						"from":  "2024-01-01",
						"store": "tsdb",
						"index": map[string]any{
							"prefix": "loki_index_",
							"period": "24h",
						},
						"object_store": "filesystem",
						"schema":       "v13",
					},
				},
			},

			"ingester": map[string]any{
				"chunk_encoding": "snappy",
			},

			"querier": map[string]any{
				"max_concurrent": 8,
			},

			"limits_config": map[string]any{
				"allow_structured_metadata": true,
			},
		},

		"singleBinary": map[string]any{
			"replicas": 1,

			"persistence": map[string]any{
				"size": "10Gi",
			},
		},

		"resultsCache": map[string]any{
			"enabled": false,
		},

		"chunksCache": map[string]any{
			"enabled": false,
		},

		"read": map[string]any{
			"replicas": 0,
		},

		"backend": map[string]any{
			"replicas": 0,
		},

		"write": map[string]any{
			"replicas": 0,
		},

		"test": map[string]any{
			"enabled": false,
		},

		"lokiCanary": map[string]any{
			"enabled": false,
		},
	}

	if err := helm.Ensure(ctx, client, namespace, name, repoURL, chartName, chartVersion, values); err != nil {
		return err
	}

	return nil
}
