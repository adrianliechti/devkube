package observability

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/devkube/pkg/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	prometheus        = "monitoring"
	prometheusChart   = "kube-prometheus-stack"
	prometheusVersion = "39.6.0"
)

func installPrometheus(ctx context.Context, kubeconfig, namespace string) error {
	client, err := kubernetes.NewFromConfig(kubeconfig)

	if err != nil {
		return err
	}

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

	if err := helm.Install(ctx, prometheus, prometheusRepo, prometheusChart, prometheusVersion, values, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace), helm.WithDefaultOutput()); err != nil {
		return err
	}

	client.CoreV1().ConfigMaps(namespace).Delete(ctx, prometheus+"-nodes-darwin", metav1.DeleteOptions{})

	return nil
}

func uninstallPrometheus(ctx context.Context, kubeconfig, namespace string) error {
	client, err := kubernetes.NewFromConfig(kubeconfig)

	if err != nil {
		return err
	}

	if err := helm.Uninstall(ctx, prometheus, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace)); err != nil {
		//return err
	}

	client.CoreV1().Secrets(namespace).Delete(ctx, prometheus+"-admission", metav1.DeleteOptions{})

	client.CoreV1().PersistentVolumeClaims(namespace).DeleteCollection(ctx, metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/instance=" + prometheus,
	})

	return nil
}
