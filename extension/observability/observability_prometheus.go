package observability

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/devkube/pkg/kubectl"
)

const (
	monitoring        = "monitoring"
	monitoringChart   = "kube-prometheus-stack"
	monitoringVersion = "40.3.1"
)

func installPrometheus(ctx context.Context, kubeconfig, namespace string) error {
	values := map[string]any{
		"nameOverride":     monitoring,
		"fullnameOverride": monitoring,

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

	if err := helm.Install(ctx, monitoring, prometheusRepo, monitoringChart, monitoringVersion, values, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace), helm.WithDefaultOutput()); err != nil {
		return err
	}

	if err := kubectl.Invoke(ctx, []string{"delete", "configmap", monitoring + "-nodes-darwin"}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(namespace)); err != nil {
		// return err
	}

	return nil
}

func uninstallPrometheus(ctx context.Context, kubeconfig, namespace string) error {
	if err := helm.Uninstall(ctx, monitoring, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace)); err != nil {
		// return err
	}

	if err := kubectl.Invoke(ctx, []string{"delete", "pvc", "-l", "app.kubernetes.io/instance=" + monitoring}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(namespace)); err != nil {
		// return err
	}

	if err := kubectl.Invoke(ctx, []string{"delete", "secret", monitoring + "-admission"}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(namespace)); err != nil {
		// return err
	}

	return nil
}
