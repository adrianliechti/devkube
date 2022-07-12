package dashboard

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/devkube/pkg/kubectl"
)

var (
	dashboard = "dashboard"

	chartRepo    = "https://kubernetes.github.io/dashboard"
	chartName    = "kubernetes-dashboard"
	chartVersion = "5.4.1"

	Images = []string{
		"kubernetesui/dashboard:v2.5.1",
		"kubernetesui/metrics-scraper:v1.0.7",
		"k8s.gcr.io/metrics-server/metrics-server:v0.5.0",
	}
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

		"metricsScraper": map[string]any{
			"enabled": true,
		},

		"metrics-server": map[string]any{
			"enabled": true,
			"args": []string{
				"--kubelet-insecure-tls",
				"--kubelet-preferred-address-types=InternalIP",
			},
		},
	}

	if err := helm.Install(ctx, kubeconfig, namespace, dashboard, chartRepo, chartName, chartVersion, values); err != nil {
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
		return err
	}

	if err := helm.Uninstall(ctx, kubeconfig, namespace, dashboard); err != nil {
		return err
	}

	return nil
}
