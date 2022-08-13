package dashboard

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/devkube/pkg/kubectl"
)

var (
	dashboardRepo = "https://kubernetes.github.io/dashboard"

	dashboard        = "dashboard"
	dashboardChart   = "kubernetes-dashboard"
	dashboardVersion = "5.8.0"
)

func Install(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	values := map[string]any{
		"nameOverride": dashboard,

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

	if err := helm.Install(ctx, kubeconfig, namespace, dashboard, dashboardRepo, dashboardChart, dashboardVersion, values); err != nil {
		return err
	}

	kubectl.Invoke(ctx, kubeconfig, "delete", "clusterrolebinding", dashboard)

	if err := kubectl.Invoke(ctx, kubeconfig, "create", "clusterrolebinding", dashboard, "--clusterrole=cluster-admin", "--serviceaccount="+namespace+":"+dashboard); err != nil {
		return err
	}

	return nil
}

func Uninstall(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	if err := kubectl.Invoke(ctx, kubeconfig, "delete", "clusterrolebinding", dashboard); err != nil {
		//return err
	}

	if err := helm.Uninstall(ctx, kubeconfig, namespace, dashboard); err != nil {
		//return err
	}

	return nil
}
