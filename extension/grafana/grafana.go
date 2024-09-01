package grafana

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/loop/pkg/kubernetes"
)

const (
	name      = "grafana"
	namespace = "platform"

	// https://artifacthub.io/packages/helm/grafana/grafana
	repoURL      = "https://grafana.github.io/helm-charts"
	chartName    = "grafana"
	chartVersion = "8.5.0"
)

func Ensure(ctx context.Context, client kubernetes.Client) error {
	grafanaHost := ""
	grafanaURL := "http://" + name + "." + namespace

	values := map[string]any{
		"adminUser":     "admin",
		"adminPassword": "admin",

		"persistence": map[string]any{
			"enabled": true,
			"size":    "10Gi",
		},

		"deploymentStrategy": map[string]any{
			"type": "Recreate",
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
			"server": map[string]any{
				"domain":   grafanaHost,
				"root_url": grafanaURL,
			},

			"analytics": map[string]any{
				"reporting_enabled": false,

				"check_for_updates":        false,
				"check_for_plugin_updates": false,

				"feedback_links_enabled": false,
			},

			"users": map[string]any{
				"allow_sign_up":        false,
				"allow_org_create":     false,
				"auto_assign_org":      true,
				"auto_assign_org_role": "Editor",
				"viewers_can_edit":     true,
				"editors_can_admin":    true,
			},

			"auth": map[string]any{
				"disable_login_form":   true,
				"disable_signout_menu": true,
			},

			"auth.basic": map[string]any{
				"enabled": true,
			},

			"auth.anonymous": map[string]any{
				"enabled":  true,
				"org_name": "Main Org.",
				"org_role": "Editor",
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
						"url":    "http://loki:3100",
						"access": "proxy",
					},
					{
						"name":   "Tempo",
						"type":   "tempo",
						"uid":    "tempo",
						"url":    "http://tempo:3100",
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
						"url":    "http://prometheus:9090",
						"access": "proxy",

						"jsonData": map[string]any{
							"httpMethod": "GET",
						},
					},
					{
						"name":   "Alertmanager",
						"type":   "alertmanager",
						"uid":    "alertmanager",
						"url":    "http://alertmanager:9093",
						"access": "proxy",

						"jsonData": map[string]any{
							"implementation": "prometheus",
						},
					},
				},
			},
		},
	}

	if err := helm.Ensure(ctx, client, namespace, name, repoURL, chartName, chartVersion, values); err != nil {
		return err
	}

	return nil
}
