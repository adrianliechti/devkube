package ingress

import (
	"context"
	_ "embed"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/devkube/pkg/kubernetes"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ingressRepo = "https://kubernetes.github.io/ingress-nginx"

	ingress        = "ingress-nginx"
	ingressChart   = "ingress-nginx"
	ingressVersion = "4.3.0"
)

var (
	//go:embed dashboard.json
	dashboard string
)

func Install(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	values := map[string]any{
		"controller": map[string]any{
			"service": map[string]any{
				"type": "ClusterIP",
			},

			"metrics": map[string]any{
				"enabled": true,
			},

			"serviceMonitor": map[string]any{
				"enabled": true,
			},
		},
	}

	if err := helm.Install(ctx, ingress, ingressRepo, ingressChart, ingressVersion, values, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace), helm.WithDefaultOutput()); err != nil {
		return err
	}

	if err := installDashboard(ctx, kubeconfig, namespace); err != nil {
		return err
	}

	return nil
}

func Uninstall(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	if err := helm.Uninstall(ctx, ingress, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace)); err != nil {
		//return err
	}

	if err := uninstallDashboard(ctx, kubeconfig, namespace); err != nil {
		// return err
	}

	return nil
}

func installDashboard(ctx context.Context, kubeconfig, namespace string) error {
	client, err := kubernetes.NewFromConfig(kubeconfig)

	if err != nil {
		return err
	}

	dashboardName := "ingress-nginx-dashboard"

	dashboard := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: dashboardName,

			Labels: map[string]string{
				"grafana_dashboard": dashboardName,
			},
		},

		Data: map[string]string{
			"dashboard.json": dashboard,
		},
	}

	client.CoreV1().ConfigMaps(namespace).Delete(ctx, dashboard.Name, metav1.DeleteOptions{})

	if _, err := client.CoreV1().ConfigMaps(namespace).Create(ctx, dashboard, metav1.CreateOptions{}); err != nil {
		return err
	}

	return nil
}

func uninstallDashboard(ctx context.Context, kubeconfig, namespace string) error {
	client, err := kubernetes.NewFromConfig(kubeconfig)

	if err != nil {
		return err
	}

	dashboardName := "ingress-nginx-dashboard"

	if err := client.CoreV1().ConfigMaps(namespace).Delete(ctx, dashboardName, metav1.DeleteOptions{}); err != nil {
		// return err
	}

	return nil
}
