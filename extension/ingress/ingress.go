package ingress

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
)

var (
	ingressRepo = "https://kubernetes.github.io/ingress-nginx"

	ingress        = "ingress-nginx"
	ingressChart   = "ingress-nginx"
	ingressVersion = "4.2.5"

	Images = []string{}
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

	return nil
}

func Uninstall(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	if err := helm.Uninstall(ctx, ingress, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace)); err != nil {
		//return err
	}

	return nil
}
