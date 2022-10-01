package observability

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
)

const (
	grafana        = "grafana"
	grafanaChart   = "grafana"
	grafanaVersion = "6.40.0"
)

func installGrafana(ctx context.Context, kubeconfig, namespace string) error {
	values := map[string]any{
		"adminUser":     "admin",
		"adminPassword": "admin",

		"rbac": map[string]any{
			"pspEnabled": false,
		},

		"persistence": map[string]any{
			"enabled": true,
			"size":    "10Gi",
		},

		"serviceMonitor": map[string]any{
			"enabled": true,
		},

		"sidecar": map[string]any{
			"dashboards": map[string]any{
				"enabled": true,

				"searchNamespace": "ALL",
			},

			"datasources": map[string]any{
				"enabled": true,

				"searchNamespace": "ALL",
			},

			"plugins": map[string]any{
				"enabled": true,

				"searchNamespace": "ALL",
			},

			"notifiers": map[string]any{
				"enabled": true,

				"searchNamespace": "ALL",
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
						"jsonData": map[string]any{
							"httpMethod": "GET",

							"tracesToLogs": map[string]any{
								"datasourceUid":      "loki",
								"mapTagNamesEnabled": true,
							},
							"tracesToMetrics": map[string]any{
								"datasourceUid": "prometheus",
							},
							"serviceMap": map[string]any{
								"datasourceUid": "prometheus",
							},
							"nodeGraph": map[string]any{
								"enabled": true,
							},
							"lokiSearch": map[string]any{
								"datasourceUid": "loki",
							},
						},
					},
					{
						"name":   "Prometheus",
						"type":   "prometheus",
						"uid":    "prometheus",
						"url":    "http://" + prometheus + "-prometheus:9090",
						"access": "proxy",
						"jsonData": map[string]any{
							"httpMethod": "GET",
						},
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

	if err := helm.Install(ctx, grafana, grafanaRepo, grafanaChart, grafanaVersion, values, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace), helm.WithDefaultOutput()); err != nil {
		return err
	}

	return nil
}

func uninstallGrafana(ctx context.Context, kubeconfig, namespace string) error {
	if err := helm.Uninstall(ctx, grafana, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace)); err != nil {
		// return err
	}

	return nil
}
