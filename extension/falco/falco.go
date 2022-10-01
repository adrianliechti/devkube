package falco

import (
	"context"
	"errors"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/devkube/pkg/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	falco        = "falco"
	falcoChart   = "falco"
	falcoVersion = "2.0.18"

	exporter        = "falco-exporter"
	exporterChart   = "falco-exporter"
	exporterVersion = "0.8.2"

	falcoRepo = "https://falcosecurity.github.io/charts"
)

func Install(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	if err := verifyNodes(ctx, kubeconfig); err != nil {
		return err
	}

	if err := installFalco(ctx, kubeconfig, namespace); err != nil {
		return err
	}

	if err := installExporter(ctx, kubeconfig, namespace); err != nil {
		return err
	}

	return nil
}

func Uninstall(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	if err := uninstallFalco(ctx, kubeconfig, namespace); err != nil {
		// return err
	}

	if err := uninstallExporter(ctx, kubeconfig, namespace); err != nil {
		// return err
	}

	return nil
}

func verifyNodes(ctx context.Context, kubeconfig string) error {
	return nil
	client, err := kubernetes.NewFromConfig(kubeconfig)

	if err != nil {
		return err
	}

	nodes, err := client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})

	if err != nil {
		return err
	}

	for _, node := range nodes.Items {
		os := node.Labels["kubernetes.io/os"]
		arch := node.Labels["kubernetes.io/arch"]

		if os == "linux" && arch == "amd64" {
			return nil
		}
	}

	return errors.New("falco is currently only supported on linux/amd64 nodes")
}

func installFalco(ctx context.Context, kubeconfig, namespace string) error {
	values := map[string]any{
		"falco": map[string]any{
			"grpc": map[string]any{
				"enabled": true,
			},

			"grpc_output": map[string]any{
				"enabled": true,
			},
		},
	}

	if err := helm.Install(ctx, falco, falcoRepo, falcoChart, falcoVersion, values, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace), helm.WithDefaultOutput()); err != nil {
		return err
	}

	return nil
}

func uninstallFalco(ctx context.Context, kubeconfig, namespace string) error {
	if err := helm.Uninstall(ctx, falco, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace)); err != nil {
		// return err
	}

	return nil
}

func installExporter(ctx context.Context, kubeconfig, namespace string) error {
	values := map[string]any{
		"serviceMonitor": map[string]any{
			"enabled": true,
		},

		"grafanaDashboard": map[string]any{
			"enabled":                  true,
			"namespace":                nil,
			"prometheusDatasourceName": "prometheus",
		},

		"prometheusRules": map[string]any{
			"enabled": true,
		},
	}

	if err := helm.Install(ctx, exporter, falcoRepo, exporterChart, exporterVersion, values, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace), helm.WithDefaultOutput()); err != nil {
		return err
	}

	return nil
}

func uninstallExporter(ctx context.Context, kubeconfig, namespace string) error {
	if err := helm.Uninstall(ctx, exporter, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace)); err != nil {
		// return err
	}

	return nil
}
