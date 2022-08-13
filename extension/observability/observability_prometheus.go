package observability

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/devkube/pkg/kubectl"
)

const (
	prometheus        = "monitoring"
	prometheusChart   = "kube-prometheus-stack"
	prometheusVersion = "39.6.0"
)

func installPrometheus(ctx context.Context, kubeconfig, namespace string) error {
	values := map[string]any{
		"nameOverride":     prometheus,
		"fullnameOverride": prometheus,

		"cleanPrometheusOperatorObjectNames": true,

		"kubeEtcd": map[string]any{
			"service": map[string]any{
				"targetPort": 2381,
			},
		},

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
				"serviceMonitorSelector":                  nil,
				"serviceMonitorSelectorNilUsesHelmValues": false,

				"podMonitorSelector":                  nil,
				"podMonitorSelectorNilUsesHelmValues": false,

				"probeSelector":                  nil,
				"probeSelectorNilUsesHelmValues": false,

				"ruleSelector":                  nil,
				"ruleSelectorNilUsesHelmValues": false,

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

	if err := helm.Install(ctx, kubeconfig, namespace, prometheus, prometheusRepo, prometheusChart, prometheusVersion, values); err != nil {
		return err
	}

	if err := kubectl.Invoke(ctx, kubeconfig, "delete", "configmap", "-n", namespace, prometheus+"-nodes-darwin"); err != nil {
		//return err
	}

	return nil
}

func uninstallPrometheus(ctx context.Context, kubeconfig, namespace string) error {
	if err := helm.Uninstall(ctx, kubeconfig, namespace, prometheus); err != nil {
		//return err
	}

	if err := kubectl.Invoke(ctx, kubeconfig, "delete", "pvc", "-n", namespace, "-l", "app.kubernetes.io/instance="+prometheus); err != nil {
		//return err
	}

	if err := kubectl.Invoke(ctx, kubeconfig, "delete", "secret", "-n", namespace, prometheus+"-admission"); err != nil {
		//return err
	}

	return nil
}
