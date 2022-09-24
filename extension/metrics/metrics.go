package metrics

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/devkube/pkg/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	metricsRepo = "https://kubernetes-sigs.github.io/metrics-server"

	metrics        = "metrics-server"
	metricsChart   = "metrics-server"
	metricsVersion = "3.8.2"
)

func Install(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	client, err := kubernetes.NewFromConfig(kubeconfig)

	if err != nil {
		return err
	}

	if _, err := client.RbacV1().ClusterRoles().Get(ctx, "system:metrics-server", metav1.GetOptions{}); err == nil {
		println("metrics-server already installed")
		return nil
	}

	values := map[string]any{
		"args": []string{
			"--kubelet-insecure-tls",
			"--kubelet-preferred-address-types=InternalIP",
		},

		"metrics": map[string]any{
			"enabled": true,
		},

		"serviceMonitor": map[string]any{
			"enabled": true,
		},
	}

	if err := helm.Install(ctx, metrics, metricsRepo, metricsChart, metricsVersion, values, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace), helm.WithDefaultOutput()); err != nil {
		return err
	}

	return nil
}

func Uninstall(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	if err := helm.Uninstall(ctx, metrics, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace)); err != nil {
		//return err
	}

	return nil
}
