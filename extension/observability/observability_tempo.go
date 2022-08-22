package observability

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/devkube/pkg/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	tempo        = "tempo"
	tempoChart   = "tempo"
	tempoVersion = "0.15.8"
)

func installTempo(ctx context.Context, kubeconfig, namespace string) error {
	values := map[string]any{
		"persistence": map[string]any{
			"enabled": true,
			"size":    "10Gi",
		},

		"tempoQuery": map[string]any{
			"enabled": false,
		},
	}

	if err := helm.Install(ctx, tempo, grafanaRepo, tempoChart, tempoVersion, values, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace), helm.WithDefaultOutput()); err != nil {
		return err
	}

	return nil
}

func uninstallTempo(ctx context.Context, kubeconfig, namespace string) error {
	client, err := kubernetes.NewFromConfig(kubeconfig)

	if err != nil {
		return err
	}

	if err := helm.Uninstall(ctx, tempo, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace)); err != nil {
		//return err
	}

	client.CoreV1().PersistentVolumeClaims(namespace).DeleteCollection(ctx, metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/instance=" + tempo,
	})

	return nil
}
