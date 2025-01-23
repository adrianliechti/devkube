package alloy

import (
	"context"
	_ "embed"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/loop/pkg/kubernetes"
)

const (
	name      = "alloy"
	namespace = "platform"

	// https://artifacthub.io/packages/helm/grafana/alloy
	repoURL      = "https://grafana.github.io/helm-charts"
	chartName    = "alloy"
	chartVersion = "0.11.0"
)

var (
	//go:embed config.alloy
	config string
)

func Ensure(ctx context.Context, client kubernetes.Client) error {
	values := map[string]any{
		"alloy": map[string]any{
			"securityContext": map[string]any{
				"runAsUser":  0,
				"privileged": true,
			},

			"mounts": map[string]any{
				"varlog":           true,
				"dockercontainers": true,
			},

			"configMap": map[string]any{
				"create":  true,
				"content": config,
			},

			"extraPorts": []map[string]any{
				{
					"name":       "otlp-grpc",
					"port":       4317,
					"protocol":   "TCP",
					"targetPort": 4317,
				},
				{
					"name":       "otlp-http",
					"port":       4318,
					"protocol":   "TCP",
					"targetPort": 4318,
				},
			},
		},

		"configReloader": map[string]any{
			"enabled": false,
		},
	}

	if err := helm.Ensure(ctx, client, namespace, name, repoURL, chartName, chartVersion, values); err != nil {
		return err
	}

	return nil
}
