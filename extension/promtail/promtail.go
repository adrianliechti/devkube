package promtail

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/loop/pkg/kubernetes"
)

const (
	name      = "promtail"
	namespace = "platform"

	// https://artifacthub.io/packages/helm/grafana/promtail
	repoURL      = "https://grafana.github.io/helm-charts"
	chartName    = "promtail"
	chartVersion = "6.16.6"
)

func Ensure(ctx context.Context, client kubernetes.Client) error {
	values := map[string]any{
		"config": map[string]any{
			"clients": []map[string]any{
				{
					"url": "http://loki:3100/loki/api/v1/push",
				},
			},
		},
	}

	if err := helm.Ensure(ctx, client, namespace, name, repoURL, chartName, chartVersion, values); err != nil {
		return err
	}

	return nil
}
