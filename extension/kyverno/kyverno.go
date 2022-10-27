package kyverno

import (
	"context"
	_ "embed"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/devkube/pkg/kubernetes"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	//go:embed dashboard.json
	dashboard string
)

const (
	kyvernoRepo = "https://kyverno.github.io/kyverno"

	kyverno        = "kyverno"
	kyvernoChart   = "kyverno"
	kyvernoVersion = "v2.5.3"

	policies        = "kyverno-policies"
	policiesChart   = "kyverno-policies"
	policiesVersion = "v2.5.5"
)

func Install(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	kyvernoValues := map[string]any{
		"excludeKyvernoNamespace": true,

		"resourceFiltersExcludeNamespaces": []string{
			namespace,
			"kube-system",
			"cert-manager",
			"linkerd",
			"linkerd-viz",
			"linkerd-jaeger",
		},

		"serviceMonitor": map[string]any{
			"enabled": true,
		},
	}

	if err := helm.Install(ctx, kyverno, kyvernoRepo, kyvernoChart, kyvernoVersion, kyvernoValues, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace), helm.WithDefaultOutput()); err != nil {
		return err
	}

	policiesValues := map[string]any{
		"failurePolicy": "Ignore",

		"validationFailureAction": "audit",
	}

	if err := helm.Install(ctx, policies, kyvernoRepo, policiesChart, policiesVersion, policiesValues, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace), helm.WithDefaultOutput()); err != nil {
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

	if err := helm.Uninstall(ctx, kyverno, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace)); err != nil {
		// return err
	}

	if err := helm.Uninstall(ctx, policies, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace)); err != nil {
		// return err
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

	dashboardName := "kyverno-dashboard"

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

	dashboardName := "kyverno-dashboard"

	if err := client.CoreV1().ConfigMaps(namespace).Delete(ctx, dashboardName, metav1.DeleteOptions{}); err != nil {
		// return err
	}

	return nil
}
