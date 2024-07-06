package tempo

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/loop/pkg/kubernetes"
)

const (
	name      = "tempo"
	namespace = "monitoring"

	repoURL      = "https://grafana.github.io/helm-charts"
	chartName    = "tempo"
	chartVersion = "1.8.0"
)

func Ensure(ctx context.Context, client kubernetes.Client) error {
	values := map[string]any{
		"tempo": map[string]any{
			"metricsGenerator": map[string]any{
				"enabled":        true,
				"remoteWriteUrl": "http://monitoring:9090/api/v1/write",
			},
		},

		"persistence": map[string]any{
			"enabled": true,
			"size":    "10Gi",
		},
	}

	if err := helm.Ensure(ctx, client, namespace, name, repoURL, chartName, chartVersion, values); err != nil {
		return err
	}

	return nil
}
