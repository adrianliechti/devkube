package monitoring

import (
	"context"
	_ "embed"
	"strings"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/loop/pkg/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	name      = "monitoring"
	namespace = "platform"

	// https://artifacthub.io/packages/helm/prometheus-community/kube-prometheus-stack
	repoURL      = "https://prometheus-community.github.io/helm-charts"
	chartName    = "kube-prometheus-stack"
	chartVersion = "61.3.0"
)

var (
	//go:embed manifest.yaml
	manifest string
)

func Ensure(ctx context.Context, client kubernetes.Client) error {
	values := map[string]any{
		"nameOverride":     name,
		"fullnameOverride": name,

		"cleanPrometheusOperatorObjectNames": true,

		"grafana": map[string]any{
			"enabled":                false,
			"forceDeployDashboards":  true,
			"forceDeployDatasources": true,
		},

		"alertmanager": map[string]any{
			"alertmanagerSpec": map[string]any{
				"storage": map[string]any{
					"volumeClaimTemplate": map[string]any{
						"spec": map[string]any{
							"accessModes": []string{"ReadWriteOnce"},
							"resources": map[string]any{
								"requests": map[string]any{
									"storage": "10Gi",
								},
							},
						},
					},
				},
			},
		},

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

				"retentionSize": "9GiB",

				"storageSpec": map[string]any{
					"volumeClaimTemplate": map[string]any{
						"spec": map[string]any{
							"accessModes": []string{"ReadWriteOnce"},
							"resources": map[string]any{
								"requests": map[string]any{
									"storage": "10Gi",
								},
							},
						},
					},
				},
			},
		},
	}

	if err := helm.Ensure(ctx, client, namespace, name, repoURL, chartName, chartVersion, values); err != nil {
		return err
	}

	if err := client.Apply(ctx, namespace, strings.NewReader(manifest)); err != nil {
		return err
	}

	for _, name := range []string{
		"monitoring-grafana-overview",
		"monitoring-prometheus",
		"monitoring-alertmanager-overview",
		"monitoring-nodes-darwin",
	} {
		client.CoreV1().ConfigMaps(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	}

	return nil
}
