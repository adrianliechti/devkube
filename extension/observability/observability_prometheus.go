package observability

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/devkube/pkg/kubectl"
)

const (
	prometheus        = "kube-prometheus-stack"
	prometheusRepo    = "https://prometheus-community.github.io/helm-charts"
	prometheusChart   = "kube-prometheus-stack"
	prometheusVersion = "39.6.0"
)

func installPrometheus(ctx context.Context, kubeconfig, namespace string) error {
	values := map[string]any{
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

	return nil
}

func uninstallPrometheus(ctx context.Context, kubeconfig, namespace string) error {
	if err := helm.Uninstall(ctx, kubeconfig, namespace, prometheus); err != nil {
		//return err
	}

	if err := kubectl.Invoke(ctx, kubeconfig, "delete", "pvc", "-n", namespace, "-l", "app.kubernetes.io/instance="+prometheus+"-alertmanager"); err != nil {
		//return err
	}

	if err := kubectl.Invoke(ctx, kubeconfig, "delete", "pvc", "-n", namespace, "-l", "app.kubernetes.io/instance="+prometheus+"-prometheus"); err != nil {
		//return err
	}

	return nil
}
