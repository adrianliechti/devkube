package ingress

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
)

var (
	ingressRepo = "https://kubernetes.github.io/ingress-nginx"

	ingress        = "ingress-nginx"
	ingressChart   = "ingress-nginx"
	ingressVersion = "4.2.1"

	Images = []string{}
)

func Install(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	values := map[string]any{
		"controller": map[string]any{
			"service": map[string]any{
				"type": "NodePort",
			},

			"metrics": map[string]any{
				"enabled": true,
			},

			"serviceMonitor": map[string]any{
				"enabled": true,
			},
		},
	}

	if err := helm.Install(ctx, kubeconfig, namespace, ingress, ingressRepo, ingressChart, ingressVersion, values); err != nil {
		return err
	}

	return nil
}

func Uninstall(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	if err := helm.Uninstall(ctx, kubeconfig, namespace, ingress); err != nil {
		//return err
	}

	return nil
}
