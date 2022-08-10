package observability

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
)

const (
	grafana        = "grafana"
	grafanaRepo    = "https://grafana.github.io/helm-charts"
	grafanaChart   = "grafana"
	grafanaVersion = "6.32.12"
)

func installGrafana(ctx context.Context, kubeconfig, namespace string) error {
	values := map[string]any{
		"nameOverride": grafana,

		"adminUser":     "admin",
		"adminPassword": "admin",

		"persistence": map[string]any{
			"enabled": true,
			"size":    "10Gi",
		},

		"sidecar": map[string]any{
			"dashboards": map[string]any{
				"enabled": true,
			},

			"datasources": map[string]any{
				"enabled": true,
			},

			"plugins": map[string]any{
				"enabled": true,
			},
		},

		"grafana.ini": map[string]any{
			"alerting": map[string]any{
				"enabled": false,
			},

			"users": map[string]any{
				"allow_sign_up":        false,
				"allow_org_create":     false,
				"auto_assign_org":      true,
				"auto_assign_org_role": "Viewer",
				"viewers_can_edit":     true,
				"editors_can_admin":    false,
			},

			"auth": map[string]any{
				"disable_login_form": true,
			},

			"auth.basic": map[string]any{
				"enabled": false,
			},

			"auth.anonymous": map[string]any{
				"enabled":  true,
				"org_name": "Main Org.",
				"org_role": "Admin",
			},
		},

		"datasources": map[string]any{
			"datasources.yaml": map[string]any{
				"apiVersion": 1,
				"datasources": []map[string]any{
					{
						"name":   "Loki",
						"type":   "loki",
						"uid":    "loki",
						"url":    "http://" + loki + ":3100",
						"access": "proxy",
					},
					{
						"name":   "Tempo",
						"type":   "tempo",
						"uid":    "tempo",
						"url":    "http://" + tempo + ":3100",
						"access": "proxy",
					},
					{
						"name":   "Prometheus",
						"type":   "prometheus",
						"uid":    "prometheus",
						"url":    "http://" + prometheus + "-prometheus:9090",
						"access": "proxy",
					},
					{
						"name":   "Alertmanager",
						"type":   "alertmanager",
						"uid":    "alertmanager",
						"url":    "http://" + prometheus + "-alertmanager:9093",
						"access": "proxy",

						"jsonData": map[string]any{
							"implementation": "prometheus",
						},
					},
				},
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
