package dashboard

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/devkube/pkg/kubectl"
)

const (
	dashboardRepo = "https://kubernetes.github.io/dashboard"

	dashboard        = "dashboard"
	dashboardChart   = "kubernetes-dashboard"
	dashboardVersion = "5.10.0"
)

func Install(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	values := map[string]any{
		"nameOverride":     dashboard,
		"fullnameOverride": dashboard,

		"extraArgs": []string{
			"--enable-skip-login",
			"--enable-insecure-login",
			"--disable-settings-authorizer",
		},

		"protocolHttp": true,

		"service": map[string]any{
			"externalPort": 80,
		},

		"serviceMonitor": map[string]any{
			"enabled": true,
		},

		"metricsScraper": map[string]any{
			"enabled": true,
		},

		"resources": nil,
	}

	if err := helm.Install(ctx, dashboard, dashboardRepo, dashboardChart, dashboardVersion, values, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace), helm.WithDefaultOutput()); err != nil {
		return err
	}

	if err := kubectl.Invoke(ctx, []string{"delete", "clusterrolebinding", dashboard}, kubectl.WithKubeconfig(kubeconfig)); err != nil {
		// ignore error
	}

	if err := kubectl.Invoke(ctx, []string{"create", "clusterrolebinding", dashboard, "--clusterrole=cluster-admin", "--serviceaccount=" + namespace + ":" + dashboard}, kubectl.WithKubeconfig(kubeconfig)); err != nil {
		return err
	}

	return nil
}

func Uninstall(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	if err := kubectl.Invoke(ctx, []string{"delete", "clusterrolebinding", dashboard}, kubectl.WithKubeconfig(kubeconfig)); err != nil {
		//return err
	}

	if err := helm.Uninstall(ctx, dashboard, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace)); err != nil {
		//return err
	}

	return nil
}
