package prometheus

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/loop/pkg/kubernetes"
)

const (
	namespace = "monitoring"

	repo    = "https://prometheus-community.github.io/helm-charts"
	chart   = "kube-prometheus-stack"
	version = "61.2.0"
)

func Ensure(ctx context.Context, client kubernetes.Client) error {
	values := map[string]any{
		"nameOverride":     "monitoring",
		"fullnameOverride": "monitoring",

		"cleanPrometheusOperatorObjectNames": true,

		"coreDns": map[string]any{
			"enabled": false,
		},

		"kubeDns": map[string]any{
			"enabled": false,
		},

		"kubeEtcd": map[string]any{
			"enabled": false,
		},

		"kubeScheduler": map[string]any{
			"enabled": false,
		},

		"kubeProxy": map[string]any{
			"enabled": false,
		},

		"kubeControllerManager": map[string]any{
			"enabled": false,
		},

		"grafana": map[string]any{
			"enabled":                false,
			"forceDeployDashboards":  true,
			"forceDeployDatasources": true,
		},

		// "alertmanager": map[string]any{
		// 	"alertmanagerSpec": map[string]any{
		// 		"storage": map[string]any{
		// 			"volumeClaimTemplate": map[string]any{
		// 				"spec": map[string]any{
		// 					"accessModes": []string{"ReadWriteOnce"},
		// 					"resources": map[string]any{
		// 						"requests": map[string]any{
		// 							"storage": "10Gi",
		// 						},
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// },

		"prometheus": map[string]any{
			"prometheusSpec": map[string]any{
				"enableAdminAPI":            true,
				"enableRemoteWriteReceiver": true,

				"serviceMonitorSelector":                  nil,
				"serviceMonitorSelectorNilUsesHelmValues": false,

				"podMonitorSelector":                  nil,
				"podMonitorSelectorNilUsesHelmValues": false,

				"probeSelector":                  nil,
				"probeSelectorNilUsesHelmValues": false,

				"ruleSelector":                  nil,
				"ruleSelectorNilUsesHelmValues": false,

				//"retentionSize": "9GiB",

				// "storageSpec": map[string]any{
				// 	"volumeClaimTemplate": map[string]any{
				// 		"spec": map[string]any{
				// 			"accessModes": []string{"ReadWriteOnce"},
				// 			"resources": map[string]any{
				// 				"requests": map[string]any{
				// 					"storage": "10Gi",
				// 				},
				// 			},
				// 		},
				// 	},
				// },
			},
		},
	}

	if err := helm.Ensure(ctx, client, namespace, "monitoring", repo, chart, version, values); err != nil {
		return err
	}

	return nil
}
